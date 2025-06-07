package devscripts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Helper function to read file content
func readFileContent(t *testing.T, path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}
	return string(content)
}

func TestGobadgeScriptCreateReadme(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock files
	functionsScript := filepath.Join(tempDir, "functions.sh")
	functionsContent := `#!/bin/bash
info() {
  echo "INFO: $1"
}

warning() {
  echo "WARNING: $1"
}
`
	err := os.WriteFile(functionsScript, []byte(functionsContent), 0755)
	if err != nil {
		t.Fatalf("Failed to write functions.sh mock: %v", err)
	}

	gomodutilsScript := filepath.Join(tempDir, "gomodutils.sh")
	gomodutilsContent := `#!/bin/bash
get_go_version() {
    echo "1.22"
}
`
	err = os.WriteFile(gomodutilsScript, []byte(gomodutilsContent), 0755)
	if err != nil {
		t.Fatalf("Failed to write gomodutils.sh mock: %v", err)
	}

	licenseScript := filepath.Join(tempDir, "license.sh")
	licenseContent := `#!/bin/bash
get_license_type() {
    echo "MIT"
}
`
	err = os.WriteFile(licenseScript, []byte(licenseContent), 0755)
	if err != nil {
		t.Fatalf("Failed to write license.sh mock: %v", err)
	}

	// Copy gobadge.sh
	originalScript, err := os.ReadFile(filepath.Join(".", "gobadge.sh"))
	if err != nil {
		t.Fatalf("Failed to read original gobadge.sh: %v", err)
	}
	err = os.WriteFile(filepath.Join(tempDir, "gobadge.sh"), originalScript, 0755)
	if err != nil {
		t.Fatalf("Failed to copy gobadge.sh to temp dir: %v", err)
	}

	// Create a runner with the temp directory
	runner := NewScriptRunner(tempDir)

	// Test when README.md doesn't exist
	exitCode, output, err := runner.ExecScript("gobadge.sh", "testmodule", "Passing", "85", "Clean", "OK", "MIT")
	if err != nil {
		t.Errorf("Failed to execute gobadge.sh: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Output: %s", exitCode, output)
	}

	// Should show warning about missing README
	if !strings.Contains(output, "README.md not found") {
		t.Errorf("Expected warning about missing README.md, got: %s", output)
	}
}

func TestGobadgeScriptWithExistingReadme(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock files (same as above)
	functionsScript := filepath.Join(tempDir, "functions.sh")
	functionsContent := `#!/bin/bash
info() {
  echo "INFO: $1"
}

warning() {
  echo "WARNING: $1"
}
`
	err := os.WriteFile(functionsScript, []byte(functionsContent), 0755)
	if err != nil {
		t.Fatalf("Failed to write functions.sh mock: %v", err)
	}

	gomodutilsScript := filepath.Join(tempDir, "gomodutils.sh")
	gomodutilsContent := `#!/bin/bash
get_go_version() {
    echo "1.22"
}
`
	err = os.WriteFile(gomodutilsScript, []byte(gomodutilsContent), 0755)
	if err != nil {
		t.Fatalf("Failed to write gomodutils.sh mock: %v", err)
	}

	licenseScript := filepath.Join(tempDir, "license.sh")
	licenseContent := `#!/bin/bash
get_license_type() {
    echo "MIT"
}
`
	err = os.WriteFile(licenseScript, []byte(licenseContent), 0755)
	if err != nil {
		t.Fatalf("Failed to write license.sh mock: %v", err)
	}

	// Copy gobadge.sh
	originalScript, err := os.ReadFile(filepath.Join(".", "gobadge.sh"))
	if err != nil {
		t.Fatalf("Failed to read original gobadge.sh: %v", err)
	}
	err = os.WriteFile(filepath.Join(tempDir, "gobadge.sh"), originalScript, 0755)
	if err != nil {
		t.Fatalf("Failed to copy gobadge.sh to temp dir: %v", err)
	}

	// Create README.md with title
	readmeContent := `# MyTestModule

This is a test module for demonstration purposes.

## Features

- Feature 1
- Feature 2
`
	readmePath := filepath.Join(tempDir, "README.md")
	err = os.WriteFile(readmePath, []byte(readmeContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create README.md: %v", err)
	}

	// Create a runner with the temp directory
	runner := NewScriptRunner(tempDir)

	// Test adding badges to existing README
	exitCode, output, err := runner.ExecScript("gobadge.sh", "testmodule", "Passing", "100", "Clean", "OK", "MIT")
	if err != nil {
		t.Errorf("Failed to execute gobadge.sh: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Output: %s", exitCode, output)
	}

	// Check that badges were added
	if !strings.Contains(output, "Adding new badges") {
		t.Errorf("Expected 'Adding new badges' message, got: %s", output)
	}

	// Verify README.md content
	updatedContent := readFileContent(t, readmePath)

	// Should contain the original title
	if !strings.Contains(updatedContent, "# MyTestModule") {
		t.Errorf("Original title should be preserved in README")
	}
	// Should contain the badge HTML
	if !strings.Contains(updatedContent, "Generated dynamically by gotest.sh") {
		t.Errorf("Badge comment should be present in README")
	}

	if !strings.Contains(updatedContent, "<div class=\"project-badges\">") {
		t.Errorf("Badge div should be present in README")
	}

	// Should contain badge values
	if !strings.Contains(updatedContent, "Passing") {
		t.Errorf("Test status should be in README")
	}
	if !strings.Contains(updatedContent, "100%") {
		t.Errorf("Coverage should be in README")
	}
	if !strings.Contains(updatedContent, "MIT") {
		t.Errorf("License should be in README")
	}

	// Should preserve original content
	if !strings.Contains(updatedContent, "This is a test module") {
		t.Errorf("Original content should be preserved")
	}
}

func TestGobadgeScriptUpdateExistingBadges(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock files (same as above test)
	functionsScript := filepath.Join(tempDir, "functions.sh")
	functionsContent := `#!/bin/bash
info() {
  echo "INFO: $1"
}

warning() {
  echo "WARNING: $1"
}
`
	err := os.WriteFile(functionsScript, []byte(functionsContent), 0755)
	if err != nil {
		t.Fatalf("Failed to write functions.sh mock: %v", err)
	}

	gomodutilsScript := filepath.Join(tempDir, "gomodutils.sh")
	gomodutilsContent := `#!/bin/bash
get_go_version() {
    echo "1.22"
}
`
	err = os.WriteFile(gomodutilsScript, []byte(gomodutilsContent), 0755)
	if err != nil {
		t.Fatalf("Failed to write gomodutils.sh mock: %v", err)
	}

	licenseScript := filepath.Join(tempDir, "license.sh")
	licenseContent := `#!/bin/bash
get_license_type() {
    echo "MIT"
}
`
	err = os.WriteFile(licenseScript, []byte(licenseContent), 0755)
	if err != nil {
		t.Fatalf("Failed to write license.sh mock: %v", err)
	}

	// Copy gobadge.sh
	originalScript, err := os.ReadFile(filepath.Join(".", "gobadge.sh"))
	if err != nil {
		t.Fatalf("Failed to read original gobadge.sh: %v", err)
	}
	err = os.WriteFile(filepath.Join(tempDir, "gobadge.sh"), originalScript, 0755)
	if err != nil {
		t.Fatalf("Failed to copy gobadge.sh to temp dir: %v", err)
	}
	// Create README.md with existing badges
	readmeWithBadges := `# MyTestModule
<!-- Generated dynamically by gotest.sh from github.com/cdvelop/devscripts -->
<link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/cdvelop/devscripts@main/badges.css">
<div class="project-badges">
    <div class="badge-group">
        <span class="badge-label">Tests</span><span class="badge-value tests-failing">Failed</span>
    </div>
    <div class="badge-group">
        <span class="badge-label">Coverage</span><span class="badge-value coverage-none">0%</span>
    </div>
</div>

This is a test module with old badges.
`
	readmePath := filepath.Join(tempDir, "README.md")
	err = os.WriteFile(readmePath, []byte(readmeWithBadges), 0644)
	if err != nil {
		t.Fatalf("Failed to create README.md with badges: %v", err)
	}

	// Create a runner with the temp directory
	runner := NewScriptRunner(tempDir)

	// Test updating existing badges
	exitCode, output, err := runner.ExecScript("gobadge.sh", "testmodule", "Passing", "95", "Clean", "OK", "Apache")
	if err != nil {
		t.Errorf("Failed to execute gobadge.sh: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Output: %s", exitCode, output)
	}

	// Check that badges were updated
	if !strings.Contains(output, "Updating existing badges") {
		t.Errorf("Expected 'Updating existing badges' message, got: %s", output)
	}

	// Verify README.md content was updated
	updatedContent := readFileContent(t, readmePath)

	// Should contain updated values
	if !strings.Contains(updatedContent, "Passing") {
		t.Errorf("Updated test status should be in README")
	}
	if !strings.Contains(updatedContent, "95%") {
		t.Errorf("Updated coverage should be in README")
	}
	if !strings.Contains(updatedContent, "Apache") {
		t.Errorf("Updated license should be in README")
	} // Should NOT contain old values
	if strings.Contains(updatedContent, "Failed") {
		t.Errorf("Old test status should be replaced")
	}
	if strings.Contains(updatedContent, "0%") {
		t.Errorf("Old coverage should be replaced")
	}

	// Should preserve other content
	if !strings.Contains(updatedContent, "This is a test module with old badges") {
		t.Errorf("Original content should be preserved")
	}
}
