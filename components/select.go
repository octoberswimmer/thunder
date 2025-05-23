package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
)

// SelectOption represents a single dropdown option.
type SelectOption struct {
	Label string
	Value string
}

// Select renders an SLDS styled dropdown with a label.
// label is the form element label text.
// options is the list of SelectOption structs.
// selected is the currently selected option value.
// onChange is the change event handler.
func Select(label string, options []SelectOption, selected string, onChange func(*masc.Event)) masc.ComponentOrHTML {
	// Form element wrapper
	selectOpts := func() []masc.MarkupOrChild {
		var opts []masc.MarkupOrChild
		for _, opt := range options {
			// Option markup: always set value; mark selected option
			var pm []masc.Applyer
			pm = append(pm, masc.Property("value", opt.Value))
			if opt.Value == selected {
				pm = append(pm, masc.Property("selected", true))
			}
			opts = append(opts,
				elem.Option(
					masc.Markup(pm...),
					masc.Text(opt.Label),
				),
			)
		}
		return opts
	}()
	markupOrChild := []masc.MarkupOrChild{
		masc.Markup(
			masc.Class("slds-select"),
			event.Change(onChange),
			masc.Property("value", selected),
		),
	}
	markupOrChild = append(markupOrChild, selectOpts...)
	wrapper := elem.Div(
		masc.Markup(masc.Class("slds-form-element", "slds-m-bottom_small")),
		// Label for the select
		elem.Label(
			masc.Markup(masc.Class("slds-form-element__label")),
			masc.Text(label),
		),
		// Control wrapper
		elem.Div(
			masc.Markup(masc.Class("slds-form-element__control")),
			// The select element
			elem.Select(markupOrChild...),
		),
	)
	return wrapper
}
