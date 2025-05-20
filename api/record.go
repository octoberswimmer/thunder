//go:build js && !dev
// +build js,!dev

package api

import (
	"fmt"
	"syscall/js"
)

func RecordId() (string, error) {
	recordId := js.Global().Get("recordId")
	if recordId.Type() == js.TypeUndefined {
		return "", fmt.Errorf("recordId is not available")
	}
	return recordId.String(), nil
}
