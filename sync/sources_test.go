package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Sources(t *testing.T) {
	config := Config{Ignore: []string{
		"/journal/**",
	}}

	tests := []struct {
		name    string
		changed []string
		sources []string
		want    []string
	}{
		{
			"all files",
			[]string{},
			[]string{
				"/journal/yyyy-mm-dd.md",
				"/german/vocabulary/s.md",
			},
			[]string{
				"/german/vocabulary/s.md",
			},
		},
		{
			"changed files",
			[]string{
				"/german/vocabulary/s.md",
			},
			[]string{
				"/journal/yyyy-mm-dd.md",
				"/german/vocabulary/f.md",
				"/german/vocabulary/p.md",
				"/german/vocabulary/s.md",
				"/german/grammar/noun.md",
			},
			[]string{
				"/german/vocabulary/f.md",
				"/german/vocabulary/p.md",
				"/german/vocabulary/s.md",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := new(MockFilesystem)
			fs.On("Sources", extensions).Return(tt.sources, nil)

			got, err := Sources(tt.changed, config, fs)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			fs.AssertExpectations(t)
		})
	}
}

func Test_uniqueDirs(t *testing.T) {
	sources := []string{
		"/journal/yyyy-mm-dd.md",
		"/german/vocabulary/f.md",
		"/german/vocabulary/p.md",
		"/german/vocabulary/s.md",
		"/german/grammar/noun.md",
	}
	want := []string{
		"/",
		"/german",
		"/german/grammar",
		"/german/vocabulary",
		"/journal",
	}

	got := uniqueDirs(sources)
	assert.Equal(t, want, got)
}
