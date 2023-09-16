package parser

import (
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"github.com/leonhfr/mochi/parser/frontmatter"
)

const noteName = "note"

type Note struct {
	parser parser.Parser
}

func NewNote() *Note {
	return &Note{
		parser: parser.NewParser(
			parser.WithBlockParsers(
				util.Prioritized(&frontmatter.Parser{}, 0),
				util.Prioritized(parser.NewATXHeadingParser(), 100),
				util.Prioritized(parser.NewHTMLBlockParser(), 200),
				util.Prioritized(parser.NewParagraphParser(), 300),
			),
			parser.WithInlineParsers(
				util.Prioritized(parser.NewLinkParser(), 100),
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
	var fmLength int
	var name string
	images := make(map[string]Image)

	doc := n.parser.Parse(text.NewReader(source))
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if fm, ok := n.(*frontmatter.Node); ok && entering {
			fmLength = frontmatterLength(fm)
			return ast.WalkSkipChildren, nil
		}

		_, isHeading := n.(*ast.Heading)
		_, isParagraph := n.(*ast.Paragraph)

		if name == "" && (isHeading || isParagraph) && entering {
			name = getName(n, source)
		}

		if img, ok := n.(*ast.Image); ok && entering {
			destination := string(img.Destination)
			altText := string(img.Text(source))
			if path, image := newImage(destination, altText); path != "" {
				images[path] = image
			}
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, err
	}

	content := trimFrontmatter(string(source), fmLength)
	content = replaceImages(content, images)
	return []Card{
		{
			Name:    name,
			Content: content,
			Fields:  map[string]string{},
			Images:  images,
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

func trimFrontmatter(content string, length int) string {
	return strings.TrimSpace(content[length:]) + "\n"
}

func frontmatterLength(node *frontmatter.Node) int {
	return node.Segment.Stop + node.DelimCount + 1
}
