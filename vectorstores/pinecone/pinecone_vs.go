package pinecone

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/deluan/flowllm"
	"github.com/deluan/flowllm/vectorstores"
	"golang.org/x/exp/slices"
)

// Options for the Pinecone vector store.
type Options struct {
	ApiKey      string
	Environment string
	Index       string
	NameSpace   string
	Pods        int
	Replicas    int
	PodType     string
	Metric      Metric
}

// VectorStore is a vector store backed by Pinecone. It requires an already created Pinecone index,
// with the same dimensionality as the embeddings used to create the store.
type VectorStore struct {
	client     *client
	embeddings flowllm.Embeddings
	textKey    string
}

// NewVectorStore creates a new Pinecone vector store.
func NewVectorStore(ctx context.Context, embeddings flowllm.Embeddings, opts Options) (*VectorStore, error) {
	if opts.ApiKey == "" {
		opts.ApiKey = os.Getenv("PINECONE_API_KEY")
	}
	if opts.Environment == "" {
		opts.Environment = os.Getenv("PINECONE_ENVIRONMENT")
	}
	if opts.Index == "" {
		opts.Index = os.Getenv("PINECONE_INDEX")
	}
	if opts.Pods == 0 {
		opts.Pods = 1
	}
	if opts.Replicas == 0 {
		opts.Replicas = 1
	}
	if opts.PodType == "" {
		opts.PodType = "s1"
	}
	if opts.Metric == "" {
		opts.Metric = Cosine
	}
	c := &client{
		index:       opts.Index,
		pods:        opts.Pods,
		replicas:    opts.Replicas,
		podType:     opts.PodType,
		metric:      string(opts.Metric),
		apiKey:      opts.ApiKey,
		environment: opts.Environment,
		namespace:   opts.NameSpace,
	}
	s := VectorStore{
		embeddings: embeddings,
		client:     c,
		textKey:    "text",
	}
	err := connect(ctx, c)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *VectorStore) AddDocuments(ctx context.Context, documents ...flowllm.Document) error {
	var texts []string
	for i := 0; i < len(documents); i++ {
		texts = append(texts, documents[i].PageContent)
	}

	vectors, err := s.embeddings.EmbedStrings(ctx, texts)
	if err != nil {
		return err
	}

	var items []pineconeItem
	for i := 0; i < len(vectors); i++ {
		curMetadata := make(map[string]string)
		for key, value := range documents[i].Metadata {
			curMetadata[key] = fmt.Sprintf("%s", value)
		}

		curMetadata[s.textKey] = documents[i].PageContent

		items = append(items, pineconeItem{
			Values:   vectors[i],
			Metadata: curMetadata,
			ID:       fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%v", curMetadata)))),
		})
	}

	return s.client.upsert(ctx, items)
}

func (s *VectorStore) SimilaritySearch(ctx context.Context, query string, k int) ([]flowllm.Document, error) {
	return vectorstores.SimilaritySearch(ctx, s, s.embeddings, query, k)
}

func (s *VectorStore) SimilaritySearchVectorWithScore(ctx context.Context, query []float32, k int) ([]flowllm.ScoredDocument, error) {
	queryResponse, err := s.client.query(ctx, query, k)
	if err != nil {
		return nil, err
	}

	var resultDocuments []flowllm.ScoredDocument
	for _, match := range queryResponse.Matches {
		pageContent, ok := match.Metadata[s.textKey]
		if !ok {
			return nil, fmt.Errorf("missing textKey %s in query response match", s.textKey)
		}

		metadata := make(map[string]any)
		for key, value := range match.Metadata {
			if key == s.textKey {
				continue
			}
			metadata[key] = value
		}

		resultDocuments = append(resultDocuments, flowllm.ScoredDocument{
			Document: flowllm.Document{
				PageContent: pageContent,
				Metadata:    metadata,
			},
			Score: match.Score,
		})
	}
	slices.SortFunc(resultDocuments, func(a, b flowllm.ScoredDocument) bool {
		return a.Score > b.Score
	})

	return resultDocuments, nil
}

type Metric string

const (
	Euclidean  Metric = "euclidean"
	Cosine     Metric = "cosine"
	DotProduct Metric = "dotproduct"
)
