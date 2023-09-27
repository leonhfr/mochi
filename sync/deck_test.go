package sync

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/leonhfr/mochi/api"
)

func Test_SynchronizeDecks(t *testing.T) {
	type want struct {
		lock *Lock
		res  DeckResult
	}

	tests := []struct {
		name        string
		sources     []string
		lock        *Lock
		config      Config
		deckCreates map[api.CreateDeckRequest]api.Deck
		deckUpdates map[string]api.UpdateDeckRequest
		want        want
	}{
		{
			name: "success",
			sources: []string{
				"/note.md",
				"/german/vocabulary/s.md",
				"/german/grammar/noun.md",
			},
			lock: &Lock{
				data: lockData{
					"id_root":              {Path: "/", Name: "Root", Cards: map[string]lockCard{}},
					"id_german":            {Path: "/german", Name: "German", Cards: map[string]lockCard{}},
					"id_german_vocabulary": {Path: "/german/vocabulary", Name: "Vocabulary", Cards: map[string]lockCard{}},
				},
			},
			config: Config{Sync: []Sync{
				{Path: "/", Name: "Notes (root)"},
				{Path: "/german/vocabulary"},
				{Path: "/german/grammar"},
			}},
			deckCreates: map[api.CreateDeckRequest]api.Deck{
				{Name: "Grammar", ParentID: "id_german"}: {Name: "Grammar", ID: "id_german_grammar"},
			},
			deckUpdates: map[string]api.UpdateDeckRequest{
				"id_root": {Name: "Notes (root)"},
			},
			want: want{
				&Lock{
					data: lockData{
						"id_root":              {Path: "/", Name: "Notes (root)", Cards: map[string]lockCard{}},
						"id_german":            {Path: "/german", Name: "German", Cards: map[string]lockCard{}},
						"id_german_grammar":    {Path: "/german/grammar", Name: "Grammar", Cards: map[string]lockCard{}},
						"id_german_vocabulary": {Path: "/german/vocabulary", Name: "Vocabulary", Cards: map[string]lockCard{}},
					},
					updated: true,
				},
				DeckResult{
					Created: 1,
					Updated: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := new(MockClient)
			for req, deck := range tt.deckCreates {
				client.On("CreateDeck", mock.Anything, req).Return(deck, nil)
			}
			for id, req := range tt.deckUpdates {
				client.On("UpdateDeck", mock.Anything, id, req).Return(api.Deck{}, nil)
			}

			got, err := SynchronizeDecks(context.Background(), tt.sources, tt.lock, tt.config, client, testLogger{})

			require.NoError(t, err)
			assert.Equal(t, tt.want.lock.data, tt.lock.data)
			assert.Equal(t, tt.want.lock.updated, tt.lock.updated)
			assert.Equal(t, tt.want.res, got)
			client.AssertExpectations(t)
		})
	}
}
