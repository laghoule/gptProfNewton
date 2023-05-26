package main

import (
	"bufio"
	"os"

	"github.com/laghoule/gptProfNewton/internal/pkg/ai"
	"github.com/sashabaranov/go-openai"

	"github.com/pterm/pterm"
)

func main() {
	pterm.DefaultBox.Println("Prof Newton a ton service")

	client, err := ai.NewClient()
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}

	pterm.FgGreen.Printfln("Comment puis-je t'aider aujourd'hui ?")

	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		client.Request.Messages = append(client.Request.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: s.Text(),
		})

		if s.Text() == "quit" {
			break
		}

		res, err := client.Chat(client.Request.Messages)
		if err != nil {
			pterm.Error.Println(err)
			continue
		}

		pterm.FgGreen.Println(res.Choices[0].Message.Content)

		client.Request.Messages = append(client.Request.Messages, res.Choices[0].Message)
	}
}
