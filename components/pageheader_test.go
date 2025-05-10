package components

import (
	"github.com/octoberswimmer/masc"
	"testing"
)

// TestPageHeaderBasic verifies that PageHeader returns a non-nil component for a title only.
func TestPageHeaderBasic(t *testing.T) {
	comp := PageHeader("Title", "")
	if comp == nil {
		t.Error("PageHeader returned nil for title only")
	}
}

// TestPageHeaderWithSubtitle verifies that including a subtitle does not break rendering.
func TestPageHeaderWithSubtitle(t *testing.T) {
	comp := PageHeader("Title", "Subtitle")
	if comp == nil {
		t.Error("PageHeader returned nil when subtitle provided")
	}
}

// TestPageHeaderWithActions verifies that including action children does not break rendering.
func TestPageHeaderWithActions(t *testing.T) {
	// Use a simple text node as a placeholder for an action
	action := masc.Text("Action")
	comp := PageHeader("Title", "Subtitle", action)
	if comp == nil {
		t.Error("PageHeader returned nil when actions provided")
	}
}
