//go:build dev
// +build dev

package api

import (
	"encoding/json"
	"fmt"
)

// GetThunderSettings retrieves Thunder Settings from thunder serve environment endpoint in dev mode
func GetThunderSettings() (*ThunderSettings, error) {
	responseData, err := Get("/api/settings")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Thunder Settings from dev server: %w", err)
	}

	var settings ThunderSettings
	if err := json.Unmarshal(responseData, &settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Thunder Settings: %w", err)
	}

	// Check if there was an error in the response
	if settings.Error {
		return nil, fmt.Errorf("Thunder Settings error: %s", settings.Message)
	}

	return &settings, nil
}
