//go:build js && !dev
// +build js,!dev

package api

import (
	"fmt"
	"syscall/js"

	"github.com/octoberswimmer/thunder"
)

func RecordId() (string, error) {
	// Get the current div element for this Thunder instance
	div := thunder.GetCurrentDiv()
	if div.IsUndefined() {
		return "", fmt.Errorf("thunder instance not initialized")
	}

	// Call the getRecordIdForDiv function to retrieve the recordId from the WeakMap
	getRecordIdForDiv := js.Global().Get("getRecordIdForDiv")
	if getRecordIdForDiv.Type() == js.TypeUndefined {
		return "", fmt.Errorf("getRecordIdForDiv function not available")
	}

	recordId := getRecordIdForDiv.Invoke(div)
	if recordId.Type() == js.TypeUndefined || recordId.Type() == js.TypeNull {
		return "", fmt.Errorf("recordId is not available")
	}
	return recordId.String(), nil
}
