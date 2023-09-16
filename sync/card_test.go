package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/leonhfr/mochi/parser"
)

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
