//go:build dev
// +build dev

package api

import (
	"bytes"
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

// Post performs an HTTP POST against the local dev server and returns the response body.
// For composite requests, it returns CompositeErrors if any sub-requests fail.
func Post(url string, body []byte) ([]byte, error) {
	fmt.Printf("POST %s %s\n", url, string(body))
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var sfErrs []struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(data, &sfErrs); err == nil && len(sfErrs) > 0 {
			return nil, fmt.Errorf("POST %s returned status %d: %s", url, resp.StatusCode, sfErrs[0].Message)
		}
		return nil, fmt.Errorf("POST %s returned status %d: %s", url, resp.StatusCode, string(data))
	}

	// Check if this is a composite request response
	if isCompositeRequest(url, body) {
		if compositeErrs, err := parseCompositeResponse(data); err == nil && compositeErrs.HasErrors() {
			return data, compositeErrs
		}
	}

	return data, nil
}

// Patch performs an HTTP PATCH against the local dev server and returns the response body.
func Patch(url string, body []byte) ([]byte, error) {
	fmt.Printf("PATCH %s %s\n", url, string(body))
	req, err := http.NewRequest("PATCH", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var sfErrs []struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(data, &sfErrs); err == nil && len(sfErrs) > 0 {
			return nil, fmt.Errorf("PATCH %s returned status %d: %s", url, resp.StatusCode, sfErrs[0].Message)
		}
		return nil, fmt.Errorf("PATCH %s returned status %d: %s", url, resp.StatusCode, string(data))
	}
	return data, nil
}

// Delete performs an HTTP DELETE against the local dev server and returns the response body.
func Delete(url string) ([]byte, error) {
	fmt.Printf("DELETE %s\n", url)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var sfErrs []struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(data, &sfErrs); err == nil && len(sfErrs) > 0 {
			return nil, fmt.Errorf("DELETE %s returned status %d: %s", url, resp.StatusCode, sfErrs[0].Message)
		}
		return nil, fmt.Errorf("DELETE %s returned status %d: %s", url, resp.StatusCode, string(data))
	}
	return data, nil
}
