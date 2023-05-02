package flowllm

import (
	"context"
	"errors"
	"io"
)

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

// Document represents a document to be stored in a VectorStore.
type Document struct {
	ID          string
	PageContent string
	Metadata    map[string]any
}

// ScoredDocument represents a document along with its similarity score.
type ScoredDocument struct {
	Document
	Score float32
}

// Splitter is a function that splits a string into a slice of strings.
type Splitter = func(string) ([]string, error)

// DocumentLoader is the interface implemented by types that can load documents.
// The LoadNext method should the next available document, or io.EOF if there are no more documents.
type DocumentLoader interface {
	LoadNext(ctx context.Context) (Document, error)
}

// DocumentLoaderFunc is an adapter to allow the use of ordinary functions as DocumentLoaders.
type DocumentLoaderFunc func(ctx context.Context) (Document, error)

func (f DocumentLoaderFunc) LoadNext(ctx context.Context) (Document, error) {
	return f(ctx)
}

// LoadDocs loads the next n documents from the given DocumentLoader.
func LoadDocs(n int, loader DocumentLoader) ([]Document, error) {
	ctx := context.Background()
	var docs []Document
	for {
		doc, err := loader.LoadNext(ctx)
		if errors.Is(err, io.EOF) {
			return docs, nil
		}
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
		n--
		if n == 0 {
			return docs, nil
		}
	}
}
