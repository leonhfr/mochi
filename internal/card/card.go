package card

import (
	"io"
	"path/filepath"

	"github.com/leonhfr/mochi/internal/converter"
	"github.com/leonhfr/mochi/internal/heap"
	"github.com/leonhfr/mochi/internal/parser"
)

// Reader represents the interface to read note files.
type Reader interface {
	Read(path string) (io.ReadCloser, error)
}

// Parser represents the interface to parse note files.
type Parser interface {
	Parse(reader parser.Reader, parser, path string) (parser.Result, error)
}

// Converter represents the interface to convert cards.
type Converter interface {
	Convert(reader converter.Reader, path string, source string) (converter.Result, error)
}

// Parse parses the note files for cards.
func Parse(r Reader, p Parser, c Converter, workspace, parserName string, filePaths []string) ([]Card, error) {
	var cards []Card
	for _, path := range filePaths {
		path = filepath.Join(workspace, path)

		deck, parsedCards, err := parseFile(r, p, parserName, path)
		if err != nil {
			return nil, err
		}

		converted, err := convertCards(r, c, deck, path, parsedCards)
		if err != nil {
			return nil, err
		}

		cards = append(cards, converted...)
	}
	return cards, nil
}

func parseFile(r Reader, p Parser, parserName, path string) (string, []parser.Card, error) {
	result, err := p.Parse(r, parserName, path)
	if err != nil {
		return "", nil, err
	}

	return result.Deck, result.Cards, nil
}

func convertCards(r Reader, c Converter, deck, path string, parsedCards []parser.Card) ([]Card, error) {
	cards := make([]Card, 0, len(parsedCards))
	for _, card := range parsedCards {
		converted, err := c.Convert(r, path, card.Content)
		if err != nil {
			return nil, err
		}
		cards = append(cards, newCard(deck, converted.Markdown, card, converted.Attachments))
	}
	return cards, nil
}

// Heap creates a card heap from cards.
func Heap(cards []Card) *heap.Heap[Card] {
	h := heap.New[Card]()
	for _, card := range cards {
		h.Push(card)
	}
	return h
}

// Card contains the data to group and prioritize a card.
type Card struct {
	base        string
	Attachments []converter.Attachment
	parser.Card
}

var _ heap.Item = &Card{}

func newCard(deck, _ string, card parser.Card, attachments []converter.Attachment) Card {
	return Card{
		base:        deck,
		Attachments: attachments,
		Card:        card,
	}
}

// Base implements the PriorityItem interface.
func (c Card) Base() string {
	return c.base
}

// Priority implements the PriorityItem interface.
func (c Card) Priority() int {
	return len(c.base)
}
