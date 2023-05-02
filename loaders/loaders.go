package loaders

import "github.com/deluan/flowllm"

func SplitDocuments(splitter flowllm.Splitter, documents []flowllm.Document) ([]flowllm.Document, error) {
	var texts []string
	var metadatas []map[string]any
	for _, document := range documents {
		texts = append(texts, document.PageContent)
		metadatas = append(metadatas, document.Metadata)
	}

	return createDocuments(splitter, texts, metadatas)
}

func createDocuments(splitter flowllm.Splitter, texts []string, metadatas []map[string]any) ([]flowllm.Document, error) {
	if len(metadatas) == 0 {
		metadatas = make([]map[string]any, len(texts))
		for i := range metadatas {
			metadatas[i] = make(map[string]any)
		}
	}

	var documents []flowllm.Document
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
			documents = append(documents, flowllm.Document{
				PageContent: chunk,
				Metadata:    metadata,
			})
		}
	}

	return documents, nil
}
