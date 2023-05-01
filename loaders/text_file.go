package loaders

import (
	"context"
	"io"
	"os"

	"github.com/deluan/pipelm"
)

func TextFile(path string, splitter ...pipelm.Splitter) pipelm.DocumentLoaderFunc {
	var docs []pipelm.Document
	var idx int
	var spl pipelm.Splitter
	if len(splitter) > 0 {
		spl = splitter[0]
	}

	return func(context.Context) (pipelm.Document, error) {
		// Return next document if already loaded
		if len(docs) > 0 {
			if idx < len(docs) {
				idx++
				return docs[idx-1], nil
			}
			return pipelm.Document{}, io.EOF
		}

		// Load file
		text, err := os.ReadFile(path)
		if err != nil {
			return pipelm.Document{}, err
		}
		metadata := map[string]any{"source": path}
		docs = []pipelm.Document{{PageContent: string(text), Metadata: metadata}}

		// Use splitter if provided
		if spl != nil {
			docs, err = SplitDocuments(splitter[0], docs)
			if err != nil {
				return pipelm.Document{}, err
			}
		}

		// Return first document
		idx++
		return docs[0], nil
	}
}
