package parser

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// table represents a table parser.
//
// Each row returns a separate card.
type table struct {
	parser parser.Parser
}

func newTable() *table {
	p := parser.NewParser(
		parser.WithBlockParsers(
			parser.DefaultBlockParsers()...,
		),
	)
	p.AddOptions(
		parser.WithParagraphTransformers(
			util.Prioritized(extension.NewTableParagraphTransformer(), 200),
		),
		parser.WithASTTransformers(
			util.Prioritized(extension.NewTableASTTransformer(), 200),
		),
	)
	return &table{
		parser: p,
	}
}

func (t *table) parse(path string, source []byte) (Result, error) {
	var headers []string
	var rows [][]string
	var parsingRows, resetRow bool

	doc := t.parser.Parse(text.NewReader(source))
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch node := n.(type) {
		case *east.TableRow:
			parsingRows = true
			resetRow = true
		case *east.TableCell:
			text := string(node.Text(source))
			if parsingRows && resetRow {
				rows = append(rows, []string{text})
			} else if parsingRows {
				rows[len(rows)-1] = append(rows[len(rows)-1], text)
			} else {
				headers = append(headers, text)
			}
			resetRow = false
		}
		return ast.WalkContinue, nil
	})

	return Result{
		Deck:  getNameFromPath(path),
		Cards: getTableCards(path, headers, rows),
	}, err
}

func getTableCards(path string, headers []string, rows [][]string) []Card {
	cards := []Card{}
	for i, row := range rows {
		if len(headers) != len(row) {
			continue
		}
		cards = append(cards, newTableCard(headers, row, path, i))
	}
	return cards
}

func newTableCard(headers, cells []string, path string, index int) Card {
	filename := getFilename(path)
	position := fmt.Sprintf("%s%04d", filename, index)
	return Card{
		Content:  tableContent(headers, cells),
		Fields:   nameFields(strings.Join(cells, "|")),
		Path:     path,
		Position: sanitizePosition(position),
	}
}

func tableContent(headers, cells []string) string {
	rows := []string{"|Headers|Values|", "|---|---|"}
	for i, header := range headers {
		rows = append(rows, fmt.Sprintf("|%s|%s|", header, cells[i]))
	}
	return fmt.Sprintf("%s\n", strings.Join(rows, "\n"))
}
