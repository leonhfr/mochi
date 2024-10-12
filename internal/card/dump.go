package card

import (
	"context"

	"github.com/leonhfr/mochi/mochi"
)

// DumpClient is the interface that should be implemented to dump cards.
type DumpClient interface {
	ListCardsInDeck(ctx context.Context, id string) ([]mochi.Card, error)
}

// DumpRequests returns the requests to dump cards.
func DumpRequests(ctx context.Context, client DumpClient, deckID string) ([]Request, error) {
	cards, err := client.ListCardsInDeck(ctx, deckID)
	if err != nil {
		return nil, err
	}
	reqs := make([]Request, 0, len(cards))
	for _, card := range cards {
		reqs = append(reqs, newDeleteCardRequest(card.ID))
	}
	return reqs, nil
}
