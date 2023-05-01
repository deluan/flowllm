package splitters

import (
	"strings"

	"github.com/deluan/pipelm"
)

var (
	defaultSeparators = []string{"\n\n", "\n", " ", ""}
	defaultChunkSize  = 1000
	defaultOverlap    = 200
)

type Options struct {
	ChunkSize    int
	ChunkOverlap int
	Separators   []string
}

// RecursiveCharacterText splits a text into chunks of a given size, trying to
// split at the given separators. If the text is smaller than the chunk size,
// it will be returned as a single chunk. If the text is larger than the chunk
// size, it will be split into chunks of the given size, trying to split at the
// given separators. If the text cannot be split at any of the given separators,
// it will be split at the last separator.
func RecursiveCharacterText(opts Options) pipelm.Splitter {
	if opts.ChunkSize == 0 {
		opts.ChunkSize = defaultChunkSize
	}
	if opts.ChunkOverlap == 0 {
		opts.ChunkOverlap = defaultOverlap
	}
	if len(opts.Separators) == 0 {
		opts.Separators = defaultSeparators
	}
	var splitter pipelm.Splitter
	splitter = func(text string) ([]string, error) {
		var separator string
		for _, s := range opts.Separators {
			if s == "" || strings.Contains(text, s) {
				separator = s
				break
			}
		}

		splits := strings.Split(text, separator)
		var finalChunks []string
		var goodSplits []string
		for _, split := range splits {
			if len(split) < opts.ChunkSize {
				goodSplits = append(goodSplits, split)
			} else {
				if len(goodSplits) > 0 {
					mergedText := mergeSplits(goodSplits, separator, opts.ChunkSize, opts.ChunkOverlap)
					finalChunks = append(finalChunks, mergedText...)
					goodSplits = nil
				}
				otherInfo, err := splitter(split)
				if err != nil {
					return nil, err
				}
				finalChunks = append(finalChunks, otherInfo...)
			}
		}

		if len(goodSplits) > 0 {
			mergedText := mergeSplits(goodSplits, separator, opts.ChunkSize, opts.ChunkOverlap)
			finalChunks = append(finalChunks, mergedText...)
		}
		return finalChunks, nil
	}
	return splitter
}
