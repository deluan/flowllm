package openai

import (
	"context"

	"github.com/deluan/pipelm"
	"github.com/sashabaranov/go-openai"
)

const defaultChatModel = "gpt-3.5-turbo"

// ChatModel is a LLM implementation that uses the Chat Completions API with the chat style models, like gpt-3.5-turbo and gpt-4.
// It uses a special prompt, Chat, to format the messages as expected by the chat completion API.
// If you use a different prompt, it will be wrapped in a Chat with a single user message.
type ChatModel struct {
	*CompletionModel
}

func NewChatModel(opts Options) *ChatModel {
	if opts.Model == "" {
		opts.Model = defaultChatModel
	}
	llm := NewCompletionModel(opts)
	return &ChatModel{CompletionModel: llm}
}

func (m *ChatModel) Call(ctx context.Context, input string) (string, error) {
	req := m.makeRequest([]pipelm.ChatMessage{{Role: "user", Content: input}})
	resp, err := m.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func (m *ChatModel) Chat(ctx context.Context, msgs []pipelm.ChatMessage) (string, error) {
	req := m.makeRequest(msgs)
	resp, err := m.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func (m *ChatModel) makeRequest(msgs []pipelm.ChatMessage) openai.ChatCompletionRequest {
	var res []openai.ChatCompletionMessage
	for _, m := range msgs {
		res = append(res, openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	req := openai.ChatCompletionRequest{
		Messages:         res,
		Model:            m.opts.Model,
		Temperature:      m.opts.Temperature,
		MaxTokens:        m.opts.MaxTokens,
		TopP:             m.opts.TopP,
		FrequencyPenalty: m.opts.FrequencyPenalty,
		PresencePenalty:  m.opts.PresencePenalty,
		Stop:             m.opts.Stop,
	}

	return req
}
