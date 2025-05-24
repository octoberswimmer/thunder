package components

import (
	"fmt"
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// GridAlignment defines how grid content is aligned.
type GridAlignment string

const (
	// AlignStart aligns content to the start (left) of the grid.
	AlignStart GridAlignment = ""
	// AlignCenter centers content in the grid.
	AlignCenter GridAlignment = "slds-align_absolute-center"
	// AlignEnd aligns content to the end (right) of the grid.
	AlignEnd GridAlignment = "slds-grid_align-end"
)

// GridOptions configures grid layout and alignment.
type GridOptions struct {
	Alignment GridAlignment
	Wrap      bool
	Gutters   bool
}

// Grid renders an SLDS grid container with gutters for spacing.
// Arrange child columns using GridColumn components.
func Grid(children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	return GridWithOptions(GridOptions{
		Wrap:    true,
		Gutters: true,
	}, children...)
}

// GridWithOptions renders an SLDS grid container with custom options.
func GridWithOptions(options GridOptions, children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	if len(children) == 0 {
		return nil
	}

	classes := []string{"slds-grid"}

	if options.Wrap {
		classes = append(classes, "slds-wrap")
	}

	if options.Gutters {
		classes = append(classes, "slds-gutters")
	}

	if options.Alignment != "" {
		classes = append(classes, string(options.Alignment))
	}

	args := []masc.MarkupOrChild{
		masc.Markup(masc.Class(classes...)),
	}
	args = append(args, children...)
	return elem.Div(args...)
}

// CenteredGrid creates a grid with center alignment.
func CenteredGrid(children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	return GridWithOptions(GridOptions{
		Alignment: AlignCenter,
		Wrap:      true,
		Gutters:   true,
	}, children...)
}

// GridColumn renders an SLDS grid column.
// size is the SLDS sizing string (e.g. "1-of-2" yields the class "slds-size_1-of-2").
// children are nested content within the column.
func GridColumn(size string, children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	classes := []string{"slds-col"}
	if size != "" {
		classes = append(classes, fmt.Sprintf("slds-size_%s", size))
	}
	args := []masc.MarkupOrChild{
		masc.Markup(masc.Class(classes...)),
	}
	args = append(args, children...)
	return elem.Div(args...)
}
