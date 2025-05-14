package components

import (
	"fmt"

	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// ProgressBar renders an SLDS horizontal progress bar.
// percent is the progress value between 0 and 100.
// It displays an assistive text announcing the progress.
func ProgressBar(percent int) masc.ComponentOrHTML {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	// Percentage string for CSS width and assistive text
	pct := fmt.Sprintf("%d%%", percent)
	// Progress bar container
	return elem.Div(
		masc.Markup(
			masc.Class("slds-progress-bar"),
			masc.Attribute("role", "progressbar"),
			masc.Property("aria-valuemin", 0),
			masc.Property("aria-valuemax", 100),
			masc.Property("aria-valuenow", percent),
		),
		elem.Span(
			masc.Markup(
				masc.Class("slds-progress-bar__value"),
				masc.Style("width", pct),
			),
			// Assistive text for screen readers
			elem.Span(
				masc.Markup(masc.Class("slds-assistive-text")),
				masc.Text(fmt.Sprintf("Progress: %s", pct)),
			),
		),
	)
}
