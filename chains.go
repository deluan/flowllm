package pipelm

import (
	"context"
	"log"
	"strings"

	"github.com/deluan/pipelm/pl"
)

// Chain is a special handler that executes a list of handlers in sequence.
// The output of each chain is passed as input to the next one.
// The output of the last chain is returned as the output of the Sequential chain.
func Chain(handlers ...Handler) HandlerFunc {
	return func(ctx context.Context, values ...Values) (Values, error) {
		variables := Values{}.Merge(values...)
		for _, chain := range handlers {
			var err error
			variables, err = chain.Call(ctx, variables)
			if err != nil {
				return nil, err
			}
		}
		return variables, nil
	}
}

// MapOutputTo renames the output of the chain to the given key
func MapOutputTo(key string) HandlerFunc {
	return func(ctx context.Context, values ...Values) (Values, error) {
		variables := Values{}.Merge(values...)
		variables[key] = variables[DefaultKey]
		delete(variables, DefaultKey)
		return variables, nil
	}
}

// TrimSpace trims all spaces from the values of the given keys
func TrimSpace(keys ...string) HandlerFunc {
	return func(ctx context.Context, values ...Values) (Values, error) {
		variables := Values{}.Merge(values...)
		for _, key := range keys {
			variables[key] = strings.TrimSpace(variables.Get(key))
		}
		return variables, nil
	}
}

// TrimSuffix trims the given suffix from the values of the given keys
func TrimSuffix(suffix string, keys ...string) HandlerFunc {
	return func(ctx context.Context, values ...Values) (Values, error) {
		variables := Values{}.Merge(values...)
		for _, key := range keys {
			variables[key] = strings.TrimSuffix(variables.Get(key), suffix)
		}
		return variables, nil
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

		variables := Values{}.Merge(values...)
		resC, errC := pl.Stage(ctx, maxParallel, chains, func(ctx context.Context, handler Handler) (Values, error) {
			return handler.Call(ctx, variables)
		})

		results := Values{}
		for res := range pl.ReadOrDone(ctx, resC) {
			results = results.Merge(res)
		}

		var fatalErr error
		for err := range pl.ReadOrDone(ctx, errC) {
			if err != nil {
				log.Printf("Error: %v", err)
				cancel()
				fatalErr = err
				break
			}
		}
		if fatalErr != nil {
			return nil, fatalErr
		}
		<-resC

		return results, nil
	}
}

type LanguageModel interface {
	Call(ctx context.Context, input string) (string, error)
}

type ChatMessage struct {
	Role    string
	Content string
}

type ChatLanguageModel interface {
	Chat(ctx context.Context, msgs []ChatMessage) (string, error)
}

func LLM(model LanguageModel) HandlerFunc {
	return func(ctx context.Context, values ...Values) (Values, error) {
		variables := Values{}.Merge(values...)
		input := variables.Get(DefaultKey)
		output, err := model.Call(ctx, input)
		if err != nil {
			return nil, err
		}
		variables[DefaultKey] = output
		return variables, nil
	}
}

func ChatLLM(model ChatLanguageModel) HandlerFunc {
	return func(ctx context.Context, values ...Values) (Values, error) {
		variables := Values{}.Merge(values...)
		msgs, ok := variables[DefaultChatKey].([]ChatMessage)
		if !ok {
			msgs = []ChatMessage{
				{Role: "user", Content: variables.Get(DefaultKey)},
			}
		}
		output, err := model.Chat(ctx, msgs)
		if err != nil {
			return nil, err
		}
		variables[DefaultKey] = output
		return variables, nil
	}
}
