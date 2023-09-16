package api

import "context"

const deckPath = "/api/decks"

type Deck struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ParentID string `json:"parent-id"`
}

type CreateDeckRequest struct {
	Name     string `json:"name"`
	ParentID string `json:"parent-id,omitempty"`
}

type UpdateDeckRequest struct {
	Name     string `json:"name"`
	ParentID string `json:"parent-id,omitempty"`
}

func (c *Client) CreateDeck(ctx context.Context, req CreateDeckRequest) (Deck, error) {
	return createItem[Deck](ctx, c, deckPath, req)
}

func (c *Client) GetDeck(ctx context.Context, id string) (Deck, error) {
	return getItem[Deck](ctx, c, deckPath, id)
}

func (c *Client) ListDecks(ctx context.Context) ([]Deck, error) {
	return listItems[Deck](ctx, c, deckPath, nil)
}

func (c *Client) UpdateDeck(ctx context.Context, id string, req UpdateDeckRequest) (Deck, error) {
	return updateItem[Deck](ctx, c, deckPath, id, req)
}

func (c *Client) DeleteDeck(ctx context.Context, id string) error {
	return deleteItem(ctx, c, deckPath, id)
}
