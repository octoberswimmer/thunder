package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// Prompt renders an SLDS modal dialog with the given title, body content, and actions.
// The modal and backdrop are always visible. To toggle visibility, wrap the
// component in masc.If or similar conditional logic.
// title is displayed in the header; content is rendered in the modal body, actions in the footer.
func Prompt(title string, actions []masc.MarkupOrChild, content ...masc.MarkupOrChild) masc.ComponentOrHTML {
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
	contentArgs = append(contentArgs, content...)
	containerArgs = append(containerArgs,
		elem.Div(contentArgs...),
	)
	// Footer
	var footerArgs []masc.MarkupOrChild
	footerArgs = append(footerArgs, masc.Markup(masc.Class("slds-modal__footer")))
	footerArgs = append(footerArgs, actions...)
	containerArgs = append(containerArgs,
		elem.Div(footerArgs...),
	)
	// Modal wrapper
	modal := elem.Div(
		masc.Markup(
			masc.Class("slds-modal", "slds-fade-in-open", "slds-modal_prompt"),
			masc.Attribute("role", "alertdialog"),
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
