package components

import (
	"strings"
	"testing"

	"github.com/gost-dom/browser/html"
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// iconModel is a simple masc.Model wrapper that renders an Icon component.
type iconModel struct {
	masc.Core
}

func (m *iconModel) Init() masc.Cmd                             { return nil }
func (m *iconModel) Update(msg masc.Msg) (masc.Model, masc.Cmd) { return m, nil }
func (m *iconModel) Render(send func(masc.Msg)) masc.ComponentOrHTML {
	// Wrap in <body> for RenderComponentIntoWithSend
	return elem.Body(
		Icon(UtilityIcon, "close", IconSmall),
	)
}

// TestIconBasic verifies that Icon returns a non-nil component for valid parameters.
func TestIconBasic(t *testing.T) {
	comp := Icon(UtilityIcon, "close", IconSmall)
	if comp == nil {
		t.Error("Icon returned nil for valid parameters")
	}
}

// TestIconRendersSVG verifies that the Icon component includes an <svg> element in the rendered output.
func TestIconRendersSVG(t *testing.T) {
	win, err := html.NewWindowReader(
		strings.NewReader("<!DOCTYPE html><html><body></body></html>"),
	)
	if err != nil {
		t.Fatal(err)
	}
	model := &iconModel{}
	body, _, err := masc.RenderComponentIntoWithSend(win, model)
	if err != nil {
		t.Fatal(err)
	}
	inner := body.InnerHTML()
	if !strings.Contains(inner, "<svg") {
		t.Errorf("expected <svg> in Icon output; got %s", inner)
	}
}
