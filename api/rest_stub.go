//go:build !js
// +build !js

package api

// Get is a stub implementation for non-WASM builds and will panic if called.
func Get(url string) []byte {
	panic("api.Get is not supported outside the WASM environment")
}
