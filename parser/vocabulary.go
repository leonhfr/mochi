package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"github.com/leonhfr/mochi/parser/example"
	"github.com/leonhfr/mochi/parser/frontmatter"
)

const (
	vocabularyName          = "vocabulary"
	vocabularyFieldWord     = "word"
	vocabularyFieldExamples = "examples"
	vocabularyFieldNotes    = "notes"
)

type Vocabulary struct {
	parser parser.Parser
}

func NewVocabulary() *Vocabulary {
	return &Vocabulary{
		parser: parser.NewParser(
			parser.WithBlockParsers(
				util.Prioritized(&frontmatter.Parser{}, 0),
				util.Prioritized(parser.NewATXHeadingParser(), 100),
				util.Prioritized(parser.NewHTMLBlockParser(), 200),
				util.Prioritized(parser.NewParagraphParser(), 300),
			),
			parser.WithInlineParsers(
				util.Prioritized(example.NewParser(), 100),
			),
		),
	}
}

func (v *Vocabulary) String() string {
	return vocabularyName
}

func (v *Vocabulary) Fields() []string {
	return []string{
		vocabularyFieldWord,
		vocabularyFieldExamples,
		vocabularyFieldNotes,
	}
}

func (v *Vocabulary) Convert(_ string, source []byte) ([]Card, error) {
	var cards []Card
	doc := v.parser.Parse(text.NewReader(source))

	if err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if paragraph, ok := n.(*ast.Paragraph); ok && entering {
			card, err := newVocabularyCard(paragraph, source)
			if err != nil {
				return ast.WalkStop, err
			}
			cards = append(cards, card)
		}

		return ast.WalkContinue, nil
	}); err != nil {
		return nil, err
	}

	return cards, nil
}

func newVocabularyCard(paragraph *ast.Paragraph, source []byte) (Card, error) {
	var word string
	var examples []string
	var notes []string

	if err := ast.Walk(paragraph, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch n.(type) {
		case *example.Node:
			example := strings.TrimSpace(string(n.Text(source)))
			examples = append(examples, example)
			return ast.WalkSkipChildren, nil
		case *ast.Text:
			if text := strings.TrimSpace(string(n.Text(source))); word == "" {
				word = text
			} else {
				notes = append(notes, text)
			}
		default:
		}

		return ast.WalkContinue, nil
	}); err != nil {
		return Card{}, err
	}

	if len(word) == 0 {
		return Card{}, errors.New("word not found")
	}

	return Card{
		Name:    word,
		Content: formatVocabularyContent(word, examples, notes),
		Fields: map[string]string{
			vocabularyFieldWord:     word,
			vocabularyFieldExamples: concatenateBlocks(examples),
			vocabularyFieldNotes:    concatenateBlocks(notes),
		},
	}, nil
}

func formatVocabularyContent(word string, examples, notes []string) string {
	blocks := []string{fmt.Sprintf("# %s", word)}
	if len(examples) > 0 {
		blocks = append(blocks, "## Examples")
		for _, example := range examples {
			if len(example) > 0 {
				blocks = append(blocks, example)
			}
		}
	}
	if len(notes) > 0 {
		blocks = append(blocks, "## Notes")
		for _, note := range notes {
			if len(note) > 0 {
				blocks = append(blocks, note)
			}
		}
	}
	return strings.Join(blocks, "\n\n") + "\n"
}

func concatenateBlocks(blocks []string) string {
	return strings.Join(blocks, "\n\n")
}
