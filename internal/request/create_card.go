package request

import (
	"context"
	"fmt"

	"github.com/leonhfr/mochi/internal/image"
	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/mochi"
)

type createRequest struct {
	deckID string
	card   parser.Card
}

// CreateCard returns a new create card request.
func CreateCard(deckID string, card parser.Card) Request {
	return &createRequest{
		deckID: deckID,
		card:   card,
	}
}

// Execute implements the Request interface.
func (r *createRequest) Execute(ctx context.Context, client Client, reader image.Reader, lf Lockfile) error {
	images := image.New(reader, r.card.Path(), r.card.Images())

	req := mochi.CreateCardRequest{
		Content: images.Replace(r.card.Content()),
		DeckID:  r.deckID,
		Fields: map[string]mochi.Field{
			"name": {ID: "name", Value: r.card.Name()},
		},
		Attachments: images.Attachments(),
		Pos:         r.card.Position(),
	}

	card, err := client.CreateCard(ctx, req)
	if err != nil {
		return err
	}

	hashMap := images.HashMap()
	lf.Lock()
	defer lf.Unlock()

	if err := lf.SetCard(r.deckID, card.ID, r.card.Filename(), hashMap); err != nil {
		return err
	}

	return nil
}

// String implements the fmt.Stringer interface.
func (r *createRequest) String() string {
	return fmt.Sprintf("create request for file %s", r.card.Filename())
}
