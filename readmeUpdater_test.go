package devscripts

import (
	"strings"
	"testing"
)

// Tests for backward compatibility API
func TestBackwardCompatibilityReadmeUpdater(t *testing.T) {
	runner := NewScriptRunner()

	t.Run("UpdateOriginalReadme", func(t *testing.T) {
		// Execute update
		update, err := runner.UpdateReadmeIfNeeded("README.md")
		if err != nil {
			t.Fatalf("Error updating README: %v", err)
		}
		t.Logf("Readme Updated: %v", update)
	})

	t.Run("GenerateReadmeSection", func(t *testing.T) {
		section, err := runner.GenerateReadmeSection()
		if err != nil {
			t.Fatalf("Error generating README section: %v", err)
		}

		if !strings.Contains(section, "## Available Scripts") {
			t.Error("Section should contain header")
		}
	})

	t.Run("GetScriptNames", func(t *testing.T) {
		scripts, err := runner.GetScriptNames()
		if err != nil {
			t.Fatalf("Error getting script names: %v", err)
		}

		if len(scripts) == 0 {
			t.Error("Should find some scripts")
		}

		// Verify scripts have .sh extension
		for _, script := range scripts {
			if !strings.HasSuffix(script, ".sh") {
				t.Errorf("Script %s should have .sh extension", script)
			}
		}
	})
}
