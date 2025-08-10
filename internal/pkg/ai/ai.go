package AI

import (
	"context"
	"errors"
	"os"
	"regexp"

	"github.com/laghoule/gptProfNewton/internal/pkg/config"

	openai "github.com/sashabaranov/go-openai"
)

type AI struct {
	client  *openai.Client
	Request *openai.ChatCompletionRequest
	Config  *config.Config
	Debug   bool
}

func NewClient(conf *config.Config, debug bool) (*AI, error) {
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

	prompt := newPrompt(conf.Student.Name, conf.Student.Grade, conf.Student.Details)

	promptSafe, err := prompt.isPromptSafe(*client)
	if err != nil {
		return nil, errors.Join(apiErr(), err)
	}

	if !promptSafe {
		return nil, flaggedTermsErr()
	}

	return &AI{
		client: client,
		Request: &openai.ChatCompletionRequest{
			Model:       model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: prompt.String(),
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

func isThinkingModel(model string) bool {
	r := regexp.MustCompile(`^o[1-9]-.*`)
	return r.MatchString(model)
}
