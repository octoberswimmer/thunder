package api

import (
	"testing"
)

// TestPlacePredictionStruct verifies PlacePrediction structure
func TestPlacePredictionStruct(t *testing.T) {
	prediction := PlacePrediction{
		PlaceID:     "ChIJ123",
		Description: "123 Main St, Anytown, State",
	}

	if prediction.PlaceID != "ChIJ123" {
		t.Error("Expected PlaceID to be set correctly")
	}

	if prediction.Description != "123 Main St, Anytown, State" {
		t.Error("Expected Description to be set correctly")
	}
}

// TestPlaceDetailsStruct verifies PlaceDetails structure
func TestPlaceDetailsStruct(t *testing.T) {
	details := PlaceDetails{
		PlaceID:          "ChIJ123",
		FormattedAddress: "123 Main St, Anytown, State 12345",
		Street:           "123 Main St",
		City:             "Anytown",
		State:            "State",
		PostalCode:       "12345",
		Country:          "Country",
		Latitude:         40.7128,
		Longitude:        -74.0060,
	}

	if details.PlaceID != "ChIJ123" {
		t.Error("Expected PlaceID to be set correctly")
	}

	if details.Street != "123 Main St" {
		t.Error("Expected Street to be set correctly")
	}

	if details.City != "Anytown" {
		t.Error("Expected City to be set correctly")
	}

	if details.Latitude != 40.7128 {
		t.Error("Expected Latitude to be set correctly")
	}

	if details.Longitude != -74.0060 {
		t.Error("Expected Longitude to be set correctly")
	}
}

// TestGetPlacesAutocompleteEmptyInput verifies behavior with empty input
func TestGetPlacesAutocompleteEmptyInput(t *testing.T) {
	predictions, err := GetPlacesAutocomplete("test-key", "")
	if err != nil {
		t.Errorf("Expected no error with empty input, got: %v", err)
	}

	if predictions != nil {
		t.Error("Expected nil predictions with empty input")
	}
}

// TestGetPlaceDetailsEmptyPlaceID verifies behavior with empty place ID
func TestGetPlaceDetailsEmptyPlaceID(t *testing.T) {
	_, err := GetPlaceDetails("test-key", "")
	if err == nil {
		t.Error("Expected error with empty place ID")
	}
}
