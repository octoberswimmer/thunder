package api

import "testing"

// TestRecordId_stub_panics verifies that calling RecordId in stub mode panics.
func TestRecordId_stub_panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when calling RecordId in stub version")
		}
	}()

	RecordId()
}
