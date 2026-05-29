package components

import (
	"strings"
	"testing"
	"time"

	"github.com/gost-dom/browser/html"
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// datepickerModel renders a single component inside <body> for DOM assertions.
type datepickerModel struct {
	masc.Core
	comp masc.ComponentOrHTML
}

func (m *datepickerModel) Init() masc.Cmd                             { return nil }
func (m *datepickerModel) Update(masc.Msg) (masc.Model, masc.Cmd)     { return m, nil }
func (m *datepickerModel) Render(func(masc.Msg)) masc.ComponentOrHTML { return elem.Body(m.comp) }

// renderComponent renders comp into a headless DOM and returns the window so
// tests can query the resulting elements.
func renderComponent(t *testing.T, comp masc.ComponentOrHTML) html.Window {
	t.Helper()
	win, err := html.NewWindowReader(
		strings.NewReader("<!DOCTYPE html><html><body></body></html>"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := masc.RenderComponentIntoWithSend(win, &datepickerModel{comp: comp}); err != nil {
		t.Fatal(err)
	}
	return win
}

// TestDatepickerBasic verifies that Datepicker returns non-nil for valid parameters.
func TestDatepickerBasic(t *testing.T) {
	comp := Datepicker("Date", time.Time{}, func(t time.Time) {})
	if comp == nil {
		t.Error("Datepicker returned nil for valid parameters")
	}
}

// TestDatepickerWithValue verifies that Datepicker works with a specific date value.
func TestDatepickerWithValue(t *testing.T) {
	date := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	comp := Datepicker("Birth Date", date, func(t time.Time) {})
	if comp == nil {
		t.Error("Datepicker returned nil with date value")
	}
}

// TestDatepickerWithZeroValue verifies that Datepicker works with zero time value.
func TestDatepickerWithZeroValue(t *testing.T) {
	comp := Datepicker("Optional Date", time.Time{}, func(t time.Time) {})
	if comp == nil {
		t.Error("Datepicker returned nil with zero time value")
	}
}

// TestDatepickerWithHandler verifies that Datepicker works with event handler.
func TestDatepickerWithHandler(t *testing.T) {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	handler := func(t time.Time) {}
	comp := Datepicker("Event Date", date, handler)
	if comp == nil {
		t.Error("Datepicker returned nil with event handler")
	}
}

// TestAlignedDatepickerBasic verifies AlignedDatepicker returns non-nil.
func TestAlignedDatepickerBasic(t *testing.T) {
	comp := AlignedDatepicker("Ends", time.Time{}, time.Time{}, time.Time{}, func(time.Time) {})
	if comp == nil {
		t.Error("AlignedDatepicker returned nil for valid parameters")
	}
}

// TestAlignedDatepickerRendersValueAndMax verifies the date input carries the
// value and the max bound, and omits min when unset.
func TestAlignedDatepickerRendersValueAndMax(t *testing.T) {
	value := time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC)
	max := time.Date(2026, 5, 28, 0, 0, 0, 0, time.UTC)
	win := renderComponent(t, AlignedDatepicker("Week Ending", value, time.Time{}, max, func(time.Time) {}))

	node, err := win.Document().QuerySelector("input[type=date]")
	if err != nil {
		t.Fatal(err)
	}
	if node == nil {
		t.Fatal("expected an input[type=date], got none")
	}
	el := node.(html.HTMLElement)
	if got, _ := el.GetAttribute("value"); got != "2026-05-20" {
		t.Errorf("value: got %q, want %q", got, "2026-05-20")
	}
	if got, _ := el.GetAttribute("max"); got != "2026-05-28" {
		t.Errorf("max: got %q, want %q", got, "2026-05-28")
	}
	if _, ok := el.GetAttribute("min"); ok {
		t.Error("min attribute should be absent when min is the zero time")
	}
}

// TestAlignedDatepickerRendersLabelSlot verifies that an empty label still
// produces a label element so the control aligns with labeled siblings.
func TestAlignedDatepickerRendersLabelSlot(t *testing.T) {
	win := renderComponent(t, AlignedDatepicker("", time.Time{}, time.Time{}, time.Time{}, func(time.Time) {}))
	node, err := win.Document().QuerySelector("label.slds-form-element__label")
	if err != nil {
		t.Fatal(err)
	}
	if node == nil {
		t.Fatal("expected an aligned label slot, got none")
	}
}
