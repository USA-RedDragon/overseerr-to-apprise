package cmd_test

import (
	"testing"

	"github.com/USA-RedDragon/overseerr-to-apprise/cmd"
)

//nolint:golint,gochecknoglobals
var requiredFlags = []string{}

func TestDefault(t *testing.T) {
	t.Parallel()
	baseCmd := cmd.NewCommand("testing", "default")
	// Avoid port conflict
	baseCmd.SetArgs(append([]string{"--http.metrics.enabled", "true", "--http.port", "8082", "--http.metrics.port", "8083", "--log_level", "error"}, requiredFlags...))
	err := baseCmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
