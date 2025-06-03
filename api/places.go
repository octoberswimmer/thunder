package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// PlacePrediction represents a Google Places API autocomplete prediction
type PlacePrediction struct {
	PlaceID     string `json:"place_id"`
	Description string `json:"description"`
}

// PlaceDetails represents detailed place information from Google Places API
type PlaceDetails struct {
	PlaceID          string  `json:"place_id"`
	FormattedAddress string  `json:"formatted_address"`
	Street           string  `json:"street"`
	City             string  `json:"city"`
	State            string  `json:"state"`
	PostalCode       string  `json:"postal_code"`
	Country          string  `json:"country"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
}

// GooglePlacesAutocompleteResponse represents Google's API response
type GooglePlacesAutocompleteResponse struct {
	Suggestions []struct {
		PlacePrediction PlacePrediction `json:"placePrediction"`
	} `json:"suggestions"`
}

// GooglePlaceDetailsResponse represents Google's place details API response
type GooglePlaceDetailsResponse struct {
	PlaceID           string `json:"id"`
	FormattedAddress  string `json:"formattedAddress"`
	AddressComponents []struct {
		LongText  string   `json:"longText"`
		ShortText string   `json:"shortText"`
		Types     []string `json:"types"`
	} `json:"addressComponents"`
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
}

// GetPlacesAutocomplete fetches address predictions from Google Places API
func GetPlacesAutocomplete(apiKey, input string) ([]PlacePrediction, error) {
	if input == "" {
		return nil, nil
	}

	// Build request body for Google Places API
	requestBody := map[string]interface{}{
		"input":                input,
		"includedPrimaryTypes": []string{"street_address", "premise", "subpremise"},
		"languageCode":         "en",
	}

	requestData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Build URL for Google Places API
	u := fmt.Sprintf("https://places.googleapis.com/v1/places:autocomplete?key=%s", url.QueryEscape(apiKey))

	// Make direct HTTP POST to Google Places API
	resp, err := http.Post(u, "application/json", bytes.NewReader(requestData))
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("places autocomplete request failed: %w", err)
	}

	var response GooglePlacesAutocompleteResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract predictions
	var predictions []PlacePrediction
	for _, suggestion := range response.Suggestions {
		predictions = append(predictions, suggestion.PlacePrediction)
	}

	return predictions, nil
}

// GetPlaceDetails fetches detailed place information from Google Places API
func GetPlaceDetails(apiKey, placeID string) (*PlaceDetails, error) {
	if placeID == "" {
		return nil, fmt.Errorf("place ID is required")
	}

	// Build URL for Google Places API
	u := fmt.Sprintf("https://places.googleapis.com/v1/places/%s?fields=id,formattedAddress,addressComponents,location&key=%s",
		url.QueryEscape(placeID), url.QueryEscape(apiKey))

	// Make direct HTTP GET to Google Places API
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("place details request failed: %w", err)
	}

	var response GooglePlaceDetailsResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Convert to our format
	details := &PlaceDetails{
		PlaceID:          response.PlaceID,
		FormattedAddress: response.FormattedAddress,
		Latitude:         response.Location.Latitude,
		Longitude:        response.Location.Longitude,
	}

	// Parse address components
	for _, component := range response.AddressComponents {
		if len(component.Types) > 0 {
			switch component.Types[0] {
			case "street_number":
				details.Street = component.LongText
			case "route":
				if details.Street != "" {
					details.Street = details.Street + " " + component.LongText
				} else {
					details.Street = component.LongText
				}
			case "locality":
				details.City = component.LongText
			case "administrative_area_level_1":
				details.State = component.ShortText
			case "postal_code":
				details.PostalCode = component.ShortText
			case "country":
				details.Country = component.LongText
			}
		}
	}

	return details, nil
}
