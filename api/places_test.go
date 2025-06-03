package api

import (
	"encoding/json"
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

// TestParseGoogleAPIError verifies error parsing functionality
func TestParseGoogleAPIError(t *testing.T) {
	// Test JSON error response
	jsonError := `{"error":{"code":400,"message":"Invalid request","status":"INVALID_REQUEST"}}`
	err := parseGoogleAPIError(400, []byte(jsonError))
	if err == nil {
		t.Error("Expected error to be parsed")
	}
	if err.Error() != "Google Places API error (400): Invalid request" {
		t.Errorf("Expected specific error message, got: %s", err.Error())
	}

	// Test HTTP status code fallback
	err = parseGoogleAPIError(401, []byte("invalid response"))
	if err == nil {
		t.Error("Expected error for 401 status")
	}
	if err.Error() != "invalid API key - please check your Google Places API key" {
		t.Errorf("Expected API key error, got: %s", err.Error())
	}

	// Test rate limiting
	err = parseGoogleAPIError(429, []byte(""))
	if err == nil {
		t.Error("Expected error for 429 status")
	}
	if err.Error() != "too many requests - please try again in a moment" {
		t.Errorf("Expected rate limit error, got: %s", err.Error())
	}
}

// TestGooglePlacesResponseParsing verifies JSON response parsing
func TestGooglePlacesResponseParsing(t *testing.T) {
	// Test with actual Google Places API response format
	jsonResponse := `{
		"suggestions": [
			{
				"placePrediction": {
					"place": "places/ChIJm2PSmds5DogR0R9sUII7y8A",
					"placeId": "ChIJm2PSmds5DogR0R9sUII7y8A",
					"text": {
						"text": "7285 West 87th Street, Bridgeview, IL, USA"
					}
				}
			},
			{
				"placePrediction": {
					"placeId": "EjA3Mjg1IFdlc3QgQXJtaXRhZ2UgQXZlbnVlLCBFbG13b29kIFBhcmssIElMLCBVU0E",
					"text": {
						"text": "7285 West Armitage Avenue, Elmwood Park, IL, USA"
					}
				}
			}
		]
	}`

	var response GooglePlacesAutocompleteResponse
	err := json.Unmarshal([]byte(jsonResponse), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if len(response.Suggestions) != 2 {
		t.Errorf("Expected 2 suggestions, got %d", len(response.Suggestions))
	}

	// Test first suggestion
	firstSuggestion := response.Suggestions[0].PlacePrediction
	if firstSuggestion.PlaceID != "ChIJm2PSmds5DogR0R9sUII7y8A" {
		t.Errorf("Expected specific place ID, got: %s", firstSuggestion.PlaceID)
	}
	if firstSuggestion.Text.Text != "7285 West 87th Street, Bridgeview, IL, USA" {
		t.Errorf("Expected specific text, got: %s", firstSuggestion.Text.Text)
	}

	// Test conversion to PlacePrediction struct
	var predictions []PlacePrediction
	for _, suggestion := range response.Suggestions {
		predictions = append(predictions, PlacePrediction{
			PlaceID:     suggestion.PlacePrediction.PlaceID,
			Description: suggestion.PlacePrediction.Text.Text,
		})
	}

	if len(predictions) != 2 {
		t.Errorf("Expected 2 converted predictions, got %d", len(predictions))
	}

	if predictions[0].PlaceID != "ChIJm2PSmds5DogR0R9sUII7y8A" {
		t.Errorf("Expected specific place ID in converted prediction, got: %s", predictions[0].PlaceID)
	}

	if predictions[0].Description != "7285 West 87th Street, Bridgeview, IL, USA" {
		t.Errorf("Expected specific description in converted prediction, got: %s", predictions[0].Description)
	}
}
