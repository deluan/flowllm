// Package tiktoken implements a wrapper around the github.com/tiktoken-go/tokenizer library.
// It is compatible with OpenAI's tiktoken Python library.
package tiktoken

import "github.com/tiktoken-go/tokenizer"

// Len returns a len function for the given model that returns the number of tokens in a string.
func Len(forModel string) func(string) int {
	enc, err := tokenizer.ForModel(tokenizer.Model(forModel))
	if err != nil {
		panic(err)
	}
	return func(s string) int {
		tokens, _, _ := enc.Encode(s)
		return len(tokens)
	}
}
