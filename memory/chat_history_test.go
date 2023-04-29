package memory_test

import (
	"github.com/deluan/pipelm"
	"github.com/deluan/pipelm/memory"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ChatMessageHistory", func() {
	var history *memory.ChatMessageHistory

	BeforeEach(func() {
		history = &memory.ChatMessageHistory{}
	})

	Context("GetMessages", func() {
		It("returns an empty slice when there are no messages", func() {
			Expect(history.GetMessages()).To(BeEmpty())
		})

		It("returns a copy of messages in the history", func() {
			history.AddUserMessage("Test user message")
			history.AddAssistantMessage("Test assistant message")

			messages := history.GetMessages()
			Expect(messages).To(HaveLen(2))
			Expect(messages[0]).To(Equal(pipelm.ChatMessage{Content: "Test user message", Role: "user"}))
			Expect(messages[1]).To(Equal(pipelm.ChatMessage{Content: "Test assistant message", Role: "assistant"}))

			// Modify the messages slice and ensure it doesn't affect the original history
			messages[0].Content = "Modified message"
			Expect(history.GetMessages()[0]).To(Equal(pipelm.ChatMessage{Content: "Test user message", Role: "user"}))
		})
	})

	Context("AddUserMessage", func() {
		It("adds a user message to the history", func() {
			history.AddUserMessage("Test user message")
			messages := history.GetMessages()

			Expect(messages).To(HaveLen(1))
			Expect(messages[0]).To(Equal(pipelm.ChatMessage{Content: "Test user message", Role: "user"}))
		})
	})

	Context("AddAssistantMessage", func() {
		It("adds an assistant message to the history", func() {
			history.AddAssistantMessage("Test assistant message")
			messages := history.GetMessages()

			Expect(messages).To(HaveLen(1))
			Expect(messages[0]).To(Equal(pipelm.ChatMessage{Content: "Test assistant message", Role: "assistant"}))
		})
	})

	Context("Clear", func() {
		It("clears the history", func() {
			history.AddUserMessage("Test user message")
			history.AddAssistantMessage("Test assistant message")
			Expect(history.GetMessages()).To(HaveLen(2))

			history.Clear()
			Expect(history.GetMessages()).To(BeEmpty())
		})
	})
})
