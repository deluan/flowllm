package bolt

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/fs"
	"time"

	"github.com/deluan/pipelm"
	"github.com/deluan/pipelm/vectorstores"
	"go.etcd.io/bbolt"
	"golang.org/x/exp/slices"
)

const (
	DefaultPath       = "vector_store.db"
	DefaultBucket     = "embeddings"
	DefaultPermission = 0600
)

type Options struct {
	Path       string
	Bucket     string
	Permission fs.FileMode
	Timeout    time.Duration
}

type VectorStore struct {
	embeddings pipelm.Embeddings
	db         *bbolt.DB
	bucket     string
}

func NewVectorStore(embeddings pipelm.Embeddings, opts Options) (*VectorStore, func(), error) {
	if opts.Path == "" {
		opts.Path = DefaultPath
	}
	if opts.Bucket == "" {
		opts.Bucket = DefaultBucket
	}
	if opts.Permission == 0 {
		opts.Permission = DefaultPermission
	}
	if opts.Timeout == 0 {
		opts.Timeout = time.Second
	}
	s := VectorStore{
		embeddings: embeddings,
		bucket:     opts.Bucket,
	}
	db, err := bbolt.Open(opts.Path, opts.Permission, &bbolt.Options{Timeout: opts.Timeout})
	if err != nil {
		return nil, func() {}, err
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			return fmt.Errorf("create bucket: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, func() {}, err
	}
	s.db = db
	return &s, func() { _ = db.Close() }, nil
}

type boltItem struct {
	Vectors  []float32              `json:"vectors"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (d boltItem) id() string {
	return fmt.Sprintf("%x", sha256.Sum256(d.Marshall()))
}

func (d boltItem) Marshall() []byte {
	buf, _ := json.Marshal(d)
	return buf
}

func (s *VectorStore) AddDocuments(ctx context.Context, documents ...pipelm.Document) error {
	texts := make([]string, len(documents))
	for i, document := range documents {
		texts[i] = document.PageContent
	}
	vectors, err := s.embeddings.EmbedStrings(ctx, texts)
	if err != nil {
		return err
	}

	return s.db.Batch(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(s.bucket))
		for i, doc := range documents {
			item := boltItem{
				Vectors:  vectors[i],
				Content:  doc.PageContent,
				Metadata: doc.Metadata,
			}
			if err := bucket.Put([]byte(item.id()), item.Marshall()); err != nil {
				return err
			}
		}
		return nil
	})
}

type match struct {
	id         []byte
	similarity float32
}

func (s *VectorStore) SimilaritySearch(ctx context.Context, query string, k int) ([]pipelm.Document, error) {
	return vectorstores.SimilaritySearch(ctx, s, s.embeddings, query, k)
}

func (s *VectorStore) SimilaritySearchVectorWithScore(_ context.Context, query []float32, k int) ([]pipelm.ScoredDocument, error) {
	var matches []match
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(s.bucket))
		return bucket.ForEach(func(k, v []byte) error {
			var item boltItem
			err := json.Unmarshal(v, &item)
			if err != nil {
				return err
			}
			similarity := vectorstores.CosineSimilarity(query, item.Vectors)
			matches = append(matches, match{id: k, similarity: similarity})
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	slices.SortFunc(matches, func(a, b match) bool {
		return a.similarity > b.similarity
	})
	k = min(k, len(matches))
	matches = matches[:k]
	var results []pipelm.ScoredDocument
	err = s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(s.bucket))
		for _, match := range matches {
			var item boltItem
			err := json.Unmarshal(bucket.Get(match.id), &item)
			if err != nil {
				return err
			}
			results = append(results, pipelm.ScoredDocument{
				Score: match.similarity,
				Document: pipelm.Document{
					PageContent: item.Content,
					Metadata:    item.Metadata,
				},
			})
		}
		return nil
	})
	return results, err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
