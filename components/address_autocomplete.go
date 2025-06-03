package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
	"github.com/octoberswimmer/thunder/api"
)

// AddressAutocomplete renders an SLDS address autocomplete field using Google Places API.
// label: field label
// value: current input text
// apiKey: Google Maps API key
// predictions: current list of address predictions
// apiError: error message to display (if any)
// onInput: called on input change
// onSelect: called when an address is selected with full details
func AddressAutocomplete(
	label string,
	value string,
	apiKey string,
	predictions []api.PlacePrediction,
	apiError string,
	onInput func(string),
	onSelect func(api.PlaceDetails),
) masc.ComponentOrHTML {

	// Build suggestion items from predictions
	var items []masc.MarkupOrChild
	for _, pred := range predictions {
		items = append(items,
			elem.ListItem(
				masc.Markup(masc.Class("slds-lookup__item")),
				elem.Button(
					masc.Markup(
						masc.Class("slds-button", "slds-lookup__item-action"),
						event.Click(func(e *masc.Event) {
							// Get place details and notify selection
							go func() {
								if details, err := api.GetPlaceDetails(apiKey, pred.PlaceID); err == nil && details != nil && onSelect != nil {
									onSelect(*details)
								}
							}()
							// Update input with selected description
							if onInput != nil {
								onInput(pred.Description)
							}
						}),
					),
					elem.Div(
						masc.Markup(masc.Class("slds-truncate")),
						masc.Text(pred.Description),
					),
				),
			),
		)
	}

	// Render the input field
	input := elem.Input(
		masc.Markup(
			masc.Class("slds-input", "slds-lookup__search-input"),
			masc.Property("type", "text"),
			masc.Property("value", value),
			masc.Property("placeholder", "Enter address..."),
			event.Input(func(e *masc.Event) {
				if onInput != nil {
					onInput(e.Target.Get("value").String())
				}
			}),
		),
	)

	wrapperClasses := []string{"slds-lookup", "slds-form-element"}
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
				// Add location icon
				elem.Span(
					masc.Markup(masc.Class("slds-icon_container", "slds-icon-utility-location", "slds-input__icon", "slds-input__icon_right")),
					Icon("utility", "location", "x-small"),
				),
			),
			// suggestions dropdown
			elem.Div(
				masc.Markup(masc.Class("slds-lookup__menu")),
				elem.UnorderedList(
					append([]masc.MarkupOrChild{
						masc.Markup(masc.Class("slds-lookup__list"), masc.Attribute("role", "listbox")),
					}, items...)...,
				),
			),
		),
		// Show API error if present
		masc.If(apiError != "",
			elem.Div(
				masc.Markup(masc.Class("slds-form-element__help", "slds-text-color_error")),
				masc.Text(apiError),
			),
		),
	)
}

// AddressAutocompleteResult is the message type returned by AddressAutocompleteCmd
type AddressAutocompleteResult struct {
	Predictions []api.PlacePrediction
	Error       error
}

// AddressAutocompleteCmd creates a command to fetch address predictions
// This follows Thunder's pattern of returning masc.Cmd for async operations
func AddressAutocompleteCmd(apiKey, input string) masc.Cmd {
	return func() masc.Msg {
		predictions, err := api.GetPlacesAutocomplete(apiKey, input)
		return AddressAutocompleteResult{
			Predictions: predictions,
			Error:       err,
		}
	}
}
