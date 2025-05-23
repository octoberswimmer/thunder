package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
)

// TextInput renders an SLDS styled text input with a label.
// label is the form element label text.
// value is the current input value.
// placeholder is optional placeholder text.
// onInput is the input event handler.
func TextInput(label, value, placeholder string, onInput func(*masc.Event)) masc.ComponentOrHTML {
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
			elem.Input(
				masc.Markup(
					masc.Class("slds-input"),
					masc.Property("type", "text"),
					masc.Property("value", value),
					masc.Property("placeholder", placeholder),
					event.Input(onInput),
				),
			),
		),
	)
}
