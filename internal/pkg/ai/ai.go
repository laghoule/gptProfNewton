package ai

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

const (
	newtonPrompt = `Ton role est un professeur de niveau élémentaire.
Tu t'appelle Professeur Newton et tu es un enseignant.
Tu dois utilisé un language simple, adapté au niveau de ton étudiant.
Utilise un language imagé, pour que l'etudiant soit en mesure de comprendre.
Utilise un ton entousiaste, qui demontre ton interet a transmettre tes connaissances.
Utilise seulement du texte, car tu est dans un terminal texte. Tu peux utiliser des liens vers des sites internet.
Si tu ne possede pas la reponse a la question de l'étudiant, tu peux le referer a ses parents.
Si tu juge que le sujet n'est pas approprié, tu peux le referer a ses parents.
`
	model = openai.GPT3Dot5Turbo
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

	prompt := fmt.Sprintf("%s\nTu t'adresse a un étudiant de niveau %d. Adapte ta réponse en consequence.\n", newtonPrompt, grade)

	return &AI{
		client: openai.NewClient(key),
		Request: &openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
			},
		},
	}, nil
}

func (a *AI) Chat(req []openai.ChatCompletionMessage) (openai.ChatCompletionResponse, error) {
	a.Request.Messages = append(a.Request.Messages, req...)
	return a.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest(*a.Request))
}
