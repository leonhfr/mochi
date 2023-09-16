package api

import "context"

const templatePath = "/api/templates"

type Template struct {
	ID      string                   `json:"id"`
	Name    string                   `json:"name"`
	Content string                   `json:"content"`
	Fields  map[string]FieldTemplate `json:"fields"`
}

type FieldTemplate struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) GetTemplate(ctx context.Context, id string) (Template, error) {
	return getItem[Template](ctx, c, templatePath, id)
}

func (c *Client) ListTemplates(ctx context.Context) ([]Template, error) {
	return listItems[Template](ctx, c, templatePath, nil)
}
