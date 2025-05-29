package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
)

// Timepicker renders an SLDS styled time picker.
// label is the form element label text.
// value is the current time value in HH:MM format (24-hour).
// onChange is the change event handler.
func Timepicker(label, value string, onChange func(*masc.Event)) masc.ComponentOrHTML {
	// Build form element
	return elem.Div(
		masc.Markup(masc.Class("slds-form-element", "slds-m-bottom_small")),
		// Label
		elem.Label(
			masc.Markup(masc.Class("slds-form-element__label")),
			masc.Text(label),
		),
		// Control wrapper
		elem.Div(
			masc.Markup(masc.Class("slds-form-element__control")),
			elem.Input(
				masc.Markup(
					masc.Class("slds-input"),
					masc.Property("type", "time"),
					masc.Property("value", value),
					event.Change(onChange),
				),
			),
		),
	)
}
