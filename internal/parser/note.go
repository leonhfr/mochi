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
	parsed := []image.Parsed{}
	doc := n.parser.Parse(text.NewReader(source))
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch node := n.(type) {
		case *ast.Image:
			parsed = append(parsed, image.Parsed{
				Destination: string(node.Destination),
				AltText:     string(node.Text(source)),
			})
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, err
	}

	name := getNameFromPath(path)
	images := image.NewMap(n.fc, path, parsed)
	return []Card{createNoteCard(name, path, source, images)}, nil
}

func createNoteCard(name, path string, source []byte, images map[string]image.Image) Card {
	return Card{
		Name:     name,
		Content:  image.Replace(images, string(source)),
		Filename: getFilename(path),
		Images:   images,
	}
}
