package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

// CompositeRequest represents a Salesforce composite request
type CompositeRequest struct {
	AllOrNone        bool                  `json:"allOrNone"`
	CompositeRequest []CompositeSubRequest `json:"compositeRequest"`
}

// CompositeSubRequest represents a single sub-request within a composite request
type CompositeSubRequest struct {
	Method      string      `json:"method"`
	URL         string      `json:"url"`
	ReferenceID string      `json:"referenceId"`
	Body        interface{} `json:"body,omitempty"`
}

// CompositeResponse represents a Salesforce composite response
type CompositeResponse struct {
	CompositeResponse []CompositeSubResponse `json:"compositeResponse"`
}

// CompositeSubResponse represents a single sub-response within a composite response
type CompositeSubResponse struct {
	Body           interface{} `json:"body"`
	HTTPStatusCode int         `json:"httpStatusCode"`
	ReferenceID    string      `json:"referenceId"`
}

// CompositeError represents an error in a composite request
type CompositeError struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

// CompositeErrors represents a collection of composite request errors
type CompositeErrors struct {
	Errors      []CompositeSubResponse `json:"errors"`
	PartialData []CompositeSubResponse `json:"partialData,omitempty"`
}

// Error implements the error interface for CompositeErrors
func (ce *CompositeErrors) Error() string {
	if len(ce.Errors) == 0 {
		return "composite request failed"
	}

	var errorMessages []string
	for _, err := range ce.Errors {
		if body, ok := err.Body.(map[string]interface{}); ok {
			if message, exists := body["message"]; exists {
				errorMessages = append(errorMessages, fmt.Sprintf("ref %s: %v", err.ReferenceID, message))
			}
		}
	}

	if len(errorMessages) > 0 {
		return fmt.Sprintf("composite request failed: %s", strings.Join(errorMessages, "; "))
	}

	return fmt.Sprintf("composite request failed with %d errors", len(ce.Errors))
}

// HasErrors returns true if the composite response contains any errors
func (ce *CompositeErrors) HasErrors() bool {
	return len(ce.Errors) > 0
}

// isCompositeRequest checks if the URL and body indicate a composite request
func isCompositeRequest(url string, body []byte) bool {
	if !strings.Contains(url, "/composite") {
		return false
	}

	// Check if the body contains composite request structure
	var req map[string]interface{}
	if err := json.Unmarshal(body, &req); err != nil {
		return false
	}

	_, hasCompositeRequest := req["compositeRequest"]
	return hasCompositeRequest
}

// parseCompositeResponse parses the response body and extracts composite errors
func parseCompositeResponse(data []byte) (*CompositeErrors, error) {
	// Try to parse as composite response first
	var compositeResp CompositeResponse
	if err := json.Unmarshal(data, &compositeResp); err == nil && len(compositeResp.CompositeResponse) > 0 {
		errors := &CompositeErrors{}

		for _, subResp := range compositeResp.CompositeResponse {
			if subResp.HTTPStatusCode >= 400 {
				errors.Errors = append(errors.Errors, subResp)
			} else {
				errors.PartialData = append(errors.PartialData, subResp)
			}
		}

		return errors, nil
	}

	// Try to parse as single composite error
	var compositeErr CompositeError
	if err := json.Unmarshal(data, &compositeErr); err == nil && compositeErr.Message != "" {
		errors := &CompositeErrors{
			Errors: []CompositeSubResponse{
				{
					Body:           compositeErr,
					HTTPStatusCode: 400,
					ReferenceID:    "composite",
				},
			},
		}
		return errors, nil
	}

	return nil, fmt.Errorf("unable to parse composite response")
}
