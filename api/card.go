package api

import "context"

const cardPath = "/api/cards"

type Card struct {
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	Content       string           `json:"content"`
	DeckID        string           `json:"deck-id"`
	TemplateID    string           `json:"template-id"`
	POS           string           `json:"pos"`
	Archived      bool             `json:"archived"`
	New           bool             `json:"new?"`
	ReviewReverse bool             `json:"review-reverse"`
	Fields        map[string]Field `json:"fields"`
	CreatedAt     Date             `json:"created-at"`
	UpdatedAt     Date             `json:"updated-at"`
}

type CreateCardRequest struct {
	Content       string           `json:"content"`
	DeckID        string           `json:"deck-id"`
	TemplateID    string           `json:"template-id,omitempty"`
	Archived      bool             `json:"archived"`
	ReviewReverse bool             `json:"review-reverse,omitempty"`
	POS           string           `json:"pos,omitempty"`
	Fields        map[string]Field `json:"fields,omitempty"`
	Attachments   []Attachment     `json:"attachments,omitempty"`
}

type UpdateCardRequest struct {
	Content       string           `json:"content,omitempty"`
	DeckID        string           `json:"deck-id,omitempty"`
	TemplateID    string           `json:"template-id,omitempty"`
	Archived      bool             `json:"archived"`
	ReviewReverse bool             `json:"review-reverse,omitempty"`
	POS           string           `json:"pos,omitempty"`
	Fields        map[string]Field `json:"fields,omitempty"`
	Attachments   []Attachment     `json:"attachments,omitempty"`
}

type Field struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type Attachment struct {
	FileName    string `json:"file-name"`    // File name must match the regex /[0-9a-zA-Z]{8,16}/. E.g. "j94fuC0R.jpg".
	ContentType string `json:"content-type"` // MIME type.
	Data        string `json:"data"`         // Base64 encoded representation of the attachment data.
}

type Date struct {
	Date string `json:"date"`
}

func (c *Client) CreateCard(ctx context.Context, req CreateCardRequest) (Card, error) {
	return createItem[Card](ctx, c, cardPath, req)
}

func (c *Client) GetCard(ctx context.Context, id string) (Card, error) {
	return getItem[Card](ctx, c, cardPath, id)
}

func (c *Client) ListCards(ctx context.Context) ([]Card, error) {
	return listItems[Card](ctx, c, cardPath, nil)
}

func (c *Client) ListCardsInDeck(ctx context.Context, id string) ([]Card, error) {
	return listItems[Card](ctx, c, cardPath, map[string][]string{"deck-id": {id}})
}

func (c *Client) UpdateCard(ctx context.Context, id string, req UpdateCardRequest) (Card, error) {
	return updateItem[Card](ctx, c, cardPath, id, req)
}

func (c *Client) DeleteCard(ctx context.Context, id string) error {
	return deleteItem(ctx, c, cardPath, id)
}
