package loaders

import (
	"context"
	"io"
	"os"

	"github.com/deluan/flowllm"
)

func TextFile(path string, splitter ...flowllm.Splitter) flowllm.DocumentLoaderFunc {
	var docs []flowllm.Document
	var idx int
	var spl flowllm.Splitter
	if len(splitter) > 0 {
		spl = splitter[0]
	}

	return func(context.Context) (flowllm.Document, error) {
		// Return next document if already loaded
		if len(docs) > 0 {
			if idx < len(docs) {
				idx++
				return docs[idx-1], nil
			}
			return flowllm.Document{}, io.EOF
		}

		// Load file
		text, err := os.ReadFile(path)
		if err != nil {
			return flowllm.Document{}, err
		}
		metadata := map[string]any{"source": path}
		docs = []flowllm.Document{{PageContent: string(text), Metadata: metadata}}

		// Use splitter if provided
		if spl != nil {
			docs, err = SplitDocuments(splitter[0], docs)
			if err != nil {
				return flowllm.Document{}, err
			}
		}

		// Return first document
		idx++
		return docs[0], nil
	}
}
