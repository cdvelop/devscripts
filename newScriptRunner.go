package devscripts

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// scriptRunner is a handler for executing different types of scripts
type scriptRunner struct {
	scriptsDir   string            // Base directory of scripts
	interpreters map[string]string // Map of file extensions to interpreter commands
}

// NewScriptRunner creates a handler for scripts using an optional scripts directory parameter.
// If no directory is provided, it uses the current working directory.
func NewScriptRunner(scriptsDir ...string) *scriptRunner {
	// Default value: scriptsDir is the current path
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}

	dir := wd
	if len(scriptsDir) > 0 && scriptsDir[0] != "" {
		dir = scriptsDir[0]
	}

	// Initialize interpreters map with default interpreters
	interpreters := map[string]string{
		".sh": "bash",
		".py": "python",
	}

	// Adjustment for Windows
	if runtime.GOOS == "windows" {
		interpreters[".sh"] = `C:\Program Files\Git\bin\bash.exe`
	}

	return &scriptRunner{
		scriptsDir:   dir,
		interpreters: interpreters,
	}
}

// ScriptChain represents a chain of scripts to be executed in sequence
type ScriptChain struct {
	runner       *scriptRunner
	scripts      []scriptExecution
	lastExitCode int
	lastOutput   string
	lastError    error
}

// scriptExecution represents a single script execution with its arguments
type scriptExecution struct {
	name string
	args []string
}

// Chain creates a new script execution chain
func (sr *scriptRunner) Chain() *ScriptChain {
	return &ScriptChain{
		runner:       sr,
		scripts:      make([]scriptExecution, 0),
		lastExitCode: 0,
		lastOutput:   "",
		lastError:    nil,
	}
}

// Then adds a script to the execution chain
func (sc *ScriptChain) Then(scriptName string, args ...string) *ScriptChain {
	sc.scripts = append(sc.scripts, scriptExecution{
		name: scriptName,
		args: args,
	})
	return sc
}

// Execute runs all scripts in the chain until one fails
func (sc *ScriptChain) Execute() (int, string, error) {
	var combinedOutput strings.Builder

	for _, script := range sc.scripts {
		exitCode, output, err := sc.runner.ExecScript(script.name, script.args...)
		combinedOutput.WriteString(output)

		sc.lastExitCode = exitCode
		sc.lastOutput = output
		sc.lastError = err

		if err != nil || exitCode != 0 {
			// Stop execution if a script fails
			return exitCode, combinedOutput.String(), err
		}
	}

	return 0, combinedOutput.String(), nil
}

// ExitCode returns the exit code of the last executed script
func (sc *ScriptChain) ExitCode() int {
	return sc.lastExitCode
}

// Output returns the output of the last executed script
func (sc *ScriptChain) Output() string {
	return sc.lastOutput
}

// Error returns the error of the last executed script
func (sc *ScriptChain) Error() error {
	return sc.lastError
}

// ExecScript executes a script and returns the exit code, output, and any error
func (sr *scriptRunner) ExecScript(scriptName string, args ...string) (int, string, error) {

	// Path to the main script in the scripts directory
	scriptPath := filepath.Join(sr.scriptsDir, scriptName)

	// Check if the script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		// List files in the directory for debugging
		files, _ := os.ReadDir(sr.scriptsDir)
		fileNames := make([]string, 0, len(files))
		for _, file := range files {
			fileNames = append(fileNames, file.Name())
		}
		return 1, "", fmt.Errorf("error: script '%s' does not exist. Available files: %v", scriptName, fileNames)
	}

	// Get the file extension
	ext := filepath.Ext(scriptName)

	// Determine the interpreter based on the file extension
	interpreter, supported := sr.interpreters[ext]
	if !supported {
		return 1, "", fmt.Errorf("unsupported script type: %s support: %v", ext, sr.interpreters)
	}

	// Ensure the script is executable
	if err := sr.makeScriptsExecutable(scriptPath); err != nil {
		return 1, "", fmt.Errorf("error making script executable: %w", err)
	}

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// Execute script with Git Bash on Windows by converting paths to Unix format
		unixPath := strings.ReplaceAll(scriptPath, "\\", "/")
		// Combine script path and arguments into a single quoted command string
		// Properly escape arguments and use bash positional parameters
		escapedArgs := make([]string, len(args))
		for i, arg := range args {
			escapedArgs[i] = fmt.Sprintf("%q", arg)
		}
		fullCommand := fmt.Sprintf("%q \"$@\"", unixPath)
		cmdArgs := []string{"-c", fullCommand, "--"}
		cmdArgs = append(cmdArgs, args...)
		cmd = exec.Command(interpreter, cmdArgs...)
	} else {
		// On other operating systems, execute directly
		cmd = exec.Command(interpreter, append([]string{scriptPath}, args...)...)
	}

	// Set the working directory to the directory where the scripts are located
	cmd.Dir = sr.scriptsDir

	// Configure environment variables to ensure stability
	env := os.Environ()
	cmd.Env = append(env, "LANG=C")

	// Execute and capture output
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	// Determine the exit code and handle errors
	if err != nil {
		var exitCode int
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
		return exitCode, outputStr, fmt.Errorf("error executing script: %w", err)
	}

	return 0, outputStr, nil
}

// makeScriptsExecutable makes the specified script executable if needed
func (sr *scriptRunner) makeScriptsExecutable(scriptPath string) error {
	// On Windows it's not necessary to make scripts executable
	if runtime.GOOS == "windows" {
		return nil
	}

	if err := os.Chmod(scriptPath, 0755); err != nil {
		return fmt.Errorf("failed to make script %s executable: %w", scriptPath, err)
	}

	return nil
}
