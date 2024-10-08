package parser

var extensions = []string{".md"}

// Card contains the card data parsed from a file.
type Card struct {
	Name     string
	Content  string
	Filename string
}

type cardParser interface {
	Convert(path string, source []byte) ([]Card, error)
}

// Parser represents a parser.
type Parser struct {
	def cardParser
}

// New returns a new parser.
func New() *Parser {
	return &Parser{
		def: newNote(),
	}
}

// Convert converts a source file into a slice of cards.
func (p *Parser) Convert(path string, source []byte) ([]Card, error) {
	return p.def.Convert(path, source)
}

// Extensions returns the list of supported extensions.
func (p *Parser) Extensions() []string {
	return extensions
}
