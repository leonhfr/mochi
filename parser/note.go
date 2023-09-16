package parser

import (
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

const noteName = "note"

type Note struct {
	parser parser.Parser
}

func NewNote() *Note {
	return &Note{
		parser: parser.NewParser(
			parser.WithBlockParsers(
				util.Prioritized(parser.NewATXHeadingParser(), 100),
				util.Prioritized(parser.NewHTMLBlockParser(), 200),
				util.Prioritized(parser.NewParagraphParser(), 300),
			),
		),
	}
}

func (n *Note) String() string {
	return noteName
}

func (n *Note) Fields() []string {
	return nil
}

func (n *Note) Convert(source []byte) ([]Card, error) {
	var name string

	doc := n.parser.Parse(text.NewReader(source))
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		_, isHeading := n.(*ast.Heading)
		_, isParagraph := n.(*ast.Paragraph)

		if name == "" && (isHeading || isParagraph) && entering {
			name = getName(n, source)
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, err
	}

	return []Card{
		{
			Name:    name,
			Content: string(source),
			Fields:  map[string]string{},
		},
	}, nil
}

func getName(node ast.Node, source []byte) string {
	var texts []string
	_ = ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if _, ok := n.(*ast.Text); ok && entering {
			texts = append(texts, string(n.Text(source)))
		}
		return ast.WalkContinue, nil
	})
	return strings.Join(texts, "")
}
