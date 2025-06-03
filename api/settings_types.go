package api

// ThunderSettings represents the Thunder Settings custom settings object
type ThunderSettings struct {
	GoogleMapsAPIKey string `json:"Google_Maps_API_Key__c"`
	Error            bool   `json:"error,omitempty"`
	Message          string `json:"message,omitempty"`
}
