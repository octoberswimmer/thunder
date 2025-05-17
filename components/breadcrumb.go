package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
)

// BreadcrumbOption represents one breadcrumb item.
// If Href is non-empty, renders a link; if OnClick is provided, renders a button;
// otherwise, renders plain text for the current page.
type BreadcrumbOption struct {
	Label   string
	Href    string
	OnClick func(*masc.Event)
}

// Breadcrumb renders an SLDS breadcrumb navigation.
// Example:
//
//	Breadcrumb([]BreadcrumbOption{
//	  {Label: "Home", Href: "/"},
//	  {Label: "Library", Href: "/library"},
//	  {Label: "Data", OnClick: handleClick},
//	})
func Breadcrumb(opts []BreadcrumbOption) masc.ComponentOrHTML {
	// Build list items
	var items []masc.MarkupOrChild
	for _, opt := range opts {
		// Determine link or text
		var child masc.ComponentOrHTML
		switch {
		case opt.Href != "":
			child = elem.Anchor(
				masc.Markup(masc.Class("slds-breadcrumb__link"), masc.Property("href", opt.Href)),
				masc.Text(opt.Label),
			)
		case opt.OnClick != nil:
			child = elem.Button(
				masc.Markup(
					masc.Class("slds-breadcrumb__link"),
					event.Click(opt.OnClick),
				),
				masc.Text(opt.Label),
			)
		default:
			child = elem.Span(
				masc.Markup(masc.Class("slds-breadcrumb__link")),
				masc.Text(opt.Label),
			)
		}
		// Wrap in list item
		li := elem.ListItem(
			masc.Markup(masc.Class("slds-breadcrumb__item"), masc.Property("role", "presentation")),
			child,
		)
		items = append(items, li)
	}
	// Ordered list for breadcrumb
	// Prepare markup and items for variadic call
	olArgs := make([]masc.MarkupOrChild, 0, len(items)+1)
	// SLDS breadcrumb list styling
	olArgs = append(olArgs,
		masc.Markup(
			// SLDS breadcrumb list styling
			masc.Class("slds-breadcrumb", "slds-list_horizontal"),
			masc.Property("role", "list"),
		),
	)
	olArgs = append(olArgs, items...)
	ol := elem.OrderedList(olArgs...)
	// Navigation container with SLDS breadcrumb styling
	nav := elem.Navigation(
		masc.Markup(
			// SLDS nav wrapper for breadcrumbs
			masc.Class("slds-breadcrumbs", "slds-m-bottom_medium"),
			masc.Property("aria-label", "Breadcrumbs"),
		),
		ol,
	)
	return nav
}
