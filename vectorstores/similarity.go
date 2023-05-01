package vectorstores

import (
	"context"
	"math"

	"github.com/deluan/pipelm"
)

func CosineSimilarity(a, b []float32) float32 {
	var p, p2, q2 float32
	for i := 0; i < len(a) && i < len(b); i++ {
		p += a[i] * b[i]
		p2 += a[i] * a[i]
		q2 += b[i] * b[i]
	}
	if p2 == 0 || q2 == 0 {
		return 0
	}
	return p / (float32(math.Sqrt(float64(p2))) * float32(math.Sqrt(float64(q2))))
}

// SimilaritySearch returns the k most similar documents to the given query. It uses the given
// vector store's SimilaritySearchVectorWithScore method to perform the search.
func SimilaritySearch(ctx context.Context, store pipelm.VectorStore, embeddings pipelm.Embeddings, query string, k int) ([]pipelm.Document, error) {
	queryVector, err := embeddings.EmbedString(ctx, query)
	if err != nil {
		return nil, err
	}
	var docs []pipelm.Document
	results, _ := store.SimilaritySearchVectorWithScore(ctx, queryVector, k)
	for _, result := range results {
		docs = append(docs, result.Document)
	}
	return docs, nil
}
