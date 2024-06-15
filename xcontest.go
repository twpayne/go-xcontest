// Package xcontest is a client for https://tools.xcontest.org/.
package xcontest

import (
	"fmt"
	"net/http"
)

const defaultBaseURL = "https://tools.xcontest.org/api"

type errExpectedHTTPStatusOK int

func (e errExpectedHTTPStatusOK) Error() string {
	return fmt.Sprintf("expected status code 200 OK, got %d %s", int(e), http.StatusText(int(e)))
}

// A Client is an XContest client.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// A ClientOption sets an option on a Client.
type ClientOption func(*Client)

// NewClient returns a new Client.
func NewClient(options ...ClientOption) *Client {
	client := &Client{
		httpClient: http.DefaultClient,
		baseURL:    defaultBaseURL,
	}
	for _, option := range options {
		option(client)
	}
	return client
}

// WithBaseURL sets the base URL on a Client.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithHTTPClient sets the http.Client used by a Client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}
