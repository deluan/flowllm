package pipelm_test

import (
	. "github.com/deluan/pipelm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Splitters", func() {

	Describe("RecursiveTextSplitter", func() {
		var (
			text           string
			splitter       Splitter
			expectedOutput []string
		)

		BeforeEach(func() {
			text = "This is a sample text for testing the RecursiveTextSplitter function."
		})

		Context("with default options", func() {
			BeforeEach(func() {
				splitter = RecursiveTextSplitter(SplitterOptions{})
				expectedOutput = []string{
					"This is a sample text for testing the RecursiveTextSplitter function.",
				}
			})

			It("splits the text into chunks based on the default chunk size and overlap", func() {
				chunks, err := splitter(text)
				Expect(err).NotTo(HaveOccurred())
				Expect(chunks).To(Equal(expectedOutput))
			})
		})

		Context("with custom options", func() {
			BeforeEach(func() {
				splitter = RecursiveTextSplitter(SplitterOptions{
					ChunkSize:    20,
					ChunkOverlap: 2,
					Separators:   []string{"\n"},
				})
				expectedOutput = []string{
					"This is a sample te",
					"text for testing th",
					"the RecursiveTextSp",
					"Splitter function.",
				}
			})

			It("splits the text into chunks based on the custom chunk size and overlap", func() {
				chunks, err := splitter(text)
				Expect(err).NotTo(HaveOccurred())
				Expect(chunks).To(Equal(expectedOutput))
			})
		})
	})

	Describe("MarkdownSplitter", func() {
		var (
			text           string
			splitter       Splitter
			expectedOutput []string
		)

		BeforeEach(func() {
			text = `
# Header 1

This is some content.

## Header 2

This is some more content.

### Header 3

This is even more content.`
		})

		Context("with markdown formatted text", func() {
			BeforeEach(func() {
				splitter = MarkdownSplitter(SplitterOptions{ChunkSize: 40, ChunkOverlap: 20})
				expectedOutput = []string{
					"# Header 1\n\nThis is some content.",
					"Header 2\n\nThis is some more content.",
					"Header 3\n\nThis is even more content.",
				}
			})

			It("splits the text into chunks based on the default chunk size and overlap", func() {
				chunks, err := splitter(text)
				Expect(err).NotTo(HaveOccurred())
				Expect(chunks).To(Equal(expectedOutput))
			})
		})
	})
})
