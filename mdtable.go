package devscripts

import (
	"strings"
)

// MdTable handles creation of markdown tables
type MdTable struct {
	headers           []string
	rows              [][]string
	minColumnWidths   map[int]int
	maxColumnWidths   map[int]int
	columnFormatters  map[int]func(string) string
	emptyPlaceholders map[int]string
}

// NewMdTable creates a new MdTable with headers
func NewMdTable(headers []string) *MdTable {
	return &MdTable{
		headers:           headers,
		rows:              make([][]string, 0),
		minColumnWidths:   make(map[int]int),
		maxColumnWidths:   make(map[int]int),
		columnFormatters:  make(map[int]func(string) string),
		emptyPlaceholders: make(map[int]string),
	}
}

// SetMinColumnWidth sets minimum width for a column (0-based index)
func (mt *MdTable) SetMinColumnWidth(colIndex, width int) {
	mt.minColumnWidths[colIndex] = width
}

// SetMaxColumnWidth sets maximum width for a column (0-based index)
func (mt *MdTable) SetMaxColumnWidth(colIndex, width int) {
	mt.maxColumnWidths[colIndex] = width
}

// SetColumnFormatter sets a formatter function for a column (0-based index)
func (mt *MdTable) SetColumnFormatter(colIndex int, formatter func(string) string) {
	mt.columnFormatters[colIndex] = formatter
}

// SetEmptyPlaceholder sets placeholder text for empty cells in a column (0-based index)
func (mt *MdTable) SetEmptyPlaceholder(colIndex int, placeholder string) {
	mt.emptyPlaceholders[colIndex] = placeholder
}

// AddRow adds a row to the table
func (mt *MdTable) AddRow(row []string) {
	mt.rows = append(mt.rows, row)
}

// SetRows sets all rows at once
func (mt *MdTable) SetRows(rows [][]string) {
	mt.rows = rows
}

// Generate creates the markdown table string
func (mt *MdTable) Generate() string {
	if len(mt.headers) == 0 {
		return ""
	}

	// Calculate column widths
	columnWidths := mt.calculateColumnWidths()

	var sb strings.Builder

	// Create header
	sb.WriteString("| ")
	for i, header := range mt.headers {
		sb.WriteString(mt.padRight(header, columnWidths[i]))
		sb.WriteString(" | ")
	}
	sb.WriteString("\n")

	// Create separator
	sb.WriteString("| ")
	for i := range mt.headers {
		sb.WriteString(strings.Repeat("-", columnWidths[i]))
		sb.WriteString(" | ")
	}
	sb.WriteString("\n")

	// Create data rows
	for _, row := range mt.rows {
		sb.WriteString("| ")
		for i := 0; i < len(mt.headers); i++ {
			cellValue := ""
			if i < len(row) {
				cellValue = row[i]
			}

			// Apply empty placeholder if needed
			if cellValue == "" {
				if placeholder, exists := mt.emptyPlaceholders[i]; exists {
					cellValue = placeholder
				}
			}

			// Apply formatter if exists
			if formatter, exists := mt.columnFormatters[i]; exists {
				cellValue = formatter(cellValue)
			}

			sb.WriteString(mt.padRight(cellValue, columnWidths[i]))
			sb.WriteString(" | ")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// calculateColumnWidths determines the width for each column
func (mt *MdTable) calculateColumnWidths() []int {
	columnWidths := make([]int, len(mt.headers))

	// Initialize with header lengths
	for i, header := range mt.headers {
		columnWidths[i] = len(header)
	}

	// Check all rows for maximum content length
	for _, row := range mt.rows {
		for i := 0; i < len(mt.headers) && i < len(row); i++ {
			cellValue := row[i]

			// Apply empty placeholder if needed for length calculation
			if cellValue == "" {
				if placeholder, exists := mt.emptyPlaceholders[i]; exists {
					cellValue = placeholder
				}
			}

			// Apply formatter if exists for length calculation
			if formatter, exists := mt.columnFormatters[i]; exists {
				cellValue = formatter(cellValue)
			}

			if len(cellValue) > columnWidths[i] {
				columnWidths[i] = len(cellValue)
			}
		}
	}

	// Apply minimum width constraints
	for colIndex, minWidth := range mt.minColumnWidths {
		if colIndex < len(columnWidths) && columnWidths[colIndex] < minWidth {
			columnWidths[colIndex] = minWidth
		}
	}

	// Apply maximum width constraints
	for colIndex, maxWidth := range mt.maxColumnWidths {
		if colIndex < len(columnWidths) && columnWidths[colIndex] > maxWidth {
			columnWidths[colIndex] = maxWidth
		}
	}

	return columnWidths
}

// padRight pads a string to the given length by adding spaces on the right
func (mt *MdTable) padRight(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(" ", length-len(s))
}

// Helper formatters that can be used with SetColumnFormatter
func AddBackticks(s string) string {
	if s == "" {
		return s
	}
	return "`" + s + "`"
}

func TrimHashPrefix(s string) string {
	s = strings.TrimSpace(s)       // First trim all spaces
	s = strings.TrimPrefix(s, "#") // Then remove # prefix
	return strings.TrimSpace(s)    // Trim again in case there were spaces after #
}
