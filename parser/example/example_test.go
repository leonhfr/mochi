package example

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
)

var _ ast.Node = &Node{}

var _ parser.InlineParser = &exampleParser{}
