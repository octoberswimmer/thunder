package components

import (
	"testing"

	"github.com/octoberswimmer/masc"
)

// TestBadgeBasic verifies that Badge returns non-nil for valid parameters.
func TestBadgeBasic(t *testing.T) {
	comp := Badge("New")
	if comp == nil {
		t.Error("Badge returned nil for valid parameters")
	}
}

// TestPillBasic verifies that Pill returns non-nil when onRemove is nil.
func TestPillBasic(t *testing.T) {
	comp := Pill("Tag", nil)
	if comp == nil {
		t.Error("Pill returned nil for valid parameters")
	}
}

// TestPillWithRemove verifies that Pill returns non-nil when onRemove is provided.
func TestPillWithRemove(t *testing.T) {
	comp := Pill("Tag", func(e *masc.Event) {})
	if comp == nil {
		t.Error("Pill returned nil when onRemove is provided")
	}
}
