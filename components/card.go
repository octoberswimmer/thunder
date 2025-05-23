// Package components provides SLDS-styled Masc components under the thunder namespace.
package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// Card renders an SLDS Card with a header and optional body content.
// title is the card header text.
// children represent the card body inner content.
func Card(title string, children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	// Prepare header arguments
	var headerArgs []masc.MarkupOrChild
	headerArgs = append(headerArgs,
		masc.Markup(
			masc.Class("slds-card__header", "slds-grid", "slds-grid_vertical-align-center"),
		),
	)
	headerArgs = append(headerArgs,
		elem.Heading2(
			masc.Markup(
				masc.Class("slds-text-heading_small", "slds-truncate"),
				masc.Property("title", title),
			),
			masc.Text(title),
		),
	)
	cardHeader := elem.Div(headerArgs...)

	// Prepare body arguments
	var bodyArgs []masc.MarkupOrChild
	bodyArgs = append(bodyArgs,
		masc.Markup(
			masc.Class("slds-card__body", "slds-card__body_inner"),
		),
	)
	bodyArgs = append(bodyArgs, children...)
	cardBody := elem.Div(bodyArgs...)

	// Assemble card container
	return elem.Article(
		masc.Markup(
			masc.Class("slds-card", "slds-m-bottom_medium"),
		),
		cardHeader,
		cardBody,
	)
}
