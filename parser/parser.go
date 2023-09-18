package parser

import "fmt"

type Parser interface {
	fmt.Stringer
	Fields() []string
	Convert(path string, source []byte) ([]Card, error)
}

type Card struct {
	Name    string
	Content string
	Fields  map[string]string
	Images  map[string]Image
}
