package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
)

// Datepicker renders an SLDS styled date picker with a label.
// label is the form element label text.
// value is the current selected date value in YYYY-MM-DD format.
// onChange is the event handler when the date value changes.
func Datepicker(label, value string, onChange func(*masc.Event)) masc.ComponentOrHTML {
	return elem.Div(
		masc.Markup(masc.Class("slds-form-element")),
		elem.Label(
			masc.Markup(masc.Class("slds-form-element__label")),
			masc.Text(label),
		),
		elem.Div(
			masc.Markup(masc.Class("slds-form-element__control")),
			elem.Input(
				masc.Markup(
					masc.Class("slds-input"),
					masc.Property("type", "date"),
					masc.Property("value", value),
					event.Change(onChange),
				),
			),
		),
	)
}
