package devscripts

import (
	"os"
	"strings"
	"testing"
)

func TestSectionUpdateCreateNewReadme(t *testing.T) {
	runner := NewScriptRunner()

	// Clean up any existing test files
	os.Remove("test_README.md")
	defer os.Remove("test_README.md")
	// Test content for badges section
	badgeContent := `<div align="center">
<img src="https://img.shields.io/badge/Go-1.22-blue" alt="Go version">
</div>` // Test creating new README with badges section using sectionUpdate
	exitCode, output, err := runner.ExecScript("sectionUpdate.sh", "BADGES_SECTION", "", badgeContent, "test_README.md")

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

func TestSectionUpdateUpdateExistingSection(t *testing.T) {
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
	exitCode, output, err := runner.ExecScript("sectionUpdate.sh", "BADGES_SECTION", "", newBadgeContent, "test_README.md")

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

func TestSectionUpdateNoChangeWhenSame(t *testing.T) {
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
	exitCode, output, err := runner.ExecScript("sectionUpdate.sh", "BADGES_SECTION", "", badgeContent, "test_README.md")

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

func TestSectionUpdatePositionAfterLine(t *testing.T) {
	runner := NewScriptRunner()

	// Clean up any existing test files
	os.Remove("test_README.md")
	defer os.Remove("test_README.md")

	// Create README with title and content
	readmeContent := `# Test Project

This is a test project.

## Section 1
Content here.
`
	err := os.WriteFile("test_README.md", []byte(readmeContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test README: %v", err)
	}

	// Test badge content
	badgeContent := `<img src=".github/badges.svg" alt="Project Badges">` // Test creating badges after line 1 (after title)
	exitCode, output, err := runner.ExecScript("sectionUpdate.sh", "BADGES_SECTION", "1", badgeContent, "test_README.md")

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d. Output: %s, Error: %v", exitCode, output, err)
	}

	// Verify README was updated
	content, err := os.ReadFile("test_README.md")
	if err != nil {
		t.Fatalf("Failed to read updated README: %v", err)
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	// Should have badges section after title (line 1)
	if len(lines) < 3 {
		t.Fatalf("README should have at least 3 lines, got %d", len(lines))
	}
	// Line 0: # Test Project
	// Line 1: <!-- START_SECTION:BADGES_SECTION -->
	if !strings.Contains(lines[1], "START_SECTION:BADGES_SECTION") {
		t.Errorf("Badges section should start at line 2, but found: %s", lines[1])
	}

	// Should contain badge content
	if !strings.Contains(contentStr, "badges.svg") {
		t.Errorf("Badge content missing from README")
	}

	// Should preserve original content
	if !strings.Contains(contentStr, "This is a test project") {
		t.Errorf("Original content should be preserved")
	}
}

func TestSectionUpdateMoveSectionToNewPosition(t *testing.T) {
	runner := NewScriptRunner()

	// Clean up any existing test files
	os.Remove("test_README.md")
	defer os.Remove("test_README.md")

	// Create README with badges at the end
	readmeContent := `# Test Project

This is a test project.

## Section 1
Content here.

<!-- START_SECTION:BADGES_SECTION -->
<img src="old-badges.svg" alt="Old Badges">
<!-- END_SECTION:BADGES_SECTION -->
`
	err := os.WriteFile("test_README.md", []byte(readmeContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test README with badges at end: %v", err)
	}

	// New badge content
	newBadgeContent := `<img src=".github/badges.svg" alt="Project Badges">`
	// Test moving badges to after line 1 (after title)
	exitCode, output, err := runner.ExecScript("sectionUpdate.sh", "BADGES_SECTION", "1", newBadgeContent, "test_README.md")

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d. Output: %s, Error: %v", exitCode, output, err)
	}

	// Verify README was updated
	content, err := os.ReadFile("test_README.md")
	if err != nil {
		t.Fatalf("Failed to read updated README: %v", err)
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	// Should have badges section after title (line 1)
	if !strings.Contains(lines[1], "START_SECTION:BADGES_SECTION") {
		t.Errorf("Badges section should start at line 2, but found: %s", lines[1])
	}

	// Should contain new badge content
	if !strings.Contains(contentStr, ".github/badges.svg") {
		t.Errorf("New badge content missing from README")
	}

	// Should NOT contain old badge content
	if strings.Contains(contentStr, "old-badges.svg") {
		t.Errorf("Old badge content should be replaced")
	}

	// Should preserve other content
	if !strings.Contains(contentStr, "This is a test project") {
		t.Errorf("Original content should be preserved")
	}

	if !strings.Contains(contentStr, "## Section 1") {
		t.Errorf("Other sections should be preserved")
	} // Should not have badges at the end anymore
	lastLines := strings.Join(lines[len(lines)-5:], "\n")
	if strings.Contains(lastLines, "BADGES_SECTION") {
		t.Errorf("Badges section should not be at the end anymore")
	}
}

func TestSectionUpdateNoDuplicateSections(t *testing.T) {
	runner := NewScriptRunner()

	// Clean up any existing test files
	os.Remove("test_duplicate_README.md")
	defer os.Remove("test_duplicate_README.md")

	// Create README with duplicate badges sections (simulating the race condition)
	readmeWithDuplicates := `# Test Project

<!-- START_SECTION:BADGES_SECTION -->
<a href="docs/img/badges.svg"><img src="docs/img/badges.svg" alt="Project Badges" title="Generated by badges.sh from github.com/cdvelop/devscripts"></a>
<!-- END_SECTION:BADGES_SECTION -->
<!-- START_SECTION:BADGES_SECTION -->
<a href="docs/img/badges.svg"><img src="docs/img/badges.svg" alt="Project Badges" title="Generated by badges.sh from github.com/cdvelop/devscripts"></a>
<!-- END_SECTION:BADGES_SECTION -->

This is a test project with duplicate sections.
`
	err := os.WriteFile("test_duplicate_README.md", []byte(readmeWithDuplicates), 0644)
	if err != nil {
		t.Fatalf("Failed to create test README with duplicates: %v", err)
	}

	// New badge content
	newBadgeContent := `<a href="docs/img/badges.svg"><img src="docs/img/badges.svg" alt="Updated Project Badges" title="Generated by badges.sh from github.com/cdvelop/devscripts"></a>` // Test updating section when duplicates exist - should remove all duplicates and create only one
	exitCode, output, err := runner.ExecScript("sectionUpdate.sh", "BADGES_SECTION", "1", newBadgeContent, "test_duplicate_README.md")

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d. Output: %s, Error: %v", exitCode, output, err)
	} // Verify README was updated
	content, err := os.ReadFile("test_duplicate_README.md")
	if err != nil {
		t.Fatalf("Failed to read updated README: %v", err)
	}

	contentStr := string(content)

	// Count occurrences of BADGES_SECTION
	startSectionCount := strings.Count(contentStr, "START_SECTION:BADGES_SECTION")
	endSectionCount := strings.Count(contentStr, "END_SECTION:BADGES_SECTION")

	// Should have exactly one start and one end section
	if startSectionCount != 1 {
		t.Errorf("Expected exactly 1 START_SECTION:BADGES_SECTION, found %d", startSectionCount)
	}

	if endSectionCount != 1 {
		t.Errorf("Expected exactly 1 END_SECTION:BADGES_SECTION, found %d", endSectionCount)
	}

	// Should contain the updated content
	if !strings.Contains(contentStr, "Updated Project Badges") {
		t.Errorf("Should contain updated badge content")
	}

	// Should preserve other content
	if !strings.Contains(contentStr, "This is a test project with duplicate sections") {
		t.Errorf("Other content should be preserved")
	}

	// Verify section is positioned correctly (after line 1)
	lines := strings.Split(contentStr, "\n")
	if len(lines) < 3 {
		t.Fatalf("README should have at least 3 lines, got %d", len(lines))
	}

	if !strings.Contains(lines[1], "START_SECTION:BADGES_SECTION") {
		t.Errorf("Badges section should start at line 2, but found: %s", lines[1])
	}
}

func TestSectionUpdateDetectMultipleSectionsWithDifferentContent(t *testing.T) {
	runner := NewScriptRunner()

	// Clean up any existing test files
	os.Remove("test_multiple_sections_README.md")
	defer os.Remove("test_multiple_sections_README.md")

	// Create README with multiple sections but different content
	readmeWithMultipleSections := `# Test Project

<!-- START_SECTION:BADGES_SECTION -->
<a href="old-badges.svg"><img src="old-badges.svg" alt="Old Badges"></a>
<!-- END_SECTION:BADGES_SECTION -->

Some content here.

<!-- START_SECTION:BADGES_SECTION -->
<a href="other-badges.svg"><img src="other-badges.svg" alt="Other Badges"></a>
<!-- END_SECTION:BADGES_SECTION -->

More content.
`
	err := os.WriteFile("test_multiple_sections_README.md", []byte(readmeWithMultipleSections), 0644)
	if err != nil {
		t.Fatalf("Failed to create test README with multiple sections: %v", err)
	}

	// New badge content
	newBadgeContent := `<a href="docs/img/badges.svg"><img src="docs/img/badges.svg" alt="Consolidated Badges"></a>`
	// Test updating section - should consolidate all duplicate sections into one
	exitCode, output, err := runner.ExecScript("sectionUpdate.sh", "BADGES_SECTION", "1", newBadgeContent, "test_multiple_sections_README.md")

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d. Output: %s, Error: %v", exitCode, output, err)
	}

	// Verify README was updated
	content, err := os.ReadFile("test_multiple_sections_README.md")
	if err != nil {
		t.Fatalf("Failed to read updated README: %v", err)
	}

	contentStr := string(content)

	// Count occurrences of BADGES_SECTION
	startSectionCount := strings.Count(contentStr, "START_SECTION:BADGES_SECTION")
	endSectionCount := strings.Count(contentStr, "END_SECTION:BADGES_SECTION")

	// Should have exactly one start and one end section
	if startSectionCount != 1 {
		t.Errorf("Expected exactly 1 START_SECTION:BADGES_SECTION, found %d", startSectionCount)
	}

	if endSectionCount != 1 {
		t.Errorf("Expected exactly 1 END_SECTION:BADGES_SECTION, found %d", endSectionCount)
	}

	// Should contain only the new content, not the old ones
	if !strings.Contains(contentStr, "Consolidated Badges") {
		t.Errorf("Should contain new consolidated badge content")
	}

	if strings.Contains(contentStr, "old-badges.svg") {
		t.Errorf("Should not contain old badge content")
	}

	if strings.Contains(contentStr, "other-badges.svg") {
		t.Errorf("Should not contain other old badge content")
	}

	// Should preserve non-section content
	if !strings.Contains(contentStr, "Some content here") {
		t.Errorf("Should preserve content between sections")
	}

	if !strings.Contains(contentStr, "More content") {
		t.Errorf("Should preserve content after sections")
	}
}
