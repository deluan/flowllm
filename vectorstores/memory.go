package vectorstores

import (
	"context"

	"github.com/deluan/pipelm"
	"golang.org/x/exp/slices"
)

type Memory struct {
	embeddings pipelm.Embeddings
	data       []memoryItem
}

type memoryItem struct {
	content  string
	vector   []float32
	metadata map[string]any
}

func NewMemoryVectorStore(embeddings pipelm.Embeddings) *Memory {
	return &Memory{
		embeddings: embeddings,
	}
}

func (m *Memory) AddDocuments(ctx context.Context, documents ...pipelm.Document) error {
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

func (m *Memory) SimilaritySearch(ctx context.Context, query string, k int) ([]pipelm.Document, error) {
	return SimilaritySearch(ctx, m, m.embeddings, query, k)
}

func (m *Memory) SimilaritySearchVectorWithScore(_ context.Context, query []float32, k int) ([]pipelm.ScoredDocument, error) {
	var results []pipelm.ScoredDocument
	for _, item := range m.data {
		similarity := CosineSimilarity(query, item.vector)
		results = append(results, pipelm.ScoredDocument{
			Document: pipelm.Document{
				PageContent: item.content,
				Metadata:    item.metadata,
			},
			Score: similarity,
		})
	}
	slices.SortFunc(results, func(a, b pipelm.ScoredDocument) bool {
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

func (m *Memory) addVectors(vectors [][]float32, documents []pipelm.Document) {
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
