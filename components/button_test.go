package components

import (
	"testing"
)

// TestAlignedFieldWrapsContent verifies AlignedField reserves an empty label
// slot above its content so the content lines up with labeled form fields.
func TestAlignedFieldWrapsContent(t *testing.T) {
	win := renderComponent(t, AlignedField(Button("Go", VariantNeutral, nil)))

	label, err := win.Document().QuerySelector("label.slds-form-element__label")
	if err != nil {
		t.Fatal(err)
	}
	if label == nil {
		t.Fatal("expected an empty label slot, got none")
	}

	btn, err := win.Document().QuerySelector(".slds-form-element__control button")
	if err != nil {
		t.Fatal(err)
	}
	if btn == nil {
		t.Fatal("expected the wrapped button inside the control slot, got none")
	}
}

// TestAlignedButtonUsesAlignedField verifies AlignedButton renders a button
// within the aligned form-element structure.
func TestAlignedButtonUsesAlignedField(t *testing.T) {
	win := renderComponent(t, AlignedButton("Next", VariantBrand, nil))

	btn, err := win.Document().QuerySelector(".slds-form-element__control button.slds-button_brand")
	if err != nil {
		t.Fatal(err)
	}
	if btn == nil {
		t.Fatal("expected an aligned brand button, got none")
	}
}
