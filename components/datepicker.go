package components

import (
	"time"

	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
)

// Datepicker renders an SLDS styled date picker with a label.
// label is the form element label text.
// value is the current selected date value as time.Time (zero value for empty).
// onChange is the event handler when the date value changes, receives the new time.Time value.
func Datepicker(label string, value time.Time, onChange func(time.Time)) masc.ComponentOrHTML {
	// Convert time.Time to string for HTML input
	var valueStr string
	if !value.IsZero() {
		valueStr = value.Format("2006-01-02")
	}
	return elem.Div(
		masc.Markup(masc.Class("slds-form-element", "slds-m-bottom_small")),
		elem.Label(
			masc.Markup(masc.Class("slds-form-element__label")),
			masc.Text(label),
		),
		elem.Div(
			masc.Markup(masc.Class("slds-form-element__control")),
			elem.Input(
				masc.Markup(
					masc.Class("slds-input"),
					masc.Property("type", "date"),
					masc.Property("value", valueStr),
					event.Change(func(e *masc.Event) {
						dateStr := e.Target.Get("value").String()
						var newValue time.Time
						if dateStr != "" {
							if parsedDate, err := time.Parse("2006-01-02", dateStr); err == nil {
								newValue = parsedDate
							}
							// If parse fails, newValue remains zero value
						}
						// Call onChange with parsed time.Time value
						onChange(newValue)
					}),
				),
			),
		),
	)
}
