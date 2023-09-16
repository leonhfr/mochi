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
						Content: "# Note 1\n\nContent 1",
					},
					{
						DeckID:  "id_root",
						ID:      "id_note_2",
						Name:    "Note 2",
						Content: "# Note 1\n\nWrong content.",
					},
				},
			},
			map[string]string{
				"/note-1.md": "# Note 1\n\nContent 1",
				"/note-2.md": "# Note 2\n\nContent 2",
				"/note-3.md": "# Note 3\n\nContent 3",
			},
			[]cardRequest{
				{
					id:   "id_note_2",
					kind: updateRequest,
					update: api.UpdateCardRequest{
						DeckID:  "id_root",
						Content: "# Note 2\n\nContent 2",
					},
				},
				{
					kind: createRequest,
					create: api.CreateCardRequest{
						DeckID:  "id_root",
						Content: "# Note 3\n\nContent 3",
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
				"/note.md": "# Note\n\nA simple note",
			},
			[]parser.Card{
				{
					Name:    "Note",
					Content: "# Note\n\nA simple note",
					Fields:  map[string]string{},
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
				"/german/vocabulary/s.md": "# s\n\nSpaziergang\nNotes notes.\n\nSpiegel",
				"/german/vocabulary/p.md": "# p\n\nPapagei",
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
