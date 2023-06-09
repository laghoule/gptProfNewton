package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/laghoule/gptProfNewton/internal/pkg/ai"
	"github.com/sashabaranov/go-openai"

	"github.com/briandowns/spinner"
	"github.com/pterm/pterm"
)

const (
	spinnerColor string = "fgGreen"
)

var (
	version = "dev"
)

func main() {
	creative := flag.Bool("creative", false, "Utiliser le modele creatif")
	debug := flag.Bool("debug", false, "Activer le mode debug")
	grade := flag.Int("grade", 4, "Grade de l'éléve (1-12)")
	model := flag.String("model", "gpt-3.5", "Modéle de l'API d'OpenAI")
	stream := flag.Bool("stream", true, "Activer le mode streaming")
	version := flag.Bool("version", false, "Afficher la version")
	flag.Parse()

	if *version {
		printVersion()
		return
	}

	if *grade < 1 || *grade > 12 {
		exitOnError(fmt.Errorf("Vous devez choisir un grade entre 1 et 12)"))
	}

	printHeader()

	ai, err := ai.NewClient(*grade, *model, *stream, *creative)
	if err != nil {
		exitOnError(err)
	}

	if err := run(ai, *debug); err != nil {
		exitOnError(err)
	}
}

func run(ai *ai.AI, debug bool) error {
	pterm.FgGreen.Printfln("Comment puis-je t'aider aujourd'hui ?")
	pterm.Italic.Printf("Pour quitter [quit], pour reinitiliser [reset]\n\n")

	spinner := spinner.New(spinner.CharSets[2], 100*time.Millisecond)
	if err := spinner.Color(spinnerColor); err != nil {
		exitOnError(err)
	}

	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT)

		ai.Request.Messages = append(ai.Request.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: s.Text(),
		})

		switch s.Text() {
		case "quit":
			pterm.FgGreen.Printf("\nAurevoir, et bonne étude!\n\n")
			return nil
		case "reset":
			pterm.FgLightGreen.Printf("\nReinitialisation de la conversation\n\n")
			ai.Reset()
			continue
		}

		switch ai.Request.Stream {
		case true:
			if err := printChatStream(ctx, ai, debug); err != nil {
				return err
			}
		case false:
			if err := printChat(ctx, ai, debug); err != nil {
				return err
			}
		}
	}

	return nil
}

func printHeader() {
	pterm.DefaultBox.Println("Prof Newton assitant scolaire")
	pterm.Printfln("")
}

func printChatStream(ctx context.Context, ai *ai.AI, debug bool) error {
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
				pterm.FgLightGreen.Printf("\n\nMessage annulé")
				ai.CancelLastMessage()
				if debug {
					printMsg(ai.Request.Messages)
				}
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

	if debug {
		fmt.Printf("\n\n")
		printMsg(ai.Request.Messages)
	} else {
		pterm.Printf("\n\n")
	}

	return nil
}

func printChat(ctx context.Context, ai *ai.AI, debug bool) error {
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
			pterm.FgLightGreen.Printf("Message annulé\n\n")
			ai.CancelLastMessage()
			if debug {
				printMsg(ai.Request.Messages)
			}
			return nil
		}
		return err
	}

	spinner.Stop()
	pterm.FgGreen.Printf("%s\n\n", res.Choices[0].Message.Content)
	ai.Request.Messages = append(ai.Request.Messages, res.Choices[0].Message)

	if debug {
		printMsg(ai.Request.Messages)
	}

	return nil
}

func printMsg(msgs []openai.ChatCompletionMessage) {
	for _, m := range msgs {
		pterm.FgLightBlue.Printf("Role: %s\n", m.Role)
		pterm.FgLightBlue.Printf("Content: %s\n", m.Content)
		pterm.Println()
	}
}

func printVersion() {
	pterm.Printfln("Version: %s\n", version)
}

func exitOnError(err error) {
	pterm.Error.Println(err)
	os.Exit(1)
}
