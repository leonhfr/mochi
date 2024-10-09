package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_List(t *testing.T) {
	want := []string{
		"/Root level card.md",
		"/journal/yyyy-mm-dd.md",
		"/headings/headings.md",
		"/lorem-ipsum/Lorem ipsum.md",
	}
	var got []string
	err := NewSystem().Walk(
		"../../testdata",
		[]string{".md"},
		func(path string) { got = append(got, path) },
	)
	assert.ElementsMatch(t, want, got)
	assert.NoError(t, err)
}

func Test_Exists(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"exists", "../../testdata/mochi.yml", true},
		{"does not exist", "../../testdata/mochi.yaml", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSystem().Exists(tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}
