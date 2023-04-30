package loaders

import "github.com/deluan/pipelm"

func SplitDocuments(splitter pipelm.Splitter, documents []pipelm.Document) ([]pipelm.Document, error) {
	var texts []string
	var metadatas []map[string]any
	for _, document := range documents {
		texts = append(texts, document.PageContent)
		metadatas = append(metadatas, document.Metadata)
	}

	return createDocuments(splitter, texts, metadatas)
}

func createDocuments(splitter pipelm.Splitter, texts []string, metadatas []map[string]any) ([]pipelm.Document, error) {
	if len(metadatas) == 0 {
		metadatas = make([]map[string]any, len(texts))
		for i := range metadatas {
			metadatas[i] = make(map[string]any)
		}
	}

	var documents []pipelm.Document
	for i, text := range texts {
		chunks, err := splitter(text)
		if err != nil {
			return nil, err
		}
		for _, chunk := range chunks {
			metadata := make(map[string]any)
			for k, v := range metadatas[i] {
				metadata[k] = v
			}
			documents = append(documents, pipelm.Document{
				PageContent: chunk,
				Metadata:    metadata,
			})
		}
	}

	return documents, nil
}
