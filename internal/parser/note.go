package parser

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"github.com/leonhfr/mochi/internal/image"
)

// note represents a note parser.
//
// The whole content of the file is returned as a card.
// The is the file name without the extension.
type note struct {
	fc     FileCheck
	parser parser.Parser
}

// newNote returns a new note parser.
func newNote(fc FileCheck) *note {
	return &note{
		fc: fc,
		parser: parser.NewParser(
			parser.WithBlockParsers(
				util.Prioritized(parser.NewParagraphParser(), 100),
			),
			parser.WithInlineParsers(
				util.Prioritized(parser.NewLinkParser(), 100),
			),
		),
	}
}

// Convert implements the cardParser interface.
func (n *note) convert(path string, source []byte) ([]Card, error) {
	images := image.New(n.fc, path)
	doc := n.parser.Parse(text.NewReader(source))
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch node := n.(type) {
		case *ast.Image:
			destination := string(node.Destination)
			altText := string(node.Text(source))
			images.Add(destination, altText)
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, err
	}

	return []Card{
		{
			Name:     getNameFromPath(path),
			Content:  images.Replace(string(source)),
			Filename: getFilename(path),
			Images:   images.Images(),
		},
	}, nil
}
