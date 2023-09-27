package sync

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/leonhfr/mochi/api"
)

func Test_ReadLock(t *testing.T) {
	tests := []struct {
		name      string
		fileReads map[string]string
		deckList  []api.Deck
		want      *Lock
	}{
		{
			"success",
			map[string]string{
				lockName: `{"deck_id":{"path":"/deck_path","name":"Deck Name"}}`,
			},
			[]api.Deck{
				{ID: "deck_id", Name: "Deck Name"},
			},
			&Lock{
				data: lockData{
					"deck_id": lockDeck{Path: "/deck_path", Name: "Deck Name"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := new(MockClient)
			client.On("ListDecks", mock.Anything).Return(tt.deckList, nil)

			fs := new(MockFilesystem)
			for path, content := range tt.fileReads {
				fs.On("Read", path).Return([]byte(content), nil)
			}

			lock, err := ReadLock(context.Background(), client, fs)
			require.NoError(t, err)
			assert.Equal(t, tt.want.data, lock.data)
			assert.False(t, lock.updated)
			client.AssertExpectations(t)
			fs.AssertExpectations(t)
		})
	}
}

func Test_Lock_Write(t *testing.T) {
	tests := []struct {
		name       string
		lock       *Lock
		fileWrites map[string]string
		want       bool
	}{
		{
			"not updated",
			&Lock{
				data: lockData{
					"deck_id": lockDeck{Path: "/deck_path", Name: "Deck Name"},
				},
				updated: false,
			},
			map[string]string{},
			false,
		},
		{
			"updated",
			&Lock{
				data: lockData{
					"deck_id": lockDeck{
						Path: "/deck_path",
						Name: "Deck Name",
						Cards: map[string]lockCard{
							"card_id": {
								Filename: "note.md",
								Images: map[string]string{
									"../images/example-1.png": "md5_hash",
								},
							},
						},
					},
				},
				updated: true,
			},
			map[string]string{
				lockName: `{"deck_id":{"path":"/deck_path","name":"Deck Name","cards":{"card_id":{"filename":"note.md","images":{"../images/example-1.png":"md5_hash"}}}}}`,
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := new(MockFilesystem)
			for path, content := range tt.fileWrites {
				fs.On("Write", path, content).Return(nil)
			}

			updated, err := tt.lock.Write(fs)

			require.NoError(t, err)
			assert.Equal(t, tt.want, updated)
			fs.AssertExpectations(t)
		})
	}
}

func Test_Lock_cleanDecks(t *testing.T) {
	tests := []struct {
		name  string
		decks []api.Deck
		lock  *Lock
		want  *Lock
	}{
		{
			"no update",
			[]api.Deck{
				{ID: "deck_id", Name: "Deck Name"},
			},
			&Lock{
				data: lockData{
					"deck_id": lockDeck{Path: "/deck_path", Name: "Deck Name"},
				},
			},
			&Lock{
				data: lockData{
					"deck_id": lockDeck{Path: "/deck_path", Name: "Deck Name"},
				},
			},
		},
		{
			"deck delete",
			[]api.Deck{},
			&Lock{
				data: lockData{
					"deck_id": lockDeck{Path: "/deck_path", Name: "Deck Name"},
				},
			},
			&Lock{
				data:    lockData{},
				updated: true,
			},
		},
		{
			"deck name update",
			[]api.Deck{
				{ID: "deck_id", Name: "Deck Name"},
			},
			&Lock{
				data: lockData{
					"deck_id": lockDeck{Path: "/deck_path", Name: "Wrong Name"},
				},
			},
			&Lock{
				data: lockData{
					"deck_id": lockDeck{Path: "/deck_path", Name: "Deck Name"},
				},
				updated: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.lock.cleanDecks(tt.decks)
			assert.Equal(t, tt.want.data, tt.lock.data)
			assert.Equal(t, tt.want.updated, tt.lock.updated)
		})
	}
}
