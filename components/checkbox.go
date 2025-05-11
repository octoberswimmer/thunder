package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
)

// Checkbox renders an SLDS styled checkbox with a label.
// label is the text to display next to the checkbox.
// checked indicates the current state of the checkbox.
// onChange is the event handler when the checkbox is toggled.
func Checkbox(label string, checked bool, onChange func(*masc.Event)) masc.ComponentOrHTML {
	return elem.Div(
		masc.Markup(masc.Class("slds-form-element")),
		elem.Label(
			masc.Markup(masc.Class("slds-checkbox")),
			elem.Input(
				masc.Markup(
					masc.Property("type", "checkbox"),
					masc.Property("checked", checked),
					event.Change(onChange),
				),
			),
			elem.Span(
				masc.Markup(masc.Class("slds-checkbox_faux")),
			),
			elem.Span(
				masc.Markup(masc.Class("slds-form-element__label")),
				masc.Text(label),
			),
		),
	)
}
