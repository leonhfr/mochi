package parser

import "fmt"

type Parser interface {
	fmt.Stringer
	Fields() []string
}
