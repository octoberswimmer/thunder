package components

import (
	"github.com/octoberswimmer/masc"
	"testing"
)

// TestGridEmpty verifies that Grid returns nil when no children are provided.
func TestGridEmpty(t *testing.T) {
	if comp := Grid(); comp != nil {
		t.Error("Grid should be nil when no children provided")
	}
}

// TestGridBasic verifies that Grid returns a non-nil component when children are provided.
func TestGridBasic(t *testing.T) {
	child := masc.Text("cell")
	if comp := Grid(child); comp == nil {
		t.Error("Grid returned nil when children provided")
	}
}

// TestGridColumnBasic verifies that GridColumn returns a non-nil component for a valid size and child.
func TestGridColumnBasic(t *testing.T) {
	if comp := GridColumn("1-of-2", masc.Text("cell")); comp == nil {
		t.Error("GridColumn returned nil for valid size and child")
	}
}

// TestGridColumnNoChildren verifies that GridColumn returns a non-nil component when no children are provided.
func TestGridColumnNoChildren(t *testing.T) {
	if comp := GridColumn("1-of-3"); comp == nil {
		t.Error("GridColumn returned nil when no children provided")
	}
}
