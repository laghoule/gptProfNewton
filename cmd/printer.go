package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/laghoule/gptProfNewton/internal/pkg/ai"

	"github.com/briandowns/spinner"
	"github.com/pterm/pterm"
	"github.com/sashabaranov/go-openai"
)

const (
	spinnerColor string = "fgGreen"
)

func printHeader() {
	pterm.DefaultBox.Println("Prof Newton assistant scolaire")
	pterm.Printfln("")
}

func printChat(ctx context.Context, ai *ai.AI) error {
	spinner := spinner.New(spinner.CharSets[2], 100*time.Millisecond)
	if err := spinner.Color(spinnerColor); err != nil {
		return fmt.Errorf("error while setting spinner color: %w", err)
	}

	pterm.Println()
	spinner.Start()

	res, err := ai.Chat(ctx)
	if err != nil {
		spinner.Stop()
		if ctx.Err() == context.Canceled {
			canceledMessage(ai, ai.Debug)
			return nil
		}
		return err
	}

	spinner.Stop()
	pterm.FgGreen.Printf("%s\n\n", res.Choices[0].Message.Content)
	ai.Request.Messages = append(ai.Request.Messages, res.Choices[0].Message)

	if ai.Debug {
		printMsg(ai.Request.Messages)
	}

	return nil
}

func printChatStream(ctx context.Context, ai *ai.AI) error {
	stream, err := ai.ChatStream(ctx)
	if err != nil {
		return err
	}
	defer stream.Close()

	pterm.Println()
	var streamedData []byte

	for {
		recv, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			if errors.Is(err, context.Canceled) {
				canceledMessage(ai, ai.Debug)
				break
			}
			return err
		}

		pterm.FgGreen.Printf("%s", recv.Choices[0].Delta.Content)
		streamedData = append(streamedData, recv.Choices[0].Delta.Content...)
	}

	ai.Request.Messages = append(ai.Request.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: string(streamedData),
	})

	if ai.Debug {
		pterm.Printf("\n\n")
		printMsg(ai.Request.Messages)
	} else {
		pterm.Printf("\n\n")
	}

	return nil
}

func printMsg(msgs []openai.ChatCompletionMessage) {
	for _, msg := range msgs {
		pterm.FgLightBlue.Printf("Role: %s\n", msg.Role)
		pterm.FgLightBlue.Printf("Content: %s\n", msg.Content)
		pterm.Println()
	}
}

func printVersion() {
	pterm.Printfln("Version: %s\n", version)
}
