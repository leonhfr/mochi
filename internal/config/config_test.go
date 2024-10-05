package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/mochi/internal/file"
)

func Test_Parse(t *testing.T) {
	tests := []struct {
		name   string
		target string
		want   Config
		err    error
	}{
		{
			name:   "no config found",
			target: "testdata/noconfig",
			want:   Config{},
			err:    ErrNoConfig,
		},
		{
			name:   "mochi.yaml",
			target: "testdata/yaml",
			want:   Config{Sync: []Sync{{Path: "lorem-ipsum"}}},
		},
		{
			name:   "mochi.yml",
			target: "testdata/yml",
			want:   Config{Sync: []Sync{{Path: "lorem-ipsum"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(file.NewSystem(), tt.target)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.err, err)
		})
	}
}
