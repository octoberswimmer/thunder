package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// PageHeader renders an SLDS page header with a title, optional subtitle, and optional actions.
// title is the main heading text.
// subtitle is optional descriptive text below the title.
// actions are optional child components (e.g., buttons) rendered in the header control area.
func PageHeader(title string, subtitle string, actions ...masc.MarkupOrChild) masc.ComponentOrHTML {
	var children []masc.MarkupOrChild

	// Title element
	children = append(children,
		elem.Heading1(
			masc.Markup(
				masc.Class("slds-page-header__title", "slds-truncate"),
				masc.Property("title", title),
			),
			masc.Text(title),
		),
	)

	// Optional subtitle
	if subtitle != "" {
		children = append(children,
			// Use Paragraph element for subtitle text
			elem.Paragraph(
				masc.Markup(
					masc.Class("slds-page-header__meta"),
				),
				masc.Text(subtitle),
			),
		)
	}

	// Optional actions container
	if len(actions) > 0 {
		// Wrap actions in control div
		var controlChildren []masc.MarkupOrChild
		controlChildren = append(controlChildren,
			masc.Markup(masc.Class("slds-page-header__control")),
		)
		controlChildren = append(controlChildren, actions...)
		children = append(children,
			elem.Div(controlChildren...),
		)
	}

	// Outer page header container
	// Combine markup and children into a single args slice for Div
	var outer []masc.MarkupOrChild
	outer = append(outer, masc.Markup(masc.Class("slds-page-header")))
	outer = append(outer, children...)
	return elem.Div(
		outer...,
	)
}
