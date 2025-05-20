//go:build dev
// +build dev

package api

import (
	"fmt"
	"syscall/js"
)

// RecordId returns the current Lightning record page ID from the URL's query parameters.
// It returns an error if the recordId parameter is not present.
func RecordId() (string, error) {
	rid := js.Global().Get("URL").New(js.Global().Get("location").Get("href")).Get("searchParams").Call("get", "recordId")
	if rid.Type() == js.TypeUndefined || rid.Type() == js.TypeNull || rid.String() == "" {
		return "", fmt.Errorf("recordId is not available")
	}
	return rid.String(), nil
}
