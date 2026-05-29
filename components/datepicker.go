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
	return datepicker(label, value, time.Time{}, time.Time{}, false, onChange)
}

// AlignedDatepicker renders a date picker that vertically aligns with adjacent
// labeled form fields (e.g. AlignedButton). When label is empty an invisible
// label keeps the control aligned. The selectable range is bounded by min and
// max; pass the zero time for either to leave that bound open.
func AlignedDatepicker(label string, value, min, max time.Time, onChange func(time.Time)) masc.ComponentOrHTML {
	return datepicker(label, value, min, max, true, onChange)
}

// datepicker is the shared implementation behind Datepicker and
// AlignedDatepicker. aligned controls the empty-label-slot behaviour and drops
// the default bottom margin so the control lines up with other aligned fields.
func datepicker(label string, value, min, max time.Time, aligned bool, onChange func(time.Time)) masc.ComponentOrHTML {
	// Convert time.Time to string for HTML input
	var valueStr string
	if !value.IsZero() {
		valueStr = value.Format("2006-01-02")
	}

	formClasses := []string{"slds-form-element"}
	if !aligned {
		formClasses = append(formClasses, "slds-m-bottom_small")
	}

	// When aligned, an empty label is rendered as a non-breaking space so the
	// control still reserves the label row and lines up with labeled siblings.
	labelText := label
	if aligned && labelText == "" {
		labelText = " "
	}

	inputMarkup := []masc.Applyer{
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
	}
	if !min.IsZero() {
		inputMarkup = append(inputMarkup, masc.Property("min", min.Format("2006-01-02")))
	}
	if !max.IsZero() {
		inputMarkup = append(inputMarkup, masc.Property("max", max.Format("2006-01-02")))
	}

	return elem.Div(
		masc.Markup(masc.Class(formClasses...)),
		elem.Label(
			masc.Markup(masc.Class("slds-form-element__label")),
			masc.Text(labelText),
		),
		elem.Div(
			masc.Markup(masc.Class("slds-form-element__control")),
			elem.Input(
				masc.Markup(inputMarkup...),
			),
		),
	)
}
