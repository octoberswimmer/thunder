//go:build dev
// +build dev

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Get performs an HTTP GET against the local dev server and returns the response body.
func Get(url string) ([]byte, error) {
	fmt.Printf("Getting %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Read the response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// If non-2xx, attempt to unmarshal Salesforce error message
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Salesforce error responses are JSON arrays of objects with 'message'
		var sfErrs []struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(data, &sfErrs); err == nil && len(sfErrs) > 0 {
			return nil, fmt.Errorf("GET %s returned status %d: %s", url, resp.StatusCode, sfErrs[0].Message)
		}
		// Fallback to raw body
		return nil, fmt.Errorf("GET %s returned status %d: %s", url, resp.StatusCode, string(data))
	}
	return data, nil
}
