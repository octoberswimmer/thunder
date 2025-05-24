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
