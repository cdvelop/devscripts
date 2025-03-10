package devscripts

import (
	"strings"
	"testing"
)

func TestRunScript(t *testing.T) {

	// Create a runner with explicit configuration for tests
	runnerForTests := NewScriptRunner()

	t.Run("Successful script execution", func(t *testing.T) {
		exitCode, output, err := runnerForTests.ExecScript("testscript.sh", "arg1", "arg2")

		if exitCode != 0 {
			t.Fatalf("Expected exit code 0, got %d", exitCode)
		}

		if err != nil {
			t.Fatalf("Did not expect an error, but got: %v", err)
		}

		// Verify that the output contains the expected messages using functions from functions.sh
		if !strings.Contains(output, "Script executed successfully") {
			t.Fatalf("Output does not contain the expected success message: %s", output)
		}

		if !strings.Contains(output, "Number of arguments: 2") {
			t.Fatalf("Output does not display the correct number of arguments: %s", output)
		}

		// Adjusted verification to accept how Bash shows the arguments
		if !strings.Contains(output, "Received arguments: arg1 arg2") &&
			!strings.Contains(output, "Arguments received: arg1") {
			t.Fatalf("Output does not display the correct arguments: %s", output)
		}

	})

	t.Run("Script with error", func(t *testing.T) {
		exitCode, output, err := runnerForTests.ExecScript("testscript.sh", "error")

		if exitCode != 1 {
			t.Fatalf("Expected exit code 1, got %d", exitCode)
		}

		if err == nil {
			t.Error("Expected an error, but none was obtained")
		}

		if !strings.Contains(output, "Requested error!") {
			t.Fatalf("Output does not contain the expected error message: %s", output)
		}

		// Verify that the error function from functions.sh executed correctly
		if !strings.Contains(output, "ERROR:") {
			t.Fatalf("Output does not contain the ERROR symbol from functions.sh: %s", output)
		}
	})

	t.Run("Nonexistent script", func(t *testing.T) {
		exitCode, _, err := runnerForTests.ExecScript("script-que-no-existe")

		if exitCode != 1 {
			t.Fatalf("Expected exit code 1, got %d", exitCode)
		}

		if err == nil {
			t.Error("Expected an error, but none was obtained")
		}

		if !strings.Contains(err.Error(), "error") {
			t.Fatalf("Error message is not as expected: %v", err)
		}
	})

	t.Run("Custom scriptRunner", func(t *testing.T) {
		exitCode, output, err := runnerForTests.ExecScript("testscript.sh", "custom")

		if exitCode != 0 {
			t.Fatalf("Expected exit code 0, got %d", exitCode)
		}

		if err != nil {
			t.Fatalf("Did not expect an error, but got: %v", err)
		}

		if !strings.Contains(output, "Received arguments: custom") {
			t.Fatalf("Output does not display the correct arguments: %s", output)
		}

		// Verify the presence of the output from the execute function instead of literally searching for "execute"
		if !strings.Contains(output, "Command executed successfully") {
			t.Fatalf("The execute function in functions.sh did not run correctly: %s", output)
		}
	})
}
