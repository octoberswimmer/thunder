package components

import (
	"strings"

	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
)

// LookupOption represents one autocomplete suggestion.
type LookupOption struct {
	Label string
	Value string
}

// Lookup renders an SLDS lookup field.
// label: field label
// suggestions: full list to filter
// value: current input text
// onInput: called on input change
// onSelect: called when a suggestion is chosen
func Lookup(label string,
	suggestions []LookupOption,
	value string,
	onInput func(string),
	onSelect func(string),
) masc.ComponentOrHTML {

	// 1. Filter suggestions by input substring
	var filtered []LookupOption
	// Filter suggestions by input substring, but exclude exact matches to close dropdown after select
	q := strings.ToLower(strings.TrimSpace(value))
	if q != "" {
		for _, opt := range suggestions {
			lower := strings.ToLower(opt.Label)
			if lower != q && strings.Contains(lower, q) {
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

	// 3. Render the input field
	input := elem.Input(
		masc.Markup(
			masc.Class("slds-input", "slds-lookup__search-input"),
			masc.Property("type", "text"),
			masc.Property("value", value),
			masc.Property("placeholder", label),
			event.Input(func(e *masc.Event) {
				if onInput != nil {
					onInput(e.Target.Get("value").String())
				}
			}),
		),
	)

	wrapperClasses := []string{"slds-lookup"}
	if len(items) > 0 {
		wrapperClasses = append(wrapperClasses, "slds-is-open")
	}
	return elem.Div(
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
			elem.Div(
				masc.Markup(masc.Class("slds-lookup__menu")),
				elem.UnorderedList(
					append([]masc.MarkupOrChild{masc.Markup(masc.Class("slds-lookup__list"), masc.Attribute("role", "listbox"))},
						items...)..., // your <li class="slds-lookup__item"> entries
				),
			),
		),
	)

}
