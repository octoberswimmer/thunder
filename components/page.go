package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// Page renders an SLDS page layout with a header and content area.
// header is a PageHeader or similar component.
// content is the main page body, typically cards or other containers.
func Page(header masc.ComponentOrHTML, content ...masc.MarkupOrChild) masc.ComponentOrHTML {
	// Wrap header and main content
	// Build content wrapper args
	var cwArgs []masc.MarkupOrChild
	cwArgs = append(cwArgs, masc.Markup(masc.Class("slds-page-content", "slds-p-around_medium")))
	cwArgs = append(cwArgs, content...)
	contentWrapper := elem.Div(cwArgs...)
	// Return header and content wrapper as a list
	return masc.List{header, contentWrapper}
}
