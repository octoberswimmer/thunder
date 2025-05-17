package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
)

// Badge renders an SLDS badge with the given label.
func Badge(label string) masc.ComponentOrHTML {
	return elem.Span(
		masc.Markup(masc.Class("slds-badge")),
		masc.Text(label),
	)
}

// Pill renders an SLDS pill container with the given label.
// If onRemove is non-nil, a remove button is rendered to the right.
func Pill(label string, onRemove func(*masc.Event)) masc.ComponentOrHTML {
	var args []masc.MarkupOrChild
	args = append(args, masc.Markup(masc.Class("slds-pill")))
	args = append(args, elem.Span(
		masc.Markup(masc.Class("slds-pill__label")),
		masc.Text(label),
	))
	if onRemove != nil {
		// Inline SVG for close icon plus assistive text (embedded path to avoid external sprite)
		svgHTML := `<svg class="slds-button__icon slds-button__icon-small" aria-hidden="true" viewBox="0 0 24 24" fill="currentColor">` +
			`<path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>` +
			`</svg><span class="slds-assistive-text">Remove</span>`
		args = append(args,
			elem.Button(
				masc.Markup(
					masc.Class("slds-button", "slds-button_icon", "slds-button_icon-small", "slds-pill__remove"),
					event.Click(onRemove),
					masc.Property("title", "Remove"),
				),
				masc.Markup(masc.UnsafeHTML(svgHTML)),
			),
		)
	}
	return elem.Span(args...)
}
