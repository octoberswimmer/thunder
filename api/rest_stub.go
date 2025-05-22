//go:build !js
// +build !js

package api

// Get is a stub implementation for non-WASM builds and will panic if called.
func Get(url string) ([]byte, error) {
	panic("api.Get is not supported outside the WASM environment")
}

// Post is a stub implementation for non-WASM builds and will panic if called.
func Post(url string, body []byte) ([]byte, error) {
	panic("api.Post is not supported outside the WASM environment")
}

// Patch is a stub implementation for non-WASM builds and will panic if called.
func Patch(url string, body []byte) ([]byte, error) {
	panic("api.Patch is not supported outside the WASM environment")
}

// Delete is a stub implementation for non-WASM builds and will panic if called.
func Delete(url string) ([]byte, error) {
	panic("api.Delete is not supported outside the WASM environment")
}
