package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
)

// TabOption represents a single tab choice with label, value, and content.
type TabOption struct {
	Label   string
	Value   string
	Content masc.ComponentOrHTML
}

// Tabs renders an SLDS default-style tabs component.
// id is the identifier prefix for tab and panel IDs.
// options is the slice of TabOption values.
// selected is the currently selected tab value.
// onChange is the callback invoked with the Value of the clicked tab.
func Tabs(id string, options []TabOption, selected string, onChange func(string)) masc.ComponentOrHTML {
	if len(options) == 0 {
		return nil
	}
	// Build navigation list
	var navItems []masc.MarkupOrChild
	navItems = append(navItems,
		masc.Markup(
			masc.Class("slds-tabs_default__nav"),
			masc.Attribute("role", "tablist"),
		),
	)
	for _, opt := range options {
		active := opt.Value == selected
		// Build class list for the tab item
		classList := []string{"slds-tabs_default__item"}
		if active {
			classList = append(classList, "slds-is-active")
		}
		tabID := id + "-" + opt.Value + "__item"
		panelID := id + "-" + opt.Value
		navItems = append(navItems,
			elem.ListItem(
				masc.Markup(
					masc.Class(classList...),
					masc.Property("title", opt.Label),
					masc.Attribute("role", "presentation"),
				),
				elem.Anchor(
					masc.Markup(
						masc.Class("slds-tabs_default__link"),
						masc.Attribute("role", "tab"),
						masc.Property("id", tabID),
						masc.Property("aria-controls", panelID),
						masc.Property("aria-selected", active),
						masc.Property("tabindex", func() string {
							if active {
								return "0"
							}
							return "-1"
						}()),
						event.Click(func(e *masc.Event) {
							if onChange != nil {
								onChange(opt.Value)
							}
						}),
					),
					masc.Text(opt.Label),
				),
			),
		)
	}
	nav := elem.UnorderedList(navItems...)

	// Build content panels
	var panels []masc.MarkupOrChild
	for _, opt := range options {
		active := opt.Value == selected
		panelID := id + "-" + opt.Value
		tabID := panelID + "__item"
		// Build class list for the content panel
		classList := []string{"slds-tabs_default__content"}
		if active {
			classList = append(classList, "slds-show")
		} else {
			classList = append(classList, "slds-hide")
		}
		panels = append(panels,
			elem.Div(
				masc.Markup(
					masc.Class(classList...),
					masc.Attribute("role", "tabpanel"),
					masc.Property("id", panelID),
					masc.Property("aria-labelledby", tabID),
				),
				opt.Content,
			),
		)
	}

	// Wrap in tabs container
	var args []masc.MarkupOrChild
	args = append(args,
		masc.Markup(masc.Class("slds-tabs_default")),
		nav,
	)
	args = append(args, panels...)
	return elem.Div(args...)
}
