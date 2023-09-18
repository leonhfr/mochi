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
				lockName: "[decks]\n\"/deck_path\" = [\"deck_id\", \"Deck Name\"]\n",
			},
			[]api.Deck{
				{ID: "deck_id", Name: "Deck Name"},
			},
			&Lock{
				Decks: map[string][2]string{
					"/deck_path": {"deck_id", "Deck Name"},
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
			assert.Equal(t, tt.want.Decks, lock.Decks)
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
				Decks: map[string][2]string{
					"/deck_path": {"deck_id", "Deck Name"},
				},
				updated: false,
			},
			map[string]string{},
			false,
		},
		{
			"updated",
			&Lock{
				Decks: map[string][2]string{
					"/deck_path": {"deck_id", "Deck Name"},
				},
				Images: map[string]map[string]string{
					"card_id": {
						"../images/example-1.png": "md5_hash",
					},
				},
				updated: true,
			},
			map[string]string{
				lockName: "[decks]\n\"/deck_path\" = [\"deck_id\", \"Deck Name\"]\n\n[images]\n[images.card_id]\n\"../images/example-1.png\" = \"md5_hash\"\n",
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

func Test_updateLock(t *testing.T) {
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
				Decks: map[string][2]string{
					"/deck_path": {"deck_id", "Deck Name"},
				},
			},
			&Lock{
				Decks: map[string][2]string{
					"/deck_path": {"deck_id", "Deck Name"},
				},
			},
		},
		{
			"deck delete",
			[]api.Deck{},
			&Lock{
				Decks: map[string][2]string{
					"/deck_path": {"deck_id", "Deck Name"},
				},
			},
			&Lock{
				Decks:   map[string][2]string{},
				updated: true,
			},
		},
		{
			"deck name update",
			[]api.Deck{
				{ID: "deck_id", Name: "Deck Name"},
			},
			&Lock{
				Decks: map[string][2]string{
					"/deck_path": {"deck_id", "Wrong Name"},
				},
			},
			&Lock{
				Decks: map[string][2]string{
					"/deck_path": {"deck_id", "Deck Name"},
				},
				updated: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateLock(tt.lock, tt.decks)
			assert.Equal(t, tt.want, tt.lock)
		})
	}
}
