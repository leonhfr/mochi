package converter

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type transformer struct{}

func newTransformer() parser.ASTTransformer {
	return &transformer{}
}

func (t *transformer) Transform(node *ast.Document, _ text.Reader, pc parser.Context) {
	reader := getReader(pc)
	path := getPath(pc)

	err := ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch node := n.(type) {
		case *ast.Image:
			attachment, err := newAttachment(reader, path, string(node.Destination))
			if err == nil {
				addAttachment(pc, attachment)
				node.Destination = attachment.destination()
			}
		}

		return ast.WalkContinue, nil
	})

	pc.Set(errorKey, err)
}
