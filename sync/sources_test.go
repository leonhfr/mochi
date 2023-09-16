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
		sources []string
		want    []string
	}{
		{
			"all files",
			[]string{
				"/journal/yyyy-mm-dd.md",
				"/german/vocabulary/s.md",
			},
			[]string{
				"/german/vocabulary/s.md",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := new(MockFilesystem)
			fs.On("Sources", extensions).Return(tt.sources, nil)

			got, err := Sources(config, fs)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			fs.AssertExpectations(t)
		})
	}
}
