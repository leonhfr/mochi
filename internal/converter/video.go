package converter

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"

	cast "github.com/leonhfr/mochi/internal/converter/ast"
)

func isVideo(node *ast.Image, source []byte) bool {
	return string(node.Text(source)) == "@video"
}

func replaceVideoNode(node *ast.Image) {
	video := cast.NewVideo(node)
	parent := node.Parent()
	parent.ReplaceChild(parent, node, video)
}

func newVideoRenderer() renderer.NodeRendererFunc {
	return func(writer util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		v := n.(*cast.Video)
		if entering {
			_, _ = writer.WriteString(videoIFrame(string(v.Destination)))
		}
		return ast.WalkContinue, nil
	}
}

func videoIFrame(destination string) string {
	return fmt.Sprintf(`<iframe src="%s?rel=0&amp;autoplay=0&amp;showinfo=0&amp;enablejsapi=0" frameborder="0" loading="lazy" gesture="media" allow="autoplay; fullscreen" allowautoplay="true" allowfullscreen="true" style="aspect-ratio:16/9;height:100%%;width:100%%;"></iframe>`, destination)
}
