package parser

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"

	"github.com/leonhfr/mochi/mochi"
)

// headings represents a headings parser.
//
// Each headings until a determined depth returns a separate card.
// The card names are formatted from the card name and the heading.
type headings struct {
	parser   parser.Parser
	maxLevel int
}

// newHeadings returns a new note parser.
func newHeadings(maxLevel int) *headings {
	return &headings{
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
func (h *headings) convert(path string, source []byte) (Result, error) {
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
			parsed[len(parsed)-1].images = append(parsed[len(parsed)-1].images, Image{
				Destination: string(node.Destination),
				AltText:     string(node.Text(source)),
			})
		}

		return ast.WalkContinue, nil
	})

	cards := getHeadingCards(path, parsed, source)

	return Result{
		Deck:  getNameFromPath(path),
		Cards: cards,
	}, err
}

func getHeadingCards(path string, headings []parsedHeading, source []byte) []Card {
	if len(headings) == 0 {
		return nil
	}

	if len(headings) == 1 && len(source) > 0 {
		name := getNameFromPath(path)
		return []Card{newNoteCard(name, path, source, headings[0].images)}
	} else if len(headings) == 1 {
		return nil
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
			start = stop
			continue
		}

		cards = append(cards, newHeadingsCard(titles, path, content, heading.images, len(cards)))
		start = stop
	}

	return cards
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
	images []Image
}

type headingsCard struct {
	name     string
	content  string
	images   []Image
	path     string
	position string
}

func newHeadingsCard(headings []string, path string, source []byte, images []Image, index int) Card {
	filename := getFilename(path)
	name := strings.ReplaceAll(strings.Join(headings, " > "), " >  > ", " > ")
	content := fmt.Sprintf("%s\n\n%s\n", name, string(source))
	return headingsCard{
		name:     name,
		content:  string(content),
		images:   images,
		path:     path,
		position: getPosition(filename, index),
	}
}

func (h headingsCard) Name() string       { return h.name }
func (h headingsCard) Content() string    { return h.content }
func (h headingsCard) Images() []Image    { return h.images }
func (h headingsCard) Path() string       { return h.path }
func (h headingsCard) Filename() string   { return getFilename(h.path) }
func (h headingsCard) Position() string   { return h.position }
func (h headingsCard) TemplateID() string { return "" }

func (h headingsCard) Fields() map[string]mochi.Field {
	return map[string]mochi.Field{
		"name": {ID: "name", Value: h.name},
	}
}

func (h headingsCard) Equals(card mochi.Card) bool {
	return card.Content == h.content && card.Pos == h.position
}
