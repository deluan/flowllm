package pinecone

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type whoAmIResponse struct {
	ProjectName string `json:"project_name"`
	UserLabel   string `json:"user_label"`
	UserName    string `json:"user_name"`
}

func (c *client) whoAmI(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://controller.%s.pinecone.io/actions/whoami", c.environment), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Api-Key", c.apiKey)

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	var response whoAmIResponse

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&response)
	return response.ProjectName, err
}
