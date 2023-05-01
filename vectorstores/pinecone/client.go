package pinecone

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type client struct {
	apiKey      string
	environment string
	namespace   string
	index       string
	projectName string
	pods        int
	podType     string
	metric      string
	replicas    int
}

func connect(ctx context.Context, c *client) error {
	if c.apiKey == "" {
		return fmt.Errorf("no value set for api key. Use WithApiKey when creating a new client")
	}

	if c.environment == "" {
		return fmt.Errorf("no value set for environment. Use WithEnvironment when creating a new client")
	}

	if c.index == "" {
		return fmt.Errorf("no value set for index name. Use WithIndexName when creating a new client")
	}

	// Get project name associated with api using the whoami command
	var err error
	c.projectName, err = c.whoAmI(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) getEndpoint() string {
	urlString := url.QueryEscape(fmt.Sprintf("%s-%s.svc.%s.pinecone.io", c.index, c.projectName, c.environment))
	return "https://" + urlString
}

func doRequest(ctx context.Context, payload any, url, apiKey string) (io.ReadCloser, int, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "text/plain")
	req.Header.Set("Api-Key", apiKey)

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	return r.Body, r.StatusCode, nil
}
