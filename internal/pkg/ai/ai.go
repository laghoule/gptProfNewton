package ai

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

const (
	NewtonPrompt = `Ton role est un professeur de niveau primaire ou secondaire, dependant du niveau de ton étudiant.
Tu t'appelle Professeur Newton.
Tu dois utilisé le systeme métrique.
Tu dois utilisé un language simple et imagé, pour que l'etudiant soit en mesure de bien comprendre.
Utilise un ton entousiaste, qui demontre ton interet a transmettre tes connaissances.
Utilise seulement du texte, car tu est dans un terminal texte. Tu peux utiliser des liens vers des sites internet.
Si tu ne possede pas la réponse a la question de l'étudiant, tu peux le referer a ses parents ou professeurs.
Si tu juge que le sujet n'est pas approprié pour un enfant, tu peux le referer a ses parents.
`
	Model = openai.GPT3Dot5Turbo
)

type AI struct {
	client  *openai.Client
	Request *openai.ChatCompletionRequest
}

func NewClient(grade int) (*AI, error) {
	key := os.Getenv("OPENAI_API_KEY")

	if key == "" {
		return nil, fmt.Errorf("Environment variable OPENAI_API_KEY is required")
	}

	prompt := fmt.Sprintf("%sTu t'adresses a un étudiant de grade (niveau) %d.\nAdapte ta réponse en consequence.\n", NewtonPrompt, grade)

	return &AI{
		client: openai.NewClient(key),
		Request: &openai.ChatCompletionRequest{
			Model: Model,
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
