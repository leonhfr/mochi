package parser

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
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
func (h *headings) parse(path string, source []byte) (Result, error) {
	parsed := []parsedHeading{{level: 0}}
	doc := h.parser.Parse(text.NewReader(source))
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch node := n.(type) {
		case *ast.Heading:
			if level := node.Level; level <= h.maxLevel && node.Lines().Len() > 0 {
				segment := node.Lines().At(0)
				parsed = append(parsed, parsedHeading{
					start: segment.Start,
					stop:  segment.Stop,
					level: node.Level,
				})
			}
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
		return []Card{newNoteCard(name, path, source)}
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

		cards = append(cards, newHeadingsCard(titles, path, content, len(cards)))
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
	level int
	start int
	stop  int
}

func newHeadingsCard(headings []string, path string, source []byte, index int) Card {
	filename := getFilename(path)
	position := fmt.Sprintf("%s%04d", filename, index)
	name := strings.ReplaceAll(strings.Join(headings, " > "), " >  > ", " > ")
	content := fmt.Sprintf("%s\n\n%s\n", name, string(source))
	return Card{
		Content:  string(content),
		Fields:   nameFields(name),
		Path:     path,
		Position: sanitizePosition(position),
	}
}
