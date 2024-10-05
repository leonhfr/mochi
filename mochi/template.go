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
//
// The callback is called with a slice of templates until all templates have been listed
// or until the callback returns an error. Each callback call makes a HTTP request.
func (c *Client) ListTemplates(ctx context.Context, cb func([]Template) error) error {
	return listItems(ctx, c, templatePath, nil, cb)
}
