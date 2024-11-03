package converter

import (
	"bytes"
	"io"

	markdown "github.com/teekennedy/goldmark-markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"

	"github.com/leonhfr/mochi/internal/converter/ast"
)

// Reader represents the interface to read files.
type Reader interface {
	Read(path string) (io.ReadCloser, error)
}

// Result represents the result of a conversion.
type Result struct {
	Markdown    string
	Attachments []Attachment
}

// Converter converts markdown to mochi markdown.
type Converter struct {
	markdown goldmark.Markdown
}

// New returns a new Converter.
func New() *Converter {
	renderer := markdown.NewRenderer()
	renderer.Register(ast.KindVideo, newVideoRenderer())

	return &Converter{
		markdown: goldmark.New(
			goldmark.WithRenderer(renderer),
			goldmark.WithParserOptions(
				parser.WithASTTransformers(
					util.Prioritized(newTransformer(), 999),
				),
			),
		),
	}
}

// Convert converts the source markdown to mochi markdown.
func (c *Converter) Convert(reader Reader, path, source string) (Result, error) {
	ctx := newContext(reader, path)
	b := bytes.NewBuffer(nil)
	err := c.markdown.Convert([]byte(source), b, parser.WithContext(ctx))
	if err != nil {
		return Result{}, err
	}
	return Result{
		Markdown:    b.String(),
		Attachments: getAttachments(ctx),
	}, getError(ctx)
}
