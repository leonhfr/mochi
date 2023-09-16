package filesystem

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ Interface = &Filesystem{}

var (
	workspace = "../test/data"
	config    = `sync:
  - path: .
    name: Notes (root)
    parser: note
    archive: true
    walk: true

ignore:
  - journal/*
`
)

func Test_Filesystem_FileExists(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"mochi.yml", true},
		{"mochi.yaml", false},
		{"german", false},
	}

	fs := New(workspace)

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := fs.FileExists(tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Filesystem_Read(t *testing.T) {
	tests := []struct {
		path string
		want []byte
	}{
		{"mochi.yml", []byte(config)},
		{"mochi.yaml", nil},
	}

	fs := New(workspace)

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got, err := fs.Read(tt.path)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
