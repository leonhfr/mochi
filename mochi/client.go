package mochi

import "net/http"

const baseURL = "https://app.mochi.cards/"

// Client manages communications with mochi.
type Client struct {
	baseURL   string
	token     string
	client    *http.Client
	transport http.RoundTripper
}

// Option represents a Client option.
type Option func(c *Client)

// New creates a new Client with default values.
func New(token string, options ...Option) *Client {
	client := &Client{
		baseURL:   baseURL,
		token:     token,
		client:    http.DefaultClient,
		transport: http.DefaultTransport,
	}
	for _, option := range options {
		option(client)
	}
	return client
}

// WithClient sets the http.Client to use for requests.
func WithClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}

// WithTransport sets the http.RoundTripper to use for requests.
func WithTransport(transport http.RoundTripper) Option {
	return func(c *Client) {
		c.transport = transport
	}
}
