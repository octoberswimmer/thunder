package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
)

// Modal renders an SLDS modal dialog with the given title and body content.
// The modal and backdrop are always visible. To toggle visibility, wrap the
// component in masc.If or similar conditional logic.
// title is displayed in the header; content is rendered in the modal body.
func Modal(title string, children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	return modal(title, "", children...)
}

// LargeModal renders a Modal at the SLDS "large" width, giving wide content
// such as multi-column tables room to lay out without horizontal scrolling.
func LargeModal(title string, children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	return modal(title, "slds-modal_large", children...)
}

// modal is the shared implementation behind Modal and LargeModal. sizeClass is
// an optional SLDS size modifier (e.g. "slds-modal_large") applied to the modal
// wrapper; an empty string yields the default width.
func modal(title, sizeClass string, children ...masc.MarkupOrChild) masc.ComponentOrHTML {
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
	modalClasses := []string{"slds-modal", "slds-fade-in-open"}
	if sizeClass != "" {
		modalClasses = append(modalClasses, sizeClass)
	}
	modal := elem.Div(
		masc.Markup(
			masc.Class(modalClasses...),
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

// ModalWithClose renders an SLDS modal dialog with a close button and backdrop click support.
// The modal and backdrop are always visible. To toggle visibility, wrap the
// component in masc.If or similar conditional logic.
// title is displayed in the header; onClose is called when close button or backdrop is clicked.
// children: first element becomes modal body content, remaining elements become footer buttons.
func ModalWithClose(title string, onClose func(*masc.Event), children ...masc.MarkupOrChild) masc.ComponentOrHTML {
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

	// Header with close button
	containerArgs = append(containerArgs,
		elem.Header(
			masc.Markup(masc.Class("slds-modal__header")),
			elem.Button(
				masc.Markup(
					masc.Class("slds-button", "slds-button_icon", "slds-modal__close", "slds-button_icon-inverse"),
					masc.Attribute("title", "Close"),
					masc.Attribute("type", "button"),
					event.Click(onClose),
				),
				elem.Span(
					masc.Markup(
						masc.Class("slds-button__icon"),
						masc.Style("font-size", "1rem"),
						masc.Style("line-height", "1"),
					),
					masc.Text("✕"),
				),
				elem.Span(
					masc.Markup(masc.Class("slds-assistive-text")),
					masc.Text("Close"),
				),
			),
			elem.Heading2(
				masc.Markup(
					masc.Class("slds-modal__title", "slds-hyphenate"),
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
			masc.Attribute("aria-labelledby", "modal-heading-01"),
		),
		elem.Div(containerArgs...),
	)

	// Backdrop with click handler
	backdrop := elem.Div(
		masc.Markup(
			masc.Class("slds-backdrop", "slds-backdrop_open"),
			event.Click(onClose),
		),
	)

	// Combine modal and backdrop into a slice literal
	return masc.List{backdrop, modal}
}
