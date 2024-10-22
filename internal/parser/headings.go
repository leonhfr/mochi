package parser

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"

	"github.com/leonhfr/mochi/internal/image"
)

// headings represents a headings parser.
//
// Each headings until a determined depth returns a separate card.
// The card names are formatted from the card name and the heading.
type headings struct {
	fc       FileCheck
	parser   parser.Parser
	maxLevel int
}

// newHeadings returns a new note parser.
func newHeadings(fc FileCheck, maxLevel int) *headings {
	return &headings{
		fc: fc,
		parser: parser.NewParser(
			parser.WithBlockParsers(
				parser.DefaultBlockParsers()...,
			),
			parser.WithInlineParsers(
				parser.DefaultInlineParsers()...,
			),
		),
		maxLevel: maxLevel,
	}
}

// convert implements the cardParser interface.
func (h *headings) convert(path string, source []byte) ([]Card, error) {
	parsed := []parsedHeading{{level: 0}}
	doc := h.parser.Parse(text.NewReader(source))

	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch node := n.(type) {
		case *ast.Heading:
			segment := node.Lines().At(0)
			if level := node.Level; level <= h.maxLevel {
				parsed = append(parsed, parsedHeading{
					start: segment.Start,
					stop:  segment.Stop,
					level: node.Level,
				})
			}
		case *ast.Image:
			parsed[len(parsed)-1].images = append(parsed[len(parsed)-1].images, image.Parsed{
				Destination: string(node.Destination),
				AltText:     string(node.Text(source)),
			})
		}

		return ast.WalkContinue, nil
	})

	cards := getHeadingCards(h.fc, path, parsed, source)

	return cards, err
}

func getHeadingCards(fc FileCheck, path string, headings []parsedHeading, source []byte) []Card {
	if len(headings) == 0 {
		return nil
	}

	if len(headings) == 1 {
		name := getNameFromPath(path)
		images := image.NewMap(fc, path, headings[0].images)
		return []Card{createNoteCard(name, path, source, images)}
	}

	cards := []Card{}
	titles := []string{}
	var start int

	for i, heading := range headings {
		switch {
		case heading.level == 0:
			titles = append(titles, getNameFromPath(path))
		case heading.level > len(titles):
			for heading.level > len(titles) {
				titles = append(titles, "")
			}
			titles = append(titles, getHeadingText(heading, source))
		default:
			for heading.level < len(titles) {
				titles = titles[:len(titles)-1]
			}
			titles = append(titles, getHeadingText(heading, source))
		}

		stop := len(source)
		if i < len(headings)-1 {
			stop = getHeadingStart(headings[i+1])
		}

		content := bytes.TrimSpace(source[start:stop])
		if !bytes.ContainsRune(content, '\n') {
			continue
		}

		cards = append(cards, createHeadingCard(fc, titles, path, content, heading.images, len(cards)))
		start = stop
	}

	return cards
}

func createHeadingCard(fc FileCheck, headings []string, path string, source []byte, images []image.Parsed, index int) Card {
	name := strings.ReplaceAll(strings.Join(headings, " | "), " |  | ", " | ")
	content := fmt.Sprintf("%s\n\n%s\n", name, string(source))
	imageMap := image.NewMap(fc, path, images)
	return Card{
		Name:     name,
		Content:  image.Replace(imageMap, string(content)),
		Filename: getFilename(path),
		Images:   imageMap,
		Index:    index,
	}
}

func getHeadingText(heading parsedHeading, source []byte) string {
	return string(source[heading.start:heading.stop])
}

func getHeadingStart(heading parsedHeading) int {
	return heading.start - heading.level - 1
}

type parsedHeading struct {
	level  int
	start  int
	stop   int
	images []image.Parsed
}
