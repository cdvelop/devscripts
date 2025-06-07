package devscripts

import (
	"os"
	"strings"
	"testing"
)

func TestGobadgeScriptCreateReadme(t *testing.T) {
	runner := NewScriptRunner()

	// Clean up any existing test files
	os.Remove("test_nonexistent_README.md")
	defer os.Remove("test_nonexistent_README.md")

	// Test when README.md doesn't exist
	exitCode, output, err := runner.ExecScript("gobadge.sh", "test_nonexistent_README.md")

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d. Output: %s, Error: %v", exitCode, output, err)
	}

	// Should show warning about missing README
	if !strings.Contains(output, "WARNING: test_nonexistent_README.md not found") {
		t.Errorf("Expected warning about missing README.md, got: %s", output)
	}
}

func TestGobadgeScriptWithExistingReadme(t *testing.T) {
	runner := NewScriptRunner()

	// Clean up any existing test files
	os.Remove("test_README.md")
	defer os.Remove("test_README.md")

	// Create README.md with title
	readmeContent := `# Test Project

This is a test project.
`
	err := os.WriteFile("test_README.md", []byte(readmeContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test README: %v", err)
	}

	// Test adding badges to existing README
	exitCode, output, err := runner.ExecScript("gobadge.sh", "test_README.md")

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d. Output: %s, Error: %v", exitCode, output, err)
	}

	// Check that badges were added
	if !strings.Contains(output, "Adding new badges") {
		t.Errorf("Expected 'Adding new badges' message, got: %s", output)
	}

	// Verify README.md content
	content, err := os.ReadFile("test_README.md")
	if err != nil {
		t.Fatalf("Failed to read updated README: %v", err)
	}

	contentStr := string(content)

	// Should contain the original title
	if !strings.Contains(contentStr, "# Test Project") {
		t.Errorf("Original title missing from README")
	}

	// Should contain the badge HTML
	if !strings.Contains(contentStr, "START_SECTION:BADGES_SECTION") {
		t.Errorf("Badge section markers missing from README")
	}

	// Should contain badge values
	if !strings.Contains(contentStr, "Go") {
		t.Errorf("Go badge missing from README")
	}

	// Should preserve original content
	if !strings.Contains(contentStr, "This is a test project") {
		t.Errorf("Original content missing from README")
	}
}

func TestGobadgeScriptUpdateExistingBadges(t *testing.T) {
	runner := NewScriptRunner()

	// Clean up any existing test files
	os.Remove("test_README.md")
	defer os.Remove("test_README.md")

	// Create README.md with existing badges
	readmeContent := `# Test Project

<!-- START_SECTION:BADGES_SECTION -->
<div align="center">
<img src="https://img.shields.io/badge/Go-1.21-blue" alt="Go version">
</div>
<!-- END_SECTION:BADGES_SECTION -->

This is a test project.
`
	err := os.WriteFile("test_README.md", []byte(readmeContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test README with existing badges: %v", err)
	}

	// Test updating existing badges
	exitCode, output, err := runner.ExecScript("gobadge.sh", "test_README.md")

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d. Output: %s, Error: %v", exitCode, output, err)
	}

	// Check that badges were updated
	if !strings.Contains(output, "Updating existing badges") {
		t.Errorf("Expected 'Updating existing badges' message, got: %s", output)
	}

	// Verify badges were updated
	content, err := os.ReadFile("test_README.md")
	if err != nil {
		t.Fatalf("Failed to read updated README: %v", err)
	}

	contentStr := string(content)
	// Should contain updated content
	if !strings.Contains(contentStr, "1.22") {
		t.Errorf("Badges were not updated with new Go version. Content: %s", contentStr)
	}
}

func TestGobadgeScriptNoChangeWhenDataSame(t *testing.T) {
	runner := NewScriptRunner()

	// Clean up any existing test files
	os.Remove("test_README.md")
	defer os.Remove("test_README.md")

	// First run to create badges
	readmeContent := `# Test Project

This is a test project.
`
	err := os.WriteFile("test_README.md", []byte(readmeContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test README: %v", err)
	}

	// First run
	_, _, err = runner.ExecScript("gobadge.sh", "test_README.md")
	if err != nil {
		t.Fatalf("First run failed: %v", err)
	}

	// Second run with same data - should not change
	exitCode, output, err := runner.ExecScript("gobadge.sh", "test_README.md")

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d. Output: %s, Error: %v", exitCode, output, err)
	}

	// Should indicate no changes needed
	if !strings.Contains(output, "Badges are already up to date") {
		t.Errorf("Expected 'Badges are already up to date' message, got: %s", output)
	}
}
