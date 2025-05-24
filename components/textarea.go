package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
)

// Textarea renders an SLDS styled textarea with a label.
// label is the form element label text.
// value is the current textarea value.
// placeholder is optional placeholder text.
// rows is optional number of visible text lines (defaults to 3 if not specified).
// onInput is the input event handler.
func Textarea(label, value, placeholder string, rows int, onInput func(*masc.Event)) masc.ComponentOrHTML {
	// Default rows if not specified
	textareaRows := 3
	if rows > 0 {
		textareaRows = rows
	}

	// Form element wrapper
	return elem.Div(
		masc.Markup(masc.Class("slds-form-element", "slds-m-bottom_small")),
		// Label
		elem.Label(
			masc.Markup(masc.Class("slds-form-element__label")),
			masc.Text(label),
		),
		// Control
		elem.Div(
			masc.Markup(masc.Class("slds-form-element__control")),
			elem.TextArea(
				masc.Markup(
					masc.Class("slds-textarea"),
					masc.Property("placeholder", placeholder),
					masc.Property("rows", textareaRows),
					event.Input(onInput),
				),
				masc.Text(value),
			),
		),
	)
}
