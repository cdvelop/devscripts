package devscripts

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MdHandler struct {
	sectionID  string
	afterLine  string
	content    string
	readmeFile string
}

func NewMdHandler(args ...string) *MdHandler {
	h := &MdHandler{readmeFile: "README.md"}
	if len(args) > 0 {
		h.sectionID = args[0]
	}
	if len(args) > 1 {
		h.afterLine = args[1]
	}
	if len(args) > 2 {
		h.content = args[2]
	}
	if len(args) > 3 {
		h.readmeFile = args[3]
	}
	return h
}

// SectionUpdate updates or creates a section in README based on identifier
func (h *MdHandler) SectionUpdate() error {
	if h.sectionID == "" || h.content == "" {
		return fmt.Errorf("section_identifier and new_content are required")
	}

	// Handle special case for BADGES -> BADGES_SECTION for backward compatibility
	if h.sectionID == "BADGES" {
		h.sectionID = "BADGES_SECTION"
	}

	existing, err := os.ReadFile(h.readmeFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error reading file: %v", err)
	}

	currentContent := string(existing)
	newContent, changed, err := h.processContent(currentContent)
	if err != nil {
		return err
	}

	if !changed {
		fmt.Println("already up to date")
		return nil
	}

	err = os.WriteFile(h.readmeFile, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	fmt.Println("Updated section successfully")
	return nil
}

func (h *MdHandler) processContent(currentContent string) (string, bool, error) {
	sectionStart := fmt.Sprintf("<!-- START_SECTION:%s -->", h.sectionID)
	sectionEnd := fmt.Sprintf("<!-- END_SECTION:%s -->", h.sectionID)

	// Create new section content
	newSection := fmt.Sprintf("%s\n%s\n%s", sectionStart, h.content, sectionEnd)

	// If file is empty or doesn't exist, create new file with section
	if currentContent == "" {
		if h.readmeFile != "README.md" || !fileExists(h.readmeFile) {
			fmt.Printf("Creating new file %s with section\n", h.readmeFile)
		}
		return newSection + "\n", true, nil
	}

	// Find all duplicate sections and their positions
	sections, err := h.findAllSections(currentContent, sectionStart, sectionEnd)
	if err != nil {
		return "", false, err
	}

	// Determine insertion position
	insertPos := h.determineInsertPosition(currentContent, sections)

	// Check if content needs updating
	if len(sections) == 1 && sections[0].content == h.content {
		// Content is the same, check position
		if h.afterLine == "" || sections[0].startLine == insertPos {
			return currentContent, false, nil // No change needed
		}
	}

	// Remove all existing sections
	contentWithoutSections := h.removeAllSections(currentContent, sections)

	// Insert new section at determined position
	newContent := h.insertSectionAtPosition(contentWithoutSections, newSection, insertPos)

	// Report what was done
	if len(sections) > 1 {
		fmt.Printf("Consolidated %d duplicate sections into one\n", len(sections))
	} else if len(sections) == 1 {
		fmt.Println("Updated existing section")
	} else {
		fmt.Println("Added new section")
	}

	return newContent, true, nil
}

type sectionInfo struct {
	startLine int
	endLine   int
	content   string
}

func (h *MdHandler) findAllSections(content, startMarker, endMarker string) ([]sectionInfo, error) {
	lines := strings.Split(content, "\n")
	var sections []sectionInfo
	currentStart := -1

	for i, line := range lines {
		if strings.TrimSpace(line) == startMarker {
			currentStart = i
		} else if strings.TrimSpace(line) == endMarker && currentStart >= 0 {
			// Extract content between markers
			var sectionContent strings.Builder
			for j := currentStart + 1; j < i; j++ {
				if j > currentStart+1 {
					sectionContent.WriteString("\n")
				}
				sectionContent.WriteString(lines[j])
			}

			sections = append(sections, sectionInfo{
				startLine: currentStart,
				endLine:   i,
				content:   sectionContent.String(),
			})
			currentStart = -1
		}
	}

	return sections, nil
}

func (h *MdHandler) determineInsertPosition(content string, sections []sectionInfo) int {
	lines := strings.Split(content, "\n")

	if h.afterLine != "" {
		// Parse after_line parameter (1-based from user input)
		if pos, err := strconv.Atoi(h.afterLine); err == nil {
			if pos >= 1 && pos <= len(lines) {
				return pos // Convert to 0-based: after line 1 means position 1
			}
		}
	}

	// Default behavior: if sections exist, use first section position
	if len(sections) > 0 {
		return sections[0].startLine
	}

	// No sections exist and no after_line specified, append at end
	return len(lines)
}

func (h *MdHandler) removeAllSections(content string, sections []sectionInfo) string {
	if len(sections) == 0 {
		return content
	}

	lines := strings.Split(content, "\n")

	// Sort sections by start line in descending order to remove from end
	for i := 0; i < len(sections); i++ {
		for j := i + 1; j < len(sections); j++ {
			if sections[i].startLine < sections[j].startLine {
				sections[i], sections[j] = sections[j], sections[i]
			}
		}
	}

	// Remove sections from end to beginning to maintain line indices
	for _, section := range sections {
		// Remove lines from endLine to startLine (inclusive)
		newLines := make([]string, 0, len(lines)-(section.endLine-section.startLine+1))
		newLines = append(newLines, lines[:section.startLine]...)
		if section.endLine+1 < len(lines) {
			newLines = append(newLines, lines[section.endLine+1:]...)
		}
		lines = newLines
	}

	return strings.Join(lines, "\n")
}

func (h *MdHandler) insertSectionAtPosition(content, section string, position int) string {
	lines := strings.Split(content, "\n")

	// Ensure position is within bounds
	if position > len(lines) {
		position = len(lines)
	}
	if position < 0 {
		position = 0
	}

	// Insert section at position
	newLines := make([]string, 0, len(lines)+3) // +3 for section lines
	newLines = append(newLines, lines[:position]...)
	newLines = append(newLines, section)
	if position < len(lines) {
		newLines = append(newLines, lines[position:]...)
	}

	return strings.Join(newLines, "\n")
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// MdUtils is the main function that follows the established pattern
func SectionUpdate(args ...string) {
	handler := NewMdHandler(args...)
	err := handler.SectionUpdate()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
