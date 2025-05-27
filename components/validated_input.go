package components

import (
	"fmt"
	"math/rand"
	"time"

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
	Tooltip      string
	Placeholder  string
}

// Convenience constructors for common ValidationState configurations

// WithTooltip creates a ValidationState with just a tooltip
func WithTooltip(tooltip string) ValidationState {
	return ValidationState{Tooltip: tooltip}
}

// WithPlaceholder creates a ValidationState with just a placeholder
func WithPlaceholder(placeholder string) ValidationState {
	return ValidationState{Placeholder: placeholder}
}

// WithTooltipAndPlaceholder creates a ValidationState with both tooltip and placeholder
func WithTooltipAndPlaceholder(tooltip, placeholder string) ValidationState {
	return ValidationState{Tooltip: tooltip, Placeholder: placeholder}
}

// Required creates a ValidationState marked as required
func Required() ValidationState {
	return ValidationState{Required: true}
}

// RequiredWithTooltip creates a required ValidationState with a tooltip
func RequiredWithTooltip(tooltip string) ValidationState {
	return ValidationState{Required: true, Tooltip: tooltip}
}

// generateTooltipID generates a unique ID for tooltips
func generateTooltipID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("tooltip-%d", rand.Intn(1000000))
}

// renderTooltip creates an SLDS tooltip element
func renderTooltip(id, text string) masc.ComponentOrHTML {
	if text == "" {
		return nil
	}

	return elem.Div(
		masc.Markup(
			masc.Class("slds-popover", "slds-popover_tooltip", "slds-fall-into-ground"),
			masc.Property("id", id),
			masc.Property("role", "tooltip"),
		),
		elem.Div(
			masc.Markup(masc.Class("slds-popover__body")),
			masc.Text(text),
		),
	)
}

// addTooltipAttributes adds necessary attributes for tooltip functionality
func addTooltipAttributes(existing []masc.Applyer, tooltipID, tooltipText string) []masc.Applyer {
	// We don't add tooltip attributes to the input since we're using help icon pattern
	return existing
}

// ValidatedTextInput renders an SLDS styled text input with validation support.
// label is the form element label text.
// value is the current input value.
// validation contains error state, required flag, messages, tooltip, and placeholder.
// onInput is the input event handler.
func ValidatedTextInput(label, value string, validation ValidationState, onInput func(*masc.Event)) masc.ComponentOrHTML {

	// Build form element classes
	formClasses := []string{"slds-form-element", "slds-m-bottom_small"}
	inputClasses := []string{"slds-input"}

	if validation.HasError {
		formClasses = append(formClasses, "slds-has-error")
	}

	// Build label with required indicator and tooltip icon
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

	// Build input properties with tooltip support
	inputProps := []masc.Applyer{
		masc.Class(inputClasses...),
		masc.Property("type", "text"),
		masc.Property("value", value),
		masc.Property("placeholder", validation.Placeholder),
		masc.Property("required", validation.Required),
		event.Input(onInput),
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

// ValidatedTextarea renders an SLDS styled textarea with validation support.
// label is the form element label text.
// value is the current textarea value.
// rows is the number of visible text lines (defaults to 3 if 0).
// validation contains error state, required flag, messages, tooltip, and placeholder.
// onInput is the input event handler.
func ValidatedTextarea(label, value string, rows int, validation ValidationState, onInput func(*masc.Event)) masc.ComponentOrHTML {

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

	// Build label with required indicator and tooltip icon
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

	// Build textarea properties with tooltip support
	textareaProps := []masc.Applyer{
		masc.Class(textareaClasses...),
		masc.Property("placeholder", validation.Placeholder),
		masc.Property("rows", textareaRows),
		masc.Property("required", validation.Required),
		event.Input(onInput),
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
				masc.Markup(textareaProps...),
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

	// Build label with required indicator and tooltip icon
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
