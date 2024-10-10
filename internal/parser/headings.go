package parser

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// headings represents a headings parser.
//
// Each headings until a determined depth returns a separate card.
// The card names are formatted from the card name and the heading.
type headings struct {
	parser parser.Parser
	level  int
}

// newHeadings returns a new note parser.
func newHeadings(level int) *headings {
	return &headings{
		parser: parser.NewParser(
			parser.WithBlockParsers(
				util.Prioritized(parser.NewATXHeadingParser(), 100),
				util.Prioritized(parser.NewHTMLBlockParser(), 200),
				util.Prioritized(parser.NewParagraphParser(), 300),
			),
		),
		level: level,
	}
}

// Convert implements the cardParser interface.
func (h *headings) convert(path string, source []byte) ([]Card, error) {
	var cards []Card
	var title string
	var paragraphs []string

	doc := h.parser.Parse(text.NewReader(source))
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch node := n.(type) {
		case *ast.Heading:
			if level := node.Level; level > h.level {
				bytes := n.Text(source)
				heading := formatHeading(string(bytes), level)
				paragraphs = append(paragraphs, heading)
				return ast.WalkContinue, nil
			}

			if len(title) > 0 {
				cards = append(cards, newHeadingsCard(path, title, h.level, paragraphs))
			}

			title = string(node.Text(source))
			paragraphs = nil
		case *ast.Paragraph:
			bytes := node.Text(source)
			paragraphs = append(paragraphs, string(bytes))
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, err
	}

	cards = append(cards, newHeadingsCard(path, title, h.level, paragraphs))
	return cards, nil
}

func newHeadingsCard(path, name string, level int, paragraphs []string) Card {
	title := formatHeading(name, level)
	content := concatenateBlocks(append([]string{title}, paragraphs...)) + "\n"
	return Card{
		Name:     name,
		Content:  content,
		Filename: getFilename(path),
	}
}

func formatHeading(text string, level int) string {
	format := strings.Repeat("#", level)
	return fmt.Sprintf("%s %s", format, text)
}

func concatenateBlocks(blocks []string) string {
	return strings.Join(blocks, "\n\n")
}
