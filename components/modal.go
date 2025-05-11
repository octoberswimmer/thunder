package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// Modal renders an SLDS modal dialog with the given title and body content.
// The modal and backdrop are always visible. To toggle visibility, wrap the
// component in masc.If or similar conditional logic.
// title is displayed in the header; content is rendered in the modal body.
// Modal renders an SLDS modal dialog with the given title and body content.
// The modal and backdrop are always visible. To toggle visibility, wrap the
// component in masc.If or similar conditional logic.
// title is displayed in the header; content is rendered in the modal body.
func Modal(title string, children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	// Split children into body (first) and footer (rest)
	var bodyChildren, footerChildren []masc.MarkupOrChild
	if len(children) > 0 {
		bodyChildren = append(bodyChildren, children[0])
	}
	if len(children) > 1 {
		footerChildren = children[1:]
	}
	// Build container arguments
	var containerArgs []masc.MarkupOrChild
	containerArgs = append(containerArgs,
		masc.Markup(masc.Class("slds-modal__container")),
	)
	// Header
	containerArgs = append(containerArgs,
		elem.Header(
			masc.Markup(masc.Class("slds-modal__header")),
			elem.Heading2(
				masc.Markup(
					masc.Class("slds-text-heading_medium", "slds-truncate"),
					masc.Property("title", title),
				),
				masc.Text(title),
			),
		),
	)
	// Content
	var contentArgs []masc.MarkupOrChild
	contentArgs = append(contentArgs,
		masc.Markup(masc.Class("slds-modal__content", "slds-p-around_medium")),
	)
	contentArgs = append(contentArgs, bodyChildren...)
	containerArgs = append(containerArgs,
		elem.Div(contentArgs...),
	)
	// Footer (optional)
	if len(footerChildren) > 0 {
		// Footer container with SLDS footer class
		var footerArgs []masc.MarkupOrChild
		footerArgs = append(footerArgs, masc.Markup(masc.Class("slds-modal__footer")))
		footerArgs = append(footerArgs, footerChildren...)
		containerArgs = append(containerArgs,
			elem.Div(footerArgs...),
		)
	}
	// Modal wrapper
	modal := elem.Div(
		masc.Markup(
			masc.Class("slds-modal", "slds-fade-in-open"),
			masc.Attribute("role", "dialog"),
			masc.Property("aria-modal", true),
		),
		elem.Div(containerArgs...),
	)
	// Backdrop
	backdrop := elem.Div(
		masc.Markup(masc.Class("slds-backdrop", "slds-backdrop_open")),
	)
	// Combine modal and backdrop into a slice literal
	return masc.List{modal, backdrop}
}
