package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
)

// ToastVariant defines the SLDS theme for a toast notification.
type ToastVariant string

const (
	// VariantSuccess is a success toast.
	VariantSuccess ToastVariant = "slds-theme_success"
	// VariantError is an error toast.
	VariantError ToastVariant = "slds-theme_error"
	// VariantWarning is a warning toast.
	VariantWarning ToastVariant = "slds-theme_warning"
	// VariantInfo is an informational toast.
	VariantInfo ToastVariant = "slds-theme_info"
)

// Toast renders an SLDS toast notification with a header, message, and close button.
// variant governs the toast theme, header is the title text, message is the body text,
// and onClose is the click handler for the close action.
func Toast(variant ToastVariant, header, message string, onClose func(*masc.Event)) masc.ComponentOrHTML {
	// Container for notify
	return elem.Div(
		masc.Markup(masc.Class("slds-notify_container")),
		elem.Div(
			masc.Markup(
				masc.Class("slds-notify", "slds-notify_toast", string(variant)),
				masc.Property("role", "status"),
			),
			// Assistive text for screen readers
			elem.Span(
				masc.Markup(masc.Class("slds-assistive-text")),
				masc.Text(header),
			),
			// Grouped content: title and message
			elem.Div(
				masc.Markup(masc.Class("slds-notify__content")),
				// Title
				elem.Heading2(
					masc.Markup(masc.Class("slds-notify__title", "slds-text-heading_small", "slds-truncate")),
					masc.Text(header),
				),
				// Message
				elem.Div(
					masc.Markup(masc.Class("slds-notify__message")),
					masc.Text(message),
				),
			),
			// Close button
			elem.Button(
				masc.Markup(
					masc.Class("slds-button", "slds-button_icon", "slds-button_icon-inverse", "slds-notify__close"),
					event.Click(onClose),
				),
				elem.Span(
					masc.Markup(masc.Class("slds-assistive-text")),
					masc.Text("Close"),
				),
			),
		),
	)
}
