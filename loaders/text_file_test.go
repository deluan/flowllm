package loaders_test

import (
	"context"
	"io"

	"github.com/deluan/flowllm"
	"github.com/deluan/flowllm/loaders"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TextFile", func() {
	var (
		ctx    context.Context
		loader flowllm.DocumentLoader
	)

	BeforeEach(func() {
		ctx = context.Background()
	})
	Context("with default options", func() {
		BeforeEach(func() {
			loader = loaders.TextFile("../testdata/state_of_the_union.txt")
		})
		It("loads the full text file into a Document", func() {
			docs, err := loader.LoadNext(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(docs.PageContent).To(ContainSubstring("the State of the Union"))
			Expect(docs.Metadata).To(HaveKeyWithValue("source", "../testdata/state_of_the_union.txt"))
			_, err = loader.LoadNext(ctx)
			Expect(err).To(MatchError(io.EOF))
		})
	})

	Context("with a splitter", func() {
		BeforeEach(func() {
			loader = loaders.TextFile("../testdata/small_text.txt", flowllm.RecursiveTextSplitter(flowllm.SplitterOptions{ChunkSize: 10}))
		})
		It("loads the text file one chunk at a time", func() {
			doc, err := loader.LoadNext(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(doc.PageContent).To(Equal("This is a"))

			doc, err = loader.LoadNext(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(doc.PageContent).To(Equal("small text"))

			doc, err = loader.LoadNext(ctx)
			Expect(err).To(MatchError(io.EOF))
		})

		It("loads all documents", func() {
			docs, err := flowllm.LoadDocs(10000, loader)
			Expect(err).NotTo(HaveOccurred())
			Expect(docs).To(HaveLen(2))
			Expect(docs[0].PageContent).To(Equal("This is a"))
			Expect(docs[1].PageContent).To(Equal("small text"))
		})
	})
})
