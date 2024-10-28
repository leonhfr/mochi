package parser

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"

	"github.com/leonhfr/mochi/mochi"
)

// vocabulary represents a vocabulary parser.
//
// Each word returns a separate card.
type vocabulary struct {
	parser     parser.Parser
	templateID string
}

func newVocabulary(templateID string) *vocabulary {
	return &vocabulary{
		parser: parser.NewParser(
			parser.WithBlockParsers(
				parser.DefaultBlockParsers()...,
			),
			parser.WithInlineParsers(
				parser.DefaultInlineParsers()...,
			),
		),
		templateID: templateID,
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
			word := string(node.Text(source))
			cards = append(cards, newVocabularyCard(word, path, v.templateID))
		}

		return ast.WalkContinue, nil
	})
	return Result{Cards: cards}, err
}

type vocabularyCard struct {
	templateID string
	word       string
	path       string
}

func newVocabularyCard(word, path, templateID string) vocabularyCard {
	return vocabularyCard{
		templateID: templateID,
		word:       word,
		path:       path,
	}
}

func (c vocabularyCard) Name() string       { return c.word }
func (c vocabularyCard) Content() string    { return "" }
func (c vocabularyCard) TemplateID() string { return c.templateID }
func (c vocabularyCard) Images() []Image    { return nil }
func (c vocabularyCard) Path() string       { return c.path }
func (c vocabularyCard) Filename() string   { return getFilename(c.path) }
func (c vocabularyCard) Position() string   { return c.word }

func (c vocabularyCard) Fields() map[string]mochi.Field {
	return map[string]mochi.Field{
		"name": {
			ID:    "name",
			Value: c.word,
		},
	}
}

func (c vocabularyCard) Equals(card mochi.Card) bool {
	return card.Name == c.Name() &&
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
