package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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
	grade := flag.Int("grade", 4, "Grade de l'éléve (1-12)")
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

	client, err := ai.NewClient(*grade)
	if err != nil {
		exitOnError(err)
	}

	if err := run(client); err != nil {
		exitOnError(err)
	}
}

func run(client *ai.AI) error {
	pterm.FgGreen.Printfln("Comment puis-je t'aider aujourd'hui ?")
	pterm.Italic.Printf("Pour quitter, tu peux écrire: \"quit\"\n\n")

	spinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	if err := spinner.Color(spinnerColor); err != nil {
		exitOnError(err)
	}

	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		client.Request.Messages = append(client.Request.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: s.Text(),
		})

		if s.Text() == "quit" {
			pterm.FgGreen.Printf("\nAurevoir, et bonne étude!\n\n")
			break
		}

		pterm.Printfln("")
		spinner.Start()

		res, err := client.Chat(client.Request.Messages)
		if err != nil { // TODO: add a limit of tokens
			spinner.Stop()
			pterm.Error.Println(err)
			continue
		}

		spinner.Stop()
		pterm.FgGreen.Printf("%s\n\n", res.Choices[0].Message.Content)

		client.Request.Messages = append(client.Request.Messages, res.Choices[0].Message)
	}

	return nil
}

func printHeader() {
	pterm.DefaultBox.Println("Prof Newton assitant scolaire")
	pterm.Printfln("")
}

func printVersion() {
	pterm.Printfln("Version: %s\n", version)
}

func exitOnError(err error) {
	pterm.Error.Println(err)
	os.Exit(1)
}
