package parser

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
)

var extensions = []string{".md"}

// Result contains the result.
type Result struct {
	Deck  string
	Cards []Card
}

// Card contains the card data parsed from a file.
type Card struct {
	Name     string
	Content  string
	Filename string
	Path     string
	Images   []Image
	Position string
}

// Image contains the parsed image data.
type Image struct {
	Destination string
	AltText     string
}

type cardParser interface {
	convert(path string, source []byte) (Result, error)
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

// Convert converts a source file into cards.
func (p *Parser) Convert(parser, path string, r io.Reader) (Result, error) {
	var matter struct {
		Parser string `yaml:"mochi-parser"`
		Skip   bool   `yaml:"mochi-skip"`
	}

	content, err := frontmatter.Parse(r, &matter)
	if err != nil {
		return Result{}, err
	}

	if matter.Skip {
		return Result{}, nil
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

func getPosition(filename string, index int) string {
	runes := make([]rune, 0, len(filename))
	for _, r := range filename {
		if ('0' <= r && r <= '9') || ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') {
			runes = append(runes, r)
		}
	}
	return fmt.Sprintf("%s%04d", string(runes), index)
}
