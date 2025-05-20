//go:build !js
// +build !js

package api

// RecordId is a stub implementation for non-WASM builds and will panic if called.
func RecordId() (string, error) {
	panic("api.RecordId is not supported outside the WASM environment")
}
