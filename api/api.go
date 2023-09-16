package api

import (
	"context"
	"net/http"

	"github.com/carlmjohnson/requests"
)

const baseURL = "https://app.mochi.cards/"

type Client struct {
	baseURL string
	token   string
	client  *http.Client
}

type Option func(c *Client)

func New(token string, options ...Option) *Client {
	c := &Client{
		baseURL: baseURL,
		token:   token,
		client:  http.DefaultClient,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

func WithClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}

func (c *Client) builder() *requests.Builder {
	return requests.
		URL(c.baseURL).
		Client(c.client).
		Accept("application/json").
		BasicAuth(c.token, "")
}

func createItem[Item any](ctx context.Context, c *Client, path string, req any) (Item, error) {
	var res Item
	err := c.builder().
		Path(path).
		Method(http.MethodPost).
		BodyJSON(req).
		ToJSON(&res).
		Fetch(ctx)
	return res, err
}

func getItem[Item any](ctx context.Context, c *Client, path, id string) (Item, error) {
	var res Item
	err := c.builder().
		Pathf("%s/%s", path, id).
		Method(http.MethodGet).
		ToJSON(&res).
		Fetch(ctx)
	return res, err
}

func listItems[Item any](ctx context.Context, c *Client, path string, params map[string][]string) ([]Item, error) {
	type response struct {
		Bookmark string `json:"bookmark"`
		Docs     []Item `json:"docs"`
	}

	var items []Item
	var bookmark string

	for {
		var res response

		b := c.builder().
			Path(path).
			Param("limit", "100")
		if len(params) > 0 {
			b = b.Params(params)
		}
		if len(bookmark) > 0 {
			b = b.Param("bookmark", bookmark)
		}
		err := b.
			Method(http.MethodGet).
			ToJSON(&res).
			Fetch(ctx)
		if err != nil {
			return nil, err
		}

		items = append(items, res.Docs...)

		if bookmark == "" || bookmark == "nil" {
			break
		}
		bookmark = res.Bookmark
	}

	return items, nil
}

func updateItem[Item any](ctx context.Context, c *Client, path, id string, req any) (Item, error) {
	var res Item
	err := c.builder().
		Pathf("%s/%s", path, id).
		Method(http.MethodPost).
		BodyJSON(req).
		ToJSON(&res).
		Fetch(ctx)
	return res, err
}

func deleteItem(ctx context.Context, c *Client, path, id string) error {
	return c.builder().
		Pathf("%s/%s", path, id).
		Method(http.MethodDelete).
		Fetch(ctx)
}
