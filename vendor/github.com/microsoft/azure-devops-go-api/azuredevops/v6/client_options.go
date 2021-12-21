package azuredevops

import (
	"net/http"
)

// ClientOptionFunc can be used customize a new AzureDevops API client.
type ClientOptionFunc func(*Client)

// WithHTTPClient can be used to configure a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOptionFunc {
	return func(c *Client) {
		c.client = httpClient
	}
}
