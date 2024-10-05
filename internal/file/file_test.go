package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_List(t *testing.T) {
	want := []string{"/lorem-ipsum/Lorem ipsum.md"}
	got, err := NewSystem().List("../../testdata", []string{".md"})
	assert.Equal(t, want, got)
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

func Test_ParseJSON(t *testing.T) {
	var v any
	err := NewSystem().ParseJSON("../../testdata/mochi-lock.json", &v)
	assert.NoError(t, err)
}

func Test_ParseYAML(t *testing.T) {
	var v any
	err := NewSystem().ParseYAML("../../testdata/mochi.yml", &v)
	assert.NoError(t, err)
}
