package main

import (
	"context"
	"fmt"

	"github.com/deluan/pipelm"
	"github.com/deluan/pipelm/llms/openai"
	"github.com/deluan/pipelm/loaders"
	"github.com/deluan/pipelm/splitters"
	"github.com/deluan/pipelm/vectorstores"
)

func init() {
	registerExample("vector_store_memory", "InMemory Vector Store", vectorStoreMemory)
}

func vectorStoreMemory() {
	ctx := context.Background()

	//Create docs with a loader
	loader := loaders.TextFile(
		"testdata/state_of_the_union.txt",
		splitters.RecursiveCharacterText(splitters.Options{ChunkSize: 100, ChunkOverlap: 10}),
	)

	// Create a vector store
	embClient, err := openai.NewEmbeddings(openai.EmbeddingsOptions{})
	if err != nil {
		panic(err)
	}
	vectorStore := vectorstores.NewMemoryVectorStore(embClient)

	// Load the first 30 documents
	docs, err := pipelm.LoadDocs(30, loader)
	if err != nil {
		panic(err)
	}

	// Add the documents to the vector store
	err = vectorStore.AddDocuments(ctx, docs...)
	if err != nil {
		panic(err)
	}

	// Embed the query
	query, err := embClient.EmbedString(ctx, "ukraine")
	if err != nil {
		panic(err)
	}

	// Search for the top 5 most similar documents, with their similarity score
	res, err := vectorStore.SimilaritySearchVectorWithScore(ctx, query, 5)
	if err != nil {
		panic(err)
	}

	// Print the results
	for _, doc := range res {
		fmt.Printf("[%4.1f]: %s\n", doc.Score*100, doc.PageContent)
	}
}
