package api

import "testing"

// TestPost_stub_panics verifies that calling Post in stub mode panics.
func TestPost_stub_panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when calling Post in stub version")
		}
	}()
	Post("", nil)
}

// TestPatch_stub_panics verifies that calling Patch in stub mode panics.
func TestPatch_stub_panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when calling Patch in stub version")
		}
	}()
	Patch("", nil)
}

// TestDelete_stub_panics verifies that calling Delete in stub mode panics.
func TestDelete_stub_panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when calling Delete in stub version")
		}
	}()
	Delete("")
}
