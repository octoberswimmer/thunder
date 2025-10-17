//go:build !js

package components

import "github.com/octoberswimmer/masc"

func addWindowResizeListener(callback func(width, height int, evt *masc.Event)) func() {
	return func() {}
}

func queueViewportMeasurement(callback func(width, height int, evt *masc.Event)) {
	if callback != nil {
		callback(0, 0, ensureResizeEvent(nil, 0, 0))
	}
}

func ensureResizeEvent(evt *masc.Event, width, height int) *masc.Event {
	if evt == nil {
		evt = &masc.Event{}
	}
	if evt.Value == nil {
		evt.Value = masc.NewObject(nil)
	}
	if setter, ok := evt.Value.(interface {
		Set(string, interface{})
	}); ok {
		setter.Set("width", width)
		setter.Set("height", height)
	}
	if evt.Target == nil {
		evt.Target = evt.Value
	}
	return evt
}
