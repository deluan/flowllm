package tiktoken_test

import (
	"os"
	"testing"

	"github.com/deluan/pipelm"
	. "github.com/deluan/pipelm/tiktoken"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTokenizer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tokenizer Tests Suite")
}

var _ = Describe("Tokenizer", func() {
	Describe("Len", func() {
		It("returns the number of tokens in a string", func() {
			lenFunc := Len("gpt-3.5-turbo")
			Expect(lenFunc("This is a test")).To(Equal(4))
		})
	})

	Describe("Splitter", func() {
		It("splits using tiktoken token sizes", func() {
			splitter := Splitter("gpt-3.5-turbo", pipelm.SplitterOptions{
				ChunkSize: 100,
			})

			text, err := os.ReadFile("../testdata/state_of_the_union.txt")
			Expect(err).ToNot(HaveOccurred())
			chunks, err := splitter(string(text))
			Expect(err).ToNot(HaveOccurred())
			Expect(chunks).To(HaveLen(98))
			Expect(chunks[0]).To(Equal(`Madam Speaker, Madam Vice President, our First Lady and Second Gentleman. Members of Congress and the Cabinet. Justices of the Supreme Court. My fellow Americans.  

Last year COVID-19 kept us apart. This year we are finally together again. 

Tonight, we meet as Democrats Republicans and Independents. But most importantly as Americans. 

With a duty to one another to the American people to the Constitution.`))
			Expect(chunks[1]).To(Equal(`And with an unwavering resolve that freedom will always triumph over tyranny. 

Six days ago, Russia’s Vladimir Putin sought to shake the foundations of the free world thinking he could make it bend to his menacing ways. But he badly miscalculated. 

He thought he could roll into Ukraine and the world would roll over. Instead he met a wall of strength he never imagined. 

He met the Ukrainian people.`))
		})
	})
})
