package parser

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

const headingsLevel = 1

// headings represents a headings parser.
//
// Each headings until a determined depth returns a separate card.
// The card names are formatted from the card name and the heading.
type headings struct {
	parser parser.Parser
}

// newHeadings returns a new note parser.
func newHeadings() *headings {
	return &headings{
		parser: parser.NewParser(
			parser.WithBlockParsers(
				util.Prioritized(parser.NewATXHeadingParser(), 100),
				util.Prioritized(parser.NewHTMLBlockParser(), 200),
				util.Prioritized(parser.NewParagraphParser(), 300),
			),
		),
	}
}

// Convert implements the cardParser interface.
func (h *headings) convert(path string, source []byte) ([]Card, error) {
	var cards []Card
	var title string
	var paragraphs []string

	doc := h.parser.Parse(text.NewReader(source))
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if heading, ok := n.(*ast.Heading); ok && entering && heading.Level <= headingsLevel {
			if len(title) > 0 {
				cards = append(cards, newHeadingsCard(path, title, paragraphs))
			}

			title = string(heading.Text(source))
			paragraphs = nil

			return ast.WalkContinue, nil
		} else if ok && entering && len(title) > 0 {
			text := string(heading.Text(source))
			paragraphs = append(paragraphs, formatHeading(text, heading.Level))
		} else if ok && entering {
			return ast.WalkStop, fmt.Errorf("malformed card: %s", path)
		}

		if paragraph, ok := n.(*ast.Paragraph); ok && entering {
			text, err := parseParagraph(paragraph, source)
			if err != nil {
				return ast.WalkContinue, err
			}
			if len(text) > 0 {
				paragraphs = append(paragraphs, text)
			}
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, err
	}

	cards = append(cards, newHeadingsCard(path, title, paragraphs))
	return cards, nil
}

func parseParagraph(paragraph *ast.Paragraph, source []byte) (string, error) {
	var lines []string

	err := ast.Walk(paragraph, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if _, ok := n.(*ast.Text); ok && entering {
			if text := string(n.Text(source)); len(text) > 0 {
				lines = append(lines, text)
			}
		}

		return ast.WalkContinue, nil
	})
	return strings.Join(lines, "\n"), err
}

func newHeadingsCard(path, name string, paragraphs []string) Card {
	title := formatHeading(name, headingsLevel)
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
