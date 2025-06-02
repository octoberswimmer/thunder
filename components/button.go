package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
)

// ButtonVariant defines the SLDS button style variant.
type ButtonVariant string

const (
	// VariantNeutral is the default neutral button style.
	VariantNeutral ButtonVariant = "slds-button_neutral"
	// VariantBrand is the brand button style.
	VariantBrand ButtonVariant = "slds-button_brand"
	// VariantDestructive is the destructive button style.
	VariantDestructive ButtonVariant = "slds-button_destructive"
)

// Button renders an SLDS button with the given label, style variant, and click handler.
// If variant is empty, VariantNeutral will be used.
func Button(label string, variant ButtonVariant, onClick func(*masc.Event)) masc.ComponentOrHTML {
	v := string(variant)
	if v == "" {
		v = string(VariantNeutral)
	}
	return elem.Button(
		masc.Markup(
			masc.Class("slds-button", v),
			event.Click(onClick),
		),
		masc.Text(label),
	)
}

// LoadingButton creates a button with a loading spinner and disabled state.
func LoadingButton(label string, variant ButtonVariant) masc.ComponentOrHTML {
	v := string(variant)
	if v == "" {
		v = string(VariantNeutral)
	}

	content := []masc.MarkupOrChild{}

	// Add spinner
	content = append(content,
		elem.Span(
			masc.Markup(masc.Class("slds-spinner", "slds-spinner_brand", "slds-spinner_x-small")),
			masc.Markup(masc.Attribute("role", "status")),
			elem.Span(
				masc.Markup(masc.Class("slds-assistive-text")),
				masc.Text("Loading"),
			),
			elem.Div(masc.Markup(masc.Class("slds-spinner__dot-a"))),
			elem.Div(masc.Markup(masc.Class("slds-spinner__dot-b"))),
		),
	)

	if label != "" {
		content = append(content, masc.Text(" "+label))
	}

	buttonContent := []masc.MarkupOrChild{
		masc.Markup(
			masc.Class("slds-button", v),
			masc.Attribute("disabled", "true"),
		),
	}
	buttonContent = append(buttonContent, content...)

	return elem.Button(buttonContent...)
}

// DisabledButton creates a button with disabled state and proper styling.
func DisabledButton(label string, variant ButtonVariant) masc.ComponentOrHTML {
	v := string(variant)
	if v == "" {
		v = string(VariantNeutral)
	}
	return elem.Button(
		masc.Markup(
			masc.Class("slds-button", v),
			masc.Attribute("disabled", "true"),
		),
		masc.Text(label),
	)
}

// ButtonGroupSpaced creates a spaced group of buttons with proper SLDS spacing.
func ButtonGroupSpaced(children ...masc.ComponentOrHTML) masc.ComponentOrHTML {
	args := []masc.MarkupOrChild{
		masc.Markup(masc.Class("slds-button-space")),
	}
	for _, child := range children {
		args = append(args, child)
	}
	return elem.Div(args...)
}

// ActionButtons creates a group of action buttons with consistent spacing.
// Typically used for Cancel/Previous/Next button groups.
func ActionButtons(children ...masc.ComponentOrHTML) masc.ComponentOrHTML {
	args := []masc.MarkupOrChild{
		masc.Markup(masc.Class("slds-m-top_large", "slds-button-space")),
	}
	for _, child := range children {
		args = append(args, child)
	}
	return elem.Div(args...)
}

// NavigationButtons creates the standard Cancel/Previous/Next button layout.
func NavigationButtons(
	onCancel func(*masc.Event),
	onPrevious func(*masc.Event),
	onNext func(*masc.Event),
	showPrevious bool,
) masc.ComponentOrHTML {
	buttons := []masc.ComponentOrHTML{
		Button("Cancel", VariantNeutral, onCancel),
	}

	if showPrevious {
		buttons = append(buttons, Button("Previous", VariantNeutral, onPrevious))
	}

	buttons = append(buttons, Button("Next", VariantBrand, onNext))

	return ActionButtons(buttons...)
}
