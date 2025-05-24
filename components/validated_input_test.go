package components

import (
	"github.com/octoberswimmer/masc"
	"testing"
)

// TestValidatedTextInputBasic verifies basic ValidatedTextInput functionality.
func TestValidatedTextInputBasic(t *testing.T) {
	validation := ValidationState{
		HasError:     false,
		Required:     false,
		ErrorMessage: "",
		HelpText:     "",
	}
	if comp := ValidatedTextInput("Test Label", "Test Value", "Test Placeholder", validation, nil); comp == nil {
		t.Error("ValidatedTextInput returned nil for basic parameters")
	}
}

// TestValidatedTextInputRequired verifies required field rendering.
func TestValidatedTextInputRequired(t *testing.T) {
	validation := ValidationState{
		HasError:     false,
		Required:     true,
		ErrorMessage: "",
		HelpText:     "This field is required",
	}
	if comp := ValidatedTextInput("Required Field", "", "", validation, nil); comp == nil {
		t.Error("ValidatedTextInput returned nil for required field")
	}
}

// TestValidatedTextInputWithError verifies error state rendering.
func TestValidatedTextInputWithError(t *testing.T) {
	validation := ValidationState{
		HasError:     true,
		Required:     true,
		ErrorMessage: "This field is required",
		HelpText:     "Help text should not show when error is present",
	}
	if comp := ValidatedTextInput("Error Field", "", "", validation, nil); comp == nil {
		t.Error("ValidatedTextInput returned nil for error state")
	}
}

// TestValidatedTextInputWithHandler verifies event handler integration.
func TestValidatedTextInputWithHandler(t *testing.T) {
	validation := ValidationState{
		HasError:     false,
		Required:     false,
		ErrorMessage: "",
		HelpText:     "Enter some text",
	}
	handler := func(e *masc.Event) {}
	if comp := ValidatedTextInput("Handler Test", "value", "placeholder", validation, handler); comp == nil {
		t.Error("ValidatedTextInput returned nil with event handler")
	}
}

// TestValidatedTextareaBasic verifies basic ValidatedTextarea functionality.
func TestValidatedTextareaBasic(t *testing.T) {
	validation := ValidationState{
		HasError:     false,
		Required:     false,
		ErrorMessage: "",
		HelpText:     "",
	}
	if comp := ValidatedTextarea("Test Label", "Test Value", "Test Placeholder", 3, validation, nil); comp == nil {
		t.Error("ValidatedTextarea returned nil for basic parameters")
	}
}

// TestValidatedTextareaRequired verifies required textarea field rendering.
func TestValidatedTextareaRequired(t *testing.T) {
	validation := ValidationState{
		HasError:     false,
		Required:     true,
		ErrorMessage: "",
		HelpText:     "This field is required",
	}
	if comp := ValidatedTextarea("Required Textarea", "", "", 5, validation, nil); comp == nil {
		t.Error("ValidatedTextarea returned nil for required field")
	}
}

// TestValidatedTextareaWithError verifies textarea error state rendering.
func TestValidatedTextareaWithError(t *testing.T) {
	validation := ValidationState{
		HasError:     true,
		Required:     true,
		ErrorMessage: "This field is required",
		HelpText:     "Help text should not show when error is present",
	}
	if comp := ValidatedTextarea("Error Textarea", "", "", 2, validation, nil); comp == nil {
		t.Error("ValidatedTextarea returned nil for error state")
	}
}

// TestValidatedSelectBasic verifies basic ValidatedSelect functionality.
func TestValidatedSelectBasic(t *testing.T) {
	options := []SelectOption{
		{Label: "Option 1", Value: "opt1"},
		{Label: "Option 2", Value: "opt2"},
	}
	validation := ValidationState{
		HasError:     false,
		Required:     false,
		ErrorMessage: "",
		HelpText:     "",
	}
	if comp := ValidatedSelect("Test Select", options, "opt1", validation, nil); comp == nil {
		t.Error("ValidatedSelect returned nil for basic parameters")
	}
}

// TestValidatedSelectRequired verifies required select field rendering.
func TestValidatedSelectRequired(t *testing.T) {
	options := []SelectOption{
		{Label: "Choose...", Value: ""},
		{Label: "Option 1", Value: "opt1"},
	}
	validation := ValidationState{
		HasError:     false,
		Required:     true,
		ErrorMessage: "",
		HelpText:     "Please select an option",
	}
	if comp := ValidatedSelect("Required Select", options, "", validation, nil); comp == nil {
		t.Error("ValidatedSelect returned nil for required field")
	}
}

// TestValidatedSelectWithError verifies select error state rendering.
func TestValidatedSelectWithError(t *testing.T) {
	options := []SelectOption{
		{Label: "Option 1", Value: "opt1"},
	}
	validation := ValidationState{
		HasError:     true,
		Required:     true,
		ErrorMessage: "Please select a valid option",
		HelpText:     "Help text hidden when error present",
	}
	if comp := ValidatedSelect("Error Select", options, "", validation, nil); comp == nil {
		t.Error("ValidatedSelect returned nil for error state")
	}
}

// TestValidatedDatepickerBasic verifies basic ValidatedDatepicker functionality.
func TestValidatedDatepickerBasic(t *testing.T) {
	validation := ValidationState{
		HasError:     false,
		Required:     false,
		ErrorMessage: "",
		HelpText:     "",
	}
	if comp := ValidatedDatepicker("Test Date", "2023-12-25", validation, nil); comp == nil {
		t.Error("ValidatedDatepicker returned nil for basic parameters")
	}
}

// TestValidatedDatepickerRequired verifies required datepicker field rendering.
func TestValidatedDatepickerRequired(t *testing.T) {
	validation := ValidationState{
		HasError:     false,
		Required:     true,
		ErrorMessage: "",
		HelpText:     "Please select a date",
	}
	if comp := ValidatedDatepicker("Required Date", "", validation, nil); comp == nil {
		t.Error("ValidatedDatepicker returned nil for required field")
	}
}

// TestValidatedDatepickerWithError verifies datepicker error state rendering.
func TestValidatedDatepickerWithError(t *testing.T) {
	validation := ValidationState{
		HasError:     true,
		Required:     true,
		ErrorMessage: "Date is required",
		HelpText:     "Help text hidden when error present",
	}
	if comp := ValidatedDatepicker("Error Date", "", validation, nil); comp == nil {
		t.Error("ValidatedDatepicker returned nil for error state")
	}
}
