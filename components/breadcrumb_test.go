package components

import (
	"strings"
	"testing"

	"github.com/gost-dom/browser/html"
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// breadcrumbModel wraps Breadcrumb for rendering into DOM.
type breadcrumbModel struct {
	masc.Core
	opts []BreadcrumbOption
}

func (m *breadcrumbModel) Init() masc.Cmd                             { return nil }
func (m *breadcrumbModel) Update(msg masc.Msg) (masc.Model, masc.Cmd) { return m, nil }
func (m *breadcrumbModel) Render(send func(masc.Msg)) masc.ComponentOrHTML {
	// wrap in <body> for RenderComponentIntoWithSend
	return elem.Body(
		Breadcrumb(m.opts),
	)
}

// TestBreadcrumbBasic verifies that Breadcrumb returns non-nil for valid options.
func TestBreadcrumbBasic(t *testing.T) {
	opts := []BreadcrumbOption{{Label: "Home", Href: "/"}, {Label: "Page", OnClick: func(e *masc.Event) {}}}
	comp := Breadcrumb(opts)
	if comp == nil {
		t.Error("Breadcrumb returned nil for valid options")
	}
}

// TestBreadcrumbRendersListItems verifies that Breadcrumb output contains <li> elements.
func TestBreadcrumbRendersListItems(t *testing.T) {
	win, err := html.NewWindowReader(
		strings.NewReader("<!DOCTYPE html><html><body></body></html>"),
	)
	if err != nil {
		t.Fatal(err)
	}
	model := &breadcrumbModel{opts: []BreadcrumbOption{{Label: "Home", Href: "/"}, {Label: "Page"}}}
	body, _, err := masc.RenderComponentIntoWithSend(win, model)
	if err != nil {
		t.Fatal(err)
	}
	inner := body.InnerHTML()
	if !strings.Contains(inner, "<li") {
		t.Errorf("expected <li> in Breadcrumb output; got %s", inner)
	}
}
