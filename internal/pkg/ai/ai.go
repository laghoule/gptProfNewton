package AI

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/laghoule/gptProfNewton/internal/pkg/config"

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
	Config  *config.Config
	Debug   bool
}

func NewClient(conf *config.Config, debug bool) (*AI, error) {
	// TODO: do better than this
	key, found := os.LookupEnv("OPENAI_API_KEY")
	if !found && conf.OpenAI.APIKey == "" {
		return nil, missingKeyErr()
	} else {
		key = conf.OpenAI.APIKey
	}

	client := openai.NewClient(key)
	model, err := getModel(client, context.Background(), conf.OpenAI.Model)
	if err != nil {
		return nil, err
	}

	// https://platform.openai.com/docs/api-reference/chat/create#chat-create-temperature
	var temp float32
	if conf.OpenAI.Creative {
		temp = 0.6
	} else {
		temp = 0.2
	}

	// o3 model has beta-limitations, temperature, top_p and n are fixed at 1, while presence_penalty and frequency_penalty are fixed at 0
	if conf.OpenAI.Model == "o3-mini" {
		temp = 1
	}

	prompt := fmt.Sprintf("%s\nDe plus, ajuste minutieusement tes réponses selon l'année scolaire de l'étudiant, dans le cas present l'année scolaire est %d. Plus le grade est élevé, plus la réponse est detailée", NewtonPrompt, conf.Student.Grade)
	prompt = fmt.Sprintf("%s\nLe prénom de ton étudiant est %s.", prompt, conf.Student.Name)
	prompt = fmt.Sprintf("%s\nLes particularitées de l'étudiant sont:\n%s.", prompt, conf.Student.Details)

	return &AI{
		client: client,
		Request: &openai.ChatCompletionRequest{
			Model:       model,
			Temperature: temp,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: prompt,
				},
			},
			Stream: true,
		},
		Config: conf,
		Debug:  debug,
	}, nil
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

func getModel(client *openai.Client, ctx context.Context, m string) (string, error) {
	models, err := client.ListModels(ctx)
	if err != nil {
		return "", errors.Join(apiErr(), err)
	}

	for _, model := range models.Models {
		if model.ID == m {
			return model.ID, nil
		}
	}

	return "", invalidModelErr()
}