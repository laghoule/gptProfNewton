package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/laghoule/gptProfNewton/internal/pkg/ai"
	"github.com/pterm/pterm"
	"github.com/sashabaranov/go-openai"
)

var (
	version = "dev"
)

func main() {
	creative := flag.Bool("creative", false, "Utiliser le modele creatif")
	debug := flag.Bool("debug", false, "Activer le mode debug")
	grade := flag.Int("grade", 4, "Grade de l'éléve (1-12)")
	model := flag.String("model", "gpt-3.5", "Modéle de l'API d'OpenAI (gpt-3.5, gpt-4o, o3-mini)")
	stream := flag.Bool("stream", true, "Activer le mode streaming")
	studentName := flag.String("studentName", "", "Nom de l'éléve")
	version := flag.Bool("version", false, "Afficher la version")
	flag.Parse()

	if *version {
		printVersion()
		return
	}

	if *grade < 1 || *grade > 12 {
		exitOnError(fmt.Errorf("Vous devez choisir un grade entre 1 et 12)"))
	}

	conf := AI.Config{
		Debug:    *debug,
		Creative: *creative,
		Stream:   *stream,
		Grade:    *grade,
		Model:    *model,
	}

	printHeader()

	ai, err := AI.NewClient(*studentName, conf)
	if err != nil {
		exitOnError(err)
	}

	if err := run(ai); err != nil {
		exitOnError(err)
	}
}

func run(ai *AI.AI) error {
	pterm.FgGreen.Printfln("Comment puis-je t'aider aujourd'hui?")
	pterm.Italic.Printf("Pour quitter /quit, pour reinitiliser /reset\n\n")

	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT)

		ai.Request.Messages = append(ai.Request.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: s.Text(),
		})

		switch stripSpace(s.Text()) {
		case "":
			continue
		case "/quit":
			pterm.FgGreen.Printf("\nAurevoir, et bonne étude!\n\n")
			return nil
		case "/reset":
			pterm.FgLightGreen.Printf("\nRéinitialisation de la conversation\n\n")
			ai.Reset()
			continue
		}

		pterm.Println()

		if err := chat(ctx, ai); err != nil {
			return err
		}
	}

	return nil
}

func chat(ctx context.Context, ai *AI.AI) error {
	switch ai.Request.Stream {
	case true:
		if err := printChatStream(ctx, ai); err != nil {
			return err
		}
	case false:
		if err := printChat(ctx, ai); err != nil {
			return err
		}
	}
	return nil
}

func stripSpace(str string) string {
	return strings.ReplaceAll(str, " ", "")
}

func canceledMessage(ai *AI.AI, debug bool) {
	pterm.FgLightGreen.Printf("\n\nMessage annulé")
	ai.CancelLastMessage()
	if debug {
		printMsg(ai.Request.Messages)
	}
}

func exitOnError(err error) {
	pterm.Error.Println(err)
	os.Exit(1) // TODO: return AIError code instead of 1
}
