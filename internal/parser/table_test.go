package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var tableSource = `
| Infinitive | Present       | Past   | Participle   | English            |
| ---------- | ------------- | ------ | ------------ | ------------------ |
| backen     | backt (bäckt) | [buk]  | hat gebacken | to bake            |
| befehlen   | befiehlt      | befahl | hat befohlen | to order, instruct |
`

func Test_table_convert(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		source string
		want   Result
	}{
		{
			name:   "should return the expected result",
			path:   "/testdata/verbs/Strong Verbs.md",
			source: tableSource,
			want: Result{
				Deck: "Strong Verbs", Cards: []Card{
					tableCard{
						name:     "backen|backt (bäckt)|[buk]|hat gebacken|to bake",
						headers:  []string{"Infinitive", "Present", "Past", "Participle", "English"},
						cells:    []string{"backen", "backt (bäckt)", "[buk]", "hat gebacken", "to bake"},
						path:     "/testdata/verbs/Strong Verbs.md",
						position: "StrongVerbsmd0000",
					},
					tableCard{
						name:     "befehlen|befiehlt|befahl|hat befohlen|to order, instruct",
						headers:  []string{"Infinitive", "Present", "Past", "Participle", "English"},
						cells:    []string{"befehlen", "befiehlt", "befahl", "hat befohlen", "to order, instruct"},
						path:     "/testdata/verbs/Strong Verbs.md",
						position: "StrongVerbsmd0001",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newTable().convert(tt.path, []byte(tt.source))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
