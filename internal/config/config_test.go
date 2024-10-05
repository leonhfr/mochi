package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			got, err := Parse(tt.target)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.err, err)
		})
	}
}
