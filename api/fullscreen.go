//go:build js
// +build js

package api

import "syscall/js"

// RequestFullscreen asks the browser to display the first element matching selector fullscreen.
// It returns false when the Fullscreen API is unavailable — most notably inside an embedding
// iframe without the allowfullscreen permission, like Experience Cloud's Visualforce wrapper —
// or when no element matches, so callers can fall back to a CSS takeover.
func RequestFullscreen(selector string) bool {
	doc := js.Global().Get("document")
	if !doc.Get("fullscreenEnabled").Truthy() {
		return false
	}
	el := doc.Call("querySelector", selector)
	if !el.Truthy() || !el.Get("requestFullscreen").Truthy() {
		return false
	}
	el.Call("requestFullscreen")
	return true
}

// ExitFullscreen leaves fullscreen if the document is in it.
func ExitFullscreen() {
	doc := js.Global().Get("document")
	if doc.Get("fullscreenElement").Truthy() {
		doc.Call("exitFullscreen")
	}
}

// NextFullscreenChange blocks until the document's fullscreen state next changes and reports
// whether fullscreen is then active. Run it from a masc Cmd goroutine to translate Esc and other
// browser-initiated exits into a message the app's Update can react to.
func NextFullscreenChange() bool {
	done := make(chan bool, 1)
	doc := js.Global().Get("document")
	var cb js.Func
	cb = js.FuncOf(func(js.Value, []js.Value) interface{} {
		doc.Call("removeEventListener", "fullscreenchange", cb)
		cb.Release()
		done <- doc.Get("fullscreenElement").Truthy()
		return nil
	})
	doc.Call("addEventListener", "fullscreenchange", cb)
	return <-done
}
