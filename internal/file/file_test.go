package file

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Walk(t *testing.T) {
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

func Test_Open_Error(t *testing.T) {
	tests := []struct {
		name string
		path string
		err  error
	}{
		{"exists", "../../testdata/mochi.yml", nil},
		{"does not exist", "../../testdata/mochi.yaml", fs.ErrNotExist},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc, err := NewSystem().Read(tt.path)
			assert.Equal(t, tt.err, err)
			defer rc.Close()
		})
	}
}
