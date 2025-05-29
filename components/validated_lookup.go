package components

import (
	"strings"

	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
)

// ValidatedLookup renders an SLDS lookup field with validation support and Escape key reset.
// label: field label
// suggestions: full list to filter
// value: current input text
// validation: validation state
// onInput: called on input change
// onSelect: called when a suggestion is chosen
// onReset: called when Escape is pressed to get the reset value
func ValidatedLookup(label string,
	suggestions []LookupOption,
	value string,
	validation ValidationState,
	onInput func(string),
	onSelect func(string),
	onReset func() string,
) masc.ComponentOrHTML {

	// 1. Filter suggestions by input substring
	var filtered []LookupOption
	// Filter suggestions by input substring
	q := strings.ToLower(strings.TrimSpace(value))
	if q != "" && len(suggestions) > 0 {
		for _, opt := range suggestions {
			lower := strings.ToLower(opt.Label)
			if strings.Contains(lower, q) {
				filtered = append(filtered, opt)
			}
		}
	}

	// 2. Build suggestion items
	var items []masc.MarkupOrChild
	for _, opt := range filtered {
		items = append(items,
			elem.ListItem(
				masc.Markup(masc.Class("slds-lookup__item")),
				elem.Button(
					masc.Markup(
						masc.Class("slds-button", "slds-lookup__item-action"),
						event.Click(func(e *masc.Event) {
							// Notify selection
							if onSelect != nil {
								onSelect(opt.Value)
							}
							// Keep selected label in input and close dropdown
							if onInput != nil {
								onInput(opt.Label)
							}
						}),
					),
					masc.Text(opt.Label),
				),
			),
		)
	}

	// 3. Build input field with validation styling
	var inputClasses []string = []string{"slds-input", "slds-lookup__search-input"}
	if validation.HasError {
		inputClasses = append(inputClasses, "slds-has-error")
	}

	// 4. Render the input field
	input := elem.Input(
		masc.Markup(
			masc.Class(inputClasses...),
			masc.Property("type", "text"),
			masc.Property("value", value),
			masc.Property("placeholder", "Search "+label+"..."),
			event.Input(func(e *masc.Event) {
				if onInput != nil {
					onInput(e.Target.Get("value").String())
				}
			}),
			event.KeyDown(func(e *masc.Event) {
				// Access key properties from the event object
				key := e.Get("key").String()
				keyCode := e.Get("keyCode").String()
				which := e.Get("which").String()
				code := e.Get("code").String()

				if key == "Escape" || keyCode == "27" || which == "27" || code == "Escape" {
					// Reset to the current selected value on Escape
					if onReset != nil && onInput != nil {
						resetValue := onReset()
						onInput(resetValue)
					}
				}
			}),
		),
	)

	// 5. Build dropdown
	var dropdown masc.ComponentOrHTML
	if len(items) > 0 {
		dropdown = elem.Div(
			masc.Markup(masc.Class("slds-lookup__menu")),
			elem.UnorderedList(
				append([]masc.MarkupOrChild{masc.Markup(masc.Class("slds-lookup__list"), masc.Attribute("role", "listbox"))},
					items...)...,
			),
		)
	}

	// 6. Build error message if validation failed
	var errorMessage masc.ComponentOrHTML
	if validation.HasError {
		errorMessage = elem.Div(
			masc.Markup(masc.Class("slds-form-element__help")),
			masc.Text(validation.ErrorMessage),
		)
	}

	// 7. Build wrapper classes
	var wrapperClasses []string = []string{"slds-lookup"}
	if len(items) > 0 {
		wrapperClasses = append(wrapperClasses, "slds-is-open")
	}

	// 8. Build form element classes
	var formElementClasses []string = []string{"slds-form-element"}
	if validation.HasError {
		formElementClasses = append(formElementClasses, "slds-has-error")
	}

	return elem.Div(
		masc.Markup(masc.Class(formElementClasses...)),
		elem.Div(
			masc.Markup(masc.Class(wrapperClasses...)),
			elem.Label(
				masc.Markup(masc.Class("slds-form-element__label")),
				masc.Text(label),
			),
			elem.Div(
				masc.Markup(masc.Class("slds-form-element__control")),
				// search input wrapper
				elem.Div(
					masc.Markup(masc.Class(
						"slds-lookup__search-input",
						"slds-input-has-icon",
						"slds-input-has-icon_right",
					)),
					input,
				),
				// suggestions dropdown
				dropdown,
			),
		),
		errorMessage,
	)
}
