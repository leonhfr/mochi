package config

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		name    string
		target  string
		parsers []string
		exists  []testExist
		read    *testRead
		want    *Config
		err     bool
	}{
		{
			name:    "no config found",
			target:  "testdata",
			parsers: []string{"note"},
			exists: []testExist{
				{"testdata/mochi.yaml", false},
				{"testdata/mochi.yml", false},
			},
			err: true,
		},
		{
			name:    "mochi.yaml",
			target:  "testdata",
			parsers: []string{"note"},
			exists: []testExist{
				{"testdata/mochi.yaml", true},
			},
			read: &testRead{
				path: "testdata/mochi.yaml",
				file: "rootName: ROOT_NAME\ndecks:\n  - path: sed-interdum-libero\n    name: Sed interdum libero\n  - path: lorem-ipsum\n    name: Lorem ipsum\n",
				err:  nil,
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
			exists: []testExist{
				{"testdata/mochi.yaml", false},
				{"testdata/mochi.yml", true},
			},
			read: &testRead{
				path: "testdata/mochi.yml",
				file: "rootName: ROOT_NAME\ndecks:\n  - path: lorem-ipsum\n",
				err:  nil,
			},
			want: &Config{RateLimit: 50, RootName: "ROOT_NAME", Decks: []Deck{{Path: "/lorem-ipsum"}}},
		},
		{
			name:    "should set default root deck name",
			target:  "testdata",
			parsers: []string{"note"},
			exists: []testExist{
				{"testdata/mochi.yaml", false},
				{"testdata/mochi.yml", true},
			},
			read: &testRead{
				path: "testdata/mochi.yml",
				file: "decks:\n  - path: lorem-ipsum\n    name: Lorem ipsum\n",
				err:  nil,
			},
			want: &Config{RateLimit: 50, RootName: "Root Deck", Decks: []Deck{{Path: "/lorem-ipsum", Name: "Lorem ipsum"}}},
		},
		{
			name:    "invalid config",
			target:  "testdata",
			parsers: []string{"note"},
			exists: []testExist{
				{"testdata/mochi.yaml", false},
				{"testdata/mochi.yml", true},
			},
			read: &testRead{
				path: "testdata/mochi.yml",
				file: "decks:\n  - name: Lorem ipsum\n",
				err:  nil,
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := new(mockFile)
			for _, te := range tt.exists {
				r.On("Exists", te.path).Return(te.exists)
			}
			if tt.read != nil {
				r.On("Read", tt.read.path).Return(tt.read.file, tt.read.err)
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

func Test_Config_GetDeck(t *testing.T) {
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
			name: "should return false",
			config: &Config{Decks: []Deck{
				{Path: "/sed-interdum-libero"},
			}},
			path: "/lorem-ipsum",
			want: Deck{},
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.config.GetDeck(tt.path)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

type mockFile struct {
	mock.Mock
}

func (m *mockFile) Exists(p string) bool {
	args := m.Mock.Called(p)
	return args.Bool(0)
}

func (m *mockFile) Read(p string) (io.ReadCloser, error) {
	args := m.Mock.Called(p)
	rc := strings.NewReader(args.String(0))
	return io.NopCloser(rc), args.Error(1)
}
