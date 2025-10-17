//go:build js

package components

import (
	"syscall/js"

	"github.com/octoberswimmer/masc"
)

func addWindowResizeListener(callback func(width, height int, evt *masc.Event)) func() {
	if callback == nil {
		return func() {}
	}
	global := js.Global()
	if !global.Truthy() {
		return func() {}
	}
	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		width, height := viewportSize(global)
		event := buildResizeEvent(args, global, width, height)
		callback(width, height, event)
		return nil
	})
	global.Call("addEventListener", "resize", handler)
	return func() {
		global.Call("removeEventListener", "resize", handler)
		handler.Release()
	}
}

func queueViewportMeasurement(callback func(width, height int, evt *masc.Event)) {
	if callback == nil {
		return
	}
	global := js.Global()
	if !global.Truthy() {
		callback(0, 0, ensureResizeEvent(nil, 0, 0))
		return
	}
	if raf := global.Get("requestAnimationFrame"); raf.Truthy() {
		var fn js.Func
		fn = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			fn.Release()
			width, height := viewportSize(global)
			callback(width, height, ensureResizeEvent(nil, width, height))
			return nil
		})
		global.Call("requestAnimationFrame", fn)
		return
	}
	width, height := viewportSize(global)
	callback(width, height, ensureResizeEvent(nil, width, height))
}

func viewportSize(global js.Value) (int, int) {
	width := 0
	height := 0
	if inner := global.Get("innerWidth"); inner.Truthy() {
		width = inner.Int()
	}
	if inner := global.Get("innerHeight"); inner.Truthy() {
		height = inner.Int()
	}
	return width, height
}

func buildResizeEvent(args []js.Value, global js.Value, width, height int) *masc.Event {
	var evt *masc.Event
	if len(args) > 0 && args[0].Truthy() {
		raw := args[0]
		target := raw.Get("target")
		if !target.Truthy() {
			target = global
		}
		evt = &masc.Event{Value: raw, Target: target}
	}
	return ensureResizeEvent(evt, width, height)
}

func ensureResizeEvent(evt *masc.Event, width, height int) *masc.Event {
	if evt == nil {
		raw := masc.NewObject(map[string]interface{}{
			"width":  width,
			"height": height,
		})
		return &masc.Event{
			Value:  raw,
			Target: raw,
		}
	}
	value := evt.Value
	if !value.Truthy() {
		value = masc.NewObject(nil)
	}
	value.Set("width", width)
	value.Set("height", height)
	evt.Value = value
	if !evt.Target.Truthy() {
		evt.Target = value
	}
	return evt
}
