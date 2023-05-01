package openai

import (
	"context"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type EmbeddingsOptions struct {
	ApiKey       string
	KeepNewLines bool
	BatchSize    int
}
type Embeddings struct {
	client *openai.Client
	opts   EmbeddingsOptions
}

func NewEmbeddings(opts EmbeddingsOptions) (*Embeddings, error) {
	if opts.ApiKey == "" {
		opts.ApiKey = os.Getenv("OPENAI_API_KEY")
	}
	if opts.BatchSize == 0 {
		opts.BatchSize = 512
	}
	e := &Embeddings{opts: opts}
	e.client = openai.NewClient(opts.ApiKey)

	return e, nil
}

type Option func(*Embeddings)

func (o *Embeddings) EmbedString(ctx context.Context, text string) ([]float32, error) {
	texts := o.prepareTexts([]string{text})
	embeddings, err := o.embedTexts(ctx, []string{texts[0]})
	if err != nil {
		return nil, err
	}
	return embeddings[0], nil
}

func (o *Embeddings) EmbedStrings(ctx context.Context, texts []string) ([][]float32, error) {
	chunks := chunkArray(o.prepareTexts(texts), o.opts.BatchSize)
	var embeddings [][]float32
	for _, input := range chunks {
		result, err := o.embedTexts(ctx, input)
		if err != nil {
			return nil, err
		}
		embeddings = append(embeddings, result...)
	}
	return embeddings, nil
}

func (o *Embeddings) prepareTexts(texts []string) []string {
	if !o.opts.KeepNewLines {
		for i, text := range texts {
			texts[i] = strings.ReplaceAll(text, "\n", " ")
		}
	}
	return texts
}

func chunkArray(arr []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(arr); i += chunkSize {
		end := i + chunkSize
		if end > len(arr) {
			end = len(arr)
		}
		chunks = append(chunks, arr[i:end])
	}
	return chunks
}

func (o *Embeddings) embedTexts(ctx context.Context, texts []string) ([][]float32, error) {
	req := openai.EmbeddingRequest{
		Input: texts,
		Model: openai.AdaEmbeddingV2,
	}
	resp, err := o.client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, err
	}
	var embeddings [][]float32
	for _, embedding := range resp.Data {
		embeddings = append(embeddings, embedding.Embedding)
	}
	return embeddings, nil
}
