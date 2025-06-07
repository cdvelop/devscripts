package devscripts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGotestScript(t *testing.T) {
	tempDir := t.TempDir()

	// Create a mock functions.sh
	functionsScript := filepath.Join(tempDir, "functions.sh")
	functionsContent := `#!/bin/bash
execute() {
  echo "Executing: $1"
  eval "$1"
  return $?
}

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

	// Create a mock gomodutils.sh
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

	// Create a mock license.sh
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
	// Create a mock gobadge.sh
	gobadgeScript := filepath.Join(tempDir, "gobadge.sh")
	gobadgeContent := `#!/bin/bash
echo "gobadge.sh called with: $@"
echo "INFO: README.md badges updated successfully"
exit 0
`
	err = os.WriteFile(gobadgeScript, []byte(gobadgeContent), 0755)
	if err != nil {
		t.Fatalf("Failed to write gobadge.sh mock: %v", err)
	}

	// Copy the actual gotest.sh to temp dir
	originalScript, err := os.ReadFile(filepath.Join(".", "gotest.sh"))
	if err != nil {
		t.Fatalf("Failed to read original gotest.sh: %v", err)
	}
	err = os.WriteFile(filepath.Join(tempDir, "gotest.sh"), originalScript, 0755)
	if err != nil {
		t.Fatalf("Failed to copy gotest.sh to temp dir: %v", err)
	}

	// Create a sample go.mod file
	goModContent := `module github.com/example/testmodule

go 1.22

require (
    github.com/stretchr/testify v1.8.0
)
`
	err = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// Create a simple test file
	testFileContent := `package main

import "testing"

func TestExample(t *testing.T) {
    if 1+1 != 2 {
        t.Error("Math is broken")
    }
}
`
	err = os.WriteFile(filepath.Join(tempDir, "main_test.go"), []byte(testFileContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a simple main.go file
	mainFileContent := `package main

func main() {
    println("Hello, World!")
}
`
	err = os.WriteFile(filepath.Join(tempDir, "main.go"), []byte(mainFileContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create main.go: %v", err)
	}

	// Create a README.md file so gobadge.sh will be called
	readmeContent := `# Test Module
This is a test module.
`
	err = os.WriteFile(filepath.Join(tempDir, "README.md"), []byte(readmeContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create README.md: %v", err)
	}

	// Create a runner with the temp directory
	runner := NewScriptRunner(tempDir)

	// Test gotest.sh execution
	exitCode, output, err := runner.ExecScript("gotest.sh")
	if err != nil {
		t.Errorf("Failed to execute gotest.sh: %v", err)
	}

	// Check that the script ran successfully
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Output: %s", exitCode, output)
	}
	// Check that badges were updated successfully
	if !strings.Contains(output, "README.md badges updated successfully") {
		t.Errorf("Expected badges to be updated successfully, but they weren't. Output: %s", output)
	}

	// Check that the output contains expected test results
	if !strings.Contains(output, "Running go vet") {
		t.Errorf("Expected 'Running go vet' in output, got: %s", output)
	}
	if !strings.Contains(output, "Running tests") {
		t.Errorf("Expected 'Running tests' in output, got: %s", output)
	}
	if !strings.Contains(output, "Running race detection") {
		t.Errorf("Expected 'Running race detection' in output, got: %s", output)
	}
	if !strings.Contains(output, "Calculating test coverage") {
		t.Errorf("Expected 'Calculating test coverage' in output, got: %s", output)
	}
}

func TestGotestScriptWithoutTestFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Create necessary mock files (similar to above test)
	functionsScript := filepath.Join(tempDir, "functions.sh")
	functionsContent := `#!/bin/bash
execute() {
  echo "Executing: $1"
  eval "$1"
  return $?
}

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
	gobadgeScript := filepath.Join(tempDir, "gobadge.sh")
	gobadgeContent := `#!/bin/bash
echo "gobadge.sh called with: $@"
echo "INFO: README.md badges updated successfully"
exit 0
`
	err = os.WriteFile(gobadgeScript, []byte(gobadgeContent), 0755)
	if err != nil {
		t.Fatalf("Failed to write gobadge.sh mock: %v", err)
	}

	// Copy gotest.sh
	originalScript, err := os.ReadFile(filepath.Join(".", "gotest.sh"))
	if err != nil {
		t.Fatalf("Failed to read original gotest.sh: %v", err)
	}
	err = os.WriteFile(filepath.Join(tempDir, "gotest.sh"), originalScript, 0755)
	if err != nil {
		t.Fatalf("Failed to copy gotest.sh to temp dir: %v", err)
	}

	// Create go.mod but NO test files
	goModContent := `module github.com/example/testmodule

go 1.22
`
	err = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}
	// Create only main.go (no test files)
	mainFileContent := `package main

func main() {
    println("Hello, World!")
}
`
	err = os.WriteFile(filepath.Join(tempDir, "main.go"), []byte(mainFileContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create main.go: %v", err)
	}

	// Create a README.md file so gobadge.sh will be called
	readmeContent := `# Test Module
This is a test module.
`
	err = os.WriteFile(filepath.Join(tempDir, "README.md"), []byte(readmeContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create README.md: %v", err)
	}

	// Create a runner with the temp directory
	runner := NewScriptRunner(tempDir)
	// Test gotest.sh execution
	exitCode, output, err := runner.ExecScript("gotest.sh")
	if err != nil {
		t.Errorf("Failed to execute gotest.sh: %v", err)
	}

	// Should still succeed even without test files
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Output: %s", exitCode, output)
	}

	// Check that it mentions no test files
	if !strings.Contains(output, "No test files found") {
		t.Errorf("Expected 'No test files found' message, got: %s", output)
	}

	// Should still update badges successfully
	if !strings.Contains(output, "README.md badges updated successfully") {
		t.Errorf("Expected badges to be updated successfully, but they weren't. Output: %s", output)
	}
}
