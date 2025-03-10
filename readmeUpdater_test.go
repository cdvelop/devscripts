package devscripts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpdateOriginalReadme(t *testing.T) {
	runner := NewScriptRunner()
	// Execute update
	update, err := runner.UpdateReadmeIfNeeded("README.md")
	if err != nil {
		t.Fatalf("Error updating README: %v", err)
	}
	t.Logf("Readme Updated: %v", update)
}

func TestUpdateFakeReadme(t *testing.T) {
	// Set up temporary directory for tests
	tmpDir := t.TempDir()

	// Create a temporary scripts directory for test files
	scriptsDir := filepath.Join(tmpDir, "bash_scripts")
	err := os.MkdirAll(scriptsDir, 0755)
	if err != nil {
		t.Fatal("Failed to create scripts directory:", err)
	}

	// Create a test script
	testScript := `#!/bin/bash
# Description: Test script for unit tests
# Usage: ./test.sh [arg1]
echo "Test script" 
`
	err = os.WriteFile(filepath.Join(scriptsDir, "test.sh"), []byte(testScript), 0644)
	if err != nil {
		t.Fatal("Failed to create test script:", err)
	}

	// Initialize the script runner with the test directory
	runner := NewScriptRunner(scriptsDir)

	t.Run("Update section while preserving existing content", func(t *testing.T) {
		// Create a test README with content outside the section
		testContent := `# Example Project

<!-- SCRIPTS_SECTION_START -->
Old content that should be replaced
<!-- SCRIPTS_SECTION_END -->

## Other Contents
This text should remain intact.`

		testPath := filepath.Join(tmpDir, "README.md")
		err := os.WriteFile(testPath, []byte(testContent), 0644)
		if err != nil {
			t.Fatal(err)
		}

		// Execute update
		updated, err := runner.UpdateReadmeIfNeeded(testPath)
		if err != nil {
			t.Fatalf("Error updating README: %v", err)
		}

		if !updated {
			t.Error("Should have detected changes and updated the README")
		}

		// Verify updated content
		content, err := os.ReadFile(testPath)
		if err != nil {
			t.Fatal(err)
		}

		// Verify that the external content is preserved
		if !strings.Contains(string(content), "## Other Contents") ||
			!strings.Contains(string(content), "This text should remain intact") {
			t.Error("Existing content outside of the script section was lost")
		}

		// Verify that the section was updated with the script information
		if !strings.Contains(string(content), "test.sh") {
			t.Error("The script section was not updated correctly")
		}
	})

	t.Run("Add section to empty README", func(t *testing.T) {
		testPath := filepath.Join(tmpDir, "EMPTY_README.md")

		// Execute update on a new file
		updated, err := runner.UpdateReadmeIfNeeded(testPath)
		if err != nil {
			t.Fatal(err)
		}

		if !updated {
			t.Error("Should have created a new README")
		}

		content, err := os.ReadFile(testPath)
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(content), "## Available Scripts") {
			t.Error("The section was not created correctly")
		}

		if !strings.Contains(string(content), "test.sh") {
			t.Error("Script information was not included in the README")
		}
	})
}
