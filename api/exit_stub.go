//go:build !js
// +build !js

package api

// ExitApp closes the Thunder application appropriately based on context.
// If running as a quick action modal, it closes the modal.
// Otherwise, it navigates to the record's standard view page.
func ExitApp() {
	// No-op for non-JS builds
}

// ExitToRecord navigates to the specified record's standard view page.
// If recordId is empty, it uses the current record context.
func ExitToRecord(recordId string) {
	// No-op for non-JS builds
}

// CloseModal explicitly closes the modal window if running in a quick action context.
func CloseModal() {
	// No-op for non-JS builds
}
