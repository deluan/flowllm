package splitters_test

import (
	"github.com/deluan/pipelm"
	"github.com/deluan/pipelm/splitters"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Markdown", func() {
	var (
		text           string
		splitter       pipelm.Splitter
		expectedOutput []string
	)

	BeforeEach(func() {
		text = `
# Header 1

This is some content.

# Header 2

This is some more content.

# Header 3

This is even more content.`
	})

	Context("with markdown formatted text", func() {
		BeforeEach(func() {
			splitter = splitters.Markdown(
				splitters.WithChunkSize(40),
				splitters.WithChunkOverlap(20),
			)
			expectedOutput = []string{
				"# Header 1\n\nThis is some content.",
				"# Header 2\n\nThis is some more content.",
				"# Header 3\n\nThis is even more content.",
			}
		})

		It("splits the text into chunks based on the default chunk size and overlap", func() {
			chunks, err := splitter(text)
			Expect(err).NotTo(HaveOccurred())
			Expect(chunks).To(Equal(expectedOutput))
		})
	})
})
