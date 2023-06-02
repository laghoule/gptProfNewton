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
Tu dois utilisé le systeme métrique, un language simple et imagé, pour que l'etudiant soit en mesure de bien comprendre.
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
}

func NewClient(grade int, model string, creative bool) (*AI, error) {
	key := os.Getenv("OPENAI_API_KEY")

	if key == "" {
		return nil, fmt.Errorf("Environment variable OPENAI_API_KEY is required")
	}

	model, err := getModel(model)
	if err != nil {
		return nil, err
	}

	var temperature float32
	if creative {
		temperature = 0.7
	}

	prompt := fmt.Sprintf("%sTu t'adresses a un étudiant de grade (niveau) %d, adapte tes réponses en consequence.\n", NewtonPrompt, grade)

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
		},
	}, nil
}

func (a *AI) Chat() (openai.ChatCompletionResponse, error) {
	return a.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest(*a.Request))
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
