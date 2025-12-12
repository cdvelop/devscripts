package devscripts

import (
	"os"
	"strings"

	"github.com/cdvelop/mdgo"
)

// DevScriptsReadmeUpdater handles updating README.md with scripts documentation
type DevScriptsReadmeUpdater struct {
	scriptsDir string
	parser     *ScriptParser
}

// NewDevScriptsReadmeUpdater creates a new DevScriptsReadmeUpdater
func NewDevScriptsReadmeUpdater(scriptsDir string) *DevScriptsReadmeUpdater {
	return &DevScriptsReadmeUpdater{
		scriptsDir: scriptsDir,
		parser:     NewScriptParser(scriptsDir),
	}
}

// GenerateScriptsSection generates a markdown section for README with scripts table
func (dru *DevScriptsReadmeUpdater) GenerateScriptsSection() (string, error) {
	scripts, err := dru.parser.ParseScripts()
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString("## Available Scripts\n")
	sb.WriteString("<small>This section is automatically generated.</small>\n\n")
	sb.WriteString(BuildMarkdownTable(scripts))

	return sb.String(), nil
}

// UpdateReadme updates the README file with the scripts section using sectionUpdate
func (dru *DevScriptsReadmeUpdater) UpdateReadme(readmePath string) error {
	// Generate the scripts section content
	scriptsSection, err := dru.GenerateScriptsSection()
	if err != nil {
		return err
	}

	// Use mdgo to handle the file update
	m := mdgo.New(".", ".", func(name string, data []byte) error {
		return os.WriteFile(name, data, 0644)
	})
	m.InputPath(readmePath, func(name string) ([]byte, error) {
		return os.ReadFile(name)
	})

	return m.UpdateSection("SCRIPTS_SECTION", scriptsSection)
}

// UpdateReadmeIfNeeded updates README and returns true if changes were made
func (dru *DevScriptsReadmeUpdater) UpdateReadmeIfNeeded(readmePath string) (bool, error) {
	// Generate the scripts section content
	scriptsSection, err := dru.GenerateScriptsSection()
	if err != nil {
		return false, err
	}

	// Read current file content
	var currentContent string
	if existing, err := os.ReadFile(readmePath); err == nil {
		currentContent = string(existing)
	}

	// Find existing section
	sectionStart := "<!-- START_SECTION:SCRIPTS_SECTION -->"
	sectionEnd := "<!-- END_SECTION:SCRIPTS_SECTION -->"

	if strings.Contains(currentContent, sectionStart) && strings.Contains(currentContent, sectionEnd) {
		// Extract current section content
		startIdx := strings.Index(currentContent, sectionStart)
		endIdx := strings.Index(currentContent, sectionEnd)
		if startIdx >= 0 && endIdx > startIdx {
			currentSectionContent := currentContent[startIdx+len(sectionStart) : endIdx]
			currentSectionContent = strings.TrimSpace(currentSectionContent)
			if currentSectionContent == strings.TrimSpace(scriptsSection) {
				return false, nil // No changes needed
			}
		}
	}

	// Update the file
	// Update the file using mdgo
	m := mdgo.New(".", ".", func(name string, data []byte) error {
		return os.WriteFile(name, data, 0644)
	})
	m.InputPath(readmePath, func(name string) ([]byte, error) {
		return os.ReadFile(name)
	})

	err = m.UpdateSection("SCRIPTS_SECTION", scriptsSection)
	return true, err
}

// BuildMarkdownTable creates a markdown table from script info using the MdTable API
// This function provides backward compatibility and a simple interface for common use cases
func BuildMarkdownTable(scripts []ScriptInfo) string {
	if len(scripts) == 0 {
		return "No scripts found.\n"
	}

	// Create table with headers (matching the expected test format)
	table := NewMdTable([]string{"Script Name", "Description", "Usage"})

	// Configure column formatting
	table.SetColumnFormatter(0, AddBackticks) // Add backticks to script names
	table.SetColumnFormatter(2, func(s string) string {
		if s == "" {
			return "-"
		}
		return "`" + s + "`"
	}) // Add backticks to usage
	table.SetEmptyPlaceholder(1, "No description available")

	// Set minimum column widths for better formatting
	table.SetMinColumnWidth(0, 12) // Script column
	table.SetMinColumnWidth(1, 20) // Description column
	table.SetMinColumnWidth(2, 8)  // Usage column

	// Add rows
	for _, script := range scripts {
		table.AddRow([]string{script.Name, script.Description, script.Usage})
	}

	return table.Generate()
}
