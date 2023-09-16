package action

import (
	"io"
	"testing"

	"github.com/sethvargo/go-githubactions"
	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/mochi/api"
)

var _ Client = &api.Client{}

func Test_GetInput(t *testing.T) {
	//nolint:gosec
	token, workspace := "mochi_cards_token", "/mochi"

	tests := []struct {
		name   string
		envMap map[string]string
		want   Input
		err    string
	}{
		{
			name: "working",
			envMap: map[string]string{
				"INPUT_API_TOKEN":  token,
				"GITHUB_WORKSPACE": workspace,
			},
			want: Input{
				APIToken:  token,
				Workspace: workspace,
			},
			err: "",
		},
		{
			name: "optional",
			envMap: map[string]string{
				"INPUT_API_TOKEN":  token,
				"GITHUB_WORKSPACE": workspace,
			},
			want: Input{
				APIToken:  token,
				Workspace: workspace,
			},
			err: "",
		},
		{
			name: "missing api token",
			envMap: map[string]string{
				"GITHUB_WORKSPACE": workspace,
			},
			want: Input{},
			err:  "api token required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getenv := func(key string) string {
				return tt.envMap[key]
			}

			gha := githubactions.New(
				githubactions.WithWriter(io.Discard),
				githubactions.WithGetenv(getenv),
			)

			got, err := GetInput(gha)

			assert.Equal(t, tt.want, got)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
			}
		})
	}
}
