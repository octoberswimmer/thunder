package components

import (
	"strings"
	"testing"

	"github.com/gost-dom/browser/html"
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

type resizeEventRecorder struct {
	masc.Core
	events []*masc.Event
}

func (m *resizeEventRecorder) Init() masc.Cmd { return nil }

func (m *resizeEventRecorder) Update(msg masc.Msg) (masc.Model, masc.Cmd) { return m, nil }

func (m *resizeEventRecorder) Render(send func(masc.Msg)) masc.ComponentOrHTML {
	wrapper := Resizeable(func(evt *masc.Event) {
		m.events = append(m.events, evt)
	}, masc.Text("content"))
	return elem.Body(wrapper)
}

// TestResizeableReturnsComponent verifies that Resizeable returns a non-nil component.
func TestResizeableReturnsComponent(t *testing.T) {
	comp := Resizeable(nil, masc.Text("demo"))
	if comp == nil {
		t.Fatal("expected Resizeable to return non-nil component")
	}
}

// TestResizeableInvokesCallback ensures that handleResize triggers the supplied callback.
func TestResizeableInvokesCallback(t *testing.T) {
	var received []*masc.Event
	comp := Resizeable(func(evt *masc.Event) {
		received = append(received, evt)
	}, masc.Text("demo"))
	rc, ok := comp.(*resizeableComponent)
	if !ok {
		t.Fatalf("expected *resizeableComponent, got %T", comp)
	}
	rc.handleResize(320, 480, nil)
	if len(received) != 1 {
		t.Fatalf("expected callback to run once; got %d", len(received))
	}
	if width := received[0].Get("width"); width != nil {
		if w := int(width.(interface{ Int() int }).Int()); w != 320 {
			t.Fatalf("expected width 320, got %d", w)
		}
	}
}

// TestResizeableDeduplicates verifies that duplicate dimensions do not trigger the callback.
func TestResizeableDeduplicates(t *testing.T) {
	var count int
	comp := Resizeable(func(evt *masc.Event) {
		count++
	}, masc.Text("demo"))
	rc := comp.(*resizeableComponent)
	rc.handleResize(200, 100, &masc.Event{})
	rc.handleResize(200, 100, &masc.Event{})
	if count != 1 {
		t.Fatalf("expected duplicate resize to be ignored; got %d callbacks", count)
	}
}

// TestResizeableIntegratesWithDOM ensures the component mounts without error in a gost-dom window.
func TestResizeableIntegratesWithDOM(t *testing.T) {
	win, err := html.NewWindowReader(strings.NewReader("<!DOCTYPE html><html><body></body></html>"))
	if err != nil {
		t.Fatalf("failed creating window: %v", err)
	}
	model := &resizeEventRecorder{}
	if _, _, err := masc.RenderComponentIntoWithSend(win, model); err != nil {
		t.Fatalf("render failed: %v", err)
	}
}
