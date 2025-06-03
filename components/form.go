package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// Form renders an SLDS styled form container.
// children are the form fields and other content.
func Form(children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	args := []masc.MarkupOrChild{masc.Markup(masc.Class("slds-form", "slds-form_stacked"))}
	args = append(args, children...)
	return elem.Form(args...)
}

// FormWithAttributes renders an SLDS styled form container with custom attributes.
// The first argument should be markup with attributes, followed by children.
func FormWithAttributes(attributesAndChildren ...masc.MarkupOrChild) masc.ComponentOrHTML {
	args := []masc.MarkupOrChild{masc.Markup(masc.Class("slds-form", "slds-form_stacked"))}
	args = append(args, attributesAndChildren...)
	return elem.Form(args...)
}

// FormSection renders a form section with a heading and grouped fields.
// title is the section heading text.
// children are the form fields within this section.
func FormSection(title string, children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	args := []masc.MarkupOrChild{
		masc.Markup(masc.Class("slds-form-element", "slds-m-top_large")),
		elem.Div(
			masc.Markup(masc.Class("slds-form-element__legend", "slds-text-heading_small")),
			masc.Text(title),
		),
	}
	args = append(args, children...)
	return elem.Div(args...)
}

// FormRow renders form fields in a horizontal row layout.
// children are the form fields to display side by side.
func FormRow(children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	args := []masc.MarkupOrChild{masc.Markup(masc.Class("slds-grid", "slds-gutters", "slds-m-top_medium"))}
	args = append(args, children...)
	return elem.Div(args...)
}

// FormColumn renders a column within a form row.
// size determines the column width (1-12, where 12 is full width).
// Common sizes: 6 (half), 4 (third), 3 (quarter).
// children are the form fields within this column.
func FormColumn(size int, children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	var sizeClass string
	switch size {
	case 1:
		sizeClass = "slds-size_1-of-12"
	case 2:
		sizeClass = "slds-size_1-of-6"
	case 3:
		sizeClass = "slds-size_1-of-4"
	case 4:
		sizeClass = "slds-size_1-of-3"
	case 6:
		sizeClass = "slds-size_1-of-2"
	case 8:
		sizeClass = "slds-size_2-of-3"
	case 9:
		sizeClass = "slds-size_3-of-4"
	case 12:
		sizeClass = "slds-size_1-of-1"
	default:
		// Default to half width if invalid size
		sizeClass = "slds-size_1-of-2"
	}

	args := []masc.MarkupOrChild{masc.Markup(masc.Class("slds-col", sizeClass))}
	args = append(args, children...)
	return elem.Div(args...)
}
