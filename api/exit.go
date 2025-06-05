//go:build js && !dev
// +build js,!dev

package api

import (
	"syscall/js"
)

// ExitApp closes the Thunder application appropriately based on context.
// If running as a quick action modal, it closes the modal.
// Otherwise, it navigates to the record's standard view page.
func ExitApp() {
	js.Global().Call("thunderExit")
}

// ExitToRecord navigates to the specified record's standard view page.
// If recordId is empty, it uses the current record context.
func ExitToRecord(recordId string) {
	if recordId == "" {
		recordId = js.Global().Get("recordId").String()
	}
	js.Global().Call("thunderExitToRecord", recordId)
}

// CloseModal explicitly closes the modal window if running in a quick action context.
func CloseModal() {
	js.Global().Call("thunderCloseModal")
}
