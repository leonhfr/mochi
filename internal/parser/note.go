package parser

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
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
	return Result{Cards: []Card{createNoteCard(name, path, source, images)}}, nil
}

func createNoteCard(name, path string, source []byte, images []Image) Card {
	return Card{
		Name:     name,
		Content:  string(source),
		Filename: getFilename(path),
		Path:     path,
		Images:   images,
	}
}
