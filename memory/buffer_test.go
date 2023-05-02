package memory_test

import (
	"context"
	"strconv"

	"github.com/deluan/flowllm"
	"github.com/deluan/flowllm/memory"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Buffer", func() {
	var ctx context.Context
	var buf *memory.Buffer

	BeforeEach(func() {
		ctx = context.Background()
	})

	It("saves user and assistant messages to chat history", func() {
		buf = memory.NewBuffer(0, nil)
		input := "User input message"
		output := "Assistant output message"
		err := buf.Save(ctx, input, output)
		Expect(err).NotTo(HaveOccurred())

		messages, err := buf.Load(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(messages).To(HaveLen(2))
		Expect(messages[0]).To(Equal(flowllm.ChatMessage{Content: "User input message", Role: "user"}))
		Expect(messages[1]).To(Equal(flowllm.ChatMessage{Content: "Assistant output message", Role: "assistant"}))
	})

	It("initializes chat history with ChatMessages", func() {
		msgs := flowllm.ChatMessages{
			{Content: "User input message 0", Role: "user"},
			{Content: "Assistant output message 0", Role: "assistant"},
		}
		buf = memory.NewBuffer(0, &msgs)
		err := buf.Save(ctx, "User input message 1", "Assistant output message 1")
		Expect(err).NotTo(HaveOccurred())

		messages, err := buf.Load(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(messages).To(HaveLen(4))
		Expect(messages[0]).To(Equal(flowllm.ChatMessage{Content: "User input message 0", Role: "user"}))
		Expect(messages[1]).To(Equal(flowllm.ChatMessage{Content: "Assistant output message 0", Role: "assistant"}))
		Expect(messages[2]).To(Equal(flowllm.ChatMessage{Content: "User input message 1", Role: "user"}))
		Expect(messages[3]).To(Equal(flowllm.ChatMessage{Content: "Assistant output message 1", Role: "assistant"}))
	})

	It("truncates history with windowSize", func() {
		buf = memory.NewBuffer(2, nil)
		for i := 1; i <= 3; i++ {
			err := buf.Save(ctx, "User message "+strconv.Itoa(i), "Assistant message "+strconv.Itoa(i))
			Expect(err).NotTo(HaveOccurred())
		}

		messages, err := buf.Load(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(messages).To(HaveLen(4))
		Expect(messages[0].Content).To(Equal("User message 2"))
		Expect(messages[1].Content).To(Equal("Assistant message 2"))
		Expect(messages[2].Content).To(Equal("User message 3"))
		Expect(messages[3].Content).To(Equal("Assistant message 3"))
	})
})
