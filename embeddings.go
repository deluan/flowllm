package flowllm

import "context"

// Embeddings can be used to create a numerical representation of textual data.
// This numerical representation is useful when searching for similar documents.
type Embeddings interface {
	// EmbedString returns the embedding for the given string
	EmbedString(context.Context, string) ([]float32, error)

	// EmbedStrings returns the embeddings for multiple strings
	EmbedStrings(context.Context, []string) ([][]float32, error)
}
