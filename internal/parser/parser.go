package parser

var extensions = []string{".md"}

// Card contains the card data parsed from a file.
type Card struct {
	Name     string
	Content  string
	Filename string
}

type cardParser interface {
	convert(path string, source []byte) ([]Card, error)
}

// Parser represents a parser.
type Parser struct {
	parsers map[string]cardParser
	def     cardParser
}

// New returns a new parser.
func New() *Parser {
	return &Parser{
		parsers: map[string]cardParser{
			"note":     newNote(),
			"headings": newHeadings(),
		},
		def: newNote(),
	}
}

// Convert converts a source file into a slice of cards.
func (p *Parser) Convert(parserName, path string, source []byte) ([]Card, error) {
	if cp, ok := p.parsers[parserName]; ok {
		return cp.convert(path, source)
	}
	return p.def.convert(path, source)
}

// List returns the list of allowed parser names.
func (p *Parser) List() []string {
	parsers := make([]string, 0, len(p.parsers))
	for parser := range p.parsers {
		parsers = append(parsers, parser)
	}
	return parsers
}

// Extensions returns the list of supported extensions.
func (p *Parser) Extensions() []string {
	return extensions
}
