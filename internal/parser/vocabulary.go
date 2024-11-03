package parser

import (
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/parser/example"
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

func (v *vocabulary) parse(path string, source []byte) (Result, error) {
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

func newVocabularyCard(word string, examples, notes []string, path string, config config.VocabularyTemplate) Card {
	return Card{
		Fields:     vocabularyFields(word, examples, notes, config),
		TemplateID: config.TemplateID,
		Path:       path,
		Position:   sanitizePosition(word),
	}
}

func vocabularyFields(word string, examples, notes []string, config config.VocabularyTemplate) map[string]string {
	fields := map[string]string{
		"name": word,
	}
	if value := strings.Join(examples, "\n\n"); config.ExamplesID != "" && value != "" {
		fields[config.ExamplesID] = value
	}
	if value := strings.Join(notes, "\n\n"); config.NotesID != "" && value != "" {
		fields[config.NotesID] = value
	}
	return fields
}
