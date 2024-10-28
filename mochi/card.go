package mochi

import (
	"context"
	"time"
)

const cardPath = "/api/cards"

// Card represents a card.
type Card struct {
	ID            string                `json:"id"`
	Name          string                `json:"name"`
	Content       string                `json:"content"`
	DeckID        string                `json:"deck-id"`
	TemplateID    string                `json:"template-id"`
	Pos           string                `json:"pos"`
	Archived      bool                  `json:"archived?"`
	New           bool                  `json:"new?"`
	ReviewReverse bool                  `json:"review-reverse?"`
	Fields        map[string]Field      `json:"fields"`
	Attachments   map[string]Attachment `json:"attachments"`
	CreatedAt     Date                  `json:"created-at"`
	UpdatedAt     Date                  `json:"updated-at"`
}

// CreateCardRequest holds the info to create a new card.
type CreateCardRequest struct {
	Content       string                 `json:"content"`
	DeckID        string                 `json:"deck-id"`
	TemplateID    string                 `json:"template-id,omitempty"`
	Archived      bool                   `json:"archived?,omitempty"`
	ReviewReverse bool                   `json:"review-reverse?,omitempty"`
	Pos           string                 `json:"pos,omitempty"`
	Fields        map[string]Field       `json:"fields,omitempty"`
	Attachments   []DeprecatedAttachment `json:"deprecated/attachments,omitempty"`
}

// UpdateCardRequest holds the info to update a card.
type UpdateCardRequest struct {
	Content       string                 `json:"content,omitempty"`
	DeckID        string                 `json:"deck-id,omitempty"`
	TemplateID    string                 `json:"template-id,omitempty"`
	Archived      bool                   `json:"archived?,omitempty"`
	ReviewReverse bool                   `json:"review-reverse?,omitempty"`
	Pos           string                 `json:"pos,omitempty"`
	Fields        map[string]Field       `json:"fields,omitempty"`
	Attachments   []DeprecatedAttachment `json:"deprecated/attachments,omitempty"`
}

// Field represents a field.
type Field struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

// Attachment represents an attachment.
type Attachment struct {
	Size int    `json:"size"` // Size in bytes.
	Type string `json:"type"` // MIME type.
}

// DeprecatedAttachment represents a deprecated attachment.
type DeprecatedAttachment struct {
	FileName    string `json:"file-name"`    // File name must match the regex /[0-9a-zA-Z]{8,16}/. E.g. "j94fuC0R.jpg".
	ContentType string `json:"content-type"` // MIME type.
	Data        string `json:"data"`         // Base64 encoded representation of the attachment data.
}

// Date represents a time.
type Date struct {
	Date time.Time `json:"date"`
}

// CreateCard creates a new card.
func (c *Client) CreateCard(ctx context.Context, req CreateCardRequest) (Card, error) {
	return createItem[Card](ctx, c, cardPath, req)
}

// GetCard gets a single card.
func (c *Client) GetCard(ctx context.Context, id string) (Card, error) {
	return getItem[Card](ctx, c, cardPath, id)
}

// ListCards lists the cards.
func (c *Client) ListCards(ctx context.Context) ([]Card, error) {
	return listItems[Card](ctx, c, cardPath, nil)
}

// ListCardsInDeck lists the cards in a deck.
func (c *Client) ListCardsInDeck(ctx context.Context, id string) ([]Card, error) {
	return listItems[Card](ctx, c, cardPath, map[string][]string{"deck-id": {id}})
}

// UpdateCard updates an existing card.
func (c *Client) UpdateCard(ctx context.Context, id string, req UpdateCardRequest) (Card, error) {
	return updateItem[Card](ctx, c, cardPath, id, req)
}

// DeleteCard deletes a card.
func (c *Client) DeleteCard(ctx context.Context, id string) error {
	return deleteItem(ctx, c, cardPath, id)
}
