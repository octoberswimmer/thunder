package components

import (
	"fmt"
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// Grid renders an SLDS grid container with gutters for spacing.
// Arrange child columns using GridColumn components.
func Grid(children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	if len(children) == 0 {
		return nil
	}
	args := []masc.MarkupOrChild{
		masc.Markup(masc.Class("slds-grid", "slds-wrap", "slds-gutters")),
	}
	args = append(args, children...)
	return elem.Div(args...)
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
