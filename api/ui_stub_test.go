package api

import "testing"

// TestGetPicklistValuesByRecordType_stub_panics verifies that calling GetPicklistValuesByRecordType in stub mode panics.
func TestGetPicklistValuesByRecordType_stub_panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when calling GetPicklistValuesByRecordType in stub version")
		}
	}()

	GetPicklistValuesByRecordType("", "")
}

// TestGetObjectInfo_stub_panics verifies that calling GetObjectInfo in stub mode panics.
func TestGetObjectInfo_stub_panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when calling GetObjectInfo in stub version")
		}
	}()

	GetObjectInfo("")
}
