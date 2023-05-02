package flowllm

import (
	"context"
	"fmt"
	"strings"

	"github.com/deluan/flowllm/pl"
)

// Handler is the interface implemented by all composable modules in the library.
type Handler interface {
	Call(ctx context.Context, values ...Values) (Values, error)
}

// HandlerFunc is a function that implements the Handler interface.
type HandlerFunc func(context.Context, ...Values) (Values, error)

func (f HandlerFunc) Call(ctx context.Context, values ...Values) (Values, error) {
	return f(ctx, values...)
}

// Chain is a special handler that executes a list of handlers in sequence.
// The output of each chain is passed as input to the next one.
// The output of the last chain is returned as the output of the Sequential chain.
func Chain(handlers ...Handler) HandlerFunc {
	return func(ctx context.Context, values ...Values) (Values, error) {
		vals := Values{}.Merge(values...)
		for _, chain := range handlers {
			var err error
			vals, err = chain.Call(ctx, vals)
			if err != nil {
				return nil, err
			}
		}
		return vals, nil
	}
}

// MapOutputTo renames the output of the chain (DefaultKey) to the given key.
func MapOutputTo(key string) HandlerFunc {
	return func(ctx context.Context, values ...Values) (Values, error) {
		vals := Values{}.Merge(values...)
		vals[key] = vals[DefaultKey]
		delete(vals, DefaultKey)
		return vals, nil
	}
}

// TrimSpace trims all spaces from the values of the given keys.
func TrimSpace(keys ...string) HandlerFunc {
	return func(ctx context.Context, values ...Values) (Values, error) {
		vals := Values{}.Merge(values...)
		for _, key := range keys {
			vals[key] = strings.TrimSpace(vals.Get(key))
		}
		return vals, nil
	}
}

// TrimSuffix trims the given suffix from the values of the given keys.
func TrimSuffix(suffix string, keys ...string) HandlerFunc {
	return func(ctx context.Context, values ...Values) (Values, error) {
		vals := Values{}.Merge(values...)
		for _, key := range keys {
			vals[key] = strings.TrimSuffix(vals.Get(key), suffix)
		}
		return vals, nil
	}
}

// ParallelChain executes a list of handlers in parallel, up to a maximum number of concurrent executions.
// If any of the handlers returns an error, the execution is stopped and the error is returned.
// The results of all handlers are merged into a single Values object.
func ParallelChain(maxParallel int, handlers ...Handler) HandlerFunc {
	return func(ctx context.Context, values ...Values) (Values, error) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		chains := pl.FromSlice(ctx, handlers)

		vals := Values{}.Merge(values...)
		resC, errC := pl.Stage(ctx, maxParallel, chains, func(ctx context.Context, handler Handler) (Values, error) {
			return handler.Call(ctx, vals)
		})

		finalErrC := make(chan error)
		go func() {
			for err := range pl.ReadOrDone(ctx, errC) {
				if err != nil {
					finalErrC <- err
					cancel()
					return
				}
			}
			finalErrC <- nil
		}()

		results := Values{}
		for res := range pl.ReadOrDone(ctx, resC) {
			results = results.Merge(res)
		}

		if err := <-finalErrC; err != nil {
			return nil, err
		}
		return results, ctx.Err()
	}
}

// LanguageModel interface is implemented by all language models.
type LanguageModel interface {
	Call(ctx context.Context, input string) (string, error)
}

// LLM is a handler that can be used to add a language model to a chain.
func LLM(model LanguageModel) HandlerFunc {
	return func(ctx context.Context, values ...Values) (Values, error) {
		vals := Values{}.Merge(values...)
		input := vals.Get(DefaultKey)
		output, err := model.Call(ctx, input)
		if err != nil {
			return nil, err
		}
		vals[DefaultKey] = output
		return vals, nil
	}
}

// ChatMessage is a struct that represents a message in a chat conversation.
type ChatMessage struct {
	Role    string
	Content string
}

// ChatMessages is a list of ChatMessage.
type ChatMessages []ChatMessage

func (m ChatMessages) String() string {
	var output []string
	for _, msg := range m {
		output = append(output, fmt.Sprintf("%s: %s", msg.Role, msg.Content))
	}
	return strings.Join(output, "\n")
}

// Last returns the last N messages from the list.
func (m ChatMessages) Last(size int) ChatMessages {
	if len(m) < size {
		return m
	}
	return m[len(m)-size:]
}

// ChatLanguageModel interface is implemented by all chat language models.
type ChatLanguageModel interface {
	Chat(ctx context.Context, msgs []ChatMessage) (string, error)
}

// ChatLLM is a handler that can be used to add a chat model to a chain.
// It is similar to the LLM handler, but it has a few differences:
// It will use the value of the DefaultChatKey key (usually set by the ChatTemplate) as input
// to the model, if available. If not, it will use the value of the DefaultKey key.
func ChatLLM(model ChatLanguageModel) HandlerFunc {
	return func(ctx context.Context, values ...Values) (Values, error) {
		vals := Values{}.Merge(values...)
		msgs, _ := vals[DefaultChatKey].(ChatMessages)
		if msgs == nil {
			text := vals.Get(DefaultKey)
			if text != "" {
				msgs = append(msgs, ChatMessage{Role: "user", Content: text})
			}
		}
		output, err := model.Chat(ctx, msgs)
		if err != nil {
			return nil, err
		}
		vals[DefaultKey] = output
		return vals, nil
	}
}

// Memory is an interface that can be used to store and retrieve previous conversations.
type Memory interface {
	// Load returns previous conversations from the memory
	Load(context.Context) (ChatMessages, error)

	// Save last question/answer to the memory
	Save(ctx context.Context, input, output string) error
}

// WithMemory is a wrapper that loads the previous conversation from the memory,
// injects it into the chain as the value of the DefaultChatKey key, calls the wrapped handler,
// and adds the last question/answer to the memory.
func WithMemory(memory Memory, handler Handler) HandlerFunc {
	return func(ctx context.Context, values ...Values) (Values, error) {
		vals := Values{}.Merge(values...)
		history, err := memory.Load(ctx)
		if err != nil {
			return nil, err
		}
		outputVals, err := handler.Call(ctx, vals.Merge(Values{DefaultChatKey: history}))
		if err != nil {
			return nil, err
		}
		input, err := getValue(vals, "")
		if err != nil {
			return nil, err
		}
		output, err := getValue(outputVals, DefaultKey)
		if err != nil {
			return nil, err
		}
		err = memory.Save(ctx, input, output)
		if err != nil {
			return nil, err
		}
		return vals.Merge(outputVals), nil
	}
}

// getValue returns the value of the given key from the given Values object.
// If the key is empty, it returns the value of the first key in the Values object.
// If the Values object has multiple keys, it returns an error.
func getValue(values Values, key string) (string, error) {
	ret := func(v any) (string, error) {
		if s, ok := v.(string); !ok {
			return "", fmt.Errorf("input value is not a string: %v", v)
		} else {
			return s, nil
		}
	}
	if key != "" {
		return ret(values[key])
	}
	keys := values.Keys()
	if len(keys) == 1 {
		return ret(values[keys[0]])
	}
	return "", fmt.Errorf("input values have multiple keys, memory only supported when one key currently: %v", keys)
}
