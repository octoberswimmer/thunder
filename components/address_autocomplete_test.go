package components

import (
	"testing"

	"github.com/octoberswimmer/thunder/api"
)

// TestAddressAutocompleteBasic verifies that AddressAutocomplete renders for basic parameters.
func TestAddressAutocompleteBasic(t *testing.T) {
	comp := AddressAutocomplete("Home Address", "", "test-api-key", nil, "", nil, nil)
	if comp == nil {
		t.Error("AddressAutocomplete returned nil for valid parameters")
	}
}

// TestAddressAutocompleteWithValue verifies that AddressAutocomplete renders with a value.
func TestAddressAutocompleteWithValue(t *testing.T) {
	comp := AddressAutocomplete("Address", "123 Main St", "test-key", nil, "", nil, nil)
	if comp == nil {
		t.Error("AddressAutocomplete returned nil with value")
	}
}

// TestAddressAutocompleteWithPredictions verifies that AddressAutocomplete renders with predictions.
func TestAddressAutocompleteWithPredictions(t *testing.T) {
	predictions := []api.PlacePrediction{
		{PlaceID: "ChIJ123", Description: "123 Main St, Anytown, State"},
		{PlaceID: "ChIJ456", Description: "456 Oak Ave, Another City, State"},
	}

	comp := AddressAutocomplete("Address", "", "test-key", predictions, "", nil, nil)
	if comp == nil {
		t.Error("AddressAutocomplete returned nil with predictions")
	}
}

// TestAddressAutocompleteWithCallbacks verifies that AddressAutocomplete renders with callbacks.
func TestAddressAutocompleteWithCallbacks(t *testing.T) {
	onInput := func(s string) {}
	onSelect := func(details api.PlaceDetails) {}

	comp := AddressAutocomplete("Address", "", "test-key", nil, "", onInput, onSelect)
	if comp == nil {
		t.Error("AddressAutocomplete returned nil with callbacks")
	}
}

// TestAddressAutocompleteWithError verifies that AddressAutocomplete renders with an error.
func TestAddressAutocompleteWithError(t *testing.T) {
	comp := AddressAutocomplete("Address", "", "test-key", nil, "API error occurred", nil, nil)
	if comp == nil {
		t.Error("AddressAutocomplete returned nil with error")
	}
}

// TestAddressAutocompleteCmd verifies that AddressAutocompleteCmd returns a command.
func TestAddressAutocompleteCmd(t *testing.T) {
	cmd := AddressAutocompleteCmd("test-key", "123 main")
	if cmd == nil {
		t.Error("AddressAutocompleteCmd returned nil")
	}
}
