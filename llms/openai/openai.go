package openai

import (
	"context"
	"errors"
	"os"

	"github.com/sashabaranov/go-openai"
)

const (
	defaultModel     = "text-ada-001"
	defaultMaxTokens = 256
)

// Options for OpenAI Completions models
type Options struct {
	ApiKey           string
	Model            string
	Temperature      float32
	MaxTokens        int
	TopP             float32
	FrequencyPenalty float32
	PresencePenalty  float32
	BestOf           int
	Stop             []string
}

// CompletionModel is a LLM implementation that uses the Completions API to generate text.
type CompletionModel struct {
	client *openai.Client
	opts   Options
}

var ErrMissingApiKey = errors.New("missing API KEY environment variable")

func NewCompletionModel(opts Options) *CompletionModel {
	if opts.Model == "" {
		opts.Model = defaultModel
	}
	if opts.MaxTokens == 0 {
		opts.MaxTokens = defaultMaxTokens
	}
	llm := &CompletionModel{opts: opts}

	if opts.ApiKey == "" {
		opts.ApiKey = os.Getenv("OPENAI_API_KEY")
	}
	llm.client = openai.NewClient(opts.ApiKey)
	return llm
}

func (m *CompletionModel) Call(ctx context.Context, input string) (string, error) {
	req := openai.CompletionRequest{
		Prompt:           input,
		Model:            m.opts.Model,
		Temperature:      m.opts.Temperature,
		MaxTokens:        m.opts.MaxTokens,
		TopP:             m.opts.TopP,
		FrequencyPenalty: m.opts.FrequencyPenalty,
		PresencePenalty:  m.opts.PresencePenalty,
		BestOf:           m.opts.BestOf,
		Stop:             m.opts.Stop,
	}
	resp, err := m.client.CreateCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Text, nil
}
