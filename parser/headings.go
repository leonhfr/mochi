package parser

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"github.com/leonhfr/mochi/parser/frontmatter"
)

const (
	headingsName  = "headings"
	headingsLevel = 1
)

type Headings struct {
	parser parser.Parser
}

func NewHeadings() *Headings {
	return &Headings{
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

func (h *Headings) String() string {
	return headingsName
}

func (h *Headings) Fields() []string {
	return nil
}

func (h *Headings) Convert(path string, source []byte) ([]Card, error) {
	var cards []Card
	var title string
	var paragraphs []string
	var images map[string]Image

	doc := h.parser.Parse(text.NewReader(source))
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if _, ok := n.(*frontmatter.Node); ok {
			return ast.WalkSkipChildren, nil
		}

		if heading, ok := n.(*ast.Heading); ok && entering && heading.Level <= headingsLevel {
			if len(title) > 0 {
				cards = append(cards, newHeadingsCard(title, paragraphs, images))
			}

			title = string(heading.Text(source))
			paragraphs = nil
			images = make(map[string]Image)

			return ast.WalkContinue, nil
		} else if ok && entering {
			text := string(heading.Text(source))
			paragraphs = append(paragraphs, formatHeading(text, heading.Level))
		}

		if paragraph, ok := n.(*ast.Paragraph); ok && entering {
			text, imagesMap, err := parseParagraph(paragraph, path, source)
			if err != nil {
				return ast.WalkContinue, err
			}
			if len(text) > 0 {
				paragraphs = append(paragraphs, text)
			}
			for k, v := range imagesMap {
				images[k] = v
			}
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, err
	}

	cards = append(cards, newHeadingsCard(title, paragraphs, images))
	return cards, nil
}

func parseParagraph(paragraph *ast.Paragraph, path string, source []byte) (string, map[string]Image, error) {
	var lines []string
	images := make(map[string]Image)
	err := ast.Walk(paragraph, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if image, ok := n.(*ast.Image); ok && entering {
			destination := string(image.Destination)
			altText := string(image.Text(source))
			if absPath, image := newImage(path, destination, altText); absPath != "" {
				images[absPath] = image
			}
			text := fmt.Sprintf("![%s](%s)", altText, destination)
			lines = append(lines, text)
			return ast.WalkSkipChildren, nil
		}

		if link, ok := n.(*ast.Link); ok && entering {
			destination := string(link.Destination)
			text := string(link.Text(source))
			lines = append(lines, fmt.Sprintf("[%s](%s)", text, destination))
			return ast.WalkSkipChildren, nil
		}

		if _, ok := n.(*ast.Text); ok && entering {
			if text := string(n.Text(source)); len(text) > 0 {
				lines = append(lines, text)
			}
		}

		return ast.WalkContinue, nil
	})
	return strings.Join(lines, "\n"), images, err
}

func newHeadingsCard(name string, paragraphs []string, images map[string]Image) Card {
	title := formatHeading(name, headingsLevel)
	content := concatenateBlocks(append([]string{title}, paragraphs...)) + "\n"
	content = replaceImages(content, images)
	content = replaceVideos(content)
	return Card{
		Name:    name,
		Content: content,
		Fields:  map[string]string{},
		Images:  images,
	}
}

func formatHeading(text string, level int) string {
	format := strings.Repeat("#", level)
	return fmt.Sprintf("%s %s", format, text)
}
