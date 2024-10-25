package deck

import (
	"io"
	"path/filepath"

	"github.com/leonhfr/mochi/internal/heap"
	"github.com/leonhfr/mochi/internal/parser"
)

// Reader represents the interface to read note files.
type Reader interface {
	Read(path string) (io.ReadCloser, error)
}

// Parser represents the interface to parse note files.
type Parser interface {
	Convert(parser, path string, source io.Reader) (parser.Result, error)
}

// Parse parses the note files for cards.
func Parse(r Reader, p Parser, workspace, parserName string, filePaths []string) ([]parser.Card, error) {
	var cards []parser.Card
	for _, path := range filePaths {
		parsed, err := parseFile(r, p, workspace, parserName, path)
		if err != nil {
			return nil, err
		}
		cards = append(cards, parsed.Cards...)
	}
	return cards, nil
}

// ParseCards parses the note files for cards.
func ParseCards(r Reader, p Parser, workspace, parserName string, filePaths []string) ([]Card, error) {
	var cards []Card
	for _, path := range filePaths {
		parsed, err := parseFile(r, p, workspace, parserName, path)
		if err != nil {
			return nil, err
		}
		cards = append(cards, newCards(parsed)...)
	}
	return cards, nil
}

func parseFile(r Reader, p Parser, workspace, parserName, path string) (parser.Result, error) {
	path = filepath.Join(workspace, path)
	bytes, err := r.Read(path)
	if err != nil {
		return parser.Result{}, err
	}
	defer bytes.Close()

	result, err := p.Convert(parserName, path, bytes)
	if err != nil {
		return parser.Result{}, err
	}

	return result, nil
}

// Card contains the data to group and prioritize a card.
type Card struct {
	base string
	card parser.Card
}

var _ heap.Item = &Card{}

func newCards(result parser.Result) []Card {
	cards := make([]Card, len(result.Cards))
	for i, card := range result.Cards {
		cards[i] = Card{base: result.Deck, card: card}
	}
	return cards
}

// Base implements the PriorityItem interface.
func (c Card) Base() string {
	return c.base
}

// Priority implements the PriorityItem interface.
func (c Card) Priority() int {
	return len(c.base)
}
