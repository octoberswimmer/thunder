package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// Stencil renders an SLDS skeleton loading placeholder with optional label.
// label is optional text to display above the skeleton.
// height is optional height override (defaults to "2.25rem").
func Stencil(label string, height ...string) masc.ComponentOrHTML {
	stencilHeight := "2.25rem"
	if len(height) > 0 && height[0] != "" {
		stencilHeight = height[0]
	}

	stencilDiv := elem.Div(
		masc.Markup(
			masc.Class("slds-skeleton"),
			masc.Style("height", stencilHeight),
			masc.Style("border-radius", "0.25rem"),
			masc.Style("background-color", "#f3f3f3"),
			masc.Style("animation", "slds-skeleton-loading 1.5s infinite ease-in-out"),
		),
	)

	// If no label is provided, return just the skeleton div
	if label == "" {
		return stencilDiv
	}

	// Return form element structure with label when label is provided
	return elem.Div(
		masc.Markup(masc.Class("slds-form-element", "slds-m-bottom_small")),
		elem.Label(
			masc.Markup(masc.Class("slds-form-element__label")),
			masc.Text(label),
		),
		elem.Div(
			masc.Markup(masc.Class("slds-form-element__control")),
			stencilDiv,
		),
	)
}
