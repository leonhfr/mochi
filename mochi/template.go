package mochi

import "context"

const templatePath = "/api/templates"

// Template represents a template.
type Template struct {
	ID      string                   `json:"id"`
	Name    string                   `json:"name"`
	Content string                   `json:"content"`
	Fields  map[string]FieldTemplate `json:"fields"`
}

// FieldTemplate represents a field template.
type FieldTemplate struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetTemplate gets a single template.
func (c *Client) GetTemplate(ctx context.Context, id string) (Template, error) {
	return getItem[Template](ctx, c, templatePath, id)
}

// ListTemplates lists the templates.
func (c *Client) ListTemplates(ctx context.Context) ([]Template, error) {
	return listItems[Template](ctx, c, templatePath, nil)
}
