package components

import (
	"testing"
)

// TestValidatedLookup verifies that ValidatedLookup returns non-nil for valid parameters.
func TestValidatedLookup(t *testing.T) {
	options := []LookupOption{
		{Label: "Option 1", Value: "opt1"},
		{Label: "Option 2", Value: "opt2"},
	}
	validation := ValidationState{Required: false}

	comp := ValidatedLookup("Test Label", options, "test", validation,
		func(string) {}, func(string) {}, func() string { return "test" })
	if comp == nil {
		t.Error("ValidatedLookup returned nil for valid parameters")
	}
}

// TestValidatedLookupWithError verifies that ValidatedLookup returns non-nil when validation has errors.
func TestValidatedLookupWithError(t *testing.T) {
	options := []LookupOption{
		{Label: "Option 1", Value: "opt1"},
	}
	validation := ValidationState{
		Required:     false,
		HasError:     true,
		ErrorMessage: "Test error",
	}

	comp := ValidatedLookup("Test Label", options, "test", validation,
		func(string) {}, func(string) {}, func() string { return "test" })
	if comp == nil {
		t.Error("ValidatedLookup returned nil for parameters with validation errors")
	}
}

// TestValidatedLookupEmpty verifies that ValidatedLookup returns non-nil for empty parameters.
func TestValidatedLookupEmpty(t *testing.T) {
	validation := ValidationState{Required: false}

	comp := ValidatedLookup("", []LookupOption{}, "", validation,
		func(string) {}, func(string) {}, func() string { return "" })
	if comp == nil {
		t.Error("ValidatedLookup returned nil for empty parameters")
	}
}

// TestValidatedLookupWithReset verifies that ValidatedLookup accepts the onReset callback.
func TestValidatedLookupWithReset(t *testing.T) {
	options := []LookupOption{
		{Label: "Selected Option", Value: "selected"},
	}
	validation := ValidationState{Required: false}
	resetValue := "Selected Option"

	// Test that onReset callback is accepted and component renders
	comp := ValidatedLookup("Test Label", options, "current input", validation,
		func(string) {}, // onInput
		func(string) {}, // onSelect
		func() string { // onReset
			return resetValue
		})

	if comp == nil {
		t.Error("ValidatedLookup returned nil when onReset callback provided")
	}
}
