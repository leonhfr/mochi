package config

import (
	"io"
	"io/fs"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Parse(t *testing.T) {
	type (
		testRead struct {
			path string
			file string
			err  error
		}
	)
	tests := []struct {
		name    string
		target  string
		parsers []string
		read    []testRead
		want    *Config
		err     bool
	}{
		{
			name:    "no config found",
			target:  "testdata",
			parsers: []string{"note"},
			read: []testRead{
				{path: "testdata/mochi.yaml", err: fs.ErrNotExist},
				{path: "testdata/mochi.yml", err: fs.ErrNotExist},
			},
			err: true,
		},
		{
			name:    "mochi.yaml",
			target:  "testdata",
			parsers: []string{"note"},
			read: []testRead{
				{
					path: "testdata/mochi.yaml",
					file: "rootName: ROOT_NAME\ndecks:\n  - path: sed-interdum-libero\n    name: Sed interdum libero\n  - path: lorem-ipsum\n    name: Lorem ipsum\n",
				},
			},
			want: &Config{RateLimit: 50, RootName: "ROOT_NAME", Decks: []Deck{
				{Path: "/sed-interdum-libero", Name: "Sed interdum libero"},
				{Path: "/lorem-ipsum", Name: "Lorem ipsum"},
			}},
		},
		{
			name:    "mochi.yml",
			target:  "testdata",
			parsers: []string{"note"},
			read: []testRead{
				{
					path: "testdata/mochi.yaml",
					err:  fs.ErrNotExist,
				},
				{
					path: "testdata/mochi.yml",
					file: "rootName: ROOT_NAME\ndecks:\n  - path: lorem-ipsum\n",
				},
			},
			want: &Config{RateLimit: 50, RootName: "ROOT_NAME", Decks: []Deck{{Path: "/lorem-ipsum"}}},
		},
		{
			name:    "should set default root deck name",
			target:  "testdata",
			parsers: []string{"note"},
			read: []testRead{
				{
					path: "testdata/mochi.yaml",
					err:  fs.ErrNotExist,
				},
				{
					path: "testdata/mochi.yml",
					file: "decks:\n  - path: lorem-ipsum\n    name: Lorem ipsum\n",
				},
			},
			want: &Config{RateLimit: 50, RootName: "Root Deck", Decks: []Deck{{Path: "/lorem-ipsum", Name: "Lorem ipsum"}}},
		},
		{
			name:    "invalid config",
			target:  "testdata",
			parsers: []string{"note"},
			read: []testRead{
				{
					path: "testdata/mochi.yaml",
					file: "decks:\n  - name: Lorem ipsum\n",
				},
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := new(mockFile)
			for _, read := range tt.read {
				r.On("Read", read.path).Return(read.file, read.err)
			}

			got, err := Parse(r, tt.target, tt.parsers)
			assert.Equal(t, tt.want, got)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			r.AssertExpectations(t)
		})
	}
}

func Test_Config_Deck(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		path   string
		want   Deck
		ok     bool
	}{
		{
			name: "should return the expected deck config",
			config: &Config{Decks: []Deck{
				{Path: "/sed-interdum-libero"},
				{Path: "/lorem-ipsum"},
			}},
			path: "/lorem-ipsum",
			want: Deck{Path: "/lorem-ipsum"},
			ok:   true,
		},
		{
			name: "should return the expected deck config",
			config: &Config{Decks: []Deck{
				{Path: "/vocabulary/catalan"},
				{Path: "/vocabulary/german"},
			}},
			path: "/vocabulary/german",
			want: Deck{Path: "/vocabulary/german"},
			ok:   true,
		},
		{
			name: "should return nothing",
			config: &Config{Decks: []Deck{
				{Path: "/vocabulary/catalan"},
				{Path: "/vocabulary/german"},
			}},
			path: "/vocabulary",
		},
		{
			name: "should not return the deck config",
			config: &Config{Decks: []Deck{
				{Path: "/sub/lorem-ipsum"},
			}},
			path: "/sub/lorem-ipsum-other",
		},
		{
			name: "should return the root deck",
			config: &Config{RootName: "ROOT_NAME", Decks: []Deck{
				{Path: "/sed-interdum-libero"},
				{Path: "/lorem-ipsum"},
			}},
			path: "/",
			want: Deck{Path: "/", Name: "ROOT_NAME"},
			ok:   true,
		},
		{
			name: "should return false (skip root deck)",
			config: &Config{RootName: "ROOT_NAME", SkipRoot: true, Decks: []Deck{
				{Path: "/sed-interdum-libero"},
				{Path: "/lorem-ipsum"},
			}},
			path: "/",
		},
		{
			name: "should return false",
			config: &Config{Decks: []Deck{
				{Path: "/sed-interdum-libero"},
			}},
			path: "/lorem-ipsum",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.config.Deck(tt.path)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

type mockFile struct {
	mock.Mock
}

func (m *mockFile) Read(p string) (io.ReadCloser, error) {
	args := m.Mock.Called(p)
	rc := strings.NewReader(args.String(0))
	return io.NopCloser(rc), args.Error(1)
}
