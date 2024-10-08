package card

import (
	"context"

	"github.com/leonhfr/mochi/mochi"
)

// Client is the interface that should be implemented to sync cards.
type Client interface {
	CreateCard(ctx context.Context, req mochi.CreateCardRequest) (mochi.Card, error)
	UpdateCard(ctx context.Context, id string, req mochi.UpdateCardRequest) (mochi.Card, error)
}

// WriteLockfile is the interface that should be implemented to update the lockfile.
type WriteLockfile interface {
	SetCard(deckID string, cardID string, filename string) error
}

// SyncRequest is the interface that should be implemented to execute a request.
type SyncRequest interface {
	Sync(ctx context.Context, client Client, lf WriteLockfile) error
}

type createCardRequest struct {
	req      mochi.CreateCardRequest
	filename string
}

func newCreateCardRequest(filename, deckID, name, content string) *createCardRequest {
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
func (r *createCardRequest) Sync(ctx context.Context, c Client, lf WriteLockfile) error {
	card, err := c.CreateCard(ctx, r.req)
	if err != nil {
		return err
	}
	return lf.SetCard(r.req.DeckID, card.ID, r.filename)
}

type updateCardRequest struct {
	req    mochi.UpdateCardRequest
	cardID string
}

func newUpdateCardRequest(cardID, content string) *updateCardRequest {
	return &updateCardRequest{
		cardID: cardID,
		req:    mochi.UpdateCardRequest{Content: content},
	}
}

// Sync implements the SyncRequest interface.
func (r *updateCardRequest) Sync(ctx context.Context, c Client, _ WriteLockfile) error {
	_, err := c.UpdateCard(ctx, r.cardID, r.req)
	return err
}

type archiveCardRequest struct {
	req    mochi.UpdateCardRequest
	cardID string
}

func newArchiveCardRequest(cardID string) *archiveCardRequest {
	return &archiveCardRequest{
		cardID: cardID,
		req:    mochi.UpdateCardRequest{Archived: true},
	}
}

// Sync implements the SyncRequest interface.
func (r *archiveCardRequest) Sync(ctx context.Context, c Client, _ WriteLockfile) error {
	_, err := c.UpdateCard(ctx, r.cardID, r.req)
	return err
}
