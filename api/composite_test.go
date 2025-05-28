package api

import (
	"encoding/json"
	"testing"
)

func TestIsCompositeRequest(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		body     string
		expected bool
	}{
		{
			name:     "composite URL with composite request body",
			url:      "/services/data/v58.0/composite",
			body:     `{"compositeRequest": [{"method": "POST", "url": "/services/data/v58.0/sobjects/Account", "referenceId": "ref1"}]}`,
			expected: true,
		},
		{
			name:     "composite URL without composite request body",
			url:      "/services/data/v58.0/composite",
			body:     `{"name": "Test Account"}`,
			expected: false,
		},
		{
			name:     "non-composite URL with composite request body",
			url:      "/services/data/v58.0/sobjects/Account",
			body:     `{"compositeRequest": [{"method": "POST", "url": "/services/data/v58.0/sobjects/Account", "referenceId": "ref1"}]}`,
			expected: false,
		},
		{
			name:     "non-composite URL and body",
			url:      "/services/data/v58.0/sobjects/Account",
			body:     `{"name": "Test Account"}`,
			expected: false,
		},
		{
			name:     "invalid JSON body",
			url:      "/services/data/v58.0/composite",
			body:     `invalid json`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isCompositeRequest(tt.url, []byte(tt.body))
			if result != tt.expected {
				t.Errorf("isCompositeRequest() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseCompositeResponse(t *testing.T) {
	tests := []struct {
		name        string
		response    string
		expectError bool
		hasErrors   bool
		errorCount  int
	}{
		{
			name: "successful composite response",
			response: `{
				"compositeResponse": [
					{
						"body": {"id": "001XX000003DHP0", "success": true},
						"httpStatusCode": 200,
						"referenceId": "ref1"
					}
				]
			}`,
			expectError: false,
			hasErrors:   false,
			errorCount:  0,
		},
		{
			name: "composite response with errors",
			response: `{
				"compositeResponse": [
					{
						"body": {"id": "001XX000003DHP0", "success": true},
						"httpStatusCode": 200,
						"referenceId": "ref1"
					},
					{
						"body": {"message": "Required fields are missing", "errorCode": "REQUIRED_FIELD_MISSING"},
						"httpStatusCode": 400,
						"referenceId": "ref2"
					}
				]
			}`,
			expectError: false,
			hasErrors:   true,
			errorCount:  1,
		},
		{
			name: "all failed composite response",
			response: `{
				"compositeResponse": [
					{
						"body": {"message": "Required fields are missing", "errorCode": "REQUIRED_FIELD_MISSING"},
						"httpStatusCode": 400,
						"referenceId": "ref1"
					},
					{
						"body": {"message": "Invalid object type", "errorCode": "INVALID_TYPE"},
						"httpStatusCode": 400,
						"referenceId": "ref2"
					}
				]
			}`,
			expectError: false,
			hasErrors:   true,
			errorCount:  2,
		},
		{
			name: "single composite error",
			response: `{
				"message": "Error processing composite request: Invalid JSON",
				"errorCode": "COMPOSITE_REQUEST_ERROR"
			}`,
			expectError: false,
			hasErrors:   true,
			errorCount:  1,
		},
		{
			name:        "invalid JSON",
			response:    `invalid json`,
			expectError: true,
			hasErrors:   false,
			errorCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compositeErrs, err := parseCompositeResponse([]byte(tt.response))

			if tt.expectError {
				if err == nil {
					t.Errorf("parseCompositeResponse() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("parseCompositeResponse() unexpected error: %v", err)
				return
			}

			if compositeErrs.HasErrors() != tt.hasErrors {
				t.Errorf("parseCompositeResponse().HasErrors() = %v, want %v", compositeErrs.HasErrors(), tt.hasErrors)
			}

			if len(compositeErrs.Errors) != tt.errorCount {
				t.Errorf("parseCompositeResponse() error count = %d, want %d", len(compositeErrs.Errors), tt.errorCount)
			}
		})
	}
}

func TestCompositeErrorsError(t *testing.T) {
	tests := []struct {
		name     string
		errors   *CompositeErrors
		expected string
	}{
		{
			name: "single error with message",
			errors: &CompositeErrors{
				Errors: []CompositeSubResponse{
					{
						Body: map[string]interface{}{
							"message":   "Required fields are missing",
							"errorCode": "REQUIRED_FIELD_MISSING",
						},
						HTTPStatusCode: 400,
						ReferenceID:    "ref1",
					},
				},
			},
			expected: "composite request failed: ref ref1: Required fields are missing",
		},
		{
			name: "multiple errors",
			errors: &CompositeErrors{
				Errors: []CompositeSubResponse{
					{
						Body: map[string]interface{}{
							"message": "Required fields are missing",
						},
						HTTPStatusCode: 400,
						ReferenceID:    "ref1",
					},
					{
						Body: map[string]interface{}{
							"message": "Invalid object type",
						},
						HTTPStatusCode: 400,
						ReferenceID:    "ref2",
					},
				},
			},
			expected: "composite request failed: ref ref1: Required fields are missing; ref ref2: Invalid object type",
		},
		{
			name: "error without message",
			errors: &CompositeErrors{
				Errors: []CompositeSubResponse{
					{
						Body:           map[string]interface{}{},
						HTTPStatusCode: 400,
						ReferenceID:    "ref1",
					},
				},
			},
			expected: "composite request failed with 1 errors",
		},
		{
			name:     "no errors",
			errors:   &CompositeErrors{},
			expected: "composite request failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.errors.Error()
			if result != tt.expected {
				t.Errorf("CompositeErrors.Error() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestCompositeRequestJSON(t *testing.T) {
	req := CompositeRequest{
		AllOrNone: true,
		CompositeRequest: []CompositeSubRequest{
			{
				Method:      "POST",
				URL:         "/services/data/v58.0/sobjects/Account",
				ReferenceID: "ref1",
				Body: map[string]interface{}{
					"Name": "Test Account",
				},
			},
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json.Marshal() error: %v", err)
	}

	var parsed CompositeRequest
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("json.Unmarshal() error: %v", err)
	}

	if parsed.AllOrNone != req.AllOrNone {
		t.Errorf("AllOrNone = %v, want %v", parsed.AllOrNone, req.AllOrNone)
	}

	if len(parsed.CompositeRequest) != len(req.CompositeRequest) {
		t.Errorf("CompositeRequest length = %d, want %d", len(parsed.CompositeRequest), len(req.CompositeRequest))
	}
}
