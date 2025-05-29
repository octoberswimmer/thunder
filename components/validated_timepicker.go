package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
)

// ValidatedTimepicker renders an SLDS styled time picker with validation support.
// label is the form element label text.
// value is the current time value in HH:MM format (24-hour).
// validation contains error state, required flag, and messages.
// onChange is the change event handler.
func ValidatedTimepicker(label, value string, validation ValidationState, onChange func(*masc.Event)) masc.ComponentOrHTML {
	// Build form element classes
	formClasses := []string{"slds-form-element", "slds-m-bottom_small"}
	inputClasses := []string{"slds-input"}

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
	if validation.Tooltip != "" {
		labelContent = append(labelContent,
			elem.Span(
				masc.Markup(
					masc.Class("slds-m-left_xx-small", "slds-text-color_weak"),
					masc.Property("title", validation.Tooltip),
				),
				masc.Text("ⓘ"),
			),
		)
	}

	// Build input properties
	inputProps := []masc.Applyer{
		masc.Class(inputClasses...),
		masc.Property("type", "time"),
		masc.Property("value", value),
		masc.Property("placeholder", validation.Placeholder),
		masc.Property("required", validation.Required),
		event.Change(onChange),
	}

	// Build the form element
	children := []masc.MarkupOrChild{
		// Label
		func() masc.ComponentOrHTML {
			args := []masc.MarkupOrChild{masc.Markup(masc.Class("slds-form-element__label"))}
			args = append(args, labelContent...)
			return elem.Label(args...)
		}(),
		// Control wrapper
		elem.Div(
			masc.Markup(masc.Class("slds-form-element__control")),
			elem.Input(
				masc.Markup(inputProps...),
			),
		),
	}

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
