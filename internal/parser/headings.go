package parser

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"github.com/leonhfr/mochi/internal/parser/image"
)

// headings represents a headings parser.
//
// Each headings until a determined depth returns a separate card.
// The card names are formatted from the card name and the heading.
type headings struct {
	fc     FileCheck
	parser parser.Parser
	level  int
}

// newHeadings returns a new note parser.
func newHeadings(fc FileCheck, level int) *headings {
	return &headings{
		fc: fc,
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
	res := newHeadingResult(h.fc, path, h.level)
	doc := h.parser.Parse(text.NewReader(source))

	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch node := n.(type) {
		case *ast.Heading:
			bytes := node.Text(source)
			if err := res.addHeading(string(bytes), node.Level); err != nil {
				return ast.WalkStop, err
			}
		case *ast.Paragraph:
			bytes := node.Text(source)
			if err := res.addParagraph(string(bytes)); err != nil {
				return ast.WalkStop, err
			}
		}

		return ast.WalkContinue, nil
	})

	return res.getCards(), err
}

type headingResult struct {
	fc         FileCheck
	level      int
	path       string
	name       string
	headings   []string
	paragraphs []string
	cards      []Card
	images     image.Map
}

func newHeadingResult(fc FileCheck, path string, level int) *headingResult {
	return &headingResult{
		fc:     fc,
		level:  level,
		path:   path,
		name:   getNameFromPath(path),
		images: image.New(fc, path),
	}
}

func (r *headingResult) addHeading(text string, level int) (err error) {
	if level > r.level {
		heading := formatHeading(text, level)
		r.paragraphs = append(r.paragraphs, heading)
		return
	}

	if len(r.headings) == 0 {
		heading := formatHeading(text, level)
		r.paragraphs = append(r.paragraphs, heading)
		r.headings = append(r.headings, text)
		return
	}

	r.flushCard()
	heading := formatHeading(text, level)
	r.paragraphs = append(r.paragraphs, heading)
	r.headings = append(r.headings, text)

	return
}

func (r *headingResult) addParagraph(text string) (err error) {
	r.paragraphs = append(r.paragraphs, text)
	return
}

func (r *headingResult) flushCard() {
	name := strings.Join(append([]string{r.name}, r.headings...), " | ")
	r.cards = append(r.cards, Card{
		Name:     name,
		Content:  strings.Join(r.paragraphs, "\n\n") + "\n",
		Filename: getFilename(r.path),
		Images:   r.images,
	})
	r.headings = nil
	r.paragraphs = nil
	r.images = image.New(r.fc, r.path)
}

func (r *headingResult) getCards() []Card {
	r.flushCard()
	return r.cards
}

func formatHeading(text string, level int) string {
	format := strings.Repeat("#", level)
	return fmt.Sprintf("%s %s", format, text)
}
