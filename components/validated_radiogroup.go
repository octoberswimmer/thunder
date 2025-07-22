package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// ValidatedRadioButtonGroup renders an SLDS styled radio button group with validation support.
// label is the form element label text.
// name is the shared name attribute for the radio inputs.
// options are the available radio button options.
// selected is the currently selected value.
// validation contains error state, required flag, and messages.
// onChange is called with the Value of the selected option when clicked.
func ValidatedRadioButtonGroup(label, name string, options []RadioButtonOption, selected string, validation ValidationState, onChange func(string)) masc.ComponentOrHTML {
	// Build form element classes
	formClasses := []string{"slds-form-element", "slds-m-bottom_small"}

	if validation.HasError {
		formClasses = append(formClasses, "slds-has-error")
	}

	// Build label with required indicator
	labelContent := []masc.MarkupOrChild{masc.Text(label)}
	if validation.Required {
		labelContent = append(labelContent,
			elem.Span(
				masc.Markup(masc.Class("slds-required")),
				masc.Text(" *"),
			),
		)
	}

	// Build child elements
	var children []masc.MarkupOrChild

	// Add label
	labelArgs := []masc.MarkupOrChild{masc.Markup(masc.Class("slds-form-element__label"))}
	labelArgs = append(labelArgs, labelContent...)
	children = append(children,
		elem.Label(labelArgs...),
	)

	// Add control container with radio button group
	children = append(children,
		elem.Div(
			masc.Markup(masc.Class("slds-form-element__control")),
			RadioButtonGroup("", name, options, selected, onChange),
		),
	)

	// Add error message if present
	if validation.HasError && validation.ErrorMessage != "" {
		children = append(children,
			elem.Div(
				masc.Markup(masc.Class("slds-form-element__help")),
				masc.Text(validation.ErrorMessage),
			),
		)
	}

	// Add help text if present (and no error showing)
	if !validation.HasError && validation.HelpText != "" {
		children = append(children,
			elem.Div(
				masc.Markup(masc.Class("slds-form-element__help")),
				masc.Text(validation.HelpText),
			),
		)
	}

	// Build final element
	args := []masc.MarkupOrChild{masc.Markup(masc.Class(formClasses...))}
	args = append(args, children...)
	return elem.Div(args...)
}
