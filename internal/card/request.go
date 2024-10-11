package card

import (
	"context"
	"fmt"

	"github.com/leonhfr/mochi/mochi"
)

// Client is the interface that should be implemented to sync cards.
type Client interface {
	CreateCard(ctx context.Context, req mochi.CreateCardRequest) (mochi.Card, error)
	UpdateCard(ctx context.Context, id string, req mochi.UpdateCardRequest) (mochi.Card, error)
	DeleteCard(ctx context.Context, id string) error
}

// Request is the interface that should be implemented to execute a request.
type Request interface {
	fmt.Stringer
	Sync(ctx context.Context, client Client, lf Lockfile) error
}

type createCardRequest struct {
	req      mochi.CreateCardRequest
	filename string
}

func newCreateCardRequest(filename, deckID, name, content string) Request {
	return &createCardRequest{
		req: mochi.CreateCardRequest{
			Content: content,
			DeckID:  deckID,
			Fields: map[string]mochi.Field{
				"name": {ID: "name", Value: name},
			},
		},
		filename: filename,
	}
}

// Sync implements the SyncRequest interface.
func (r *createCardRequest) Sync(ctx context.Context, c Client, lf Lockfile) error {
	card, err := c.CreateCard(ctx, r.req)
	if err != nil {
		return err
	}
	return lf.SetCard(r.req.DeckID, card.ID, r.filename)
}

// String implements the fmt.Stringer interface.
func (r *createCardRequest) String() string {
	return fmt.Sprintf("create request for file %s", r.filename)
}

type updateCardRequest struct {
	req    mochi.UpdateCardRequest
	cardID string
}

func newUpdateCardRequest(cardID, content string) Request {
	return &updateCardRequest{
		cardID: cardID,
		req:    mochi.UpdateCardRequest{Content: content},
	}
}

// Sync implements the SyncRequest interface.
func (r *updateCardRequest) Sync(ctx context.Context, c Client, _ Lockfile) error {
	_, err := c.UpdateCard(ctx, r.cardID, r.req)
	return err
}

// String implements the fmt.Stringer interface.
func (r *updateCardRequest) String() string {
	return fmt.Sprintf("update request for card ID %s", r.cardID)
}

type archiveCardRequest struct {
	req    mochi.UpdateCardRequest
	cardID string
}

func newArchiveCardRequest(cardID string) Request {
	return &archiveCardRequest{
		cardID: cardID,
		req:    mochi.UpdateCardRequest{Archived: true},
	}
}

// Sync implements the SyncRequest interface.
func (r *archiveCardRequest) Sync(ctx context.Context, c Client, _ Lockfile) error {
	_, err := c.UpdateCard(ctx, r.cardID, r.req)
	return err
}

// String implements the fmt.Stringer interface.
func (r *archiveCardRequest) String() string {
	return fmt.Sprintf("archive request for card ID %s", r.cardID)
}

type deleteCardRequest struct {
	cardID string
}

func newDeleteCardRequest(cardID string) Request {
	return &deleteCardRequest{cardID: cardID}
}

// Sync implements the SyncRequest interface.
func (r *deleteCardRequest) Sync(ctx context.Context, c Client, _ Lockfile) error {
	return c.DeleteCard(ctx, r.cardID)
}

// String implements the fmt.Stringer interface.
func (r *deleteCardRequest) String() string {
	return fmt.Sprintf("delete request for card ID %s", r.cardID)
}
