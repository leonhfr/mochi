package parser

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"

	"github.com/leonhfr/mochi/internal/image"
)

var extensions = []string{".md"}

// FileCheck is the interface implemented to check file existence.
type FileCheck interface {
	Exists(path string) bool
}

// Card contains the card data parsed from a file.
type Card struct {
	Name     string
	Content  string
	Filename string
	Images   map[string]image.Image
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
func New(fc FileCheck) *Parser {
	return &Parser{
		cardParser: newNote(fc),
		parsers: map[string]cardParser{
			"note":     newNote(fc),
			"headings": newHeadings(fc, 1),
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
