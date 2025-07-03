package devscripts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestBadgesScript(t *testing.T) {
	runner := NewScriptRunner()

	t.Run("Successful badges generation", func(t *testing.T) {
		// Use custom output file for testing
		testFile := "test_badges.svg"
		defer os.Remove(testFile)

		// Use custom README file for testing
		testReadme := "test_success_readme.md"
		os.Remove(testReadme)
		defer os.Remove(testReadme)

		exitCode, output, err := runner.ExecScript("badges.sh",
			"output_svgfile:"+testFile,
			"readmefile:"+testReadme,
			"license:MIT:#007acc",
			"go:1.22:#00add8",
			"coverage:85%:#28a745")

		if exitCode != 0 {
			t.Fatalf("Expected exit code 0, got %d. Output: %s, Error: %v", exitCode, output, err)
		}
		// Verify success message
		if !strings.Contains(output, "Badges saved to "+testFile) {
			t.Errorf("Expected success message about badges saved, got: %s", output)
		}

		// Verify file was created
		if _, err := os.Stat(testFile); os.IsNotExist(err) {
			t.Errorf("Expected %s to be created", testFile)
		}

		// Verify SVG content
		content, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("Failed to read generated SVG: %v", err)
		}

		svgContent := string(content)

		// Should contain SVG header
		if !strings.Contains(svgContent, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>") {
			t.Error("SVG should contain XML header")
		}

		// Should contain badge labels
		if !strings.Contains(svgContent, "license") {
			t.Error("SVG should contain license badge")
		}

		if !strings.Contains(svgContent, "go") {
			t.Error("SVG should contain go badge")
		}

		if !strings.Contains(svgContent, "coverage") {
			t.Error("SVG should contain coverage badge")
		}

		// Should contain badge values
		if !strings.Contains(svgContent, "MIT") {
			t.Error("SVG should contain MIT value")
		}

		if !strings.Contains(svgContent, "1.22") {
			t.Error("SVG should contain Go version value")
		}

		if !strings.Contains(svgContent, "85%") {
			t.Error("SVG should contain coverage percentage")
		}
	})
	t.Run("Single badge generation", func(t *testing.T) {
		// Use custom output file for testing
		testFile := "test_single_badge.svg"
		defer os.Remove(testFile)

		// Use custom README file for testing
		testReadme := "test_single_readme.md"
		os.Remove(testReadme)
		defer os.Remove(testReadme)

		exitCode, output, err := runner.ExecScript("badges.sh",
			"output_svgfile:"+testFile,
			"readmefile:"+testReadme,
			"license:MIT:#007acc")
		if exitCode != 0 {
			t.Fatalf("Expected exit code 0, got %d. Output: %s, Error: %v", exitCode, output, err)
		}

		// Verify file was created
		if _, err := os.Stat(testFile); os.IsNotExist(err) {
			t.Errorf("Expected %s to be created", testFile)
		}
	})

	t.Run("Error handling - no parameters", func(t *testing.T) {
		exitCode, output, err := runner.ExecScript("badges.sh")

		if exitCode != 1 {
			t.Fatalf("Expected exit code 1, got %d", exitCode)
		}

		if err == nil {
			t.Error("Expected an error, but none was obtained")
		}

		// Should show usage error
		if !strings.Contains(output, "No badges specified") {
			t.Errorf("Expected 'No badges specified' error, got: %s", output)
		}

		if !strings.Contains(output, "Usage: badges.sh") {
			t.Errorf("Expected usage message, got: %s", output)
		}
	})
	t.Run("Error handling - invalid format", func(t *testing.T) {
		// Use custom output file for testing
		testFile := "test_invalid_format.svg"
		defer os.Remove(testFile)

		// Use custom README file for testing
		testReadme := "test_invalid_readme.md"
		os.Remove(testReadme)
		defer os.Remove(testReadme)

		exitCode, output, err := runner.ExecScript("badges.sh",
			"output_svgfile:"+testFile,
			"readmefile:"+testReadme,
			"license:MIT:#007acc",
			"invalid-format",
			"coverage:85%:#28a745")

		if exitCode != 0 {
			t.Fatalf("Expected exit code 0 (should continue with valid badges), got %d. Output: %s, Error: %v", exitCode, output, err)
		}
		// Should show error for invalid format but continue
		if !strings.Contains(output, "Invalid badge format: invalid-format") {
			t.Errorf("Expected error message for invalid format, got: %s", output)
		}

		// Verify file was still created
		if _, err := os.Stat(testFile); os.IsNotExist(err) {
			t.Error("Expected test file to be created despite one invalid badge")
		}
	})
	t.Run("Error handling - empty fields", func(t *testing.T) {
		// Use custom output file for testing
		testFile := "test_empty_fields.svg"
		defer os.Remove(testFile)

		// Use custom README file for testing
		testReadme := "test_empty_readme.md"
		os.Remove(testReadme)
		defer os.Remove(testReadme)

		exitCode, output, err := runner.ExecScript("badges.sh",
			"output_svgfile:"+testFile,
			"readmefile:"+testReadme,
			"license:MIT:#007acc",
			"::#ffffff",
			"coverage:85%:#28a745")

		if exitCode != 0 {
			t.Fatalf("Expected exit code 0 (should continue with valid badges), got %d. Output: %s, Error: %v", exitCode, output, err)
		}
		// Should show error for empty fields
		if !strings.Contains(output, "Empty fields in badge") {
			t.Errorf("Expected error message for empty fields, got: %s", output)
		}
	})

	t.Run("Error handling - all invalid badges", func(t *testing.T) {
		exitCode, output, err := runner.ExecScript("badges.sh",
			"invalid1",
			"invalid2",
			"invalid3")

		if exitCode != 1 {
			t.Fatalf("Expected exit code 1, got %d", exitCode)
		}
		if err == nil {
			t.Error("Expected an error, but none was obtained")
		}
		// Should show error for no valid badges
		if !strings.Contains(output, "No valid badges to generate") {
			t.Errorf("Expected 'No valid badges to generate' error, got: %s", output)
		}
	})

	t.Run("Git directory requirement", func(t *testing.T) {
		// This test specifically checks that .git directory is required
		// Save current directory
		currentDir, _ := os.Getwd()

		// Create a temporary directory without .git
		tempDir, err := os.MkdirTemp("", "test_no_git")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)
		// Change to temp directory
		err = os.Chdir(tempDir)
		if err != nil {
			t.Fatalf("Failed to change to temp directory: %v", err)
		}
		defer os.Chdir(currentDir) // Copy necessary scripts to temp directory
		scriptsToChopy := []string{"badges.sh", "functions.sh", "sectionUpdate.sh"}
		for _, script := range scriptsToChopy {
			sourceFile := filepath.Join(currentDir, script)
			destFile := filepath.Join(tempDir, script)

			content, err := os.ReadFile(sourceFile)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", script, err)
			}

			err = os.WriteFile(destFile, content, 0755)
			if err != nil {
				t.Fatalf("Failed to write %s to temp dir: %v", script, err)
			}
		}

		// Create runner in temp directory
		tempRunner := NewScriptRunner(tempDir)

		// Use a test README file
		testReadme := "test_git_readme.md"
		defer os.Remove(testReadme)

		exitCode, output, err := tempRunner.ExecScript("badges.sh", "readmefile:"+testReadme, "test:value:#ffffff")

		if exitCode != 1 {
			t.Fatalf("Expected exit code 1 when .git directory doesn't exist, got %d. Output: %s", exitCode, output)
		}

		if err == nil {
			t.Error("Expected an error when .git directory doesn't exist")
		}

		// Should show error about git repository not found
		if !strings.Contains(output, "Git repository not found") {
			t.Errorf("Expected 'Git repository not found' error, got: %s", output)
		}
	})
	t.Run("Default directory - badges created in docs/img", func(t *testing.T) {
		// This test checks that badges are created in docs/img when .git directory exists
		// Always clean up the badges.svg file and directory
		os.RemoveAll("docs/img")
		defer os.RemoveAll("docs/img")

		// Use a test README file to avoid modifying the main README.md
		testReadme := "test_default_dir_readme.md"
		os.Remove(testReadme)
		defer os.Remove(testReadme)

		exitCode, output, err := runner.ExecScript("badges.sh", "readmefile:"+testReadme, "test:value:#ffffff")

		if exitCode != 0 {
			t.Fatalf("Expected exit code 0, got %d. Output: %s, Error: %v", exitCode, output, err)
		}

		// Should show success message
		if !strings.Contains(output, "Badges saved to docs/img/badges.svg") {
			t.Errorf("Expected success message about badges saved to docs/img/badges.svg, got: %s", output)
		}

		// Verify .git directory exists (should already exist in a git repo)
		if _, err := os.Stat(".git"); os.IsNotExist(err) {
			t.Error("Expected .git directory to exist")
		}

		// Verify docs/img directory was created
		if _, err := os.Stat("docs/img"); os.IsNotExist(err) {
			t.Error("Expected docs/img directory to be created")
		}

		// Verify badges.svg was created in docs/img directory
		if _, err := os.Stat("docs/img/badges.svg"); os.IsNotExist(err) {
			t.Error("Expected docs/img/badges.svg to be created")
		}
	})
	t.Run("SVG content unchanged - file not modified", func(t *testing.T) {
		// Use custom output file for testing
		testFile := "test_unchanged_content.svg"
		defer os.Remove(testFile)

		// Use custom README file for testing
		testReadme := "test_unchanged_readme.md"
		os.Remove(testReadme)
		defer os.Remove(testReadme)

		// First generation
		exitCode, output, err := runner.ExecScript("badges.sh",
			"output_svgfile:"+testFile,
			"readmefile:"+testReadme,
			"license:MIT:#007acc",
			"go:1.22:#00add8")

		if exitCode != 0 {
			t.Fatalf("First generation failed - exit code %d. Output: %s, Error: %v", exitCode, output, err)
		}

		// Verify file was created
		if _, err := os.Stat(testFile); os.IsNotExist(err) {
			t.Fatalf("Expected %s to be created", testFile)
		}

		// Get modification time
		fileInfo1, err := os.Stat(testFile)
		if err != nil {
			t.Fatalf("Failed to get file info: %v", err)
		}
		modTime1 := fileInfo1.ModTime()

		// Wait a bit to ensure different timestamps if file is modified
		time.Sleep(10 * time.Millisecond)
		// Second generation with same content
		exitCode, output, err = runner.ExecScript("badges.sh",
			"output_svgfile:"+testFile,
			"readmefile:"+testReadme,
			"license:MIT:#007acc",
			"go:1.22:#00add8")

		if exitCode != 0 {
			t.Fatalf("Second generation failed - exit code %d. Output: %s, Error: %v", exitCode, output, err)
		}

		// Verify the script detected unchanged content
		if !strings.Contains(output, "SVG content is already up to date") {
			t.Errorf("Expected 'SVG content is already up to date' message, got: %s", output)
		}

		// Verify modification time didn't change
		fileInfo2, err := os.Stat(testFile)
		if err != nil {
			t.Fatalf("Failed to get file info after second generation: %v", err)
		}
		modTime2 := fileInfo2.ModTime()

		if !modTime1.Equal(modTime2) {
			t.Errorf("File modification time changed when content was the same. Before: %v, After: %v", modTime1, modTime2)
		}
	})
	t.Run("SVG content changed - file modified", func(t *testing.T) {
		// Use custom output file for testing
		testFile := "test_changed_content.svg"
		defer os.Remove(testFile)

		// Use custom README file for testing
		testReadme := "test_changed_readme.md"
		os.Remove(testReadme)
		defer os.Remove(testReadme)

		// First generation
		exitCode, output, err := runner.ExecScript("badges.sh",
			"output_svgfile:"+testFile,
			"readmefile:"+testReadme,
			"license:MIT:#007acc")

		if exitCode != 0 {
			t.Fatalf("First generation failed - exit code %d. Output: %s, Error: %v", exitCode, output, err)
		}

		// Get modification time
		fileInfo1, err := os.Stat(testFile)
		if err != nil {
			t.Fatalf("Failed to get file info: %v", err)
		}
		modTime1 := fileInfo1.ModTime()

		// Wait a bit to ensure different timestamps
		time.Sleep(10 * time.Millisecond)
		// Second generation with different content
		exitCode, output, err = runner.ExecScript("badges.sh",
			"output_svgfile:"+testFile,
			"readmefile:"+testReadme,
			"license:MIT:#007acc",
			"go:1.22:#00add8")

		if exitCode != 0 {
			t.Fatalf("Second generation failed - exit code %d. Output: %s, Error: %v", exitCode, output, err)
		}
		// Verify file was updated
		if !strings.Contains(output, "Badges saved to "+testFile) {
			t.Errorf("Expected 'Badges saved to' message for changed content, got: %s", output)
		}

		// Verify modification time changed
		fileInfo2, err := os.Stat(testFile)
		if err != nil {
			t.Fatalf("Failed to get file info after second generation: %v", err)
		}
		modTime2 := fileInfo2.ModTime()

		if modTime1.Equal(modTime2) || modTime2.Before(modTime1) {
			t.Errorf("File modification time should have changed when content was different. Before: %v, After: %v", modTime1, modTime2)
		}
	})
}
