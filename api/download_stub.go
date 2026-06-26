//go:build !js
// +build !js

package api

// Download is a no-op outside the browser (host builds / tests). The browser implementation lives
// in download.go.
func Download(filename, mimeType string, data []byte) {}
