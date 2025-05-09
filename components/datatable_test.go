package components

import (
	"testing"
)

// TestDataTableEmptyHeaders verifies that DataTable returns nil when no headers are provided.
func TestDataTableEmptyHeaders(t *testing.T) {
	if comp := DataTable([]string{}, []map[string]string{}); comp != nil {
		t.Errorf("DataTable expected nil for empty headers, got %v", comp)
	}
}

// TestDataTableBasic verifies that DataTable returns a component for valid headers and rows.
func TestDataTableBasic(t *testing.T) {
	headers := []string{"Col"}
	rows := []map[string]string{{"Col": "Value"}}
	if comp := DataTable(headers, rows); comp == nil {
		t.Error("DataTable returned nil for valid headers and rows")
	}
}
