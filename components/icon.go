package components

import (
	"fmt"
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// IconCategory defines SLDS icon categories for sprite lookup.
type IconCategory string

const (
	// UtilityIcon corresponds to the SLDS utility sprite.
	UtilityIcon IconCategory = "utility"
	// ActionIcon corresponds to the SLDS action sprite.
	ActionIcon IconCategory = "action"
	// StandardIcon corresponds to the SLDS standard sprite.
	StandardIcon IconCategory = "standard"
)

// IconSize defines SLDS icon sizes.
type IconSize string

const (
	IconSmall  IconSize = "small"
	IconMedium IconSize = "medium"
	IconLarge  IconSize = "large"
)

// Icon renders an SLDS icon from the given category, name, and size.
// For example: Icon(UtilityIcon, "close", IconSmall).
func Icon(category IconCategory, name string, size IconSize) masc.ComponentOrHTML {
	// Inline <use> referencing injected SVG sprite symbols by ID
	svg := fmt.Sprintf(
		`<svg class="slds-icon slds-icon_%s slds-icon-%s-%s" aria-hidden="true">`+
			`<use xlink:href="#%s"></use></svg>`,
		size,
		category,
		name,
		name,
	)
	// Wrap SVG HTML in a span to insert via innerHTML
	return elem.Span(
		masc.Markup(masc.UnsafeHTML(svg)),
	)
}
