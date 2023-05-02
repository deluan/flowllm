package vectorstores

import (
	"context"

	"github.com/deluan/flowllm"
	"golang.org/x/exp/slices"
)

// Memory is a simple in-memory vector store. It implements the VectorStore interface and
// stores the vectors in memory. It is not meant to be used in production, but it is useful
// for testing and as an example of how to implement a VectorStore.
type Memory struct {
	embeddings flowllm.Embeddings
	data       []memoryItem
}

type memoryItem struct {
	content  string
	vector   []float32
	metadata map[string]any
}

// NewMemoryVectorStore creates a new Memory vector store.
func NewMemoryVectorStore(embeddings flowllm.Embeddings) *Memory {
	return &Memory{
		embeddings: embeddings,
	}
}

func (m *Memory) AddDocuments(ctx context.Context, documents ...flowllm.Document) error {
	texts := make([]string, len(documents))
	for i, document := range documents {
		texts[i] = document.PageContent
	}
	vectors, err := m.embeddings.EmbedStrings(ctx, texts)
	if err != nil {
		return err
	}
	m.addVectors(vectors, documents)
	return nil
}

func (m *Memory) SimilaritySearch(ctx context.Context, query string, k int) ([]flowllm.Document, error) {
	return SimilaritySearch(ctx, m, m.embeddings, query, k)
}

func (m *Memory) SimilaritySearchVectorWithScore(_ context.Context, query []float32, k int) ([]flowllm.ScoredDocument, error) {
	var results []flowllm.ScoredDocument
	for _, item := range m.data {
		similarity := CosineSimilarity(query, item.vector)
		results = append(results, flowllm.ScoredDocument{
			Document: flowllm.Document{
				PageContent: item.content,
				Metadata:    item.metadata,
			},
			Score: similarity,
		})
	}
	slices.SortFunc(results, func(a, b flowllm.ScoredDocument) bool {
		return a.Score > b.Score
	})
	k = min(k, len(results))
	return results[0:k], nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m *Memory) addVectors(vectors [][]float32, documents []flowllm.Document) {
	var memoryVectors []memoryItem
	for i, vector := range vectors {
		memoryVectors = append(memoryVectors, memoryItem{
			content:  documents[i].PageContent,
			vector:   vector,
			metadata: documents[i].Metadata,
		})
	}
	m.data = append(m.data, memoryVectors...)
}
