package deck

import (
	"context"

	"github.com/leonhfr/mochi/mochi"
)

// LeafDecks returns the ids of the leaf decks.
func LeafDecks(decks []mochi.Deck) []string {
	deckMap := make(map[string][]string)
	for _, deck := range decks {
		deckMap[deck.ParentID] = append(deckMap[deck.ParentID], deck.ID)
	}
	var leaves []string
	for _, deck := range decks {
		if _, ok := deckMap[deck.ID]; !ok {
			leaves = append(leaves, deck.ID)
		}
	}
	return leaves
}

// DeleteEmptyClient is the interface to clean mochi decks.
type DeleteEmptyClient interface {
	ListCardsInDeck(ctx context.Context, id string) ([]mochi.Card, error)
	DeleteDeck(ctx context.Context, id string) error
}

// DeleteEmpty deletes the deck if it does not contain any cards.
func DeleteEmpty(ctx context.Context, client DeleteEmptyClient, deckID string) (bool, error) {
	cards, err := client.ListCardsInDeck(ctx, deckID)
	if err != nil {
		return false, err
	}

	if len(cards) > 0 {
		return false, nil
	}

	return true, client.DeleteDeck(ctx, deckID)
}
