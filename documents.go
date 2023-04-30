package pipelm

import "context"

// VectorStore is a particular type of database optimized for storing documents and their embeddings,
// and then fetching of the most relevant documents for a particular query, i.e. those whose embeddings
// are most similar to the embedding of the query.
type VectorStore interface {
	// AddDocuments adds the given documents to the store
	AddDocuments(context.Context, ...Document) error
	// SimilaritySearch returns the k most similar documents to the query
	SimilaritySearch(ctx context.Context, query string, k int) ([]Document, error)
	// SimilaritySearchVectorWithScore returns the k most similar documents to the query, along with their similarity score
	SimilaritySearchVectorWithScore(ctx context.Context, query []float32, k int) ([]ScoredDocument, error)
}

type Document struct {
	ID          string
	PageContent string
	Metadata    map[string]any
}

type ScoredDocument struct {
	Document
	Score float32
}

type Splitter = func(string) ([]string, error)

type DocumentLoader interface {
	LoadNext(ctx context.Context) (Document, error)
}

type DocumentLoaderFunc func(ctx context.Context) (Document, error)

func (f DocumentLoaderFunc) LoadNext(ctx context.Context) (Document, error) {
	return f(ctx)
}
