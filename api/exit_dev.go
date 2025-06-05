//go:build js && dev
// +build js,dev

package api

import (
	"fmt"
)

// ExitApp closes the Thunder application appropriately based on context.
// If running as a quick action modal, it closes the modal.
// Otherwise, it navigates to the record's standard view page.
func ExitApp() {
	fmt.Println("Thunder: ExitApp called (dev mode)")
}

// ExitToRecord navigates to the specified record's standard view page.
// If recordId is empty, it uses the current record context.
func ExitToRecord(recordId string) {
	fmt.Printf("Thunder: ExitToRecord called with recordId: %s (dev mode)\n", recordId)
}

// CloseModal explicitly closes the modal window if running in a quick action context.
func CloseModal() {
	fmt.Println("Thunder: CloseModal called (dev mode)")
}
