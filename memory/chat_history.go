package memory

import (
	"github.com/deluan/pipelm"
)

type ChatMessageHistory struct {
	messages []pipelm.ChatMessage
}

func (h *ChatMessageHistory) GetMessages() pipelm.ChatMessages {
	copyMessages := make(pipelm.ChatMessages, len(h.messages))
	copy(copyMessages, h.messages)
	return copyMessages
}

func (h *ChatMessageHistory) AddUserMessage(message string) {
	h.messages = append(h.messages, pipelm.ChatMessage{Content: message, Role: "user"})
}

func (h *ChatMessageHistory) AddAssistantMessage(message string) {
	h.messages = append(h.messages, pipelm.ChatMessage{Content: message, Role: "assistant"})
}

func (h *ChatMessageHistory) Clear() {
	h.messages = pipelm.ChatMessages{}
}
