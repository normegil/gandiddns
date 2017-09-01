package gandidnsclient

import (
	"io"
	"net/http"
)

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey}
}

const defaultBaseUrl = "https://dns.beta.gandi.net/api/v5/"

func (c Client) request(method string, urlPart string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, defaultBaseUrl+urlPart, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Api-Key", c.apiKey)
	return req, nil
}
