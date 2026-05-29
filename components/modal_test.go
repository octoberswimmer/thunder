package components

import (
	"github.com/octoberswimmer/masc"
	"testing"
)

// TestModalBasic verifies that Modal returns a non-nil component for a title and body.
func TestModalBasic(t *testing.T) {
	comp := Modal("Test Modal", masc.Text("Body content"))
	if comp == nil {
		t.Error("Modal returned nil for valid title and body")
	}
}

// TestModalDefaultWidthOmitsLargeClass verifies the default Modal is not large.
func TestModalDefaultWidthOmitsLargeClass(t *testing.T) {
	win := renderComponent(t, Modal("Test", masc.Text("Body")))
	node, err := win.Document().QuerySelector(".slds-modal_large")
	if err != nil {
		t.Fatal(err)
	}
	if node != nil {
		t.Error("default Modal should not carry the slds-modal_large class")
	}
}

// TestLargeModalAddsLargeClass verifies LargeModal renders at the large width.
func TestLargeModalAddsLargeClass(t *testing.T) {
	win := renderComponent(t, LargeModal("Test", masc.Text("Body")))
	node, err := win.Document().QuerySelector(".slds-modal.slds-modal_large")
	if err != nil {
		t.Fatal(err)
	}
	if node == nil {
		t.Fatal("LargeModal should add the slds-modal_large class to the modal wrapper")
	}
}
