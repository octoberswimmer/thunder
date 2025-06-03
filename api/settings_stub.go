//go:build !js && !dev
// +build !js,!dev

package api

import (
	"os"
)

// GetThunderSettings stub implementation for non-WASM, non-dev builds
func GetThunderSettings() (*ThunderSettings, error) {
	settings := &ThunderSettings{
		GoogleMapsAPIKey: os.Getenv("GOOGLE_MAPS_API_KEY"),
	}

	return settings, nil
}
