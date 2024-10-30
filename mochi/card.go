package mochi

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
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
	Attachments   map[string]Attachment `json:"attachments"` // key is filename with extension.
	CreatedAt     Date                  `json:"created-at"`
	UpdatedAt     Date                  `json:"updated-at"`
}

// CreateCardRequest holds the info to create a new card.
type CreateCardRequest struct {
	Content       string           `json:"content"`
	DeckID        string           `json:"deck-id"`
	TemplateID    string           `json:"template-id,omitempty"`
	Archived      bool             `json:"archived?,omitempty"`
	ReviewReverse bool             `json:"review-reverse?,omitempty"`
	Pos           string           `json:"pos,omitempty"`
	Fields        map[string]Field `json:"fields,omitempty"`
}

// UpdateCardRequest holds the info to update a card.
type UpdateCardRequest struct {
	Content       string           `json:"content,omitempty"`
	DeckID        string           `json:"deck-id,omitempty"`
	TemplateID    string           `json:"template-id,omitempty"`
	Archived      bool             `json:"archived?,omitempty"`
	ReviewReverse bool             `json:"review-reverse?,omitempty"`
	Pos           string           `json:"pos,omitempty"`
	Fields        map[string]Field `json:"fields,omitempty"`
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

// AddAttachment adds an attachment to a card.
func (c *Client) AddAttachment(ctx context.Context, cardID, filename string, data []byte) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, bytes.NewReader(data))
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	rb := buildRequest(c).
		Pathf("%s/%s/attachments/%s", cardPath, cardID, filename).
		Method(http.MethodPost).
		Header("Content-Type", writer.FormDataContentType()).
		BodyReader(body)
	err = executeRequest(ctx, rb)
	return err
}
