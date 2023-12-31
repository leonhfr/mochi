package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_vocabulary_convert(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   []Card
	}{
		{
			"ok",
			`---
title: YAML Frontmatter
---

<!-- This file was generated by github.com/leonhfr/vocabulary-action and is susceptible to be modified by automations. -->

# s

Schuh

Spaziergang "Example phrase." "Another example."

Spiegel
A note about the word.
Can be multiline.
`,
			[]Card{
				{
					Name:    "Schuh",
					Content: "# Schuh\n",
					Fields: map[string]string{
						"examples": "",
						"notes":    "",
						"word":     "Schuh",
					},
				},
				{
					Name:    "Spaziergang",
					Content: "# Spaziergang\n\n## Examples\n\nExample phrase.\n\nAnother example.\n\n## Notes\n",
					Fields: map[string]string{
						"examples": "Example phrase.\n\nAnother example.",
						"notes":    "",
						"word":     "Spaziergang",
					},
				},
				{
					Name:    "Spiegel",
					Content: "# Spiegel\n\n## Notes\n\nA note about the word.\n\nCan be multiline.\n",
					Fields: map[string]string{
						"examples": "",
						"notes":    "A note about the word.\n\nCan be multiline.",
						"word":     "Spiegel",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewVocabulary().Convert("", []byte(tt.source))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
