package api

import "encoding/json"

// PicklistFieldValue holds the details for a single picklist field.
type PicklistFieldValue struct {
	ControllerValues map[string]int  `json:"controllerValues"`
	DefaultValue     *PicklistValue  `json:"defaultValue"`
	ETag             string          `json:"eTag"`
	URL              string          `json:"url"`
	Values           []PicklistValue `json:"values"`
}

// PicklistValue represents a single option within a picklist field.
type PicklistValue struct {
	Attributes interface{} `json:"attributes"`
	Label      string      `json:"label"`
	ValidFor   []int       `json:"validFor"`
	Value      string      `json:"value"`
}

// picklistValuesResponse is used for unmarshalling the top-level picklist-values JSON.
type picklistValuesResponse struct {
	PicklistFieldValues map[string]PicklistFieldValue `json:"picklistFieldValues"`
}

// UnmarshalPicklistFieldValues parses picklist-values JSON into a map of field names to PicklistFieldValue.
func UnmarshalPicklistFieldValues(data []byte) (map[string]PicklistFieldValue, error) {
	var resp picklistValuesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.PicklistFieldValues, nil
}
