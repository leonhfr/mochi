package parser

import (
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/parser/example"
	"github.com/leonhfr/mochi/mochi"
)

// vocabulary represents a vocabulary parser.
//
// Each word returns a separate card.
type vocabulary struct {
	parser parser.Parser
	config config.VocabularyTemplate
}

func newVocabulary(config config.VocabularyTemplate) *vocabulary {
	return &vocabulary{
		parser: parser.NewParser(
			parser.WithBlockParsers(
				parser.DefaultBlockParsers()...,
			),
			parser.WithInlineParsers(
				util.Prioritized(example.NewParser(), 100),
			),
		),
		config: config,
	}
}

func (v *vocabulary) convert(path string, source []byte) (Result, error) {
	cards := []Card{}
	doc := v.parser.Parse(text.NewReader(source))
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch node := n.(type) {
		case *ast.Paragraph:
			word, examples, notes, err := parseParagraph(node, source)
			if err != nil {
				return ast.WalkStop, err
			}
			cards = append(cards, newVocabularyCard(word, examples, notes, path, v.config))
		}

		return ast.WalkContinue, nil
	})
	return Result{Cards: cards}, err
}

func parseParagraph(paragraph *ast.Paragraph, source []byte) (string, []string, []string, error) {
	var word string
	var examples []string
	var notes []string

	err := ast.Walk(paragraph, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch n.(type) {
		case *example.Node:
			example := strings.TrimSpace(string(n.Text(source)))
			examples = append(examples, example)
			return ast.WalkSkipChildren, nil
		case *ast.Text:
			text := strings.TrimSpace(string(n.Text(source)))
			if text == "" {
				return ast.WalkSkipChildren, nil
			}
			if word == "" {
				word = text
			} else {
				notes = append(notes, text)
			}
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return "", nil, nil, err
	}

	return word, examples, notes, nil
}

type vocabularyCard struct {
	config   config.VocabularyTemplate
	word     string
	examples []string
	notes    []string
	path     string
}

func newVocabularyCard(word string, examples, notes []string, path string, config config.VocabularyTemplate) vocabularyCard {
	return vocabularyCard{
		config:   config,
		word:     word,
		examples: examples,
		notes:    notes,
		path:     path,
	}
}

func (c vocabularyCard) Content() string    { return "" }
func (c vocabularyCard) TemplateID() string { return c.config.TemplateID }
func (c vocabularyCard) Images() []Image    { return nil }
func (c vocabularyCard) Path() string       { return c.path }
func (c vocabularyCard) Filename() string   { return getFilename(c.path) }
func (c vocabularyCard) Position() string   { return sanitizePosition(c.word) }

func (c vocabularyCard) Is(card mochi.Card) bool {
	return nameEquals(card.Fields, c.word)
}

func (c vocabularyCard) Fields() map[string]mochi.Field {
	fields := map[string]mochi.Field{
		"name": {
			ID:    "name",
			Value: c.word,
		},
	}
	if c.config.ExamplesID != "" {
		fields[c.config.ExamplesID] = mochi.Field{
			ID:    c.config.ExamplesID,
			Value: strings.Join(c.examples, "\n\n"),
		}
	}
	if c.config.NotesID != "" {
		fields[c.config.NotesID] = mochi.Field{
			ID:    c.config.NotesID,
			Value: strings.Join(c.notes, "\n\n"),
		}
	}
	return fields
}

func (c vocabularyCard) Equals(card mochi.Card) bool {
	return card.Name == c.word &&
		card.TemplateID == c.TemplateID() &&
		mapsEqual(card.Fields, c.Fields())
}

func mapsEqual[T comparable](m1, m2 map[string]T) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok || v1 != v2 {
			return false
		}
	}
	return true
}
