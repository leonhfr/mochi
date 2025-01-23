package obsidian

import (
	"regexp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	east "github.com/leonhfr/mochi/internal/converter/obsidian/ast"
)

type calloutExtension struct{}

// NewCallout returns a new Callout extension.
func NewCallout() goldmark.Extender {
	return &calloutExtension{}
}

// Extend implements goldmark.Extender.
func (e *calloutExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(defaultCalloutASTTransformer, 701),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(defaultCalloutRenderer, 701),
		),
	)
}

type calloutASTTransformer struct{}

var (
	defaultCalloutASTTransformer = &calloutASTTransformer{}

	_ parser.ASTTransformer = (*calloutASTTransformer)(nil)

	calloutRegex = regexp.MustCompile(`(?im)^\[\!(.*)\]([\+\-]?) (.*)$`)
)

// NewCalloutASTTransformer returns a new Callout AST transformer.
func NewCalloutASTTransformer() parser.ASTTransformer {
	return defaultCalloutASTTransformer
}

// Transform implements parser.ASTTransformer.
func (t *calloutASTTransformer) Transform(node *ast.Document, reader text.Reader, _ parser.Context) {
	replaceCallout := func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if n.Kind() != ast.KindBlockquote {
			return ast.WalkContinue, nil
		}

		blockquote := n.(*ast.Blockquote)
		if blockquote.ChildCount() == 0 || blockquote.FirstChild().Lines().Len() == 0 {
			return ast.WalkContinue, nil
		}

		segment := blockquote.FirstChild().Lines().At(0)
		line := segment.Value(reader.Source())
		matches := calloutRegex.FindAllSubmatch(line, -1)

		if len(matches) == 0 || len(matches[0]) < 4 {
			return ast.WalkContinue, nil
		}

		calloutType := east.ParseCalloutType(matches[0][1])
		fold := east.ParseFold(matches[0][2])
		callout := east.NewCallout(matches[0][3], calloutType, fold)
		segments := text.NewSegments()
		for i := 1; i < blockquote.Lines().Len(); i++ {
			segments.Append(blockquote.Lines().At(i))
		}
		callout.SetLines(segments)

		callout.Dump(reader.Source(), 0)

		blockquote.Parent().ReplaceChild(blockquote.Parent(), blockquote, callout)

		return ast.WalkContinue, nil
	}

	_ = ast.Walk(node, replaceCallout)
}

type calloutRenderer struct{}

var (
	defaultCalloutRenderer = &calloutRenderer{}

	_ renderer.NodeRenderer = (*calloutRenderer)(nil)
)

// NewCalloutRenderer returns a new Callout renderer.
func NewCalloutRenderer() renderer.NodeRenderer {
	return defaultCalloutRenderer
}

// RegisterFuncs implements renderer.NodeRenderer.
func (r *calloutRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(east.KindCallout, r.render)
}

func (r *calloutRenderer) render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		return ast.WalkContinue, nil
	}

	callout := node.(*east.Callout)

	// callout.Dump(source, 0)

	for i := 0; i < callout.Lines().Len(); i++ {
		line := callout.Lines().At(i)
		_, _ = w.Write(line.Value(source))
	}

	return ast.WalkContinue, nil
}

// pencil <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-pencil"><path d="M21.174 6.812a1 1 0 0 0-3.986-3.987L3.842 16.174a2 2 0 0 0-.5.83l-1.321 4.352a.5.5 0 0 0 .623.622l4.353-1.32a2 2 0 0 0 .83-.497z"/><path d="m15 5 4 4"/></svg>
