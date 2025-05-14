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
