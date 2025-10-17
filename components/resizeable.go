package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// Resizeable wraps child content and invokes the provided callback whenever the browser
// window size changes. The callback receives the DOM event; the event's Value exposes
// the underlying ResizeEvent. Width and height are available via the `width` and
// `height` properties on the event object (for example: `evt.Get("width").Int()`).
//
// The function returns a masc.Component that renders a <div> wrapping the given content.
func Resizeable(onResize func(*masc.Event), content ...masc.MarkupOrChild) masc.ComponentOrHTML {
	return &resizeableComponent{
		OnResize: onResize,
		Content:  content,
	}
}

type resizeableComponent struct {
	masc.Core

	OnResize       func(*masc.Event)    `masc:"prop"`
	Content        []masc.MarkupOrChild `masc:"prop"`
	detachListener func()
	initialQueued  bool
	lastWidth      int
	lastHeight     int
	hasDimensions  bool
}

func (r *resizeableComponent) Render(send func(masc.Msg)) masc.ComponentOrHTML {
	if r.detachListener == nil {
		r.detachListener = addWindowResizeListener(r.handleResize)
	}
	if !r.initialQueued {
		r.initialQueued = true
		queueViewportMeasurement(r.handleResize)
	}
	return elem.Div(r.Content...)
}

func (r *resizeableComponent) Unmount() {
	if r.detachListener != nil {
		r.detachListener()
		r.detachListener = nil
	}
	r.initialQueued = false
	r.hasDimensions = false
}

func (r *resizeableComponent) handleResize(width, height int, evt *masc.Event) {
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	if r.hasDimensions && width == r.lastWidth && height == r.lastHeight {
		return
	}
	r.lastWidth = width
	r.lastHeight = height
	r.hasDimensions = true

	if r.OnResize != nil {
		event := ensureResizeEvent(evt, width, height)
		r.OnResize(event)
	}
}
