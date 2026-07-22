//go:build !js
// +build !js

package api

import "testing"

func Test_request_fullscreen_reports_unavailable_on_host(t *testing.T) {
	if RequestFullscreen(".anything") {
		t.Error("expected RequestFullscreen to report the Fullscreen API unavailable on the host")
	}
}

func Test_next_fullscreen_change_reports_inactive_without_blocking_on_host(t *testing.T) {
	if NextFullscreenChange() {
		t.Error("expected NextFullscreenChange to report fullscreen inactive on the host")
	}
}
