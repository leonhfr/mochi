package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Note_Convert(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   []Card
	}{
		{
			"comment",
			"<!-- Comment. -->\n\n# Noun\n\n## Gender\n\nSome stuff about genders.\n\n- der\n- die\n- das\n",
			[]Card{
				{
					Name:    "Noun",
					Content: "<!-- Comment. -->\n\n# Noun\n\n## Gender\n\nSome stuff about genders.\n\n- der\n- die\n- das\n",
					Fields:  map[string]string{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNote().Convert([]byte(tt.source))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
