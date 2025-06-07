package devscripts

import (
	"os"
	"strings"
	"testing"
)

func TestReadmeUtilsCreateNewReadme(t *testing.T) {
	runner := NewScriptRunner()

	// Clean up any existing test files
	os.Remove("test_README.md")
	defer os.Remove("test_README.md")

	// Test content for badges section
	badgeContent := `<div align="center">
<img src="https://img.shields.io/badge/Go-1.22-blue" alt="Go version">
</div>`

	// Test creating new README with badges section using section_update
	exitCode, output, err := runner.ExecScript("readmeutils.sh", "section_update", "BADGES_SECTION", badgeContent, "test_README.md")

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d. Output: %s, Error: %v", exitCode, output, err)
	}

	// Verify README was created
	if _, err := os.Stat("test_README.md"); os.IsNotExist(err) {
		t.Fatal("README.md should have been created")
	}

	content, err := os.ReadFile("test_README.md")
	if err != nil {
		t.Fatalf("Failed to read created README: %v", err)
	}

	contentStr := string(content)

	// Should contain badge content
	if !strings.Contains(contentStr, "Go-1.22") {
		t.Errorf("Badge content missing from README")
	}

	// Should contain section markers
	if !strings.Contains(contentStr, "START_SECTION:BADGES_SECTION") {
		t.Errorf("Section start marker missing")
	}

	if !strings.Contains(contentStr, "END_SECTION:BADGES_SECTION") {
		t.Errorf("Section end marker missing")
	}
}

func TestReadmeUtilsUpdateExistingSection(t *testing.T) {
	runner := NewScriptRunner()

	// Clean up any existing test files
	os.Remove("test_README.md")
	defer os.Remove("test_README.md")

	// Create existing README with old badges
	existingContent := `# Test Project

<!-- START_SECTION:BADGES_SECTION -->
<div align="center">
<img src="https://img.shields.io/badge/Go-1.21-blue" alt="Go version">
</div>
<!-- END_SECTION:BADGES_SECTION -->

Some other content here.
`
	err := os.WriteFile("test_README.md", []byte(existingContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create existing README: %v", err)
	}

	// New badge content
	newBadgeContent := `<div align="center">
<img src="https://img.shields.io/badge/Go-1.22-blue" alt="Go version">
<img src="https://img.shields.io/badge/Tests-Passing-green" alt="Tests">
</div>`

	// Test updating existing section
	exitCode, output, err := runner.ExecScript("readmeutils.sh", "section_update", "BADGES_SECTION", newBadgeContent, "test_README.md")

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d. Output: %s, Error: %v", exitCode, output, err)
	}

	// Verify README was updated
	content, err := os.ReadFile("test_README.md")
	if err != nil {
		t.Fatalf("Failed to read updated README: %v", err)
	}

	contentStr := string(content)

	// Should contain new content
	if !strings.Contains(contentStr, "Go-1.22") {
		t.Errorf("New badge content missing from README")
	}

	if !strings.Contains(contentStr, "Tests-Passing") {
		t.Errorf("New test badge missing from README")
	}

	// Should NOT contain old content
	if strings.Contains(contentStr, "Go-1.21") {
		t.Errorf("Old badge content should be replaced")
	}

	// Should preserve other content
	if !strings.Contains(contentStr, "# Test Project") {
		t.Errorf("Original title should be preserved")
	}

	if !strings.Contains(contentStr, "Some other content here") {
		t.Errorf("Other content should be preserved")
	}
}

func TestReadmeUtilsNoChangeWhenSame(t *testing.T) {
	runner := NewScriptRunner()

	// Clean up any existing test files
	os.Remove("test_README.md")
	defer os.Remove("test_README.md")

	// Create README with badges
	badgeContent := `<div align="center">
<img src="https://img.shields.io/badge/Go-1.22-blue" alt="Go version">
</div>`

	existingContent := `# Test Project

<!-- START_SECTION:BADGES_SECTION -->
` + badgeContent + `
<!-- END_SECTION:BADGES_SECTION -->

Some content.
`
	err := os.WriteFile("test_README.md", []byte(existingContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create README: %v", err)
	}

	// Get file info before
	fileInfoBefore, err := os.Stat("test_README.md")
	if err != nil {
		t.Fatalf("Failed to get file info before: %v", err)
	}

	// Test with same content - should not change
	exitCode, output, err := runner.ExecScript("readmeutils.sh", "section_update", "BADGES_SECTION", badgeContent, "test_README.md")

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d. Output: %s, Error: %v", exitCode, output, err)
	}
	// Should indicate no changes needed
	if !strings.Contains(output, "already up to date") {
		t.Errorf("Expected 'already up to date' message, got: %s", output)
	}

	// Verify file was NOT modified
	fileInfoAfter, err := os.Stat("test_README.md")
	if err != nil {
		t.Fatalf("Failed to get file info after: %v", err)
	}

	// File modification time should be the same
	if !fileInfoBefore.ModTime().Equal(fileInfoAfter.ModTime()) {
		t.Errorf("File should not have been modified when content is the same")
	}
}
