package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// SpaceSize defines spacing size variants
type SpaceSize string

const (
	SpaceNone   SpaceSize = "none"
	SpaceXSmall SpaceSize = "x-small"
	SpaceSmall  SpaceSize = "small"
	SpaceMedium SpaceSize = "medium"
	SpaceLarge  SpaceSize = "large"
	SpaceXLarge SpaceSize = "x-large"
)

// SpaceOptions configures margin and padding for a layout container
type SpaceOptions struct {
	MarginTop         SpaceSize
	MarginBottom      SpaceSize
	MarginLeft        SpaceSize
	MarginRight       SpaceSize
	MarginVertical    SpaceSize
	MarginHorizontal  SpaceSize
	PaddingTop        SpaceSize
	PaddingBottom     SpaceSize
	PaddingLeft       SpaceSize
	PaddingRight      SpaceSize
	PaddingVertical   SpaceSize
	PaddingHorizontal SpaceSize
	PaddingAround     SpaceSize
}

// Spacer creates a container with specified spacing options.
func Spacer(options SpaceOptions, children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	var classes []string

	// Margin classes
	if options.MarginTop != "" {
		classes = append(classes, "slds-m-top_"+string(options.MarginTop))
	}
	if options.MarginBottom != "" {
		classes = append(classes, "slds-m-bottom_"+string(options.MarginBottom))
	}
	if options.MarginLeft != "" {
		classes = append(classes, "slds-m-left_"+string(options.MarginLeft))
	}
	if options.MarginRight != "" {
		classes = append(classes, "slds-m-right_"+string(options.MarginRight))
	}
	if options.MarginVertical != "" {
		classes = append(classes, "slds-m-vertical_"+string(options.MarginVertical))
	}
	if options.MarginHorizontal != "" {
		classes = append(classes, "slds-m-horizontal_"+string(options.MarginHorizontal))
	}

	// Padding classes
	if options.PaddingTop != "" {
		classes = append(classes, "slds-p-top_"+string(options.PaddingTop))
	}
	if options.PaddingBottom != "" {
		classes = append(classes, "slds-p-bottom_"+string(options.PaddingBottom))
	}
	if options.PaddingLeft != "" {
		classes = append(classes, "slds-p-left_"+string(options.PaddingLeft))
	}
	if options.PaddingRight != "" {
		classes = append(classes, "slds-p-right_"+string(options.PaddingRight))
	}
	if options.PaddingVertical != "" {
		classes = append(classes, "slds-p-vertical_"+string(options.PaddingVertical))
	}
	if options.PaddingHorizontal != "" {
		classes = append(classes, "slds-p-horizontal_"+string(options.PaddingHorizontal))
	}
	if options.PaddingAround != "" {
		classes = append(classes, "slds-p-around_"+string(options.PaddingAround))
	}

	if len(classes) == 0 {
		return Container(children...)
	}

	args := []masc.MarkupOrChild{
		masc.Markup(masc.Class(classes...)),
	}
	args = append(args, children...)
	return elem.Div(args...)
}

// MarginTop creates a container with top margin.
func MarginTop(size SpaceSize, children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	return Spacer(SpaceOptions{MarginTop: size}, children...)
}

// MarginBottom creates a container with bottom margin.
func MarginBottom(size SpaceSize, children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	return Spacer(SpaceOptions{MarginBottom: size}, children...)
}

// PaddingHorizontal creates a container with horizontal padding.
func PaddingHorizontal(size SpaceSize, children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	return Spacer(SpaceOptions{PaddingHorizontal: size}, children...)
}

// PaddingAround creates a container with padding on all sides.
func PaddingAround(size SpaceSize, children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	return Spacer(SpaceOptions{PaddingAround: size}, children...)
}

// LoadingSpinner creates a centered loading spinner container.
func LoadingSpinner(size string) masc.ComponentOrHTML {
	return CenteredGrid(
		GridColumn("",
			Spinner(size),
		),
	)
}
