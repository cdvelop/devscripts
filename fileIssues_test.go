package devscripts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Helper function to create a test issues.md file with specified content
func createTestDoingMd(t *testing.T, dir, content string) string {
	filePath := filepath.Join(dir, "issues.md")
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test issues.md: %v", err)
	}
	return filePath
}

// Helper function to read file content
func readFile(t *testing.T, path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}
	return string(content)
}

func TestGetCommitMessageFromDoingMd(t *testing.T) {
	tempDir := t.TempDir()

	// Create test wrapper script that calls get_commit_message_from_issue_md
	testScript := filepath.Join(tempDir, "test_get_msg.sh")
	scriptContent := `#!/bin/bash
source "$(dirname "$0")/fileIssues.sh"
get_commit_message_from_issue_md "$@"
`
	err := os.WriteFile(testScript, []byte(scriptContent), 0755)
	if err != nil {
		t.Fatalf("Failed to write test script: %v", err)
	}

	// Copy the actual fileIssues.sh to temp dir
	originalScript, err := os.ReadFile(filepath.Join(".", "fileIssues.sh"))
	if err != nil {
		t.Fatalf("Failed to read original fileIssues.sh: %v", err)
	}
	err = os.WriteFile(filepath.Join(tempDir, "fileIssues.sh"), originalScript, 0755)
	if err != nil {
		t.Fatalf("Failed to copy fileIssues.sh to temp dir: %v", err)
	}

	// Create sample issues.md with mixed content
	doingContent := `# My Tasks
[x] Completed task 1
[ ] Incomplete task
[x] Completed task 2
Some random notes
`
	createTestDoingMd(t, tempDir, doingContent)

	// Create a runner with the temp directory
	runner := NewScriptRunner(tempDir)

	// Test with no initial message
	exitCode, output, err := runner.ExecScript("test_get_msg.sh", "")
	if err != nil {
		t.Errorf("Failed to execute script: %v", err)
	}
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	expectedOutput := "Completed task 1, Completed task 2"
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected output to contain %q, got %q", expectedOutput, output)
	}

	// Test with initial message
	exitCode, output, err = runner.ExecScript("test_get_msg.sh", "Initial commit")
	if err != nil {
		t.Errorf("Failed to execute script: %v", err)
	}
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	expectedOutput = "Initial commit: Completed task 1, Completed task 2"
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected output to contain %q, got %q", expectedOutput, output)
	}

	// Test with empty issues.md
	createTestDoingMd(t, tempDir, "")
	exitCode, output, err = runner.ExecScript("test_get_msg.sh", "Just message")
	if err != nil {
		t.Errorf("Failed to execute script: %v", err)
	}
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	expectedOutput = "Just message"
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected output to be %q, got %q", expectedOutput, output)
	}
}

func TestCreateDoingMdFile(t *testing.T) {
	tempDir := t.TempDir()

	// Create a functions.sh mock that implements the execute function
	functionsScript := filepath.Join(tempDir, "functions.sh")
	functionsContent := `#!/bin/bash
execute() {
  echo "Executing: $1"
  eval "$1"
  return $?
}
`
	err := os.WriteFile(functionsScript, []byte(functionsContent), 0755)
	if err != nil {
		t.Fatalf("Failed to write functions.sh mock: %v", err)
	}

	// Copy the actual fileIssues.sh to temp dir
	originalScript, err := os.ReadFile(filepath.Join(".", "fileIssues.sh"))
	if err != nil {
		t.Fatalf("Failed to read original fileIssues.sh: %v", err)
	}
	err = os.WriteFile(filepath.Join(tempDir, "fileIssues.sh"), originalScript, 0755)
	if err != nil {
		t.Fatalf("Failed to copy fileIssues.sh to temp dir: %v", err)
	}

	// Create test wrapper script
	testScript := filepath.Join(tempDir, "test_create.sh")
	scriptContent := `#!/bin/bash
source "$(dirname "$0")/functions.sh"
source "$(dirname "$0")/fileIssues.sh"
create_issue_md_file
cat issues.md
`
	err = os.WriteFile(testScript, []byte(scriptContent), 0755)
	if err != nil {
		t.Fatalf("Failed to write test script: %v", err)
	}

	// Create a runner with the temp directory
	runner := NewScriptRunner(tempDir)

	// Test creating the issues.md file
	exitCode, output, err := runner.ExecScript("test_create.sh")
	if err != nil {
		t.Errorf("Failed to execute script: %v", err)
	}
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	// Check file was created with expected content
	if !strings.Contains(output, "[x] init code") {
		t.Errorf("Expected output to contain init code, got %q", output)
	}
	if !strings.Contains(output, "[ ] task 1") {
		t.Errorf("Expected output to contain task 1, got %q", output)
	}

	// Verify file was actually created in the temp directory
	doingPath := filepath.Join(tempDir, "issues.md")
	if _, err := os.Stat(doingPath); os.IsNotExist(err) {
		t.Errorf("issues.md file was not created in %s", tempDir)
	}
}

func TestDeleteChangesDoingFile(t *testing.T) {
	tempDir := t.TempDir()

	// Copy the actual fileIssues.sh to temp dir
	originalScript, err := os.ReadFile(filepath.Join(".", "fileIssues.sh"))
	if err != nil {
		t.Fatalf("Failed to read original fileIssues.sh: %v", err)
	}
	err = os.WriteFile(filepath.Join(tempDir, "fileIssues.sh"), originalScript, 0755)
	if err != nil {
		t.Fatalf("Failed to copy fileIssues.sh to temp dir: %v", err)
	}

	// Create sample issues.md with mixed content
	doingContent := `# My Tasks
[x] Completed task 1
[ ] Incomplete task 1
[x] Completed task 2
[ ] Incomplete task 2
Some random notes
`
	doingPath := createTestDoingMd(t, tempDir, doingContent)

	// Create test wrapper script
	testScript := filepath.Join(tempDir, "test_delete.sh")
	scriptContent := `#!/bin/bash
source "$(dirname "$0")/fileIssues.sh"
deleteChangesIssueFile
cat issues.md
`
	err = os.WriteFile(testScript, []byte(scriptContent), 0755)
	if err != nil {
		t.Fatalf("Failed to write test script: %v", err)
	}

	// Create a runner with the temp directory
	runner := NewScriptRunner(tempDir)

	// Test deleting completed tasks
	exitCode, output, err := runner.ExecScript("test_delete.sh")
	if err != nil {
		t.Errorf("Failed to execute script: %v", err)
	}
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	// Check that completed tasks were removed
	if strings.Contains(output, "[x] Completed task 1") {
		t.Errorf("Completed task 1 should have been removed, but was found in: %q", output)
	}
	if strings.Contains(output, "[x] Completed task 2") {
		t.Errorf("Completed task 2 should have been removed, but was found in: %q", output)
	}

	// Check that incomplete tasks and other content were kept
	if !strings.Contains(output, "[ ] Incomplete task 1") {
		t.Errorf("Incomplete task 1 should have been kept, but was not found in: %q", output)
	}
	if !strings.Contains(output, "[ ] Incomplete task 2") {
		t.Errorf("Incomplete task 2 should have been kept, but was not found in: %q", output)
	}
	if !strings.Contains(output, "# My Tasks") {
		t.Errorf("Header should have been kept, but was not found in: %q", output)
	}
	if !strings.Contains(output, "Some random notes") {
		t.Errorf("Random notes should have been kept, but were not found in: %q", output)
	}

	// Also verify the actual file content
	fileContent := readFile(t, doingPath)
	if strings.Contains(fileContent, "[x] Completed task 1") || strings.Contains(fileContent, "[x] Completed task 2") {
		t.Errorf("Completed tasks should have been removed from file, but were found in: %q", fileContent)
	}
}
