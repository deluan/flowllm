package loaders_test

import (
	"context"
	"errors"
	"io"

	"github.com/deluan/pipelm"
	"github.com/deluan/pipelm/loaders"
	"github.com/deluan/pipelm/splitters"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TextFile", func() {
	var (
		ctx    context.Context
		loader pipelm.DocumentLoader
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
			loader = loaders.TextFile("../testdata/state_of_the_union.txt", splitters.RecursiveCharacterText())
		})
		It("loads a text file into multiple Documents", func() {
			doc, err := loader.LoadNext(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(doc.PageContent).To(Equal("Madam Speaker, Madam Vice President, our First Lady and Second Gentleman. Members of Congress and the Cabinet. Justices of the Supreme Court. My fellow Americans.  \n\nLast year COVID-19 kept us apart. This year we are finally together again. \n\nTonight, we meet as Democrats Republicans and Independents. But most importantly as Americans. \n\nWith a duty to one another to the American people to the Constitution. \n\nAnd with an unwavering resolve that freedom will always triumph over tyranny. \n\nSix days ago, Russia’s Vladimir Putin sought to shake the foundations of the free world thinking he could make it bend to his menacing ways. But he badly miscalculated. \n\nHe thought he could roll into Ukraine and the world would roll over. Instead he met a wall of strength he never imagined. \n\nHe met the Ukrainian people. \n\nFrom President Zelenskyy to every Ukrainian, their fearlessness, their courage, their determination, inspires the world."))
			total := 1
			for {
				doc, err = loader.LoadNext(ctx)
				if errors.Is(err, io.EOF) {
					break
				}
				total++
				Expect(doc.Metadata).To(HaveKeyWithValue("source", "../testdata/state_of_the_union.txt"))
			}
			Expect(total).To(Equal(49))
		})
	})
})
