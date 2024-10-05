package mochi

import (
	"context"
	"net/http"
	"testing"
)

func Test_CreateDeck(t *testing.T) {
	tests := []struct {
		name string
		test createItemTestCase[CreateDeckRequest]
	}{
		{
			name: "should create a deck",
			test: createItemTestCase[CreateDeckRequest]{
				status: http.StatusCreated,
				req:    CreateDeckRequest{Name: "DeckName", ParentID: "PARENT_ID"},
				res:    Deck{ID: "DECK_ID", Name: "DeckName", ParentID: "PARENT_ID"},
				want:   Deck{ID: "DECK_ID", Name: "DeckName", ParentID: "PARENT_ID"},
				err:    "",
			},
		},
		{
			name: "should return an error",
			test: createItemTestCase[CreateDeckRequest]{
				status: http.StatusBadRequest,
				req:    CreateDeckRequest{Name: "DeckName", ParentID: "PARENT_ID"},
				res:    `{"errors":["ERROR_MESSAGE"]}`,
				want:   Deck{},
				err:    "mochi: ERROR_MESSAGE",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testCreateItem("/api/decks", tt.test, func(client *Client, req CreateDeckRequest) (any, error) {
			return client.CreateDeck(context.Background(), req)
		}))
	}
}

func Test_GetDeck(t *testing.T) {
	tests := []struct {
		name string
		test getItemTestCase
	}{
		{
			name: "should get a deck",
			test: getItemTestCase{
				status: http.StatusOK,
				id:     "DECK_ID",
				res:    Deck{ID: "DECK_ID", Name: "DeckName", ParentID: "PARENT_ID"},
				want:   Deck{ID: "DECK_ID", Name: "DeckName", ParentID: "PARENT_ID"},
				err:    "",
			},
		},
		{
			name: "should return an error",
			test: getItemTestCase{
				status: http.StatusBadRequest,
				id:     "DECK_ID",
				res:    `{"errors":["ERROR_MESSAGE"]}`,
				want:   Deck{},
				err:    "mochi: ERROR_MESSAGE",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testGetItem("/api/decks", tt.test, func(client *Client, id string) (any, error) {
			return client.GetDeck(context.Background(), id)
		}))
	}
}

func Test_ListDecks(t *testing.T) {
	tests := []struct {
		name string
		test listItemTestCase
	}{
		{
			name: "should call the callback once",
			test: listItemTestCase{
				responses: []listItemTestCaseResponse{
					{
						status: http.StatusOK,
						params: map[string]string{"limit": "100"},
						res: listResponse[Deck]{
							Docs: []Deck{{ID: "DECK_ID", Name: "DeckName", ParentID: "PARENT_ID"}},
						},
						want: []Deck{
							{ID: "DECK_ID", Name: "DeckName", ParentID: "PARENT_ID"},
						},
					},
				},
				total: 1,
			},
		},
		{
			name: "should call the callback several times",
			test: listItemTestCase{
				responses: []listItemTestCaseResponse{
					{
						status: http.StatusOK,
						params: map[string]string{"limit": "100"},
						res: listResponse[Deck]{
							Docs:     []Deck{{ID: "DECK_ID_1", Name: "DeckName1", ParentID: "PARENT_ID"}},
							Bookmark: "BOOKMARK_1",
						},
						want: []Deck{
							{ID: "DECK_ID_1", Name: "DeckName1", ParentID: "PARENT_ID"},
						},
					},
					{
						status: http.StatusOK,
						params: map[string]string{"limit": "100", "bookmark": "BOOKMARK_1"},
						res: listResponse[Deck]{
							Docs: []Deck{{ID: "DECK_ID_2", Name: "DeckName2", ParentID: "PARENT_ID"}},
						},
						want: []Deck{
							{ID: "DECK_ID_2", Name: "DeckName2", ParentID: "PARENT_ID"},
						},
					},
				},
				total: 2,
			},
		},
		{
			name: "should return an error",
			test: listItemTestCase{
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
		t.Run(tt.name, testListItem("/api/decks", tt.test, func(client *Client, _ string, cb func([]Deck) error) error {
			return client.ListDecks(context.Background(), cb)
		}))
	}
}

func Test_UpdateDeck(t *testing.T) {
	tests := []struct {
		name string
		test updateItemTestCase[UpdateDeckRequest]
	}{
		{
			name: "should update a deck",
			test: updateItemTestCase[UpdateDeckRequest]{
				status: http.StatusCreated,
				req:    UpdateDeckRequest{Name: "DeckName", ParentID: "PARENT_ID"},
				res:    Deck{ID: "DECK_ID", Name: "DeckName", ParentID: "PARENT_ID"},
				want:   Deck{ID: "DECK_ID", Name: "DeckName", ParentID: "PARENT_ID"},
				err:    "",
			},
		},
		{
			name: "should return an error",
			test: updateItemTestCase[UpdateDeckRequest]{
				status: http.StatusBadRequest,
				req:    UpdateDeckRequest{Name: "DeckName", ParentID: "PARENT_ID"},
				res:    `{"errors":["ERROR_MESSAGE"]}`,
				want:   Deck{},
				err:    "mochi: ERROR_MESSAGE",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testUpdateItem("/api/decks", tt.test, func(client *Client, req UpdateDeckRequest) (any, error) {
			return client.UpdateDeck(context.Background(), tt.test.id, req)
		}))
	}
}

func Test_DeleteDeck(t *testing.T) {
	tests := []struct {
		name string
		test deleteItemTestCase
	}{
		{
			name: "should delete a deck",
			test: deleteItemTestCase{
				status: http.StatusOK,
				id:     "DECK_ID",
				res:    nil,
				err:    "",
			},
		},
		{
			name: "should return an error",
			test: deleteItemTestCase{
				status: http.StatusBadRequest,
				id:     "DECK_ID",
				res:    `{"errors":["ERROR_MESSAGE"]}`,
				err:    "mochi: ERROR_MESSAGE",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testDeleteItem("/api/decks", tt.test, func(client *Client, id string) error {
			return client.DeleteDeck(context.Background(), id)
		}))
	}
}
