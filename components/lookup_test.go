package components

import "testing"

// TestLookupBasic verifies that Lookup renders for basic parameters.
func TestLookupBasic(t *testing.T) {
	opts := []LookupOption{{Label: "A", Value: "a"}}
	comp := Lookup("Search", opts, "", nil, nil)
	if comp == nil {
		t.Error("Lookup returned nil for valid parameters")
	}
}

// TestLookupEmptySuggestions verifies that Lookup without suggestions still renders.
func TestLookupEmptySuggestions(t *testing.T) {
	comp := Lookup("Search", nil, "", nil, nil)
	if comp == nil {
		t.Error("Lookup returned nil when suggestions empty")
	}
}
