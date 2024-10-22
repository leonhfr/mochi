package parser

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
)

var extensions = []string{".md"}

// Card contains the card data parsed from a file.
type Card struct {
	Name     string
	Content  string
	Filename string
	Path     string
	Images   []Image
	Index    int
}

// Image contains the parsed image data.
type Image struct {
	Destination string
	AltText     string
}

type cardParser interface {
	convert(path string, source []byte) ([]Card, error)
}

// Parser represents a parser.
type Parser struct {
	cardParser
	parsers map[string]cardParser
}

// New returns a new parser.
func New() *Parser {
	return &Parser{
		cardParser: newNote(),
		parsers: map[string]cardParser{
			"note":      newNote(),
			"headings":  newHeadings(1),
			"headings1": newHeadings(1),
			"headings2": newHeadings(2),
			"headings3": newHeadings(3),
		},
	}
}

// Convert converts a source file into a slice of cards.
func (p *Parser) Convert(parser, path string, r io.Reader) ([]Card, error) {
	var matter struct {
		Parser string `yaml:"mochi-parser"`
		Skip   bool   `yaml:"mochi-skip"`
	}

	content, err := frontmatter.Parse(r, &matter)
	if err != nil {
		return nil, err
	}

	if matter.Skip {
		return nil, nil
	}

	if matter.Parser != "" {
		parser = matter.Parser
	}

	if cp, ok := p.parsers[parser]; ok {
		return cp.convert(path, content)
	}

	return p.convert(path, content)
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

func getFilename(path string) string {
	return filepath.Base(path)
}

func getNameFromPath(path string) string {
	base := filepath.Base(path)
	for _, ext := range extensions {
		base = strings.TrimSuffix(base, ext)
	}
	return base
}
