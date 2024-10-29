package parser

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"

	"github.com/leonhfr/mochi/mochi"
)

// note represents a note parser.
//
// The whole content of the file is returned as a card.
// The is the file name without the extension.
type note struct {
	parser parser.Parser
}

// newNote returns a new note parser.
func newNote() *note {
	return &note{
		parser: parser.NewParser(
			parser.WithBlockParsers(
				parser.DefaultBlockParsers()...,
			),
			parser.WithInlineParsers(
				parser.DefaultInlineParsers()...,
			),
		),
	}
}

// Convert implements the cardParser interface.
func (n *note) convert(path string, source []byte) (Result, error) {
	images := []Image{}
	doc := n.parser.Parse(text.NewReader(source))
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch node := n.(type) {
		case *ast.Image:
			images = append(images, Image{
				Destination: string(node.Destination),
				AltText:     string(node.Text(source)),
			})
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return Result{}, err
	}

	name := getNameFromPath(path)
	return Result{Cards: []Card{newNoteCard(name, path, source, images)}}, nil
}

type noteCard struct {
	name    string
	content string
	images  []Image
	path    string
}

func newNoteCard(name, path string, source []byte, images []Image) Card {
	content := fmt.Sprintf("# %s\n\n%s", name, string(source))
	return noteCard{
		name:    name,
		content: content,
		path:    path,
		images:  images,
	}
}

func (n noteCard) Content() string    { return n.content }
func (n noteCard) Images() []Image    { return n.images }
func (n noteCard) Path() string       { return n.path }
func (n noteCard) Filename() string   { return getFilename(n.path) }
func (n noteCard) Position() string   { return "" }
func (n noteCard) TemplateID() string { return "" }

func (n noteCard) Is(card mochi.Card) bool {
	return nameEquals(card.Fields, n.name)
}

func (n noteCard) Fields() map[string]mochi.Field {
	return map[string]mochi.Field{
		"name": {ID: "name", Value: n.name},
	}
}

func (n noteCard) Equals(card mochi.Card) bool {
	return card.Content == n.content
}

func nameEquals(fields map[string]mochi.Field, name string) bool {
	f, ok := fields["name"]
	if !ok {
		return false
	}
	return f.Value == name
}
