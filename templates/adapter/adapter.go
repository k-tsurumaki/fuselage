package adapter

import (
	"encoding/json"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ExternalAPIAdapter struct {
	client  HTTPClient
	baseURL string
}

func NewExternalAPIAdapter(client HTTPClient, baseURL string) *ExternalAPIAdapter {
	return &ExternalAPIAdapter{
		client:  client,
		baseURL: baseURL,
	}
}

func (a *ExternalAPIAdapter) FetchData(endpoint string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", a.baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}