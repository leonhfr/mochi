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
	type image struct {
		content []byte
		hash    string
	}

	tests := []struct {
		name     string
		job      *deckJob
		lock     *Lock
		cards    map[string][]api.Card
		markdown map[string]string
		images   map[string]image
		want     []cardRequest
	}{
		{
			"generate card requests",
			&deckJob{
				id: "id_root",
				sources: []string{
					"/note-1.md",
					"/note-2.md",
					"/note-3.md",
					"/images.md",
				},
				parser: parser.NewNote(),
			},
			&Lock{
				Decks:  map[string][2]string{},
				Images: map[string]map[string]string{},
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
				"/images.md": "# Images\n\n![Image 1](path/to/image-1.jpg)\n\n![Image 1](another/path/to/image-2.jpg)",
			},
			map[string]image{
				"/path/to/image-1.jpg":         {[]byte("Image 1 content."), "image_hash_1"},
				"/another/path/to/image-2.jpg": {[]byte("Image 2 content."), "image_hash_2"},
			},
			[]cardRequest{
				{
					id:      "id_note_2",
					kind:    updateRequest,
					deckID:  "id_root",
					content: "# Note 2\n\nContent 2\n",
				},
				{
					kind:    createRequest,
					deckID:  "id_root",
					content: "# Note 3\n\nContent 3\n",
				},
				{
					kind:    createRequest,
					deckID:  "id_root",
					content: "# Images\n\n![Image 1](@media/c1816e0497517666.jpg)\n\n![Image 1](@media/5ac642a4b61d6ca1.jpg)\n",
					images: []syncImage{
						{
							attachment: api.Attachment{
								FileName:    "c1816e0497517666.jpg",
								ContentType: "image/jpg",
								Data:        "Image 1 content.",
							},
							path: "/path/to/image-1.jpg",
							hash: "image_hash_1",
						},
						{
							attachment: api.Attachment{
								FileName:    "5ac642a4b61d6ca1.jpg",
								ContentType: "image/jpg",
								Data:        "Image 2 content.",
							},
							path: "/another/path/to/image-2.jpg",
							hash: "image_hash_2",
						},
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
			for path, content := range tt.markdown {
				fs.On("Read", path).Return([]byte(content), nil)
			}
			for path, image := range tt.images {
				fs.On("FileExists", path).Return(true)
				fs.On("Image", path).Return(image.content, image.hash, nil)
			}

			got, err := generateCardRequests(context.Background(), tt.job, tt.lock, client, fs)

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
				"/note.md": "# Note\n\nA simple note\n\n![Image](../images/image.jpg)\n",
			},
			[]parser.Card{
				{
					Name:    "Note",
					Content: "# Note\n\nA simple note\n\n![Image](@media/b7e04c679d3e44ec.jpg)\n",
					Fields:  map[string]string{},
					Images: map[string]parser.Image{
						"/images/image.jpg": {
							Destination: "../images/image.jpg",
							FileName:    "b7e04c679d3e44ec",
							Extension:   "jpg",
							ContentType: "image/jpg",
							AltText:     "Image",
						},
					},
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
