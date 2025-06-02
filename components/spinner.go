package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// Spinner renders an SLDS loading spinner.
// size is optional: "small", "medium", or "large". Default is "medium".
func Spinner(size string) masc.ComponentOrHTML {
	// Base spinner class
	spinnerClass := "slds-spinner_medium"
	if size != "" {
		spinnerClass = "slds-spinner_" + size
	}
	return elem.Div(
		masc.Markup(
			masc.Class("slds-spinner", "slds-spinner_brand", spinnerClass),
			masc.Attribute("role", "status"),
		),
		// Assistive text
		elem.Span(
			masc.Markup(masc.Class("slds-assistive-text")),
			masc.Text("Loading"),
		),
		// Spinner dots
		elem.Div(masc.Markup(masc.Class("slds-spinner__dot-a"))),
		elem.Div(masc.Markup(masc.Class("slds-spinner__dot-b"))),
	)
}

// SmallSpinner creates a small inline spinner for loading states.
func SmallSpinner() masc.ComponentOrHTML {
	return elem.Span(
		masc.Markup(
			masc.Class("slds-spinner", "slds-spinner_brand", "slds-spinner_x-small"),
			masc.Attribute("role", "status"),
		),
		elem.Span(
			masc.Markup(masc.Class("slds-assistive-text")),
			masc.Text("Loading"),
		),
		elem.Div(masc.Markup(masc.Class("slds-spinner__dot-a"))),
		elem.Div(masc.Markup(masc.Class("slds-spinner__dot-b"))),
	)
}

// LoadingCard creates a card container with a centered loading spinner.
func LoadingCard(title string) masc.ComponentOrHTML {
	return Card(title, CenteredSpinner("medium"))
}
