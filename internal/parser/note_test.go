package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/mochi/internal/parser/image"
)

var _ cardParser = &note{}

func Test_Note_Convert(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		source string
		want   []Card
	}{
		{
			name:   "simple note",
			path:   "/testdata/lorem-ipsum/Lorem ipsum.md",
			source: "# Title 1\nParagraph.\n",
			want: []Card{{
				Name:     "Lorem ipsum",
				Content:  "# Title 1\nParagraph.\n",
				Filename: "Lorem ipsum.md",
				Images:   map[string]image.Image{},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := newMockFileChecker(nil)
			got, err := newNote(fc).convert(tt.path, []byte(tt.source))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
