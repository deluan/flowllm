package pinecone

import (
	"bytes"
	"context"
	"fmt"
	"io"
)

type pineconeItem struct {
	Values   []float32         `json:"values"`
	Metadata map[string]string `json:"metadata"`
	ID       string            `json:"id"`
}

type upsertPayload struct {
	Vectors   []pineconeItem `json:"vectors"`
	Namespace string         `json:"namespace,omitempty"`
}

func errorMessageFromErrorResponse(task string, body io.Reader) error {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, body)
	if err != nil {
		return fmt.Errorf("error reading body of error message: %w", err)
	}

	return fmt.Errorf("error %s: body: %s", task, buf.String())
}

func (c *client) upsert(ctx context.Context, vectors []pineconeItem) error {
	payload := upsertPayload{
		Vectors:   vectors,
		Namespace: c.namespace,
	}

	body, status, err := doRequest(ctx, payload, c.getEndpoint()+"/vectors/upsert", c.apiKey)
	if err != nil {
		return err
	}
	defer body.Close()

	if status == 200 {
		return nil
	}

	return errorMessageFromErrorResponse("upserting vectors", body)
}
