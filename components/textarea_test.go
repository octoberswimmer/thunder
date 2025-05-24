package components

import (
	"testing"
	"github.com/octoberswimmer/masc"
)

// TestTextareaBasic verifies that Textarea returns a component with basic parameters.
func TestTextareaBasic(t *testing.T) {
	if comp := Textarea("Test Label", "Test Value", "Test Placeholder", 0, nil); comp == nil {
		t.Error("Textarea returned nil for basic parameters")
	}
}

// TestTextareaWithRows verifies that Textarea handles custom rows.
func TestTextareaWithRows(t *testing.T) {
	if comp := Textarea("Test Label", "Test Value", "Test Placeholder", 5, nil); comp == nil {
		t.Error("Textarea returned nil with custom rows")
	}
}

// TestTextareaEmptyValues verifies that Textarea handles empty values.
func TestTextareaEmptyValues(t *testing.T) {
	if comp := Textarea("", "", "", 0, nil); comp == nil {
		t.Error("Textarea returned nil with empty values")
	}
}

// TestTextareaWithHandler verifies that Textarea works with event handler.
func TestTextareaWithHandler(t *testing.T) {
	handler := func(e *masc.Event) {}
	if comp := Textarea("Test Label", "Test Value", "", 0, handler); comp == nil {
		t.Error("Textarea returned nil with event handler")
	}
}