package sync

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/leonhfr/mochi/api"
	"github.com/leonhfr/mochi/parser"
)

func Test_generateCardRequests(t *testing.T) {
	tests := []struct {
		name  string
		job   *deckJob
		cards map[string][]api.Card
		files map[string]string
		want  []cardRequest
	}{
		{
			"generate card requests",
			&deckJob{
				id: "id_root",
				sources: []string{
					"/note-1.md",
					"/note-2.md",
					"/note-3.md",
				},
				parser: parser.NewNote(),
			},
			map[string][]api.Card{
				"id_root": {
					{
						DeckID:  "id_root",
						ID:      "id_note_1",
						Name:    "Note 1",
						Content: "# Note 1\n\nContent 1\n",
					},
					{
						DeckID:  "id_root",
						ID:      "id_note_2",
						Name:    "Note 2",
						Content: "# Note 1\n\nWrong content.\n",
					},
				},
			},
			map[string]string{
				"/note-1.md": "# Note 1\n\nContent 1\n",
				"/note-2.md": "# Note 2\n\nContent 2\n",
				"/note-3.md": "# Note 3\n\nContent 3\n",
			},
			[]cardRequest{
				{
					id:   "id_note_2",
					kind: updateRequest,
					update: api.UpdateCardRequest{
						DeckID:  "id_root",
						Content: "# Note 2\n\nContent 2\n",
					},
				},
				{
					kind: createRequest,
					create: api.CreateCardRequest{
						DeckID:  "id_root",
						Content: "# Note 3\n\nContent 3\n",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Client
			client := new(MockClient)
			for id, cards := range tt.cards {
				client.On("ListCardsInDeck", mock.Anything, id).Return(cards, nil)
			}

			// Filesystem
			fs := new(MockFilesystem)
			for path, content := range tt.files {
				fs.On("Read", path).Return([]byte(content), nil)
			}

			got, err := generateCardRequests(context.Background(), tt.job, client, fs)

			require.NoError(t, err)
			assert.ElementsMatch(t, tt.want, got)
			client.AssertExpectations(t)
			fs.AssertExpectations(t)
		})
	}
}

func Test_parseCards(t *testing.T) {
	tests := []struct {
		name  string
		job   *deckJob
		files map[string]string
		want  []parser.Card
	}{
		{
			"note",
			&deckJob{
				sources: []string{
					"/note.md",
				},
				parser: parser.NewNote(),
			},
			map[string]string{
				"/note.md": "# Note\n\nA simple note\n",
			},
			[]parser.Card{
				{
					Name:    "Note",
					Content: "# Note\n\nA simple note\n",
					Fields:  map[string]string{},
					Images:  map[string]parser.Image{},
				},
			},
		},
		{
			"vocabulary",
			&deckJob{
				sources: []string{
					"/german/vocabulary/s.md",
					"/german/vocabulary/p.md",
				},
				parser: parser.NewVocabulary(),
			},
			map[string]string{
				"/german/vocabulary/s.md": "# s\n\nSpaziergang\nNotes notes.\n\nSpiegel\n",
				"/german/vocabulary/p.md": "# p\n\nPapagei\n",
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := new(MockFilesystem)
			for path, content := range tt.files {
				fs.On("Read", path).Return([]byte(content), nil)
			}

			cards, err := parseCards(tt.job, fs)

			require.NoError(t, err)
			assert.Equal(t, tt.want, cards)
			fs.AssertExpectations(t)
		})
	}
}
