package components

import (
	"fmt"
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
)

// AccordionSection represents a single section in an accordion
type AccordionSection struct {
	ID       string
	Title    string
	Content  masc.ComponentOrHTML
	Expanded bool
}

// Accordion renders an SLDS Accordion component with multiple collapsible sections
// sections is a slice of AccordionSection structs
// allowMultiple determines if multiple sections can be open at once
// onToggle is a callback function that receives the section ID when toggled
func Accordion(sections []AccordionSection, allowMultiple bool, onToggle func(string)) masc.ComponentOrHTML {
	var sectionElements []masc.ComponentOrHTML

	for _, section := range sections {
		sectionElements = append(sectionElements, renderAccordionSection(section, onToggle))
	}

	// Use div container for accordion
	return elem.Div(
		masc.Markup(
			masc.Class("slds-accordion"),
		),
		masc.List(sectionElements),
	)
}

// renderAccordionSection renders a single accordion section
func renderAccordionSection(section AccordionSection, onToggle func(string)) masc.ComponentOrHTML {
	// Build section classes properly - masc.Class requires separate arguments
	var sectionClasses []string
	sectionClasses = append(sectionClasses, "slds-accordion__section")
	if section.Expanded {
		sectionClasses = append(sectionClasses, "slds-is-open")
	}

	summaryID := fmt.Sprintf("accordion-summary-%s", section.ID)
	contentID := fmt.Sprintf("accordion-details-%s", section.ID)

	// Build the section content
	var sectionContent []masc.MarkupOrChild

	// Add the header/summary
	sectionContent = append(sectionContent,
		elem.Div(
			masc.Markup(
				masc.Class("slds-accordion__summary"),
			),
			elem.Heading3(
				masc.Markup(
					masc.Class("slds-accordion__summary-heading"),
				),
				elem.Button(
					masc.Markup(
						masc.Class("slds-button", "slds-button_reset", "slds-accordion__summary-action"),
						masc.Property("type", "button"),
						masc.Property("aria-controls", contentID),
						masc.Property("aria-expanded", fmt.Sprintf("%t", section.Expanded)),
						masc.Property("title", section.Title),
						masc.Property("id", summaryID),
						event.Click(func(e *masc.Event) {
							if onToggle != nil {
								onToggle(section.ID)
							}
						}),
					),
					// Icon for expand/collapse - fallback to Unicode symbols
					elem.Span(
						masc.Markup(
							masc.Class("slds-accordion__summary-action-icon", "slds-button__icon", "slds-button__icon_left"),
							masc.Property("aria-hidden", "true"),
						),
						// Use Unicode chevron symbols - more reliable than SVG
						masc.If(section.Expanded, masc.Text("▼")),
						masc.If(!section.Expanded, masc.Text("▶")),
					),
					// Title text
					elem.Span(
						masc.Markup(
							masc.Class("slds-accordion__summary-content"),
						),
						masc.Text(section.Title),
					),
				),
			),
		),
	)

	// Add the content if expanded
	if section.Expanded {
		sectionContent = append(sectionContent,
			elem.Div(
				masc.Markup(
					masc.Class("slds-accordion__content"),
					masc.Property("id", contentID),
					masc.Property("aria-labelledby", summaryID),
				),
				section.Content,
			),
		)
	}

	// Build the section - add all content directly
	var args []masc.MarkupOrChild
	// Pass each class as a separate argument to masc.Class
	args = append(args, masc.Markup(masc.Class(sectionClasses...)))
	args = append(args, sectionContent...)

	return elem.Div(args...)
}

// SimpleAccordion creates an accordion with a single toggle callback for all sections
// Useful when you want to manage state externally
func SimpleAccordion(sections []AccordionSection, onToggle func(string)) masc.ComponentOrHTML {
	return Accordion(sections, true, onToggle)
}

// SingleAccordion creates an accordion where only one section can be open at a time
func SingleAccordion(sections []AccordionSection, onToggle func(string)) masc.ComponentOrHTML {
	return Accordion(sections, false, onToggle)
}
