package filesystem

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/leonhfr/mochi/test/data/base64"
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

func Test_Filesystem_Image(t *testing.T) {
	tests := []struct {
		path   string
		hash   string
		base64 []byte
	}{
		{"images/scream.png", "637b04d6cbd2a4a365fe57c16c90a046", bytes.TrimSpace(base64.Scream)},
	}

	fs := New(workspace)

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			base64, hash, err := fs.Image(tt.path)
			require.NoError(t, err)
			assert.Equal(t, tt.base64, base64)
			assert.Equal(t, tt.hash, hash)
		})
	}
}
