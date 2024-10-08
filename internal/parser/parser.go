package parser

// Extension is the supported markdown file extension.
const Extension = ".md"

// Parser is the interface implemented by a parser.
type Parser interface {
	Convert(path string, source []byte) ([]Card, error)
}

// Card contains the card data parsed from a file.
type Card struct {
	Name    string
	Content string
}
