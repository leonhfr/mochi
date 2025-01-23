package obsidian

import "github.com/yuin/goldmark"

type obsidianExtension struct{}

// New returns a new Obsidian extension.
func New() goldmark.Extender {
	return &obsidianExtension{}
}

// Extend implements goldmark.Extender.
func (e *obsidianExtension) Extend(m goldmark.Markdown) {
	NewCallout().Extend(m)
}
