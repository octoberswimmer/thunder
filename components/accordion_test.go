package components

import (
	"testing"

	"github.com/octoberswimmer/masc"
)

func TestAccordion(t *testing.T) {
	sections := []AccordionSection{
		{
			ID:       "section1",
			Title:    "Section 1",
			Content:  masc.Text("Content 1"),
			Expanded: true,
		},
		{
			ID:       "section2",
			Title:    "Section 2",
			Content:  masc.Text("Content 2"),
			Expanded: false,
		},
	}

	onToggle := func(id string) {
		// Callback for testing - no action needed here
	}

	result := Accordion(sections, true, onToggle)
	if result == nil {
		t.Error("Expected Accordion to return non-nil component")
	}
}

func TestAccordionSection(t *testing.T) {
	section := AccordionSection{
		ID:       "test-section",
		Title:    "Test Section",
		Content:  masc.Text("Test Content"),
		Expanded: true,
	}

	onToggle := func(id string) {
		// Callback for testing - no action needed here
	}

	result := renderAccordionSection(section, onToggle)
	if result == nil {
		t.Error("Expected renderAccordionSection to return non-nil component")
	}
}

func TestSimpleAccordion(t *testing.T) {
	sections := []AccordionSection{
		{
			ID:       "simple1",
			Title:    "Simple Section",
			Content:  masc.Text("Simple Content"),
			Expanded: false,
		},
	}

	result := SimpleAccordion(sections, nil)
	if result == nil {
		t.Error("Expected SimpleAccordion to return non-nil component")
	}
}

func TestSingleAccordion(t *testing.T) {
	sections := []AccordionSection{
		{
			ID:       "single1",
			Title:    "Single Section",
			Content:  masc.Text("Single Content"),
			Expanded: false,
		},
	}

	result := SingleAccordion(sections, nil)
	if result == nil {
		t.Error("Expected SingleAccordion to return non-nil component")
	}
}
