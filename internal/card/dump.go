package card

import (
	"context"

	"github.com/leonhfr/mochi/mochi"
)

// DumpClient is the interface that should be implemented to dump cards.
type DumpClient interface {
	Client
	ListCards(ctx context.Context) ([]mochi.Card, error)
	ListCardsInDeck(ctx context.Context, id string) ([]mochi.Card, error)
}

// DumpRequests returns the requests to dump cards.
func DumpRequests(ctx context.Context, client DumpClient, deckID *string) ([]Request, error) {
	cards, err := fetchCards(ctx, client, deckID)
	if err != nil {
		return nil, err
	}
	reqs := make([]Request, 0, len(cards))
	for _, card := range cards {
		reqs = append(reqs, newDeleteCardRequest(card.ID))
	}
	return reqs, nil
}

func fetchCards(ctx context.Context, client DumpClient, deckID *string) ([]mochi.Card, error) {
	if deckID == nil {
		return client.ListCards(ctx)
	}
	return client.ListCardsInDeck(ctx, *deckID)
}
