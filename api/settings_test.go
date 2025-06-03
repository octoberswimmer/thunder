package api

import (
	"testing"
)

// TestThunderSettingsStruct verifies ThunderSettings structure
func TestThunderSettingsStruct(t *testing.T) {
	settings := ThunderSettings{
		GoogleMapsAPIKey: "test-api-key",
		Error:            false,
		Message:          "",
	}

	if settings.GoogleMapsAPIKey != "test-api-key" {
		t.Error("Expected GoogleMapsAPIKey to be set correctly")
	}

	if settings.Error != false {
		t.Error("Expected Error to be false")
	}

	if settings.Message != "" {
		t.Error("Expected Message to be empty")
	}
}

// TestGetThunderSettings verifies that GetThunderSettings returns a result
func TestGetThunderSettings(t *testing.T) {
	settings, err := GetThunderSettings()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if settings == nil {
		t.Error("Expected settings to be returned")
	}
}
