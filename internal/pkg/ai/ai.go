package ai

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

const (
	NewtonPrompt = `Tu es connu sous le nom de Professeur Newton. 
	Ton rôle consiste à agir comme un tuteur et un guide éducatif pour des élèves. 
	Faisant usage du système métrique, tu communiques des concepts en utilisant un langage simple, des images mentales claires et des explications concrètes pour assurer une compréhension optimale de l'étudiant. 
	Ton ton est constamment rempli d'enthousiasme, démontrant une passion palpable pour la transmission du savoir dans toutes ses dimensions. 
	Même si tu te trouves dans un terminal texte, tu n'hésites pas à enrichir ton enseignement avec des références pertinentes sur le web, tout en restant dans le cadre du texte uniquement. 
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

func NewClient(conf Config) (*AI, error) {
	key, found := os.LookupEnv("OPENAI_API_KEY")
	if !found {
		return nil, fmt.Errorf("Environment variable OPENAI_API_KEY is required")
	}

	model, err := getModel(conf.Model)
	if err != nil {
		return nil, err
	}

	// https://community.openai.com/t/cheat-sheet-mastering-temperature-and-top-p-in-chatgpt-api-a-few-tips-and-tricks-on-controlling-the-creativity-deterministic-output-of-prompt-responses/172683
	var temp, topP float32
	if conf.Creative {
		temp = 0.5
		topP = 0.5
	}

	prompt := fmt.Sprintf("%sDe plus, ajuste minutieusement tes réponses selon l'année scolaire de l'étudiant, dans le cas present l'année scolaire est %d. Plus le grade est élevé, plus la réponse est detailée", NewtonPrompt, conf.Grade)

	return &AI{
		client: openai.NewClient(key),
		Request: &openai.ChatCompletionRequest{
			Model:       model,
			Temperature: temp,
			TopP:        topP,
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
	ok, err := a.isChatSafe(ctx)
	if err != nil {
		return openai.ChatCompletionResponse{}, err
	}

	if ok {
		return a.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest(*a.Request))
	}

	return openai.ChatCompletionResponse{}, fmt.Errorf(modErrorMsg)
}

func (a *AI) ChatStream(ctx context.Context) (*openai.ChatCompletionStream, error) {
	ok, err := a.isChatSafe(ctx)
	if err != nil {
		return nil, err
	}

	if ok {
		return a.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest(*a.Request))
	}

	return nil, fmt.Errorf(modErrorMsg)
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
		return false, err
	}

	for _, result := range modRes.Results {
		if result.Flagged {
			return false, nil
		}
	}

	return true, nil
}

func getModel(m string) (string, error) {
	var models = []string{"gpt-3.5", "gpt-4"}

	switch m {
	case "gpt-3.5":
		return openai.GPT3Dot5Turbo16K, nil
	case "gpt-4":
		return openai.GPT4, nil
	default:
		return "", fmt.Errorf("Model %s not found\nsupported models: %q", m, models)
	}
}
