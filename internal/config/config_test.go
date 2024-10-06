package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/mochi/internal/test"
)

func Test_Parse(t *testing.T) {
	type (
		testExist struct {
			path   string
			exists bool
		}
		testRead struct {
			path string
			file string
			err  error
		}
	)
	tests := []struct {
		name   string
		target string
		exists []testExist
		read   *testRead
		want   *Config
		err    error
	}{
		{
			name:   "no config found",
			target: "testdata",
			exists: []testExist{
				{"testdata/mochi.yaml", false},
				{"testdata/mochi.yml", false},
			},
			want: nil,
			err:  ErrNoConfig,
		},
		{
			name:   "mochi.yaml",
			target: "testdata",
			exists: []testExist{
				{"testdata/mochi.yaml", true},
			},
			read: &testRead{
				path: "testdata/mochi.yaml",
				file: "decks:\n  - path: sed-interdum-libero\n    name: Sed interdum libero\n  - path: lorem-ipsum\n    name: Lorem ipsum\n",
				err:  nil,
			},
			want: &Config{Decks: []Deck{
				{Path: "/sed-interdum-libero", Name: "Sed interdum libero"},
				{Path: "/lorem-ipsum", Name: "Lorem ipsum"},
			}},
		},
		{
			name:   "mochi.yml",
			target: "testdata",
			exists: []testExist{
				{"testdata/mochi.yaml", false},
				{"testdata/mochi.yml", true},
			},
			read: &testRead{
				path: "testdata/mochi.yml",
				file: "decks:\n  - path: lorem-ipsum\n    name: Lorem ipsum\n",
				err:  nil,
			},
			want: &Config{Decks: []Deck{{Path: "/lorem-ipsum", Name: "Lorem ipsum"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := new(test.MockFile)
			for _, te := range tt.exists {
				r.On("Exists", te.path).Return(te.exists)
			}
			if tt.read != nil {
				r.On("Read", tt.read.path).Return(tt.read.file, tt.read.err)
			}

			got, err := Parse(r, tt.target)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.err, err)
			r.AssertExpectations(t)
		})
	}
}

func Test_Config_Deck(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		base   string
		want   Deck
		ok     bool
	}{
		{
			name: "should return the expected deck config",
			config: &Config{Decks: []Deck{
				{Path: "/sed-interdum-libero"},
				{Path: "/lorem-ipsum"},
			}},
			base: "/lorem-ipsum",
			want: Deck{Path: "/lorem-ipsum"},
			ok:   true,
		},
		{
			name: "should return false",
			config: &Config{Decks: []Deck{
				{Path: "/sed-interdum-libero"},
			}},
			base: "/lorem-ipsum",
			want: Deck{},
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.config.GetDeck(tt.base)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}
