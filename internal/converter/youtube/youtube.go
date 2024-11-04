package youtube

import (
	"fmt"
	"net/url"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type youtubeExtension struct{}

// New returns a new YouTube extension.
func New() goldmark.Extender {
	return &youtubeExtension{}
}

func (e *youtubeExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(defaultASTTransformer, 500),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(NewHTMLRenderer(), 500),
		),
	)
}

// YouTube struct represents a YouTube Video embed of the Markdown text.
type YouTube struct {
	ast.Image
	Video string
}

// NewYouTube returns a new YouTube node.
func NewYouTube(img *ast.Image, v string) *YouTube {
	c := &YouTube{
		Image: *img,
		Video: v,
	}
	c.Destination = img.Destination
	c.Title = img.Title

	return c
}

// KindYouTube is a NodeKind of the YouTube node.
var KindYouTube = ast.NewNodeKind("YouTube")

// Kind implements Node.Kind.
func (n *YouTube) Kind() ast.NodeKind {
	return KindYouTube
}

type astTransformer struct{}

var defaultASTTransformer = &astTransformer{}

func (a *astTransformer) Transform(node *ast.Document, _ text.Reader, _ parser.Context) {
	replaceImages := func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if n.Kind() != ast.KindImage {
			return ast.WalkContinue, nil
		}

		img := n.(*ast.Image)
		u, err := url.Parse(string(img.Destination))
		if err != nil {
			msg := ast.NewString([]byte(fmt.Sprintf("<!-- %s -->", err)))
			msg.SetCode(true)
			n.Parent().InsertAfter(n.Parent(), n, msg)
			return ast.WalkContinue, nil
		}

		if u.Host != "www.youtube.com" || u.Path != "/watch" {
			return ast.WalkContinue, nil
		}
		v := u.Query().Get("v")
		if v == "" {
			return ast.WalkContinue, nil
		}
		yt := NewYouTube(img, v)
		n.Parent().ReplaceChild(n.Parent(), n, yt)

		return ast.WalkContinue, nil
	}

	_ = ast.Walk(node, replaceImages)
}

// HTMLRenderer struct is a renderer.NodeRenderer implementation for the extension.
type HTMLRenderer struct{}

// NewHTMLRenderer builds a new HTMLRenderer with given options and returns it.
func NewHTMLRenderer() renderer.NodeRenderer {
	return &HTMLRenderer{}
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs.
func (r *HTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindYouTube, r.renderYouTubeVideo)
}

func (r *HTMLRenderer) renderYouTubeVideo(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		return ast.WalkContinue, nil
	}

	yt := node.(*YouTube)

	_, _ = w.Write([]byte(videoIFrame(yt.Video)))
	return ast.WalkContinue, nil
}

func videoIFrame(video string) string {
	return fmt.Sprintf(`<iframe src="https://www.youtube.com/embed/%s?rel=0&amp;autoplay=0&amp;showinfo=0&amp;enablejsapi=0" frameborder="0" loading="lazy" gesture="media" allow="autoplay; fullscreen" allowautoplay="true" allowfullscreen="true" style="aspect-ratio:16/9;height:100%%;width:100%%;"></iframe>`, video)
}
