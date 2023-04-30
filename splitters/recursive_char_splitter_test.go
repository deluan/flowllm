package splitters_test

import (
	"github.com/deluan/pipelm"
	"github.com/deluan/pipelm/splitters"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RecursiveCharacterText", func() {
	var (
		text           string
		splitter       pipelm.Splitter
		expectedOutput []string
	)

	BeforeEach(func() {
		text = "This is a sample text for testing the RecursiveCharacterText function."
	})

	Context("with default options", func() {
		BeforeEach(func() {
			splitter = splitters.RecursiveCharacterText()
			expectedOutput = []string{
				"This is a sample text for testing the RecursiveCharacterText function.",
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
			splitter = splitters.RecursiveCharacterText(
				splitters.WithChunkSize(20),
				splitters.WithChunkOverlap(2),
				splitters.WithSeparators("\n"),
			)
			expectedOutput = []string{
				"This is a sample te",
				"text for testing th",
				"the RecursiveCharac",
				"acterText function.",
			}
		})

		It("splits the text into chunks based on the custom chunk size and overlap", func() {
			chunks, err := splitter(text)
			Expect(err).NotTo(HaveOccurred())
			Expect(chunks).To(Equal(expectedOutput))
		})
	})
})
