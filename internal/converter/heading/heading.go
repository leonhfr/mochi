package heading

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type headerExtension struct{}

// New returns a new Header extension.
func New() goldmark.Extender {
	return &headerExtension{}
}

// Extend implements goldmark.Extender.
func (e *headerExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(defaultASTTransformer, 601),
		),
	)
}

type astTransformer struct{}

var (
	defaultASTTransformer = &astTransformer{}

	_ parser.ASTTransformer = (*astTransformer)(nil)
)

// Transform implements parser.ASTTransformer.
func (t *astTransformer) Transform(node *ast.Document, reader text.Reader, _ parser.Context) {
	removeHeadingNumbering := func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if n.Kind() != ast.KindHeading {
			return ast.WalkContinue, nil
		}

		heading := n.(*ast.Heading)

		lines := heading.Lines()
		for i := 0; i < lines.Len(); i++ {
			segment := lines.At(i)
			index := bytes.IndexFunc(segment.Value(reader.Source()), func(r rune) bool {
				return (r < '0' || '9' < r) && r != '.'
			})

			if index > 0 {
				txt := ast.NewText()
				txt.Segment = text.NewSegment(segment.Start+index+1, segment.Stop)
				heading.ReplaceChild(heading, heading.FirstChild(), txt)
			}
		}

		return ast.WalkContinue, nil
	}

	_ = ast.Walk(node, removeHeadingNumbering)
}
