package components

import (
	"strings"
	"testing"

	"github.com/gost-dom/browser/html"
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// selectModel is a simple masc.Model wrapper that renders a Select component.
type selectModel struct {
	masc.Core
	opts []SelectOption
	sel  string
}

// Init does nothing.
func (m *selectModel) Init() masc.Cmd { return nil }

// Update does nothing and returns the same model.
func (m *selectModel) Update(msg masc.Msg) (masc.Model, masc.Cmd) { return m, nil }

// Render returns the Select component wrapped in a <body> tag for DOM insertion.
func (m *selectModel) Render(send func(masc.Msg)) masc.ComponentOrHTML {
	// Wrap the <select> inside <body> so RenderIntoNode's root matches document body
	return elem.Body(
		Select("Choose", m.opts, m.sel, nil),
	)
}

// TestSelectBasic verifies that Select returns a non-nil component for valid parameters.
func TestSelectBasic(t *testing.T) {
	opts := []SelectOption{
		{Label: "One", Value: "1"},
		{Label: "Two", Value: "2"},
	}
	comp := Select("Choose an option", opts, "1", func(e *masc.Event) {})
	if comp == nil {
		t.Error("Select returned nil for valid parameters")
	}
}

// TestSelectRetainsSelectedValue verifies that passing various selected values yields a non-nil component.
func TestSelectRetainsSelectedValue(t *testing.T) {
	opts := []SelectOption{
		{Label: "One", Value: "1"},
		{Label: "Two", Value: "2"},
		{Label: "Three", Value: "3"},
	}
	for _, sel := range []string{"1", "2", "3"} {
		comp := Select("Choose", opts, sel, nil)
		if comp == nil {
			t.Errorf("Select returned nil when selected=%s", sel)
		}
	}
}

// TestSelectOptionSelected verifies that the corresponding <option> is rendered as selected in the DOM.
func TestSelectOptionSelected(t *testing.T) {
	opts := []SelectOption{
		{Label: "One", Value: "1"},
		{Label: "Two", Value: "2"},
		{Label: "Three", Value: "3"},
	}
	win, err := html.NewWindowReader(
		strings.NewReader("<!DOCTYPE html><html><body></body></html>"),
	)
	if err != nil {
		t.Fatal(err)
	}
	model := &selectModel{opts: opts, sel: "2"}
	body, _, err := masc.RenderComponentIntoWithSend(win, model)
	if err != nil {
		t.Fatal(err)
	}
	// Query the selected option
	node, err := win.Document().QuerySelector("option[selected]")
	if err != nil {
		t.Fatal(err)
	}
	if node == nil {
		t.Fatal("expected an option[selected] element, got none")
	}
	// Cast to html.HTMLElement to access attributes
	el, ok := node.(html.HTMLElement)
	if !ok {
		t.Fatal("node is not an HTMLElement")
	}
	// Check its value attribute
	value, _ := el.GetAttribute("value")
	if value != "2" {
		t.Errorf("selected attribute on option: got %s, want %s", value, "2")
	}
	_ = body // avoid unused
}
