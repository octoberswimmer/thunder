package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
)

// RadioOption represents a single radio button choice.
type RadioOption struct {
	Label string
	Value string
}

// RadioGroup renders an SLDS styled group of radio buttons within a form element.
// name is the shared name attribute for the radio inputs.
// legend is the group label text.
// options is the slice of RadioOption values.
// selected is the currently selected option value.
// onChange is called with the Value of the selected option when clicked.
func RadioGroup(name, legend string, options []RadioOption, selected string, onChange func(string)) masc.ComponentOrHTML {
	// Build argument list for the form-element container, attach change listener to fieldset
	var args []masc.MarkupOrChild
	args = append(args,
		masc.Markup(
			masc.Class("slds-form-element"),
		),
	)
	// Legend (group label)
	args = append(args,
		elem.Legend(
			masc.Markup(masc.Class("slds-form-element__legend", "slds-form-element__label")),
			masc.Text(legend),
		),
	)
	// Build control container and radio options
	var controlArgs []masc.MarkupOrChild
	controlArgs = append(controlArgs, masc.Markup(masc.Class("slds-form-element__control")))
	for _, opt := range options {
		id := name + "-" + opt.Value
		// Each radio option
		controlArgs = append(controlArgs,
			elem.Div(
				masc.Markup(
					masc.Class("slds-radio"),
					event.Click(func(e *masc.Event) {
						if onChange != nil {
							onChange(opt.Value)
						}
					}),
				),
				// Actual radio input
				elem.Input(
					masc.Markup(
						masc.Property("type", "radio"),
						masc.Property("name", name),
						masc.Property("value", opt.Value),
						masc.Property("id", id),
						masc.Property("checked", opt.Value == selected),
					),
				),
				// Visible label for the radio
				elem.Label(
					masc.Markup(
						masc.Class("slds-radio__label"),
						masc.Property("for", id),
					),
					elem.Span(masc.Markup(masc.Class("slds-radio_faux"))),
					elem.Span(masc.Markup(masc.Class("slds-form-element__label")), masc.Text(opt.Label)),
				),
			),
		)
	}
	// Append control container to form-element container
	args = append(args, elem.Div(controlArgs...))
	// Return the assembled group wrapped in a fieldset for proper semantics
	return elem.FieldSet(args...)
}

// RadioButtonOption represents a single radio button choice for RadioButtonGroup.
type RadioButtonOption struct {
	Label string
	Value string
}

// RadioButtonGroup renders an SLDS styled group of radio buttons using the radio button group pattern.
// This follows the SLDS Radio Button Group component specification.
// label is the optional form element label text (can be empty string for no label).
// name is the shared name attribute for the radio inputs.
// options is the slice of RadioButtonOption values.
// selected is the currently selected option value.
// onChange is called with the Value of the selected option when clicked.
func RadioButtonGroup(label, name string, options []RadioButtonOption, selected string, onChange func(string)) masc.ComponentOrHTML {
	// If a label is provided, wrap the radio button group in a form element
	if label != "" {
		return elem.Div(
			masc.Markup(masc.Class("slds-form-element")),
			elem.Label(
				masc.Markup(masc.Class("slds-form-element__label")),
				masc.Text(label),
			),
			elem.Div(
				masc.Markup(masc.Class("slds-form-element__control")),
				renderRadioButtonGroup(name, options, selected, onChange),
			),
		)
	}

	// No label, just return the radio button group
	return renderRadioButtonGroup(name, options, selected, onChange)
}

// renderRadioButtonGroup renders the actual radio button group without label wrapping
func renderRadioButtonGroup(name string, options []RadioButtonOption, selected string, onChange func(string)) masc.ComponentOrHTML {
	var args []masc.MarkupOrChild
	args = append(args,
		masc.Markup(
			masc.Class("slds-radio_button-group"),
			masc.Property("role", "radiogroup"),
		),
	)

	for _, opt := range options {
		id := name + "-" + opt.Value
		// Each radio button following SLDS radio button group pattern
		args = append(args,
			elem.Span(
				masc.Markup(
					masc.Class("slds-radio_button"),
					event.Click(func(e *masc.Event) {
						if onChange != nil {
							onChange(opt.Value)
						}
					}),
				),
				// Actual radio input
				elem.Input(
					masc.Markup(
						masc.Property("type", "radio"),
						masc.Property("name", name),
						masc.Property("value", opt.Value),
						masc.Property("id", id),
						masc.Property("checked", opt.Value == selected),
					),
				),
				// Label for the radio button with styled faux element inside
				elem.Label(
					masc.Markup(
						masc.Class("slds-radio_button__label"),
						masc.Property("for", id),
						masc.Style("user-select", "none"),
						masc.Style("padding-bottom", "0.25rem"),
					),
					elem.Span(
						masc.Markup(masc.Class("slds-radio_faux")),
						masc.Text(opt.Label),
					),
				),
			),
		)
	}

	return elem.Div(args...)
}
