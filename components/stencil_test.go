package components

import "testing"

// TestStencilNoLabel verifies that Stencil returns a component without label.
func TestStencilNoLabel(t *testing.T) {
	if comp := Stencil(""); comp == nil {
		t.Error("Stencil returned nil for no label")
	}
}

// TestStencilWithLabel verifies that Stencil returns a component with label.
func TestStencilWithLabel(t *testing.T) {
	if comp := Stencil("Test Label"); comp == nil {
		t.Error("Stencil returned nil with label")
	}
}

// TestStencilWithHeight verifies that Stencil handles custom height.
func TestStencilWithHeight(t *testing.T) {
	if comp := Stencil("", "3rem"); comp == nil {
		t.Error("Stencil returned nil with custom height")
	}
}

// TestStencilWithLabelAndHeight verifies that Stencil handles both label and height.
func TestStencilWithLabelAndHeight(t *testing.T) {
	if comp := Stencil("Test Label", "4rem"); comp == nil {
		t.Error("Stencil returned nil with label and height")
	}
}
