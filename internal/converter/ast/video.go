package ast

import "github.com/yuin/goldmark/ast"

// Video is a struct that represents an embedded video.
type Video struct {
	ast.BaseInline
	Destination []byte
}

// Dump implements ast.Node.Dump.
func (n *Video) Dump(source []byte, level int) {
	m := map[string]string{}
	m["Destination"] = string(n.Destination)
	ast.DumpHelper(n, source, level, m, nil)
}

// KindVideo is a NodeKind of the Video node.
var KindVideo = ast.NewNodeKind("Video")

// Kind implements ast.Node.Kind.
func (n *Video) Kind() ast.NodeKind {
	return KindVideo
}

// NewVideo returns a new Video node.
func NewVideo(link *ast.Image) *Video {
	return &Video{
		BaseInline:  ast.BaseInline{},
		Destination: link.Destination,
	}
}
