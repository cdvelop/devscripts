package devscripts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDevScriptsReadmeUpdater(t *testing.T) {
	// Set up temporary directory for tests
	tmpDir := t.TempDir()

	// Create test scripts
	testScript := `#!/bin/bash
# Description: Test script for unit tests
# Usage: ./test.sh [arg1]
echo "Test script" 
`
	err := os.WriteFile(filepath.Join(tmpDir, "test.sh"), []byte(testScript), 0644)
	if err != nil {
		t.Fatal("Failed to create test script:", err)
	}

	updater := NewDevScriptsReadmeUpdater(tmpDir)

	t.Run("GenerateScriptsSection", func(t *testing.T) {
		section, err := updater.GenerateScriptsSection()
		if err != nil {
			t.Fatalf("GenerateScriptsSection failed: %v", err)
		}

		// Check that section contains expected elements
		if !strings.Contains(section, "## Available Scripts") {
			t.Error("Section should contain '## Available Scripts' header")
		}

		if !strings.Contains(section, "automatically generated") {
			t.Error("Section should contain 'automatically generated' note")
		}

		if !strings.Contains(section, "test.sh") {
			t.Error("Section should contain script name")
		}

		if !strings.Contains(section, "Test script for unit tests") {
			t.Error("Section should contain script description")
		}

		if !strings.Contains(section, "./test.sh [arg1]") {
			t.Error("Section should contain script usage")
		}

		// Check markdown table structure
		if !strings.Contains(section, "| Script Name") {
			t.Error("Section should contain table header")
		}

		if !strings.Contains(section, "| ---") {
			t.Error("Section should contain table separator")
		}
	})

	t.Run("UpdateReadme creates new file", func(t *testing.T) {
		testPath := filepath.Join(tmpDir, "NEW_README.md")

		err := updater.UpdateReadme(testPath)
		if err != nil {
			t.Fatalf("UpdateReadme failed: %v", err)
		}

		// Verify file was created
		if _, err := os.Stat(testPath); os.IsNotExist(err) {
			t.Error("README file should have been created")
		}

		// Read and verify content
		content, err := os.ReadFile(testPath)
		if err != nil {
			t.Fatalf("Failed to read created README: %v", err)
		}

		contentStr := string(content)

		// Check for new section format
		if !strings.Contains(contentStr, "<!-- START_SECTION:SCRIPTS_SECTION -->") {
			t.Error("File should contain new section start marker")
		}

		if !strings.Contains(contentStr, "<!-- END_SECTION:SCRIPTS_SECTION -->") {
			t.Error("File should contain new section end marker")
		}

		if !strings.Contains(contentStr, "## Available Scripts") {
			t.Error("File should contain scripts section")
		}

		// Clean up
		os.Remove(testPath)
	})

	t.Run("UpdateReadmeIfNeeded detects no changes", func(t *testing.T) {
		testPath := filepath.Join(tmpDir, "NOCHANGE_README.md")

		// First update
		changed1, err := updater.UpdateReadmeIfNeeded(testPath)
		if err != nil {
			t.Fatalf("First UpdateReadmeIfNeeded failed: %v", err)
		}

		if !changed1 {
			t.Error("First update should detect changes")
		}

		// Second update should detect no changes
		changed2, err := updater.UpdateReadmeIfNeeded(testPath)
		if err != nil {
			t.Fatalf("Second UpdateReadmeIfNeeded failed: %v", err)
		}

		if changed2 {
			t.Error("Second update should not detect changes")
		}

		// Clean up
		os.Remove(testPath)
	})

	t.Run("UpdateReadmeIfNeeded detects changes when content differs", func(t *testing.T) {
		testPath := filepath.Join(tmpDir, "CHANGE_README.md")

		// Create initial file with different content
		initialContent := `# Test README

<!-- START_SECTION:SCRIPTS_SECTION -->
Old content that should be replaced
<!-- END_SECTION:SCRIPTS_SECTION -->

## Other Content
This should remain.`

		err := os.WriteFile(testPath, []byte(initialContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create initial README: %v", err)
		}

		// Update should detect changes
		changed, err := updater.UpdateReadmeIfNeeded(testPath)
		if err != nil {
			t.Fatalf("UpdateReadmeIfNeeded failed: %v", err)
		}

		if !changed {
			t.Error("Update should detect changes when content differs")
		}

		// Verify content was updated
		content, err := os.ReadFile(testPath)
		if err != nil {
			t.Fatalf("Failed to read updated README: %v", err)
		}

		contentStr := string(content)

		// Should contain the new script content
		if !strings.Contains(contentStr, "test.sh") {
			t.Error("Updated content should contain script information")
		}

		// Should preserve other content
		if !strings.Contains(contentStr, "## Other Content") {
			t.Error("Other content should be preserved")
		}

		// Clean up
		os.Remove(testPath)
	})
}

func TestBuildMarkdownTable(t *testing.T) {
	t.Run("Build table with script info", func(t *testing.T) {
		scripts := []ScriptInfo{
			{
				Name:        "test1.sh",
				Description: "First test script",
				Usage:       "./test1.sh [arg]",
			},
			{
				Name:        "test2.sh",
				Description: "Second test script",
				Usage:       "",
			},
		}

		result := BuildMarkdownTable(scripts)

		// Check that table contains expected elements
		if !strings.Contains(result, "Script Name") {
			t.Error("Table should contain 'Script Name' header")
		}

		if !strings.Contains(result, "Description") {
			t.Error("Table should contain 'Description' header")
		}

		if !strings.Contains(result, "Usage") {
			t.Error("Table should contain 'Usage' header")
		}

		if !strings.Contains(result, "`test1.sh`") {
			t.Error("Table should contain formatted script name")
		}

		if !strings.Contains(result, "First test script") {
			t.Error("Table should contain script description")
		}

		if !strings.Contains(result, "`./test1.sh [arg]`") {
			t.Error("Table should contain formatted usage")
		}

		// Check markdown table structure
		lines := strings.Split(result, "\n")
		if len(lines) < 4 {
			t.Errorf("Expected at least 4 lines (header, separator, 2 data rows), got %d", len(lines))
		}

		// Check separator line contains dashes
		if !strings.Contains(lines[1], "---") {
			t.Error("Second line should be table separator with dashes")
		}
	})

	t.Run("Handle empty description and usage", func(t *testing.T) {
		scripts := []ScriptInfo{
			{
				Name:        "empty.sh",
				Description: "",
				Usage:       "",
			},
		}

		result := BuildMarkdownTable(scripts)

		if !strings.Contains(result, "No description available") {
			t.Error("Should show placeholder for empty description")
		}

		if !strings.Contains(result, "-") {
			t.Error("Should show placeholder for empty usage")
		}
	})

	t.Run("Handle empty scripts list", func(t *testing.T) {
		var scripts []ScriptInfo
		result := BuildMarkdownTable(scripts)

		expected := "No scripts found.\n"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})
}
