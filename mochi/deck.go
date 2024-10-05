package mochi

import "context"

const deckPath = "/api/decks"

// Deck represents a deck.
type Deck struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ParentID string `json:"parent-id,omitempty"`
}

// CreateDeckRequest holds the info to create a new deck.
type CreateDeckRequest struct {
	Name     string `json:"name"`
	ParentID string `json:"parent-id,omitempty"`
}

// UpdateDeckRequest holds the info to update a deck.
type UpdateDeckRequest struct {
	Name     string `json:"name"`
	ParentID string `json:"parent-id,omitempty"`
}

// CreateDeck creates a new deck.
func (c *Client) CreateDeck(ctx context.Context, req CreateDeckRequest) (Deck, error) {
	return createItem[Deck](ctx, c, deckPath, req)
}

// GetDeck gets a single deck.
func (c *Client) GetDeck(ctx context.Context, id string) (Deck, error) {
	return getItem[Deck](ctx, c, deckPath, id)
}

// ListDecks lists the decks.
//
// The callback is called with a slice of decks until all decks have been listed
// or until the callback returns an error. Each callback call makes a HTTP request.
func (c *Client) ListDecks(ctx context.Context, cb func([]Deck) error) error {
	return listItems(ctx, c, deckPath, nil, cb)
}

// UpdateDeck updates an existing deck.
func (c *Client) UpdateDeck(ctx context.Context, id string, req UpdateDeckRequest) (Deck, error) {
	return updateItem[Deck](ctx, c, deckPath, id, req)
}

// DeleteDeck deletes a deck.
func (c *Client) DeleteDeck(ctx context.Context, id string) error {
	return deleteItem(ctx, c, deckPath, id)
}
