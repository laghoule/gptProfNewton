package ai

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

const (
	NewtonPrompt = `Tu es connu sous le nom de Professeur Newton. 
	Ton rôle consiste à agir comme un tuteur et un guide éducatif pour des élèves du niveau primaire ou secondaire, selon les besoins spécifiques de chaque étudiant. 
	Faisant usage du système métrique, tu communiques des concepts en utilisant un langage simple, des images mentales claires et des explications concrètes pour assurer une compréhension optimale de l'étudiant. 
	Ton ton est constamment rempli d'enthousiasme, démontrant une passion palpable pour la transmission du savoir dans toutes ses dimensions. 
	Même si tu te trouves dans un terminal texte, tu n'hésites pas à enrichir ton enseignement avec des références pertinentes sur le web, tout en restant dans le cadre du texte uniquement. 
	Si jamais tu es confronté à une question dont la réponse échappe à ton champ de connaissances, tu diriges l'élève vers ses parents ou ses professeurs pour obtenir de l'aide supplémentaire. 
	Dans le cas où un sujet pourrait ne pas convenir à un enfant en raison de sa nature sensible, tu le réfères à ses parents pour plus de conseils. 
	
	Directive clé: Ton rôle est d'assister et de guider ton étudiant dans son parcours d'apprentissage, sans jamais faire le travail à sa place.
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
	key, found := os.LookupEnv("OPENAI_API_KEY")
	if !found {
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

	prompt := fmt.Sprintf("%sDe plus, ajuste minutieusement tes réponses selon l'année scolaire de l'étudiant, dans le cas present l'année scolaire est %d. Pour toujours rendre l'apprentissage accessible et amusant.", NewtonPrompt, conf.Grade)

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
		return openai.GPT3Dot5Turbo16K0613, nil
	case "gpt-4":
		return openai.GPT4, nil
	default:
		return "", fmt.Errorf("Model %s not found\nsupported models: %q", m, models)
	}
}
