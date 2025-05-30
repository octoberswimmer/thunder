package components

import (
	"testing"
	"time"

	"github.com/octoberswimmer/masc"
)

// TestValidatedTextInputBasic verifies basic ValidatedTextInput functionality.
func TestValidatedTextInputBasic(t *testing.T) {
	validation := ValidationState{
		HasError:     false,
		Required:     false,
		ErrorMessage: "",
		HelpText:     "",
		Placeholder:  "Test Placeholder",
	}
	if comp := ValidatedTextInput("Test Label", "Test Value", validation, nil); comp == nil {
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
	if comp := ValidatedTextInput("Required Field", "", validation, nil); comp == nil {
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
	if comp := ValidatedTextInput("Error Field", "", validation, nil); comp == nil {
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
		Placeholder:  "placeholder",
	}
	handler := func(e *masc.Event) {}
	if comp := ValidatedTextInput("Handler Test", "value", validation, handler); comp == nil {
		t.Error("ValidatedTextInput returned nil with event handler")
	}
}

// TestValidatedTextInputWithTooltip verifies tooltip functionality.
func TestValidatedTextInputWithTooltip(t *testing.T) {
	validation := ValidationState{
		HasError:     false,
		Required:     false,
		ErrorMessage: "",
		HelpText:     "",
		Tooltip:      "This is a helpful tooltip",
		Placeholder:  "Enter text here",
	}
	if comp := ValidatedTextInput("Tooltip Field", "", validation, nil); comp == nil {
		t.Error("ValidatedTextInput returned nil with tooltip")
	}
}

// TestValidatedTextInputWithTooltipAndRequired verifies tooltip with required field.
func TestValidatedTextInputWithTooltipAndRequired(t *testing.T) {
	validation := ValidationState{
		HasError:     false,
		Required:     true,
		ErrorMessage: "",
		HelpText:     "",
		Tooltip:      "Required field with tooltip",
	}
	if comp := ValidatedTextInput("Required Tooltip Field", "", validation, nil); comp == nil {
		t.Error("ValidatedTextInput returned nil with tooltip and required")
	}
}

// TestValidatedTextInputConvenienceConstructors verifies convenience constructor functions.
func TestValidatedTextInputConvenienceConstructors(t *testing.T) {
	// Test WithTooltip
	tooltipValidation := WithTooltip("Test tooltip")
	if comp := ValidatedTextInput("Tooltip Test", "", tooltipValidation, nil); comp == nil {
		t.Error("ValidatedTextInput returned nil with WithTooltip constructor")
	}

	// Test WithPlaceholder
	placeholderValidation := WithPlaceholder("Test placeholder")
	if comp := ValidatedTextInput("Placeholder Test", "", placeholderValidation, nil); comp == nil {
		t.Error("ValidatedTextInput returned nil with WithPlaceholder constructor")
	}

	// Test WithTooltipAndPlaceholder
	bothValidation := WithTooltipAndPlaceholder("Test tooltip", "Test placeholder")
	if comp := ValidatedTextInput("Both Test", "", bothValidation, nil); comp == nil {
		t.Error("ValidatedTextInput returned nil with WithTooltipAndPlaceholder constructor")
	}

	// Test Required
	requiredValidation := Required()
	if comp := ValidatedTextInput("Required Test", "", requiredValidation, nil); comp == nil {
		t.Error("ValidatedTextInput returned nil with Required constructor")
	}

	// Test RequiredWithTooltip
	requiredTooltipValidation := RequiredWithTooltip("Required tooltip")
	if comp := ValidatedTextInput("Required Tooltip Test", "", requiredTooltipValidation, nil); comp == nil {
		t.Error("ValidatedTextInput returned nil with RequiredWithTooltip constructor")
	}
}

// TestValidatedTextareaBasic verifies basic ValidatedTextarea functionality.
func TestValidatedTextareaBasic(t *testing.T) {
	validation := ValidationState{
		HasError:     false,
		Required:     false,
		ErrorMessage: "",
		HelpText:     "",
		Placeholder:  "Test Placeholder",
	}
	if comp := ValidatedTextarea("Test Label", "Test Value", 3, validation, nil); comp == nil {
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
	if comp := ValidatedTextarea("Required Textarea", "", 5, validation, nil); comp == nil {
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
	if comp := ValidatedTextarea("Error Textarea", "", 2, validation, nil); comp == nil {
		t.Error("ValidatedTextarea returned nil for error state")
	}
}

// TestValidatedTextareaWithTooltip verifies textarea tooltip functionality.
func TestValidatedTextareaWithTooltip(t *testing.T) {
	validation := ValidationState{
		HasError:     false,
		Required:     false,
		ErrorMessage: "",
		HelpText:     "",
		Tooltip:      "This textarea has a tooltip",
		Placeholder:  "Enter detailed text here",
	}
	if comp := ValidatedTextarea("Tooltip Textarea", "", 4, validation, nil); comp == nil {
		t.Error("ValidatedTextarea returned nil with tooltip")
	}
}

// TestIsEmptyOrWhitespace verifies the whitespace validation helper function.
func TestIsEmptyOrWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty string", "", true},
		{"single space", " ", true},
		{"multiple spaces", "   ", true},
		{"tabs", "\t\t", true},
		{"newlines", "\n\n", true},
		{"mixed whitespace", " \t\n ", true},
		{"text with spaces", " hello ", false},
		{"just text", "hello", false},
		{"single char", "a", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEmptyOrWhitespace(tt.input)
			if result != tt.expected {
				t.Errorf("IsEmptyOrWhitespace(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestValidateRequired verifies the required field validation helper.
func TestValidateRequired(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		fieldName string
		hasError  bool
		required  bool
	}{
		{"valid value", "John", "First Name", false, true},
		{"empty string", "", "First Name", true, true},
		{"whitespace only", "   ", "First Name", true, true},
		{"tabs only", "\t\t", "First Name", true, true},
		{"mixed whitespace", " \n\t ", "First Name", true, true},
		{"value with spaces", " John ", "First Name", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateRequired(tt.value, tt.fieldName)
			if result.HasError != tt.hasError {
				t.Errorf("ValidateRequired(%q, %q).HasError = %v, expected %v", tt.value, tt.fieldName, result.HasError, tt.hasError)
			}
			if result.Required != tt.required {
				t.Errorf("ValidateRequired(%q, %q).Required = %v, expected %v", tt.value, tt.fieldName, result.Required, tt.required)
			}
			if tt.hasError && result.ErrorMessage == "" {
				t.Errorf("ValidateRequired(%q, %q) should have error message when hasError is true", tt.value, tt.fieldName)
			}
		})
	}
}

// TestValidateRequiredWithTooltip verifies the tooltip validation helper.
func TestValidateRequiredWithTooltip(t *testing.T) {
	tooltip := "This field needs a value"

	tests := []struct {
		name      string
		value     string
		fieldName string
		hasError  bool
	}{
		{"valid value", "John", "First Name", false},
		{"empty string", "", "First Name", true},
		{"whitespace only", "   ", "First Name", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateRequiredWithTooltip(tt.value, tt.fieldName, tooltip)
			if result.HasError != tt.hasError {
				t.Errorf("ValidateRequiredWithTooltip(%q, %q, %q).HasError = %v, expected %v", tt.value, tt.fieldName, tooltip, result.HasError, tt.hasError)
			}
			if result.Tooltip != tooltip {
				t.Errorf("ValidateRequiredWithTooltip(%q, %q, %q).Tooltip = %q, expected %q", tt.value, tt.fieldName, tooltip, result.Tooltip, tooltip)
			}
			if result.Required != true {
				t.Errorf("ValidateRequiredWithTooltip(%q, %q, %q).Required = %v, expected true", tt.value, tt.fieldName, tooltip, result.Required)
			}
		})
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

// TestValidatedSelectWithTooltip verifies select tooltip functionality.
func TestValidatedSelectWithTooltip(t *testing.T) {
	options := []SelectOption{
		{Label: "Choose...", Value: ""},
		{Label: "Option 1", Value: "opt1"},
		{Label: "Option 2", Value: "opt2"},
	}
	validation := ValidationState{
		HasError:     false,
		Required:     false,
		ErrorMessage: "",
		HelpText:     "",
		Tooltip:      "Select one of the available options",
	}
	if comp := ValidatedSelect("Tooltip Select", options, "", validation, nil); comp == nil {
		t.Error("ValidatedSelect returned nil with tooltip")
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
	if comp := ValidatedDatepicker("Test Date", time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC), validation, func(t time.Time) {}); comp == nil {
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
	if comp := ValidatedDatepicker("Required Date", time.Time{}, validation, func(t time.Time) {}); comp == nil {
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
	if comp := ValidatedDatepicker("Error Date", time.Time{}, validation, func(t time.Time) {}); comp == nil {
		t.Error("ValidatedDatepicker returned nil for error state")
	}
}

// TestValidatedDatepickerWithValue verifies datepicker with specific date value.
func TestValidatedDatepickerWithValue(t *testing.T) {
	date := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	validation := ValidationState{
		HasError:     false,
		Required:     false,
		ErrorMessage: "",
		HelpText:     "Select your preferred date",
	}
	if comp := ValidatedDatepicker("Appointment Date", date, validation, func(t time.Time) {}); comp == nil {
		t.Error("ValidatedDatepicker returned nil with date value")
	}
}

// TestValidatedDatepickerWithHandler verifies datepicker with event handler.
func TestValidatedDatepickerWithHandler(t *testing.T) {
	date := time.Date(2023, 3, 10, 0, 0, 0, 0, time.UTC)
	validation := ValidationState{
		HasError:     false,
		Required:     true,
		ErrorMessage: "",
		HelpText:     "",
	}
	handler := func(t time.Time) {}
	if comp := ValidatedDatepicker("Event Date", date, validation, handler); comp == nil {
		t.Error("ValidatedDatepicker returned nil with event handler")
	}
}

// TestValidatedDatepickerZeroValue verifies datepicker handles zero time value correctly.
func TestValidatedDatepickerZeroValue(t *testing.T) {
	validation := ValidationState{
		HasError:     false,
		Required:     false,
		ErrorMessage: "",
		HelpText:     "Optional date field",
	}
	if comp := ValidatedDatepicker("Optional Date", time.Time{}, validation, func(t time.Time) {}); comp == nil {
		t.Error("ValidatedDatepicker returned nil with zero time value")
	}
}
