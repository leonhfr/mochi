package request

import (
	"context"
	"fmt"

	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/mochi"
)

type update struct {
	cardID string
	card   parser.Card
}

// NewUpdate returns a new update card request.
func NewUpdate(cardID string, card parser.Card) Request {
	return &update{
		cardID: cardID,
		card:   card,
	}
}

// Sync implements the SyncRequest interface.
func (r *update) Sync(ctx context.Context, c Client, _ Reader, _ Lockfile) error {
	req := mochi.UpdateCardRequest{Content: r.card.Content}
	_, err := c.UpdateCard(ctx, r.cardID, req)
	return err
}

// String implements the fmt.Stringer interface.
func (r *update) String() string {
	return fmt.Sprintf("update request for card ID %s", r.cardID)
}
