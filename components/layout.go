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
	return CenteredContainer(
		Spinner(size),
	)
}

// Section creates a container with top and bottom margins for sectioning content.
func Section(size SpaceSize, children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	return Spacer(SpaceOptions{MarginVertical: size}, children...)
}

// CenteredContainer creates a centered content container with text alignment.
func CenteredContainer(children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	args := []masc.MarkupOrChild{
		masc.Markup(masc.Class("slds-text-align_center")),
	}
	args = append(args, children...)
	return elem.Div(args...)
}

// ContentContainer creates a container with medium padding for page content.
func ContentContainer(children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	return PaddingAround(SpaceMedium, children...)
}

// ButtonGroup creates a container for grouping action buttons with proper spacing.
func ButtonGroup(children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	args := []masc.MarkupOrChild{
		masc.Markup(masc.Class("slds-button-group")),
		masc.Markup(masc.Attribute("role", "group")),
	}
	args = append(args, children...)
	return elem.Div(args...)
}

// ButtonSpacer creates spacing between buttons in a button group.
func ButtonSpacer() masc.ComponentOrHTML {
	return Spacer(SpaceOptions{MarginLeft: SpaceSmall})
}

// ActionContainer creates a container for action buttons with bottom margin.
func ActionContainer(children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	return MarginBottom(SpaceMedium, children...)
}

// CenteredSpinner creates a centered container with a spinner.
func CenteredSpinner(size string) masc.ComponentOrHTML {
	return CenteredContainer(
		Spinner(size),
	)
}

// Container creates a simple div container.
func Container(children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	return elem.Div(children...)
}

// TableElement creates a styled table element with SLDS classes.
func TableElement(children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	args := []masc.MarkupOrChild{
		masc.Markup(masc.Class("slds-table", "slds-table_cell-buffer", "slds-table_bordered", "slds-table_striped")),
	}
	args = append(args, children...)
	return elem.Table(args...)
}

// TruncatedCell creates a table cell content with truncation styling.
func TruncatedCell(content string, title ...string) masc.ComponentOrHTML {
	cellTitle := content
	if len(title) > 0 && title[0] != "" {
		cellTitle = title[0]
	}
	return elem.Div(
		masc.Markup(masc.Class("slds-truncate"), masc.Property("title", cellTitle)),
		masc.Text(content),
	)
}

// TruncatedCellDiv creates a truncated cell container for custom content.
func TruncatedCellDiv(children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	args := []masc.MarkupOrChild{masc.Markup(masc.Class("slds-truncate"))}
	args = append(args, children...)
	return elem.Div(args...)
}

// TableHeaderRow creates a table header row with proper styling.
func TableHeaderRow(children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	args := []masc.MarkupOrChild{masc.Markup(masc.Class("slds-line-height_reset"))}
	args = append(args, children...)
	return elem.TableRow(args...)
}

// TableDataRow creates a table data row with proper styling.
func TableDataRow(children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	args := []masc.MarkupOrChild{masc.Markup(masc.Class("slds-hint-parent"))}
	args = append(args, children...)
	return elem.TableRow(args...)
}

// PageContentWrapper creates a page content container with proper SLDS styling.
// This replaces the common pattern of manual div creation with SLDS classes.
func PageContentWrapper(children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	args := []masc.MarkupOrChild{
		masc.Markup(masc.Class("slds-page-content", "slds-p-around_medium")),
	}
	args = append(args, children...)
	return elem.Div(args...)
}

// DisabledOverlay creates a disabled overlay for saving states.
// This provides a consistent overlay pattern used in forms.
func DisabledOverlay() masc.ComponentOrHTML {
	return elem.Div(
		masc.Markup(
			masc.Class("slds-backdrop", "slds-backdrop_open"),
			masc.Attribute("style", "position: absolute; top: 0; left: 0; right: 0; bottom: 0; z-index: 9000; background-color: rgba(255, 255, 255, 0.5); pointer-events: auto; cursor: not-allowed;"),
		),
	)
}

// RelativeContainer creates a relatively positioned container.
// Useful for positioning overlays and spinners.
func RelativeContainer(children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	args := []masc.MarkupOrChild{
		masc.Markup(masc.Class("slds-is-relative")),
	}
	args = append(args, children...)
	return elem.Div(args...)
}

// LoadingOverlay creates a full overlay with centered spinner for loading states.
func LoadingOverlay() masc.ComponentOrHTML {
	return elem.Div(
		masc.Markup(
			masc.Class("slds-backdrop", "slds-backdrop_open"),
			masc.Attribute("style", "position: fixed; top: 0; left: 0; right: 0; bottom: 0; z-index: 9000; background-color: rgba(255, 255, 255, 0.8); display: flex; align-items: center; justify-content: center;"),
		),
		Spinner("large"),
	)
}
