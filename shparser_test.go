package devscripts

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScriptParser(t *testing.T) {
	// Set up temporary directory for tests
	tmpDir := t.TempDir()

	// Create test scripts
	testScripts := map[string]string{
		"test1.sh": `#!/bin/bash
# Description: Test script with description
# Usage: ./test1.sh [arg1]
echo "Test script 1"`, "test2.sh": `#!/bin/bash
# Description: Another test script
# Usage:
echo "Test script 2"`,
		"empty.sh": ``, "nogit.sh": `#!/bin/bash
# Description: No git operations script
echo "No git operations"`,
		"gitrepo.sh": `#!/bin/bash
echo "Git repository operations"`,
	}

	for name, content := range testScripts {
		err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test script %s: %v", name, err)
		}
	}

	parser := NewScriptParser(tmpDir)

	t.Run("GetScriptNames", func(t *testing.T) {
		scripts, err := parser.GetScriptNames()
		if err != nil {
			t.Fatalf("GetScriptNames failed: %v", err)
		}

		expectedCount := 5
		if len(scripts) != expectedCount {
			t.Errorf("Expected %d scripts, got %d", expectedCount, len(scripts))
		}

		// Check that all test scripts are found
		scriptMap := make(map[string]bool)
		for _, script := range scripts {
			scriptMap[script] = true
		}

		for expectedScript := range testScripts {
			if !scriptMap[expectedScript] {
				t.Errorf("Expected script %s not found in results", expectedScript)
			}
		}
	})

	t.Run("ParseScripts", func(t *testing.T) {
		scripts, err := parser.ParseScripts()
		if err != nil {
			t.Fatalf("ParseScripts failed: %v", err)
		}

		if len(scripts) != 5 {
			t.Errorf("Expected 5 parsed scripts, got %d", len(scripts))
		}

		// Find specific scripts and check their content
		scriptMap := make(map[string]ScriptInfo)
		for _, script := range scripts {
			scriptMap[script.Name] = script
		}

		// Test script with full description and usage
		if script, found := scriptMap["test1.sh"]; found {
			if script.Description != "Test script with description" {
				t.Errorf("Expected description 'Test script with description', got '%s'", script.Description)
			}
			if script.Usage != "./test1.sh [arg1]" {
				t.Errorf("Expected usage './test1.sh [arg1]', got '%s'", script.Usage)
			}
		} else {
			t.Error("test1.sh not found in parsed scripts")
		}

		// Test script with description but no usage
		if script, found := scriptMap["test2.sh"]; found {
			if script.Description != "Another test script" {
				t.Errorf("Expected description 'Another test script', got '%s'", script.Description)
			}
			if script.Usage != "" {
				t.Errorf("Expected empty usage, got '%s'", script.Usage)
			}
		} else {
			t.Error("test2.sh not found in parsed scripts")
		}

		// Test empty script
		if script, found := scriptMap["empty.sh"]; found {
			if script.Description != "Empty script file" {
				t.Errorf("Expected description 'Empty script file', got '%s'", script.Description)
			}
		} else {
			t.Error("empty.sh not found in parsed scripts")
		}

		// Test auto-generated description for git script
		if script, found := scriptMap["gitrepo.sh"]; found {
			if script.Description != "Git operations, Repository management" {
				t.Errorf("Expected auto-generated description for git script, got '%s'", script.Description)
			}
		} else {
			t.Error("gitrepo.sh not found in parsed scripts")
		}
		// Test script with no keywords
		if script, found := scriptMap["nogit.sh"]; found {
			if script.Description != "No git operations script" {
				t.Errorf("Expected description 'No git operations script', got '%s'", script.Description)
			}
		} else {
			t.Error("nogit.sh not found in parsed scripts")
		}
	})
}

func TestGenerateAutoDescription(t *testing.T) {
	parser := NewScriptParser("/tmp") // Directory doesn't matter for this test

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "gitutils.sh",
			content:  "echo git operations",
			expected: "Git operations",
		},
		{
			name:     "repomanager.sh",
			content:  "echo repository stuff",
			expected: "Repository management",
		},
		{
			name:     "gitrepo.sh",
			content:  "echo both",
			expected: "Git operations, Repository management",
		},
		{
			name:     "setup.sh",
			content:  "echo setup",
			expected: "System setup/config",
		},
		{
			name:     "goupdate.sh",
			content:  "echo go update",
			expected: "Dependency updates, Go language utilities",
		},
		{
			name:     "random.sh",
			content:  "echo nothing special",
			expected: "Shell script utility",
		},
	}

	for _, test := range tests {
		result := parser.generateAutoDescription(test.name, test.content)
		if result != test.expected {
			t.Errorf("generateAutoDescription(%q, %q) = %q, expected %q",
				test.name, test.content, result, test.expected)
		}
	}
}
