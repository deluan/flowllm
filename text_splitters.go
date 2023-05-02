package flowllm

import (
	"log"
	"strings"
)

var (
	defaultSplitterChunkSize  = 1000
	defaultSplitterLenFunc    = func(s string) int { return len(s) }
	defaultSplitterSeparators = []string{"\n\n", "\n", " ", ""}
)

// SplitterOptions for the RecursiveTextSplitter splitter
type SplitterOptions struct {
	// ChunkSize is the maximum size of each chunk
	ChunkSize int
	// ChunkOverlap is the number of characters that will be repeated in each
	ChunkOverlap int
	// LenFunc is the length function to be used to calculate the chunk size
	LenFunc func(string) int
	// Separators is a list of strings that will be used to split the text
	Separators []string
}

// RecursiveTextSplitter splits a text into chunks of a given size, trying to
// split at the given separators. If the text is smaller than the chunk size,
// it will be returned as a single chunk. If the text is larger than the chunk
// size, it will be split into chunks of the given size, trying to split at the
// given separators. If the text cannot be split at any of the given separators,
// it will be split at the last separator.
func RecursiveTextSplitter(opts SplitterOptions) Splitter {
	if opts.ChunkSize == 0 {
		opts.ChunkSize = defaultSplitterChunkSize
	}
	if opts.LenFunc == nil {
		opts.LenFunc = defaultSplitterLenFunc
	}
	if len(opts.Separators) == 0 {
		opts.Separators = defaultSplitterSeparators
	}
	var splitter Splitter
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
			if opts.LenFunc(split) < opts.ChunkSize { // Use LenFunc here
				goodSplits = append(goodSplits, split)
			} else {
				if len(goodSplits) > 0 {
					mergedText := mergeSplits(goodSplits, separator, opts.ChunkSize, opts.ChunkOverlap, opts.LenFunc) // Pass LenFunc
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
			mergedText := mergeSplits(goodSplits, separator, opts.ChunkSize, opts.ChunkOverlap, opts.LenFunc) // Pass LenFunc
			finalChunks = append(finalChunks, mergedText...)
		}
		return finalChunks, nil
	}
	return splitter
}

// MarkdownSplitter returns a Splitter that splits a document into chunks using a set
// of MarkdownSplitter-specific separators. It is a recursive splitter, meaning that
// it will split each chunk into smaller chunks using the same separators.
func MarkdownSplitter(opts SplitterOptions) Splitter {
	opts.Separators = []string{
		// First, try to split along MarkdownSplitter headings (starting with level 2)
		"\n## ",
		"\n### ",
		"\n#### ",
		"\n##### ",
		"\n###### ",
		// Note the alternative syntax for headings (below) is not handled here
		// Heading level 2
		// ---------------
		// End of code block
		"```\n\n",
		// Horizontal lines
		"\n\n***\n\n",
		"\n\n---\n\n",
		"\n\n___\n\n",
		// Note that this splitter doesn't handle horizontal lines defined
		// by *three or more* of ***, ---, or ___, but this is not handled
		"\n\n",
		"\n",
		" ",
		"",
	}
	return RecursiveTextSplitter(opts)
}

func joinDocs(docs []string, separator string) string {
	return strings.TrimSpace(strings.Join(docs, separator))
}

func mergeSplits(splits []string, separator string, chunkSize int, chunkOverlap int, lenFunc func(string) int) []string {
	var docs []string
	var currentDoc []string
	total := 0

	for _, d := range splits {
		length := lenFunc(d) // Use LenFunc here
		if total+length >= chunkSize {
			if total > chunkSize {
				log.Printf("Created a chunk of size %d, which is longer than the specified %d\n", total, chunkSize)
			}
			if len(currentDoc) > 0 {
				doc := joinDocs(currentDoc, separator)
				if doc != "" {
					docs = append(docs, doc)
				}
				for total > chunkOverlap || (total+length > chunkSize && total > 0) {
					total -= lenFunc(currentDoc[0]) // Use LenFunc here
					currentDoc = currentDoc[1:]
				}
			}
		}
		currentDoc = append(currentDoc, d)
		total += length
	}

	doc := joinDocs(currentDoc, separator)
	if doc != "" {
		docs = append(docs, doc)
	}
	return docs
}
