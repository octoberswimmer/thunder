package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
)

// ValidationState holds validation information for form fields.
type ValidationState struct {
	HasError     bool
	Required     bool
	ErrorMessage string
	HelpText     string
}

// ValidatedTextInput renders an SLDS styled text input with validation support.
// label is the form element label text.
// value is the current input value.
// placeholder is optional placeholder text.
// validation contains error state, required flag, and messages.
// onInput is the input event handler.
func ValidatedTextInput(label, value, placeholder string, validation ValidationState, onInput func(*masc.Event)) masc.ComponentOrHTML {
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
				masc.Markup(
					masc.Class(inputClasses...),
					masc.Property("type", "text"),
					masc.Property("value", value),
					masc.Property("placeholder", placeholder),
					masc.Property("required", validation.Required),
					event.Input(onInput),
				),
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

// ValidatedTextarea renders an SLDS styled textarea with validation support.
// label is the form element label text.
// value is the current textarea value.
// placeholder is optional placeholder text.
// rows is the number of visible text lines (defaults to 3 if 0).
// validation contains error state, required flag, and messages.
// onInput is the input event handler.
func ValidatedTextarea(label, value, placeholder string, rows int, validation ValidationState, onInput func(*masc.Event)) masc.ComponentOrHTML {
	// Default rows if not specified
	textareaRows := 3
	if rows > 0 {
		textareaRows = rows
	}

	// Build form element classes
	formClasses := []string{"slds-form-element", "slds-m-bottom_small"}
	textareaClasses := []string{"slds-textarea"}

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
			elem.TextArea(
				masc.Markup(
					masc.Class(textareaClasses...),
					masc.Property("placeholder", placeholder),
					masc.Property("rows", textareaRows),
					masc.Property("required", validation.Required),
					event.Input(onInput),
				),
				masc.Text(value),
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

// ValidatedSelect renders an SLDS styled select dropdown with validation support.
// label is the form element label text.
// options are the available select options.
// selected is the currently selected value.
// validation contains error state, required flag, and messages.
// onChange is the change event handler.
func ValidatedSelect(label string, options []SelectOption, selected string, validation ValidationState, onChange func(*masc.Event)) masc.ComponentOrHTML {
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

	// Build select options
	selectOpts := func() []masc.MarkupOrChild {
		var opts []masc.MarkupOrChild
		for _, opt := range options {
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

	// Build select element
	selectArgs := []masc.MarkupOrChild{
		masc.Markup(
			masc.Class("slds-select"),
			masc.Property("value", selected),
			masc.Property("required", validation.Required),
			event.Change(onChange),
		),
	}
	selectArgs = append(selectArgs, selectOpts...)

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
			elem.Div(
				masc.Markup(masc.Class("slds-select_container")),
				elem.Select(selectArgs...),
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

// ValidatedDatepicker renders an SLDS styled date picker with validation support.
// label is the form element label text.
// value is the current date value in YYYY-MM-DD format.
// validation contains error state, required flag, and messages.
// onChange is the change event handler.
func ValidatedDatepicker(label, value string, validation ValidationState, onChange func(*masc.Event)) masc.ComponentOrHTML {
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
				masc.Markup(
					masc.Class(inputClasses...),
					masc.Property("type", "date"),
					masc.Property("value", value),
					masc.Property("required", validation.Required),
					event.Change(onChange),
				),
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
