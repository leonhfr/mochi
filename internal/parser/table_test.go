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

func Test_table_parse(t *testing.T) {
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
					{
						Content:  "|Headers|Values|\n|---|---|\n|Infinitive|backen|\n|Present|backt (bäckt)|\n|Past|[buk]|\n|Participle|hat gebacken|\n|English|to bake|\n",
						Fields:   map[string]string{"name": "backen|backt (bäckt)|[buk]|hat gebacken|to bake"},
						Path:     "/testdata/verbs/Strong Verbs.md",
						Position: "StrongVerbsmd0000",
					},
					{
						Content:  "|Headers|Values|\n|---|---|\n|Infinitive|befehlen|\n|Present|befiehlt|\n|Past|befahl|\n|Participle|hat befohlen|\n|English|to order, instruct|\n",
						Fields:   map[string]string{"name": "befehlen|befiehlt|befahl|hat befohlen|to order, instruct"},
						Path:     "/testdata/verbs/Strong Verbs.md",
						Position: "StrongVerbsmd0001",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newTable().parse(tt.path, []byte(tt.source))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
