package pinecone

import (
	"context"
	"encoding/json"
	"net/http"
)

type SparseValues struct {
	Indices []int     `json:"indices"`
	Values  []float32 `json:"values"`
}

type match struct {
	pineconeItem
	Score        float32      `json:"score"`
	SparseValues SparseValues `json:"sparseValues"`
}

type queriesResponse struct {
	Matches   []match `json:"matches"`
	Namespace string  `json:"namespace"`
}

type queryPayload struct {
	IncludeValues   bool      `json:"includeValues"`
	IncludeMetadata bool      `json:"includeMetadata"`
	Vector          []float32 `json:"vector"`
	TopK            int       `json:"topK"`
	Namespace       string    `json:"namespace"`
}

func (c *client) query(ctx context.Context, vector []float32, numVectors int) (queriesResponse, error) {
	payload := queryPayload{
		IncludeValues:   true,
		IncludeMetadata: true,
		Vector:          vector,
		TopK:            numVectors,
		Namespace:       c.namespace,
	}

	body, statusCode, err := doRequest(ctx, payload, c.getEndpoint()+"/query", c.apiKey)
	if err != nil {
		return queriesResponse{}, err
	}
	defer body.Close()

	if statusCode != http.StatusOK {
		return queriesResponse{}, errorMessageFromErrorResponse("querying index", body)
	}

	var response queriesResponse

	decoder := json.NewDecoder(body)
	err = decoder.Decode(&response)
	return response, err
}
