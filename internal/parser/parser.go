package parser

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"

	"github.com/leonhfr/mochi/internal/config"
)

// Reader represents the interface to read files.
type Reader interface {
	Read(path string) (io.ReadCloser, error)
}

var extensions = []string{".md"}

// Result contains the result.
type Result struct {
	Deck  string
	Cards []Card
}

// Card represents a card.
type Card struct {
	Content    string
	Fields     map[string]string
	TemplateID string
	Images     []Image
	Path       string
	Position   string
}

// Filename returns the filename.
func (c Card) Filename() string { return getFilename(c.Path) }

// Image contains the parsed image data.
type Image struct {
	Destination string
	AltText     string
}

type cardParser interface {
	convert(path string, source []byte) (Result, error)
}

func defaultParsers() map[string]cardParser {
	return map[string]cardParser{
		"note":      newNote(),
		"headings":  newHeadings(1),
		"headings1": newHeadings(1),
		"headings2": newHeadings(2),
		"headings3": newHeadings(3),
		"table":     newTable(),
	}
}

// Parser represents a parser.
type Parser struct {
	cardParser
	parsers map[string]cardParser
}

// New returns a new parser.
func New(options ...Option) (*Parser, error) {
	p := &Parser{
		cardParser: newNote(),
		parsers:    defaultParsers(),
	}
	for _, option := range options {
		if err := option(p); err != nil {
			return nil, err
		}
	}
	return p, nil
}

// Option represents an option for the parser.
type Option func(*Parser) error

// WithVocabulary adds the vocabulary templates.
func WithVocabulary(vocabulary map[string]config.VocabularyTemplate) Option {
	return func(p *Parser) error {
		for name, templateID := range vocabulary {
			if _, ok := p.parsers[name]; ok {
				return fmt.Errorf("vocabulary template: cannot overwrite default parser %s", name)
			}
			p.parsers[name] = newVocabulary(templateID)
		}
		return nil
	}
}

// Convert converts a source file into cards.
func (p *Parser) Convert(reader Reader, parser, path string) (Result, error) {
	content, matter, err := parseFrontmatter(reader, path)
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

type matter struct {
	Parser string `yaml:"mochi-parser"`
	Skip   bool   `yaml:"mochi-skip"`
}

func parseFrontmatter(reader Reader, path string) ([]byte, matter, error) {
	bytes, err := reader.Read(path)
	if err != nil {
		return nil, matter{}, err
	}
	defer bytes.Close()

	var fm matter
	content, err := frontmatter.Parse(bytes, &fm)
	if err != nil {
		return nil, matter{}, err
	}

	return content, fm, nil
}

// Names returns the list of allowed parser names.
func Names() []string {
	parsers := defaultParsers()
	names := make([]string, 0, len(parsers))
	for name := range parsers {
		names = append(names, name)
	}
	return names
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

func sanitizePosition(position string) string {
	runes := make([]rune, 0, len(position))
	for _, r := range position {
		if ('0' <= r && r <= '9') || ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') {
			runes = append(runes, r)
		}
	}
	return string(runes)
}
