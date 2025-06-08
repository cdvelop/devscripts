package devscripts

import (
	"os"
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
		if !strings.Contains(output, "Generated badges SVG: "+testFile) {
			t.Errorf("Expected success message about SVG generation, got: %s", output)
		}

		// Verify badge count
		if !strings.Contains(output, "Valid badges generated: 3") {
			t.Errorf("Expected 3 valid badges generated, got: %s", output)
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

		// Verify single badge count
		if !strings.Contains(output, "Valid badges generated: 1") {
			t.Errorf("Expected 1 valid badge generated, got: %s", output)
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

		// Should still generate valid badges (2 out of 3)
		if !strings.Contains(output, "Valid badges generated: 2") {
			t.Errorf("Expected 2 valid badges generated, got: %s", output)
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

		// Should still generate valid badges (2 out of 3)
		if !strings.Contains(output, "Valid badges generated: 2") {
			t.Errorf("Expected 2 valid badges generated, got: %s", output)
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

	t.Run("GitHub directory creation", func(t *testing.T) {
		// This test specifically checks directory creation behavior
		// Always clean up the badges.svg file
		os.Remove(".github/badges.svg")
		defer os.Remove(".github/badges.svg")

		// Use a test README file to avoid modifying the main README.md
		testReadme := "test_github_readme.md"
		os.Remove(testReadme)
		defer os.Remove(testReadme)

		exitCode, output, err := runner.ExecScript("badges.sh", "readmefile:"+testReadme, "test:value:#ffffff")

		if exitCode != 0 {
			t.Fatalf("Expected exit code 0, got %d. Output: %s, Error: %v", exitCode, output, err)
		}

		// If .github directory didn't exist before, should show creation message
		// If it already existed, should still work without error
		if !strings.Contains(output, "Generated badges SVG") {
			t.Errorf("Expected success message about SVG generation, got: %s", output)
		}

		// Verify directory exists after execution
		if _, err := os.Stat(".github"); os.IsNotExist(err) {
			t.Error("Expected .github directory to exist after script execution")
		}

		// Verify badges.svg was created
		if _, err := os.Stat(".github/badges.svg"); os.IsNotExist(err) {
			t.Error("Expected .github/badges.svg to be created")
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
		if !strings.Contains(output, "Generated badges SVG: "+testFile) {
			t.Errorf("Expected 'Generated badges SVG' message for changed content, got: %s", output)
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
