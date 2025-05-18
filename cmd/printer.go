package main

import (
	"context"
	"errors"
	"io"

	"github.com/laghoule/gptProfNewton/internal/pkg/ai"

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

func printChatStream(ctx context.Context, ai *AI.AI) error {
	stream, err := ai.ChatStream(ctx)
	if err != nil {
		switch err {
		case AI.FlaggedTermsErr:
			ai.CancelLastMessage()
			pterm.Println()
			printFlaggedTerms()
			return nil
		case context.Canceled:
			canceledMessage(ai, ai.Debug)
			return nil
		}
		return err
	}
	defer stream.Close()

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
			return errors.Join(AI.ApiErr, err)
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

func printFlaggedTerms() {
	pterm.FgRed.Printf(AI.FlaggedTermsErrMsg)
	pterm.Italic.Printf("\nVotre dernier message a été annulé.\n\n")
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
