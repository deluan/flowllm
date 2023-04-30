package splitters

import (
	"log"
	"strings"
)

func joinDocs(docs []string, separator string) string {
	return strings.TrimSpace(strings.Join(docs, separator))
}

func mergeSplits(splits []string, separator string, chunkSize int, chunkOverlap int) []string {
	var docs []string
	var currentDoc []string
	total := 0

	for _, d := range splits {
		length := len(d)
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
					total -= len(currentDoc[0])
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
