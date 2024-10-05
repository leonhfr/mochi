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
				file: "sync:\n  - path: lorem-ipsum\n",
				err:  nil,
			},
			want: &Config{Sync: []Sync{{Path: "lorem-ipsum"}}},
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
				file: "sync:\n  - path: lorem-ipsum\n",
				err:  nil,
			},
			want: &Config{Sync: []Sync{{Path: "lorem-ipsum"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := new(mockReader)
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

type mockReader struct {
	mock.Mock
}

func (m *mockReader) Exists(p string) bool {
	args := m.Mock.Called(p)
	return args.Bool(0)
}

func (m *mockReader) Read(p string) (io.ReadCloser, error) {
	args := m.Mock.Called(p)
	rc := strings.NewReader(args.String(0))
	return io.NopCloser(rc), args.Error(1)
}
