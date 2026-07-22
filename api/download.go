//go:build js
// +build js

package api

import "syscall/js"

// Download triggers a browser download of data as a file named filename with the given MIME type.
// It builds a Blob and clicks a temporary object-URL anchor — a purely client-side browser action,
// so it works the same in dev and production and needs no Salesforce round-trip.
//
// The anchor is deliberately never attached to the DOM. Experience Cloud injects a document-level
// click interceptor into community Visualforce pages that hijacks anchor navigation and rewrites
// the blob: URL into a broken /s/sfdcpage/ community route; clicking a detached anchor keeps the
// event from ever reaching that interceptor while still triggering the download.
func Download(filename, mimeType string, data []byte) {
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	u8 := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(u8, data)
	parts := js.Global().Get("Array").New(1)
	parts.SetIndex(0, u8)
	blob := js.Global().Get("Blob").New(parts, map[string]interface{}{"type": mimeType})
	url := js.Global().Get("URL").Call("createObjectURL", blob)
	defer js.Global().Get("URL").Call("revokeObjectURL", url)

	doc := js.Global().Get("document")
	a := doc.Call("createElement", "a")
	a.Set("href", url)
	a.Set("download", filename)
	click := js.Global().Get("MouseEvent").New("click")
	a.Call("dispatchEvent", click)
}
