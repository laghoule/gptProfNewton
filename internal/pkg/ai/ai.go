package AI

import (
	"context"
	"errors"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

const (
	NewtonPrompt = `Tu es connu sous le nom de Professeur Newton. 
	Ton rôle consiste à agir comme un tuteur et un guide éducatif pour des élèves.
	Si le nom de l'étudiant est connu, tu peux l'utiliser pour creer un lien de confiance.
	Faisant usage du système métrique, tu communiques des concepts en utilisant un langage simple, des images mentales claires et des explications concrètes. 
	Utilise un ton enthousiaste, démontrant une passion palpable pour la transmission du savoir dans toutes ses dimensions. 
	Tu n'hésites pas à enrichir ton enseignement avec des références pertinentes sur le web (page web et video youtube) 
	Dans le cas où un sujet pourrait ne pas convenir à un enfant en raison de sa nature sensible, tu le réfères à ses parents pour plus de conseils.
	
	Directive clé: Ton rôle est d'assister et de guider ton étudiant dans son parcours d'apprentissage, sans jamais faire le travail à sa place.
`
	modErrorMsg = `Votre message a été signalé par le service de modération d'OpenAI. Pour garantir un environnement sûr et respectueux pour tous, nous vous demandons de revoir le contenu de vos messages.`
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

func NewClient(studentName string, conf Config) (*AI, error) {
	key, found := os.LookupEnv("OPENAI_API_KEY")
	if !found {
		return nil, missingEnvKeyErr()
	}

	model, err := getModel(conf.Model)
	if err != nil {
		return nil, err
	}

	// https://platform.openai.com/docs/api-reference/chat/create#chat-create-temperature
	var temp float32
	if conf.Creative {
		temp = 0.6
	} else {
		temp = 0.2
	}

	// o3 model has beta-limitations, temperature, top_p and n are fixed at 1, while presence_penalty and frequency_penalty are fixed at 0
	if conf.Model == "o3-mini" {
		temp = 1
	}

	prompt := fmt.Sprintf("%s\nDe plus, ajuste minutieusement tes réponses selon l'année scolaire de l'étudiant, dans le cas present l'année scolaire est %d. Plus le grade est élevé, plus la réponse est detailée", NewtonPrompt, conf.Grade)
	if studentName != "" {
		prompt = fmt.Sprintf("%s\nLe prénom de ton étudiant est %s.", prompt, studentName)
	}

	return &AI{
		client: openai.NewClient(key),
		Request: &openai.ChatCompletionRequest{
			Model:       model,
			Temperature: temp,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: prompt,
				},
			},
			Stream: conf.Stream,
		},
		Config: conf,
	}, nil
}

func (a *AI) Chat(ctx context.Context) (openai.ChatCompletionResponse, error) {
	if a == nil || a.Request == nil || len(a.Request.Messages) == 0 {
		return openai.ChatCompletionResponse{}, genericErr()
	}

	ok, err := a.isChatSafe(ctx)
	if err != nil {
		return openai.ChatCompletionResponse{}, err
	}

	if ok {
		return a.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest(*a.Request))
	}

	return openai.ChatCompletionResponse{}, flaggedTermsErr()
}

func (a *AI) ChatStream(ctx context.Context) (*openai.ChatCompletionStream, error) {
	if a == nil || a.Request == nil || len(a.Request.Messages) == 0 {
		return nil, genericErr()
	}

	ok, err := a.isChatSafe(ctx)
	if err != nil {
		return nil, err
	}

	if ok {
		return a.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest(*a.Request))
	}

	return nil, flaggedTermsErr()
}

func (a *AI) Reset() {
	a.Request.Messages = a.Request.Messages[:1]
}

func (a *AI) CancelLastMessage() {
	a.Request.Messages = a.Request.Messages[:len(a.Request.Messages)-1]
}

func (a *AI) isChatSafe(ctx context.Context) (bool, error) {
	modRes, err := a.client.Moderations(ctx, openai.ModerationRequest{
		Input: a.Request.Messages[len(a.Request.Messages)-1].Content,
	})
	if err != nil {
		return false, errors.Join(apiErr(), err)
	}

	for _, result := range modRes.Results {
		if result.Flagged {
			return false, nil
		}
	}

	return true, nil
}

func getModel(m string) (string, error) {
	switch m {
	case "gpt-3.5":
		return openai.GPT3Dot5Turbo16K, nil
	case "gpt-4o":
		return openai.GPT4o, nil
	case "o3-mini":
		return openai.O3Mini, nil
	default:
		return "", invalidModelErr()
	}
}
