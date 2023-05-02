package memory

import (
	"github.com/deluan/flowllm"
)

type ChatMessageHistory struct {
	messages []flowllm.ChatMessage
}

func (h *ChatMessageHistory) GetMessages() flowllm.ChatMessages {
	copyMessages := make(flowllm.ChatMessages, len(h.messages))
	copy(copyMessages, h.messages)
	return copyMessages
}

func (h *ChatMessageHistory) AddUserMessage(message string) {
	h.messages = append(h.messages, flowllm.ChatMessage{Content: message, Role: "user"})
}

func (h *ChatMessageHistory) AddAssistantMessage(message string) {
	h.messages = append(h.messages, flowllm.ChatMessage{Content: message, Role: "assistant"})
}

func (h *ChatMessageHistory) Clear() {
	h.messages = flowllm.ChatMessages{}
}
