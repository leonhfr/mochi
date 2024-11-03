package request

import (
	"context"
	"fmt"

	"github.com/sourcegraph/conc/pool"

	"github.com/leonhfr/mochi/internal/card"
	"github.com/leonhfr/mochi/mochi"
)

type createRequest struct {
	deckID string
	card   card.Card
}

// CreateCard returns a new create card request.
func CreateCard(deckID string, card card.Card) Request {
	return &createRequest{
		deckID: deckID,
		card:   card,
	}
}

// Execute implements the Request interface.
func (r *createRequest) Execute(ctx context.Context, client Client, lf Lockfile) error {
	req := mochi.CreateCardRequest{
		Content:    r.card.Content,
		DeckID:     r.deckID,
		TemplateID: r.card.TemplateID,
		Fields:     mochiFields(r.card.Fields),
		Pos:        r.card.Position,
	}

	card, err := client.CreateCard(ctx, req)
	if err != nil {
		return err
	}

	p := pool.New().WithContext(ctx)
	for _, attachment := range r.card.Attachments {
		attachment := attachment
		p.Go(func(ctx context.Context) error {
			return client.AddAttachment(ctx, card.ID, attachment.Filename, attachment.Bytes)
		})
	}

	err = p.Wait()
	if err != nil {
		return err
	}

	lf.Lock()
	defer lf.Unlock()

	if err := lf.SetCard(r.deckID, card.ID, r.card.Filename()); err != nil {
		return err
	}

	return nil
}

// String implements the fmt.Stringer interface.
func (r *createRequest) String() string {
	return fmt.Sprintf("create request for file %s", r.card.Filename())
}
