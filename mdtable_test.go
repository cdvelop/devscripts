package devscripts

import (
	"strings"
	"testing"
)

func TestMdTableConfiguration(t *testing.T) {
	table := NewMdTable([]string{"Name", "Description", "Usage"})

	// Test column formatters
	table.SetColumnFormatter(0, AddBackticks)
	table.SetColumnFormatter(2, func(s string) string {
		if s == "" {
			return "-"
		}
		return "`" + s + "`"
	})

	// Test empty placeholders
	table.SetEmptyPlaceholder(1, "No description available")

	// Test column width constraints
	table.SetMinColumnWidth(0, 10)
	table.SetMaxColumnWidth(1, 30)

	// Add test data
	table.AddRow([]string{"test.sh", "Test script", "./test.sh"})
	table.AddRow([]string{"empty.sh", "", ""})

	result := table.Generate()

	// Verify the result contains expected formatting
	if !strings.Contains(result, "`test.sh`") {
		t.Error("Expected script name to be wrapped in backticks")
	}
	if !strings.Contains(result, "No description available") {
		t.Error("Expected empty placeholder for description")
	}
	if !strings.Contains(result, "`./test.sh`") {
		t.Error("Expected usage to be wrapped in backticks")
	}
}

func TestHelperFormatters(t *testing.T) {
	// Test AddBackticks
	if AddBackticks("test") != "`test`" {
		t.Error("AddBackticks should wrap text in backticks")
	}
	if AddBackticks("") != "" {
		t.Error("AddBackticks should return empty string for empty input")
	}
	// Test TrimHashPrefix
	if TrimHashPrefix("# Test") != "Test" {
		t.Error("TrimHashPrefix should remove # prefix and trim spaces")
	}
	if TrimHashPrefix("  # Trimmed  ") != "Trimmed" {
		t.Errorf("TrimHashPrefix should handle extra spaces, got: '%s'", TrimHashPrefix("  # Trimmed  "))
	}
}
