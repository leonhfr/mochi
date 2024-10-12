package request

import (
	"context"
	"fmt"

	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/mochi"
)

type create struct {
	filename string
	deckID   string
	card     parser.Card
}

// NewCreate returns a new create card request.
func NewCreate(filename, deckID string, card parser.Card) Request {
	return &create{
		filename: filename,
		deckID:   deckID,
		card:     card,
	}
}

// Sync implements the SyncRequest interface.
func (r *create) Sync(ctx context.Context, c Client, _ Reader, lf Lockfile) error {
	req := mochi.CreateCardRequest{
		Content: r.card.Content,
		DeckID:  r.deckID,
		Fields: map[string]mochi.Field{
			"name": {ID: "name", Value: r.card.Name},
		},
	}

	card, err := c.CreateCard(ctx, req)
	if err != nil {
		return err
	}
	return lf.SetCard(r.deckID, card.ID, r.filename)
}

// String implements the fmt.Stringer interface.
func (r *create) String() string {
	return fmt.Sprintf("create request for file %s", r.filename)
}
