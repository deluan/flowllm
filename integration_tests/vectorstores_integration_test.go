package integration_tests

import (
	"context"
	"os"
	"strconv"

	"github.com/deluan/pipelm"
	"github.com/deluan/pipelm/vectorstores"
	"github.com/deluan/pipelm/vectorstores/bolt"
	"github.com/deluan/pipelm/vectorstores/pinecone"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vector Stores Integration Tests", func() {
	var (
		boltVS         pipelm.VectorStore
		memoryVS       pipelm.VectorStore
		pineconeVS     pipelm.VectorStore
		ctx            context.Context
		mockEmbeddings *FakeEmbeddings
	)

	BeforeEach(func() {
		ctx = context.Background()
		var err error
		mockEmbeddings = &FakeEmbeddings{}

		// Create a Memory VectorStore
		memoryVS = vectorstores.NewMemoryVectorStore(mockEmbeddings)

		// Create a BoltDB VectorStore
		boltTmpDB, err := os.CreateTemp("", "pipelm_bolt_*_.db")
		Expect(err).ToNot(HaveOccurred())
		_ = boltTmpDB.Close()
		var closeDB func()
		boltVS, closeDB, err = bolt.NewVectorStore(mockEmbeddings, bolt.Options{Path: boltTmpDB.Name()})
		Expect(err).ToNot(HaveOccurred())
		DeferCleanup(closeDB)
		DeferCleanup(func() { _ = os.Remove(boltTmpDB.Name()) })

		if os.Getenv("PINECONE_API_KEY") != "" {
			// Create a Pinecone VectorStore
			pineconeVS, err = pinecone.NewVectorStore(ctx, mockEmbeddings,
				pinecone.Options{
					Index:     os.Getenv("PINECONE_INDEX_INTEGRATION_TEST"),
					NameSpace: "pipelm-integration-tests-" + uuid.NewString(),
				},
			)
			Expect(err).ToNot(HaveOccurred())
		}
	})

	DescribeTable("It should perform a similarity search using the query string and return correct results",
		func(getStore func() pipelm.VectorStore) {
			store := getStore()
			if store == nil {
				Skip("Skipping test. No VectorStore found.")
			}
			documents := []pipelm.Document{
				{
					PageContent: "first document",
					Metadata:    map[string]any{"key1": "value1"},
				},
				{
					PageContent: "second document",
					Metadata:    map[string]any{"key2": "value2"},
				},
			}
			Expect(store.AddDocuments(ctx, documents...)).To(Succeed())

			query := "2"
			k := 2
			similarDocs, err := store.SimilaritySearch(ctx, query, k)

			Expect(err).ToNot(HaveOccurred())
			Expect(similarDocs).To(HaveLen(k))
			Expect(similarDocs[0].PageContent).To(Equal(documents[1].PageContent))
			Expect(similarDocs[0].Metadata).To(Equal(documents[1].Metadata))
			Expect(similarDocs[1].PageContent).To(Equal(documents[0].PageContent))
			Expect(similarDocs[1].Metadata).To(Equal(documents[0].Metadata))
		},
		Entry("Memory", func() pipelm.VectorStore { return memoryVS }),
		Entry("Bolt", func() pipelm.VectorStore { return boltVS }),
		Entry("Pinecone", func() pipelm.VectorStore { return pineconeVS }),
	)

	DescribeTable("It should perform a similarity search using the query vector and return correct results with scores",
		func(getStore func() pipelm.VectorStore) {
			store := getStore()
			if store == nil {
				Skip("Skipping test. No VectorStore found.")
			}
			documents := []pipelm.Document{
				{PageContent: "first document"},
				{PageContent: "second document"},
			}
			Expect(store.AddDocuments(ctx, documents...)).To(Succeed())

			queryVector, _ := mockEmbeddings.EmbedString(ctx, "1")
			k := 2
			scoredDocs, err := store.SimilaritySearchVectorWithScore(ctx, queryVector, k)

			Expect(err).ToNot(HaveOccurred())
			Expect(scoredDocs).To(HaveLen(k))
			Expect(scoredDocs[0].Document.PageContent).To(Equal(documents[0].PageContent))
			Expect(scoredDocs[1].Document.PageContent).To(Equal(documents[1].PageContent))
			Expect(scoredDocs[0].Score).To(BeNumerically(">", scoredDocs[1].Score))
		},
		Entry("Memory", func() pipelm.VectorStore { return memoryVS }),
		Entry("Bolt", func() pipelm.VectorStore { return boltVS }),
		Entry("Pinecone", func() pipelm.VectorStore { return pineconeVS }),
	)

	DescribeTable("It should return all documents when k is greater than the number of documents in the vector store",
		func(getStore func() pipelm.VectorStore) {
			store := getStore()
			if store == nil {
				Skip("Skipping test. No VectorStore found.")
			}
			documents := []pipelm.Document{
				{PageContent: "first document"},
				{PageContent: "second document"},
			}
			Expect(store.AddDocuments(ctx, documents...)).To(Succeed())

			query := "1"
			k := 3
			similarDocs, err := store.SimilaritySearch(ctx, query, k)

			Expect(err).ToNot(HaveOccurred())
			Expect(similarDocs).To(HaveLen(len(documents)))
			Expect(similarDocs[0].PageContent).To(Equal(documents[0].PageContent))
			Expect(similarDocs[1].PageContent).To(Equal(documents[1].PageContent))
		},
		Entry("Memory", func() pipelm.VectorStore { return memoryVS }),
		Entry("Bolt", func() pipelm.VectorStore { return boltVS }),
		Entry("Pinecone", func() pipelm.VectorStore { return pineconeVS }),
	)

	DescribeTable("It should return an empty result when performing a similarity search on an empty vector store",
		func(getStore func() pipelm.VectorStore) {
			store := getStore()
			if store == nil {
				Skip("Skipping test. No VectorStore found.")
			}
			query := "test query"
			k := 2
			similarDocs, err := store.SimilaritySearch(ctx, query, k)

			Expect(err).ToNot(HaveOccurred())
			Expect(similarDocs).To(BeEmpty())
		},
		Entry("Memory", func() pipelm.VectorStore { return memoryVS }),
		Entry("Bolt", func() pipelm.VectorStore { return boltVS }),
		Entry("Pinecone", func() pipelm.VectorStore { return pineconeVS }),
	)
})

type FakeEmbeddings struct{}

func (m *FakeEmbeddings) EmbedString(_ context.Context, query string) ([]float32, error) {
	num, _ := strconv.Atoi(query)
	vectors := m.fakeEmbed(num, 1, 1)
	return vectors, nil
}

func (m *FakeEmbeddings) EmbedStrings(_ context.Context, texts []string) ([][]float32, error) {
	vectors := make([][]float32, len(texts))
	for i := range texts {
		vectors[i] = m.fakeEmbed(i+1, 1, 1)
	}
	return vectors, nil
}

func (m *FakeEmbeddings) fakeEmbed(nums ...int) []float32 {
	vector := make([]float32, 1536)
	for i, n := range nums {
		vector[i] = float32(n)
	}
	return vector
}
