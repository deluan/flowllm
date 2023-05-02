package memory

import (
	"context"

	"github.com/deluan/flowllm"
)

type Buffer struct {
	chatHistory *ChatMessageHistory
	windowSize  int
}

func NewBuffer(windowSize int, history *flowllm.ChatMessages) *Buffer {
	chatHistory := &ChatMessageHistory{}
	if history != nil {
		chatHistory.messages = *history
	}
	return &Buffer{windowSize: windowSize, chatHistory: chatHistory}
}

func (b *Buffer) Load(_ context.Context) (flowllm.ChatMessages, error) {
	messages := b.chatHistory.GetMessages()
	if b.windowSize > 0 {
		messages = messages.Last(b.windowSize * 2)
	}
	return messages, nil
}

func (b *Buffer) Save(_ context.Context, input, output string) error {
	b.chatHistory.AddUserMessage(input)
	b.chatHistory.AddAssistantMessage(output)
	return nil
}
