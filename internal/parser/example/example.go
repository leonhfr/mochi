package example

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// KindExample is the NodeKind of the Example node.
var KindExample = ast.NewNodeKind("Example")

// Node is a struct that represents an example node.
type Node struct {
	ast.BaseInline
}

// New returns a new example node.
func New() *Node {
	return &Node{
		BaseInline: ast.BaseInline{},
	}
}

// Dump implements Node.Dump.
func (n *Node) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// Kind implements Node.Kind.
func (n *Node) Kind() ast.NodeKind {
	return KindExample
}

// isBlank returns true if this node consists of spaces, otherwise false.
func (n *Node) isBlank(source []byte) bool {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		text := c.(*ast.Text).Segment
		if !util.IsBlank(text.Value(source)) {
			return false
		}
	}
	return true
}

type exampleParser struct{}

// NewParser return a new InlineParser that parses examples surrounded by '"'.
func NewParser() parser.InlineParser {
	return &exampleParser{}
}

func (s *exampleParser) Trigger() []byte {
	return []byte{'"'}
}

func (s *exampleParser) Parse(_ ast.Node, block text.Reader, _ parser.Context) ast.Node {
	line, startSegment := block.PeekLine()
	opener := 0
	//nolint:revive
	for ; opener < len(line) && line[opener] == '"'; opener++ {
	}
	block.Advance(opener)
	l, pos := block.Position()
	node := New()
	for {
		line, segment := block.PeekLine()
		if line == nil {
			block.SetPosition(l, pos)
			return ast.NewTextSegment(startSegment.WithStop(startSegment.Start + opener))
		}
		for i := 0; i < len(line); i++ {
			c := line[i]
			if c == '"' {
				oldi := i
				//nolint:revive
				for ; i < len(line) && line[i] == '"'; i++ {
				}
				closure := i - oldi
				if closure == opener && (i >= len(line) || line[i] != '"') {
					segment = segment.WithStop(segment.Start + i - closure)
					if !segment.IsEmpty() {
						node.AppendChild(node, ast.NewRawTextSegment(segment))
					}
					block.Advance(i)
					goto end
				}
			}
		}
		node.AppendChild(node, ast.NewRawTextSegment(segment))
		block.AdvanceLine()
	}
end:
	if !node.isBlank(block.Source()) {
		// trim first halfspace and last halfspace
		segment := node.FirstChild().(*ast.Text).Segment
		shouldTrimmed := true
		if !(!segment.IsEmpty() && isSpaceOrNewline(block.Source()[segment.Start])) {
			shouldTrimmed = false
		}
		segment = node.LastChild().(*ast.Text).Segment
		if !(!segment.IsEmpty() && isSpaceOrNewline(block.Source()[segment.Stop-1])) {
			shouldTrimmed = false
		}
		if shouldTrimmed {
			t := node.FirstChild().(*ast.Text)
			segment := t.Segment
			t.Segment = segment.WithStart(segment.Start + 1)
			t = node.LastChild().(*ast.Text)
			segment = node.LastChild().(*ast.Text).Segment
			t.Segment = segment.WithStop(segment.Stop - 1)
		}
	}
	return node
}

func isSpaceOrNewline(c byte) bool {
	return c == ' ' || c == '\n'
}
