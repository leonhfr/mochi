package github

import (
	"fmt"

	"github.com/sethvargo/go-githubactions"
)

const lockfileUpdatedOutput = "lockfile_updated"

// SetOutput sets the action output.
func SetOutput(gha *githubactions.Action, updated bool) {
	gha.SetOutput(lockfileUpdatedOutput, fmt.Sprintf("%t", updated))
}
