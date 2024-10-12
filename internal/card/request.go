package card

import (
	"context"
	"fmt"

	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/mochi"
)

// Client is the interface that should be implemented to sync cards.
type Client interface {
	CreateCard(ctx context.Context, req mochi.CreateCardRequest) (mochi.Card, error)
	UpdateCard(ctx context.Context, id string, req mochi.UpdateCardRequest) (mochi.Card, error)
	DeleteCard(ctx context.Context, id string) error
}

// RequestLockfile is the interface the lockfile implement to sync cards.
type RequestLockfile interface {
	SetCard(deckID string, cardID string, filename string) error
}

// Request is the interface that should be implemented to execute a request.
type Request interface {
	fmt.Stringer
	Sync(ctx context.Context, client Client, lf RequestLockfile) error
}

type createCardRequest struct {
	filename string
	deckID   string
	card     parser.Card
}

func newCreateCardRequest(filename, deckID string, card parser.Card) Request {
	return &createCardRequest{
		filename: filename,
		deckID:   deckID,
		card:     card,
	}
}

// Sync implements the SyncRequest interface.
func (r *createCardRequest) Sync(ctx context.Context, c Client, lf RequestLockfile) error {
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
func (r *createCardRequest) String() string {
	return fmt.Sprintf("create request for file %s", r.filename)
}

type updateCardRequest struct {
	cardID string
	card   parser.Card
}

func newUpdateCardRequest(cardID string, card parser.Card) Request {
	return &updateCardRequest{
		cardID: cardID,
		card:   card,
	}
}

// Sync implements the SyncRequest interface.
func (r *updateCardRequest) Sync(ctx context.Context, c Client, _ RequestLockfile) error {
	req := mochi.UpdateCardRequest{Content: r.card.Content}
	_, err := c.UpdateCard(ctx, r.cardID, req)
	return err
}

// String implements the fmt.Stringer interface.
func (r *updateCardRequest) String() string {
	return fmt.Sprintf("update request for card ID %s", r.cardID)
}

type archiveCardRequest struct {
	cardID string
}

func newArchiveCardRequest(cardID string) Request {
	return &archiveCardRequest{
		cardID: cardID,
	}
}

// Sync implements the SyncRequest interface.
func (r *archiveCardRequest) Sync(ctx context.Context, c Client, _ RequestLockfile) error {
	req := mochi.UpdateCardRequest{Archived: true}
	_, err := c.UpdateCard(ctx, r.cardID, req)
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
func (r *deleteCardRequest) Sync(ctx context.Context, c Client, _ RequestLockfile) error {
	return c.DeleteCard(ctx, r.cardID)
}

// String implements the fmt.Stringer interface.
func (r *deleteCardRequest) String() string {
	return fmt.Sprintf("delete request for card ID %s", r.cardID)
}
