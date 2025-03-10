package devscripts

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type scriptInfo struct {
	Name        string
	Description string
	Usage       string
}

// getScriptDescriptions obtiene las descripciones de los scripts
func (sr *scriptRunner) getScriptDescriptions() ([]scriptInfo, error) {
	scripts, err := sr.GetScriptNames()
	if err != nil {
		return nil, err
	}

	var scriptsStruct []scriptInfo

	for _, script := range scripts {
		content, err := os.ReadFile(filepath.Join(sr.scriptsDir, script))
		if err != nil {
			return nil, err
		}

		if len(content) == 0 {
			scriptsStruct = append(scriptsStruct, scriptInfo{
				Name:        script,
				Description: "Empty script file",
				Usage:       "",
			})
			continue
		}

		lines := strings.Split(string(content), "\n")

		if len(lines) >= 3 {
			// Extract from line 2 (index 1)
			descLine := strings.TrimSpace(lines[1])
			descLine = strings.TrimPrefix(descLine, "# Description:")
			descLine = strings.TrimSpace(descLine)

			// Extract from line 3 (index 2)
			usageLine := strings.TrimSpace(lines[2])
			usageLine = strings.TrimPrefix(usageLine, "# Usage:")
			usageLine = strings.TrimSpace(usageLine)

			// Combine description and usage
			if descLine != "" {
				if usageLine != "" {
					scriptsStruct = append(scriptsStruct, scriptInfo{
						Name:        script,
						Description: descLine,
						Usage:       usageLine,
					})
				} else {
					scriptsStruct = append(scriptsStruct, scriptInfo{
						Name:        script,
						Description: descLine,
						Usage:       "",
					})
				}
			} else {
				scriptsStruct = append(scriptsStruct, scriptInfo{
					Name:        script,
					Description: sr.generateAutoDescription(script, strings.Join(lines, "\n")),
					Usage:       "",
				})
			}
		} else {
			scriptsStruct = append(scriptsStruct, scriptInfo{
				Name:        script,
				Description: sr.generateAutoDescription(script, strings.Join(lines, "\n")),
				Usage:       "",
			})
		}
	}

	return scriptsStruct, nil
}

// GetScriptNames obtiene los nombres de los scripts .sh en el directorio
func (sr *scriptRunner) GetScriptNames() ([]string, error) {
	files, err := os.ReadDir(sr.scriptsDir)
	if err != nil {
		return nil, err
	}

	var scripts []string
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".sh" {
			scripts = append(scripts, f.Name())
		}
	}

	return scripts, nil
}

func (sr *scriptRunner) generateAutoDescription(name, content string) string {
	keywords := map[string]string{
		"git":    "Git operations",
		"repo":   "Repository management",
		"setup":  "System setup/config",
		"update": "Dependency updates",
		"go":     "Go language utilities",
	}

	var desc []string
	for k, v := range keywords {
		if strings.Contains(strings.ToLower(name), k) {
			desc = append(desc, v)
		}
	}

	if len(desc) == 0 {
		return "Shell script utility"
	}
	return strings.Join(desc, ", ")
}

// formatDescription cleans and formats the description text
func formatDescription(desc string) string {
	// Remove leading # if present
	desc = strings.TrimPrefix(desc, "#")
	desc = strings.TrimSpace(desc)

	// If it's empty, return a placeholder
	if desc == "" {
		return "No description available"
	}
	return desc
}

// formatUsage cleans and formats the usage text
func formatUsage(usage string) string {
	// Remove leading # if present
	usage = strings.TrimPrefix(usage, "#")
	usage = strings.TrimSpace(usage)

	// If it's empty, return a placeholder
	if usage == "" {
		return "-"
	}
	return "`" + usage + "`"
}

// buildMarkdownTable creates a formatted markdown table from script info
func buildMarkdownTable(scripts []scriptInfo) string {
	// Find maximum lengths for each column to properly format the table
	maxNameLen := len("Script Name")
	maxDescLen := len("Description")
	maxUsageLen := len("Usage")

	for _, s := range scripts {
		nameLen := len(s.Name) + 2 // +2 for the backticks
		descLen := len(formatDescription(s.Description))
		usageLen := len(formatUsage(s.Usage))

		if nameLen > maxNameLen {
			maxNameLen = nameLen
		}
		if descLen > maxDescLen {
			maxDescLen = descLen
		}
		if usageLen > maxUsageLen {
			maxUsageLen = usageLen
		}
	}

	// Build the table with proper spacing
	var sb strings.Builder

	// Create header with proper column width
	sb.WriteString("| ")
	sb.WriteString(padRight("Script Name", maxNameLen))
	sb.WriteString(" | ")
	sb.WriteString(padRight("Description", maxDescLen))
	sb.WriteString(" | ")
	sb.WriteString(padRight("Usage", maxUsageLen))
	sb.WriteString(" |\n")

	// Create separator row with proper width
	sb.WriteString("| ")
	sb.WriteString(strings.Repeat("-", maxNameLen))
	sb.WriteString(" | ")
	sb.WriteString(strings.Repeat("-", maxDescLen))
	sb.WriteString(" | ")
	sb.WriteString(strings.Repeat("-", maxUsageLen))
	sb.WriteString(" |\n")

	// Add data rows
	for _, s := range scripts {
		formattedName := "`" + s.Name + "`"
		formattedDesc := formatDescription(s.Description)
		formattedUsage := formatUsage(s.Usage)

		sb.WriteString("| ")
		sb.WriteString(padRight(formattedName, maxNameLen))
		sb.WriteString(" | ")
		sb.WriteString(padRight(formattedDesc, maxDescLen))
		sb.WriteString(" | ")
		sb.WriteString(padRight(formattedUsage, maxUsageLen))
		sb.WriteString(" |\n")
	}

	return sb.String()
}

// padRight pads a string to the given length by adding spaces on the right
func padRight(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(" ", length-len(s))
}

// GenerateReadmeSection generates a markdown section for README
func (sr *scriptRunner) GenerateReadmeSection() (string, error) {
	scripts, err := sr.getScriptDescriptions()
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString("## Available Scripts\n\n")
	sb.WriteString(buildMarkdownTable(scripts))

	return sb.String(), nil
}

// UpdateReadmeIfNeeded updates the README file with the scripts section if there are changes
func (sr *scriptRunner) UpdateReadmeIfNeeded(readmePath string) (bool, error) {
	scriptsSection, err := sr.GenerateReadmeSection()
	if err != nil {
		return false, err
	}

	existing, err := os.ReadFile(readmePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return false, err
	}

	sectionStart := regexp.QuoteMeta("<!-- SCRIPTS_SECTION_START -->")
	sectionEnd := regexp.QuoteMeta("<!-- SCRIPTS_SECTION_END -->")
	pattern := regexp.MustCompile(`(?s)` + sectionStart + `\s*[\S\s]*?\s*` + sectionEnd)

	newSection := fmt.Sprintf("%s\n%s\n%s",
		"<!-- SCRIPTS_SECTION_START -->",
		scriptsSection,
		"<!-- SCRIPTS_SECTION_END -->")

	var newContent string
	currentContent := string(existing)

	if currentContent == "" {
		newContent = newSection + "\n"
	} else if pattern.MatchString(currentContent) {
		newContent = pattern.ReplaceAllString(currentContent, newSection)
	} else {
		newContent = strings.TrimSpace(currentContent) + "\n\n" + newSection
	}

	existingHash := md5.Sum(existing)
	newHash := md5.Sum([]byte(newContent))

	if existingHash != newHash {
		return true, os.WriteFile(readmePath, []byte(newContent), 0644)
	}

	return false, nil
}
