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
	chunkSize    int
	chunkOverlap int
	separators   []string
}

type Option func(*Options)

// RecursiveCharacterText splits a text into chunks of a given size, trying to
// split at the given separators. If the text is smaller than the chunk size,
// it will be returned as a single chunk. If the text is larger than the chunk
// size, it will be split into chunks of the given size, trying to split at the
// given separators. If the text cannot be split at any of the given separators,
// it will be split at the last separator.
func RecursiveCharacterText(opts ...Option) pipelm.Splitter {
	options := Options{
		chunkSize:    defaultChunkSize,
		chunkOverlap: defaultOverlap,
		separators:   defaultSeparators,
	}
	for _, opt := range opts {
		opt(&options)
	}
	var splitter pipelm.Splitter
	splitter = func(text string) ([]string, error) {
		var separator string
		for _, s := range options.separators {
			if s == "" || strings.Contains(text, s) {
				separator = s
				break
			}
		}

		splits := strings.Split(text, separator)
		var finalChunks []string
		var goodSplits []string
		for _, split := range splits {
			if len(split) < options.chunkSize {
				goodSplits = append(goodSplits, split)
			} else {
				if len(goodSplits) > 0 {
					mergedText := mergeSplits(goodSplits, separator, options.chunkSize, options.chunkOverlap)
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
			mergedText := mergeSplits(goodSplits, separator, options.chunkSize, options.chunkOverlap)
			finalChunks = append(finalChunks, mergedText...)
		}
		return finalChunks, nil
	}
	return splitter
}

func WithChunkSize(chunkSize int) Option {
	return func(o *Options) {
		o.chunkSize = chunkSize
	}
}

func WithChunkOverlap(chunkOverlap int) Option {
	return func(o *Options) {
		o.chunkOverlap = chunkOverlap
	}
}

func WithSeparators(separators ...string) Option {
	return func(o *Options) {
		o.separators = separators
	}
}
