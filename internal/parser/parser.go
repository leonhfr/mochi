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
	Index    int
}

type cardParser interface {
	convert(fc FileCheck, path string, source []byte) ([]Card, error)
}

// Parser represents a parser.
type Parser struct {
	cardParser
	fc      FileCheck
	parsers map[string]cardParser
}

// New returns a new parser.
func New(fc FileCheck) *Parser {
	return &Parser{
		fc:         fc,
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
		return cp.convert(p.fc, path, content)
	}

	return p.convert(p.fc, path, content)
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
