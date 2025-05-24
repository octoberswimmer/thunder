package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// TextSize defines text size variants
type TextSize string

const (
	TextSmall   TextSize = "slds-text-body_small"
	TextRegular TextSize = "slds-text-body_regular"
	TextLarge   TextSize = "slds-text-body_large"
)

// HeadingSize defines heading size variants
type HeadingSize string

const (
	HeadingSmall  HeadingSize = "slds-text-heading_small"
	HeadingMedium HeadingSize = "slds-text-heading_medium"
	HeadingLarge  HeadingSize = "slds-text-heading_large"
)

// Text renders styled text content.
func Text(content string, size ...TextSize) masc.ComponentOrHTML {
	s := TextRegular
	if len(size) > 0 {
		s = size[0]
	}
	return elem.Span(
		masc.Markup(masc.Class(string(s))),
		masc.Text(content),
	)
}

// Paragraph renders a paragraph with optional text size.
func Paragraph(content string, size ...TextSize) masc.ComponentOrHTML {
	s := TextRegular
	if len(size) > 0 {
		s = size[0]
	}
	return elem.Paragraph(
		masc.Markup(masc.Class(string(s))),
		masc.Text(content),
	)
}

// Heading renders a heading with specified size.
func Heading(content string, size HeadingSize) masc.ComponentOrHTML {
	switch size {
	case HeadingLarge:
		return elem.Heading1(
			masc.Markup(masc.Class(string(size))),
			masc.Text(content),
		)
	case HeadingMedium:
		return elem.Heading2(
			masc.Markup(masc.Class(string(size))),
			masc.Text(content),
		)
	default: // HeadingSmall
		return elem.Heading3(
			masc.Markup(masc.Class(string(size))),
			masc.Text(content),
		)
	}
}
