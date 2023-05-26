package ai

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

const (
	newtonPrompt = `Ton role est un professeur de niveau élémentaire.
Tu dois utilisé un language simple, tout en prenant compte que tu t'adresse a des enfants.
Utile un language imager, pour que l'etudiant soit en mesure de comprendre.
`
	model = openai.GPT3Dot5Turbo
)

type Request openai.ChatCompletionRequest

type AI struct {
	client  *openai.Client
	Request *Request
}

func NewClient() (*AI, error) {
	key := os.Getenv("OPENAI_API_KEY")

	if key == "" {
		return nil, fmt.Errorf("Environment variable OPENAI_API_KEY is required")
	}

	return &AI{
		client: openai.NewClient(key),
		Request: &Request{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: newtonPrompt,
				},
			},
		},
	}, nil
}

func (a *AI) Chat(req string) (openai.ChatCompletionResponse, error) {
	ccRequest := openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: newtonPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: req,
			},
		},
	}
	return a.client.CreateChatCompletion(context.Background(), ccRequest)
}

func (a *AI) Chat2(req []openai.ChatCompletionMessage) (openai.ChatCompletionResponse, error) {
	a.Request.Messages = append(a.Request.Messages, req...)
	return a.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest(*a.Request))
}
