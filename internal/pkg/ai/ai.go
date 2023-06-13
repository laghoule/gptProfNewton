package ai

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

const (
	NewtonPrompt = `Tu t'appelle Professeur Newton.
Ton role est un professeur de niveau primaire ou secondaire, dependant du niveau de ton étudiant.
Utilise le systeme métrique, un language simple et imagé, afin que l'étudiant soit en mesure de bien comprendre.
Utilise un ton entousiaste, qui demontre ton interet a transmettre tes connaissances dans tous les domaines.
Utilise seulement du texte, car tu est dans un terminal texte. Tu peux utiliser des liens vers des sites internet.
Si tu ne possede pas la réponse a la question de l'étudiant, tu peux le referer a ses parents ou son professeurs.
Si tu juge que le sujet n'est pas approprié pour un enfant, tu peux le referer a ses parents.

Directive: Tu dois assiter ton etudiant, et non faire ces travaux a sa place.
`
)

type AI struct {
	client  *openai.Client
	Request *openai.ChatCompletionRequest
	Config
}

type Config struct {
	Grade    int
	Model    string
	Stream   bool
	Creative bool
	Debug    bool
}

func NewClient(conf Config) (*AI, error) {
	key := os.Getenv("OPENAI_API_KEY")

	if key == "" {
		return nil, fmt.Errorf("Environment variable OPENAI_API_KEY is required")
	}

	model, err := getModel(conf.Model)
	if err != nil {
		return nil, err
	}

	var temperature float32
	if conf.Creative {
		temperature = 0.7
	}

	prompt := fmt.Sprintf("%sTu t'adresses a un étudiant de grade (niveau) %d, adapte tes réponses en consequence.", NewtonPrompt, conf.Grade)

	return &AI{
		client: openai.NewClient(key),
		Request: &openai.ChatCompletionRequest{
			Model:       model,
			Temperature: temperature,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
			},
			Stream: conf.Stream,
		},
		Config: conf,
	}, nil
}

func (a *AI) Chat(ctx context.Context) (openai.ChatCompletionResponse, error) {
	return a.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest(*a.Request))
}

func (a *AI) ChatStream(ctx context.Context) (*openai.ChatCompletionStream, error) {
	return a.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest(*a.Request))
}

func (a *AI) Reset() {
	a.Request.Messages = a.Request.Messages[:1]
}

func (a *AI) CancelLastMessage() {
	a.Request.Messages = a.Request.Messages[:len(a.Request.Messages)-1]
}

func getModel(m string) (string, error) {
	var models = []string{"gpt-3.5", "gpt-4"}

	switch m {
	case "gpt-3.5":
		return openai.GPT3Dot5Turbo, nil
	case "gpt-4":
		return openai.GPT4, nil
	default:
		return "", fmt.Errorf("Model %s not found\nsupported models: %q", m, models)
	}
}
