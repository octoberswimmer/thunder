package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// ProgressStep represents a single step in a vertical progress indicator.
type ProgressStep struct {
	Name        string
	IsActive    bool
	IsCompleted bool
}

// VerticalProgress renders an SLDS vertical progress indicator with steps.
// It displays a vertical list of steps with visual indicators for active/completed states.
func VerticalProgress(steps []ProgressStep) masc.ComponentOrHTML {
	var progressItems []masc.MarkupOrChild

	for _, step := range steps {
		var iconCategory IconCategory
		var iconName string
		var textClasses []string
		var itemClasses []string

		if step.IsCompleted {
			iconCategory = UtilityIcon
			iconName = "check"
			textClasses = []string{"slds-text-color_success"}
			itemClasses = []string{"slds-progress__item", "slds-is-completed"}
		} else if step.IsActive {
			iconCategory = UtilityIcon
			iconName = "record"
			textClasses = []string{"slds-text-color_brand", "slds-text-heading_small"}
			itemClasses = []string{"slds-progress__item", "slds-is-active"}
		} else {
			iconCategory = UtilityIcon
			iconName = "record"
			textClasses = []string{"slds-text-color_weak"}
			itemClasses = []string{"slds-progress__item"}
		}

		progressItem := elem.ListItem(
			masc.Markup(masc.Class(itemClasses...)),
			elem.Div(
				masc.Markup(masc.Class("slds-progress__marker")),
				Icon(iconCategory, iconName, IconSmall),
			),
			elem.Div(
				masc.Markup(masc.Class("slds-progress__item_content", "slds-grid", "slds-grid_align-center")),
				elem.Span(
					masc.Markup(masc.Class(textClasses...)),
					masc.Text(step.Name),
				),
			),
		)

		progressItems = append(progressItems, progressItem)
	}

	// Prepare arguments for OrderedList
	olArgs := make([]masc.MarkupOrChild, 0, len(progressItems)+1)
	olArgs = append(olArgs, masc.Markup(masc.Class("slds-progress__list")))
	olArgs = append(olArgs, progressItems...)

	return elem.Div(
		masc.Markup(masc.Class("slds-progress", "slds-progress_vertical")),
		elem.OrderedList(olArgs...),
	)
}
