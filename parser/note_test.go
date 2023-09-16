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
					Images:  map[string]Image{},
				},
			},
		},
		{
			"front matter",
			"---\nfoo: bar\n---\n\n<!-- Comment. -->\n\n# Noun\n\n## Gender\n\nSome stuff about genders.\n\n- der\n- die\n- das\n",
			[]Card{
				{
					Name:    "Noun",
					Content: "<!-- Comment. -->\n\n# Noun\n\n## Gender\n\nSome stuff about genders.\n\n- der\n- die\n- das\n",
					Fields:  map[string]string{},
					Images:  map[string]Image{},
				},
			},
		},
		{
			"images",
			"# Noun\n\n## Gender\n\n![Example 1](../images/example-1.png)\n\nInline image: ![Example 2](./example-2.png)",
			[]Card{
				{
					Name:    "Noun",
					Content: "# Noun\n\n## Gender\n\n![Example 1](@media/ZErEWCjZqZiR61Nn.png)\n\nInline image: ![Example 2](@media/zKtslBV8vJnCFDAo.png)\n",
					Fields:  map[string]string{},
					Images: map[string]Image{
						"../images/example-1.png": {
							FileName:    "ZErEWCjZqZiR61Nn",
							Extension:   "png",
							ContentType: "image/png",
							AltText:     "Example 1",
						},
						"./example-2.png": {
							FileName:    "zKtslBV8vJnCFDAo",
							Extension:   "png",
							ContentType: "image/png",
							AltText:     "Example 2",
						},
					},
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
