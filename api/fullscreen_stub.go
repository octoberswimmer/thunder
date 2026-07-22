//go:build !js
// +build !js

package api

// RequestFullscreen always reports the Fullscreen API unavailable outside the browser (host
// builds / tests), so callers take their fallback path. The browser implementation lives in
// fullscreen.go.
func RequestFullscreen(selector string) bool { return false }

// ExitFullscreen is a no-op outside the browser.
func ExitFullscreen() {}

// NextFullscreenChange reports fullscreen inactive outside the browser. Unlike the browser
// implementation it does not block: there is no fullscreen state to change on the host.
func NextFullscreenChange() bool { return false }
