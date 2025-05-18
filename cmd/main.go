package main

import (
	"bufio"
	"context"
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/laghoule/gptProfNewton/internal/pkg/ai"
	"github.com/laghoule/gptProfNewton/internal/pkg/config"
	"github.com/pterm/pterm"
	"github.com/sashabaranov/go-openai"
)

var (
	version = "dev"
)

func main() {
	debug := flag.Bool("debug", false, "Activer le mode debug")
	cfgFile := flag.String("config", "", "Chemin vers le fichier de configuration")
	version := flag.Bool("version", false, "Afficher la version")
	flag.Parse()

	if *version {
		printVersion()
		return
	}

	cfg, err := config.New(*cfgFile)
	if err != nil {
		exitOnError(err)
	}

	printHeader()

	ai, err := AI.NewClient(cfg, *debug)
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

		if err := printChatStream(ctx, ai); err != nil {
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
