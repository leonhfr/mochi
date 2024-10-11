package mochi

import (
	"context"
	"net/http"
	"testing"
)

func Test_CreateCard(t *testing.T) {
	tests := []struct {
		name string
		test createItemTestCase[CreateCardRequest]
	}{
		{
			name: "should create a card",
			test: createItemTestCase[CreateCardRequest]{
				status: http.StatusCreated,
				req:    CreateCardRequest{Content: "Card content", DeckID: "DECK_ID"},
				res:    Card{ID: "CARD_ID", Name: "Card Name", Content: "Card content", DeckID: "DECK_ID"},
				want:   Card{ID: "CARD_ID", Name: "Card Name", Content: "Card content", DeckID: "DECK_ID"},
				err:    "",
			},
		},
		{
			name: "should return an error",
			test: createItemTestCase[CreateCardRequest]{
				status: http.StatusBadRequest,
				req:    CreateCardRequest{Content: "Card content", DeckID: "DECK_ID"},
				res:    `{"errors":["ERROR_MESSAGE"]}`,
				want:   Card{},
				err:    "mochi: ERROR_MESSAGE",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testCreateItem("/api/cards", tt.test, func(client *Client, req CreateCardRequest) (any, error) {
			return client.CreateCard(context.Background(), req)
		}))
	}
}

func Test_GetCard(t *testing.T) {
	tests := []struct {
		name string
		test getItemTestCase
	}{
		{
			name: "should get a card",
			test: getItemTestCase{
				status: http.StatusOK,
				id:     "CARD_ID",
				res:    Card{ID: "CARD_ID", Name: "Card Name", Content: "Card content", DeckID: "DECK_ID"},
				want:   Card{ID: "CARD_ID", Name: "Card Name", Content: "Card content", DeckID: "DECK_ID"},
				err:    "",
			},
		},
		{
			name: "should return an error",
			test: getItemTestCase{
				status: http.StatusBadRequest,
				id:     "CARD_ID",
				res:    `{"errors":["ERROR_MESSAGE"]}`,
				want:   Card{},
				err:    "mochi: ERROR_MESSAGE",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testGetItem("/api/cards", tt.test, func(client *Client, id string) (any, error) {
			return client.GetCard(context.Background(), id)
		}))
	}
}

func Test_ListCards(t *testing.T) {
	tests := []struct {
		name string
		test listItemTestCase[Card]
	}{
		{
			name: "should call the endpoint once",
			test: listItemTestCase[Card]{
				responses: []listItemTestCaseResponse{
					{
						status: http.StatusOK,
						params: map[string]string{"limit": "100"},
						res: listResponse[Card]{
							Docs: []Card{{ID: "CARD_ID", Name: "CardName", Content: "Card content"}},
						},
					},
				},
				want: []Card{
					{ID: "CARD_ID", Name: "CardName", Content: "Card content"},
				},
			},
		},
		{
			name: "should call the endpoint several times",
			test: listItemTestCase[Card]{
				responses: []listItemTestCaseResponse{
					{
						status: http.StatusOK,
						params: map[string]string{"limit": "100"},
						res: listResponse[Card]{
							Docs:     []Card{{ID: "CARD_ID_1", Name: "CardName1", Content: "Card content"}},
							Bookmark: "BOOKMARK_1",
						},
					},
					{
						status: http.StatusOK,
						params: map[string]string{"limit": "100", "bookmark": "BOOKMARK_1"},
						res: listResponse[Card]{
							Docs: []Card{{ID: "CARD_ID_2", Name: "CardName2", Content: "Card content"}},
						},
					},
				},
				want: []Card{
					{ID: "CARD_ID_1", Name: "CardName1", Content: "Card content"},
					{ID: "CARD_ID_2", Name: "CardName2", Content: "Card content"},
				},
			},
		},
		{
			name: "should return an error",
			test: listItemTestCase[Card]{
				responses: []listItemTestCaseResponse{
					{
						status: http.StatusBadRequest,
						params: map[string]string{"limit": "100"},
						res:    `{"errors":["ERROR_MESSAGE"]}`,
					},
				},
				err: "mochi: ERROR_MESSAGE",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testListItem("/api/cards", tt.test, func(client *Client, _ string) ([]Card, error) {
			return client.ListCards(context.Background())
		}))
	}
}

func Test_ListCardsInDeck(t *testing.T) {
	tests := []struct {
		name string
		test listItemTestCase[Card]
	}{
		{
			name: "should call the endpoint once",
			test: listItemTestCase[Card]{
				id:     "DECK_ID",
				params: map[string]string{"deck-id": "DECK_ID"},
				responses: []listItemTestCaseResponse{
					{
						status: http.StatusOK,
						params: map[string]string{"limit": "100"},
						res: listResponse[Card]{
							Docs: []Card{{ID: "CARD_ID", Name: "CardName", Content: "Card content"}},
						},
					},
				},
				want: []Card{
					{ID: "CARD_ID", Name: "CardName", Content: "Card content"},
				},
			},
		},
		{
			name: "should call the endpoint several times",
			test: listItemTestCase[Card]{
				id:     "DECK_ID",
				params: map[string]string{"deck-id": "DECK_ID"},
				responses: []listItemTestCaseResponse{
					{
						status: http.StatusOK,
						params: map[string]string{"limit": "100"},
						res: listResponse[Card]{
							Docs:     []Card{{ID: "CARD_ID_1", Name: "CardName1", Content: "Card content"}},
							Bookmark: "BOOKMARK_1",
						},
					},
					{
						status: http.StatusOK,
						params: map[string]string{"limit": "100", "bookmark": "BOOKMARK_1"},
						res: listResponse[Card]{
							Docs: []Card{{ID: "CARD_ID_2", Name: "CardName2", Content: "Card content"}},
						},
					},
				},
				want: []Card{
					{ID: "CARD_ID_1", Name: "CardName1", Content: "Card content"},
					{ID: "CARD_ID_2", Name: "CardName2", Content: "Card content"},
				},
			},
		},
		{
			name: "should return an error",
			test: listItemTestCase[Card]{
				id:     "DECK_ID",
				params: map[string]string{"deck-id": "DECK_ID"},
				responses: []listItemTestCaseResponse{
					{
						status: http.StatusBadRequest,
						params: map[string]string{"limit": "100"},
						res:    `{"errors":["ERROR_MESSAGE"]}`,
					},
				},
				err: "mochi: ERROR_MESSAGE",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testListItem("/api/cards", tt.test, func(client *Client, id string) ([]Card, error) {
			return client.ListCardsInDeck(context.Background(), id)
		}))
	}
}

func Test_UpdateCard(t *testing.T) {
	tests := []struct {
		name string
		test updateItemTestCase[UpdateCardRequest]
	}{
		{
			name: "should update a deck",
			test: updateItemTestCase[UpdateCardRequest]{
				status: http.StatusCreated,
				req:    UpdateCardRequest{Content: "Card content", DeckID: "DECK_ID"},
				res:    Card{ID: "CARD_ID", Name: "Card Name", Content: "Card content", DeckID: "DECK_ID"},
				want:   Card{ID: "CARD_ID", Name: "Card Name", Content: "Card content", DeckID: "DECK_ID"},
				err:    "",
			},
		},
		{
			name: "should return an error",
			test: updateItemTestCase[UpdateCardRequest]{
				status: http.StatusBadRequest,
				req:    UpdateCardRequest{Content: "Card content", DeckID: "DECK_ID"},
				res:    `{"errors":["ERROR_MESSAGE"]}`,
				want:   Card{},
				err:    "mochi: ERROR_MESSAGE",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testUpdateItem("/api/cards", tt.test, func(client *Client, req UpdateCardRequest) (any, error) {
			return client.UpdateCard(context.Background(), tt.test.id, req)
		}))
	}
}

func Test_DeleteCard(t *testing.T) {
	tests := []struct {
		name string
		test deleteItemTestCase
	}{
		{
			name: "should delete a card",
			test: deleteItemTestCase{
				status: http.StatusOK,
				id:     "CARD_ID",
				res:    nil,
				err:    "",
			},
		},
		{
			name: "should return an error",
			test: deleteItemTestCase{
				status: http.StatusBadRequest,
				id:     "CARD_ID",
				res:    `{"errors":["ERROR_MESSAGE"]}`,
				err:    "mochi: ERROR_MESSAGE",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testDeleteItem("/api/cards", tt.test, func(client *Client, id string) error {
			return client.DeleteCard(context.Background(), id)
		}))
	}
}
