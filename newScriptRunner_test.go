package devscripts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewScriptRunner(t *testing.T) {
	// Test with default directory
	runner := NewScriptRunner()
	wd, _ := os.Getwd()
	if runner.scriptsDir != wd {
		t.Errorf("Expected scriptsDir to be %s, got %s", wd, runner.scriptsDir)
	}

	// Test with custom directory
	customDir := "/custom/dir"
	runner = NewScriptRunner(customDir)
	if runner.scriptsDir != customDir {
		t.Errorf("Expected scriptsDir to be %s, got %s", customDir, runner.scriptsDir)
	}

	// Verify interpreters are set correctly
	if _, exists := runner.interpreters[".sh"]; !exists {
		t.Error("Expected .sh interpreter to be set")
	}
	if _, exists := runner.interpreters[".py"]; !exists {
		t.Error("Expected .py interpreter to be set")
	}
}

func TestExecScript(t *testing.T) {
	// Use t.TempDir() which automatically cleans up when the test finishes
	tempDir := t.TempDir()

	// Create a test script
	testScript := filepath.Join(tempDir, "test.sh")
	scriptContent := `#!/bin/bash
echo "Hello $1!"
exit $2
`
	err := os.WriteFile(testScript, []byte(scriptContent), 0755)
	if err != nil {
		t.Fatalf("Failed to write test script: %v", err)
	}

	// Create a runner with the temp directory
	runner := NewScriptRunner(tempDir)

	// Test successful execution
	exitCode, output, err := runner.ExecScript("test.sh", "world", "0")
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
	expectedOutput := "Hello world!"
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected output to contain %q, got %q", expectedOutput, output)
	}
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test execution with non-zero exit code
	exitCode, output, err = runner.ExecScript("test.sh", "world", "1")
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected output to contain %q, got %q", expectedOutput, output)
	}
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Test execution with non-existent script
	exitCode, _, err = runner.ExecScript("nonexistent.sh")
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
	if err == nil {
		t.Error("Expected error for non-existent script, got nil")
	}
}

func TestScriptChain(t *testing.T) {
	// Use t.TempDir() which automatically cleans up when the test finishes
	tempDir := t.TempDir()

	// Create test scripts
	script1 := filepath.Join(tempDir, "script1.sh")
	err := os.WriteFile(script1, []byte(`#!/bin/bash
echo "Script 1 with arg: $1"
exit 0
`), 0755)
	if err != nil {
		t.Fatalf("Failed to create script1.sh: %v", err)
	}

	script2 := filepath.Join(tempDir, "script2.sh")
	err = os.WriteFile(script2, []byte(`#!/bin/bash
echo "Script 2 with arg: $1"
exit $2
`), 0755)
	if err != nil {
		t.Fatalf("Failed to create script2.sh: %v", err)
	}

	script3 := filepath.Join(tempDir, "script3.sh")
	err = os.WriteFile(script3, []byte(`#!/bin/bash
echo "Script 3 with arg: $1"
exit 0
`), 0755)
	if err != nil {
		t.Fatalf("Failed to create script3.sh: %v", err)
	}

	// Create a runner with the temp directory
	runner := NewScriptRunner(tempDir)

	// Test successful chain execution
	chain := runner.Chain().
		Then("script1.sh", "arg1").
		Then("script2.sh", "arg2", "0").
		Then("script3.sh", "arg3")

	exitCode, output, err := chain.Execute()
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check that all script outputs are present, using contains instead of exact match
	if !strings.Contains(output, "Script 1 with arg: arg1") {
		t.Errorf("Output missing script1 results: %q", output)
	}
	if !strings.Contains(output, "Script 2 with arg: arg2") {
		t.Errorf("Output missing script2 results: %q", output)
	}
	if !strings.Contains(output, "Script 3 with arg: arg3") {
		t.Errorf("Output missing script3 results: %q", output)
	}

	// Test chain breaking on error
	chain = runner.Chain().
		Then("script1.sh", "arg1").
		Then("script2.sh", "arg2", "1"). // This should fail
		Then("script3.sh", "arg3")       // This shouldn't execute

	exitCode, output, err = chain.Execute()
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Only the output from script1 and script2 should be present
	if !strings.Contains(output, "Script 1 with arg: arg1") {
		t.Errorf("Output missing script1 results: %q", output)
	}
	if !strings.Contains(output, "Script 2 with arg: arg2") {
		t.Errorf("Output missing script2 results: %q", output)
	}
	if strings.Contains(output, "Script 3 with arg: arg3") {
		t.Error("Script 3 should not have executed, but its output is present")
	}

	// Test access to last results
	if chain.ExitCode() != 1 {
		t.Errorf("Expected ExitCode() to return 1, got %d", chain.ExitCode())
	}
	if chain.Error() == nil {
		t.Error("Expected Error() to return error, got nil")
	}
	if !strings.Contains(chain.Output(), "Script 2 with arg: arg2") {
		t.Errorf("Expected output to contain script 2 output, got: %q", chain.Output())
	}
}
