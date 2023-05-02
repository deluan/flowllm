package flowllm_test

import (
	"context"
	"errors"
	"time"

	. "github.com/deluan/flowllm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handlers", func() {
	Describe("Chain", func() {
		It("should execute handlers in sequence and return the output of the last handler", func() {
			handler1 := HandlerFunc(func(ctx context.Context, values ...Values) (Values, error) {
				return Values{"key1": "value1"}, nil
			})
			handler2 := HandlerFunc(func(ctx context.Context, values ...Values) (Values, error) {
				vals := Values{}.Merge(values...)
				vals["key2"] = "value2"
				return vals, nil
			})
			chain := Chain(handler1, handler2)
			result, err := chain.Call(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(Values{"key1": "value1", "key2": "value2"}))
		})

		It("should return an error if any of the handlers returns an error", func() {
			handler1 := HandlerFunc(func(ctx context.Context, values ...Values) (Values, error) {
				return Values{"key1": "value1"}, nil
			})
			handler2 := HandlerFunc(func(ctx context.Context, values ...Values) (Values, error) {
				return nil, errors.New("handler error")
			})
			chain := Chain(handler1, handler2)
			_, err := chain.Call(context.Background())
			Expect(err).To(MatchError("handler error"))
		})
	})

	Describe("MapOutputTo", func() {
		It("should rename the output of the chain to the given key", func() {
			handler := HandlerFunc(func(ctx context.Context, values ...Values) (Values, error) {
				return Values{DefaultKey: "value"}, nil
			})
			chain := Chain(handler, MapOutputTo("newKey"))
			result, err := chain.Call(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(Values{"newKey": "value"}))
		})
	})

	Describe("TrimSpace", func() {
		It("should trim all spaces from the values of the given keys", func() {
			handler := HandlerFunc(func(ctx context.Context, values ...Values) (Values, error) {
				return Values{"key1": " value1 ", "key2": " value2 "}, nil
			})
			chain := Chain(handler, TrimSpace("key1", "key2"))
			result, err := chain.Call(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(Values{"key1": "value1", "key2": "value2"}))
		})
	})

	Describe("TrimSuffix", func() {
		It("should trim the given suffix from the values of the given keys", func() {
			handler := HandlerFunc(func(ctx context.Context, values ...Values) (Values, error) {
				return Values{"key1": "value1_suffix", "key2": "value2_suffix"}, nil
			})
			chain := Chain(handler, TrimSuffix("_suffix", "key1", "key2"))
			result, err := chain.Call(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(Values{"key1": "value1", "key2": "value2"}))
		})
	})

	Describe("ParallelChain", func() {
		It("should execute handlers in parallel and merge the results", func() {
			started := make(chan struct{})
			finish := make(chan struct{})

			handler1 := HandlerFunc(func(ctx context.Context, values ...Values) (Values, error) {
				started <- struct{}{}
				<-finish
				return Values{"key1": "value1"}, nil
			})
			handler2 := HandlerFunc(func(ctx context.Context, values ...Values) (Values, error) {
				started <- struct{}{}
				<-finish
				return Values{"key2": "value2"}, nil
			})

			chain := ParallelChain(2, handler1, handler2)

			go func() {
				// Wait for both handlers to start
				<-started
				<-started
				// Allow both handlers to finish
				close(finish)
			}()

			result, err := chain.Call(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(Values{"key1": "value1", "key2": "value2"}))

			Eventually(finish).Should(BeClosed())
		})

		It("should return an error if any of the handlers returns an error", func() {
			handler1 := HandlerFunc(func(ctx context.Context, values ...Values) (Values, error) {
				return Values{"key1": "value1"}, nil
			})
			handler2 := HandlerFunc(func(ctx context.Context, values ...Values) (Values, error) {
				return nil, errors.New("handler error")
			})
			chain := ParallelChain(2, handler1, handler2)
			_, err := chain.Call(context.Background())
			Expect(err).To(MatchError("handler error"))
		})
		It("should return context deadline error when context times out", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			defer cancel()
			handler1 := HandlerFunc(func(ctx context.Context, values ...Values) (Values, error) {
				select {
				case <-time.After(1 * time.Second):
				case <-ctx.Done():
				}

				return Values{"key1": "value1"}, nil
			})
			handler2 := HandlerFunc(func(ctx context.Context, values ...Values) (Values, error) {
				return Values{"key2": "value2"}, nil
			})
			chain := ParallelChain(2, handler1, handler2)
			_, err := chain.Call(ctx)
			Expect(err).To(MatchError(context.DeadlineExceeded))
		})
	})

	Describe("WithMemory", func() {
		It("should load previous conversations, call the wrapped handler, and save the last question/answer", func() {
			memory := &fakeMemory{
				ChatMessages: ChatMessages{{"user", "previous conversation"}},
			}
			handler := HandlerFunc(func(ctx context.Context, values ...Values) (Values, error) {
				Expect(values).To(HaveLen(1))
				Expect(values[0]).To(HaveKeyWithValue(DefaultKey, "input"))
				Expect(values[0]).To(HaveKeyWithValue(DefaultChatKey, ChatMessages{{"user", "previous conversation"}}))
				return Values{DefaultKey: "output"}, nil
			})
			chain := WithMemory(memory, handler)

			result, err := chain.Call(context.Background(), Values{DefaultKey: "input"})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(HaveKeyWithValue(DefaultKey, "output"))
			Expect(memory.ChatMessages).To(Equal(ChatMessages{
				{"user", "previous conversation"},
				{"user", "input"},
				{"assistant", "output"},
			}))
		})

		It("should return an error if loading from memory fails", func() {
			memory := &fakeMemory{
				LoadErr: errors.New("load error"),
			}
			handler := HandlerFunc(func(ctx context.Context, values ...Values) (Values, error) {
				return Values{DefaultKey: "output"}, nil
			})
			chain := WithMemory(memory, handler)

			result, err := chain.Call(context.Background(),
				Values{DefaultKey: "input"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("load error"))
			Expect(result).To(BeNil())
		})

		It("should return an error if saving to memory fails", func() {
			memory := &fakeMemory{
				SaveErr: errors.New("save error"),
			}
			handler := HandlerFunc(func(ctx context.Context, values ...Values) (Values, error) {
				return Values{DefaultKey: "output"}, nil
			})
			chain := WithMemory(memory, handler)

			result, err := chain.Call(context.Background(), Values{DefaultKey: "input"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("save error"))
			Expect(result).To(BeNil())
		})

		It("should return an error if input value is not a string", func() {
			memory := &fakeMemory{}
			handler := HandlerFunc(func(ctx context.Context, values ...Values) (Values, error) {
				return Values{DefaultKey: "output"}, nil
			})
			chain := WithMemory(memory, handler)

			result, err := chain.Call(context.Background(), Values{DefaultKey: 123})
			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
		})

		It("should return an error if output value is not a string", func() {
			memory := &fakeMemory{}
			handler := HandlerFunc(func(ctx context.Context, values ...Values) (Values, error) {
				return Values{DefaultKey: 123}, nil
			})
			chain := WithMemory(memory, handler)

			result, err := chain.Call(context.Background(), Values{DefaultKey: "input"})
			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
		})
	})
})

type fakeMemory struct {
	ChatMessages ChatMessages
	SaveErr      error
	LoadErr      error
}

func (m *fakeMemory) Load(context.Context) (ChatMessages, error) {
	if m.LoadErr != nil {
		return nil, m.LoadErr
	}
	return m.ChatMessages, nil
}

func (m *fakeMemory) Save(_ context.Context, input, output string) error {
	if m.SaveErr != nil {
		return m.SaveErr
	}
	m.ChatMessages = append(m.ChatMessages, ChatMessage{Role: "user", Content: input}, ChatMessage{Role: "assistant", Content: output})
	return nil
}
