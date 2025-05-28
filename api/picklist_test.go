package api

import (
	"testing"
)

func TestUnmarshalPicklistFieldValues(t *testing.T) {
	// Sample picklist values JSON for Account.AccountSource field
	picklistJSON := `{
		"picklistFieldValues": {
			"AccountSource": {
				"controllerValues": {},
				"defaultValue": null,
				"eTag": "b2b9f8c7-3a1d-4e5f-9c7a-1a2b3c4d5e6f",
				"url": "/services/data/v63.0/ui-api/object-info/Account/picklist-values/012QP000000DsyiYAC/AccountSource",
				"values": [
					{
						"attributes": null,
						"label": "Advertisement",
						"validFor": [],
						"value": "Advertisement"
					},
					{
						"attributes": null,
						"label": "Cold Call",
						"validFor": [],
						"value": "Cold Call"
					},
					{
						"attributes": null,
						"label": "Employee Referral",
						"validFor": [],
						"value": "Employee Referral"
					},
					{
						"attributes": null,
						"label": "External Referral",
						"validFor": [],
						"value": "External Referral"
					},
					{
						"attributes": null,
						"label": "Partner",
						"validFor": [],
						"value": "Partner"
					},
					{
						"attributes": null,
						"label": "Public Relations",
						"validFor": [],
						"value": "Public Relations"
					},
					{
						"attributes": null,
						"label": "Seminar - Internal",
						"validFor": [],
						"value": "Seminar - Internal"
					},
					{
						"attributes": null,
						"label": "Seminar - Partner",
						"validFor": [],
						"value": "Seminar - Partner"
					},
					{
						"attributes": null,
						"label": "Trade Show",
						"validFor": [],
						"value": "Trade Show"
					},
					{
						"attributes": null,
						"label": "Web",
						"validFor": [],
						"value": "Web"
					},
					{
						"attributes": null,
						"label": "Word of mouth",
						"validFor": [],
						"value": "Word of mouth"
					},
					{
						"attributes": null,
						"label": "Other",
						"validFor": [],
						"value": "Other"
					}
				]
			}
		}
	}`

	m, err := UnmarshalPicklistFieldValues([]byte(picklistJSON))
	if err != nil {
		t.Fatalf("UnmarshalPicklistFieldValues returned error: %v", err)
	}
	if len(m) == 0 {
		t.Fatal("expected at least one picklist field, got none")
	}
	pfv, ok := m["AccountSource"]
	if !ok {
		t.Fatal("expected AccountSource in picklist field map")
	}
	if pfv.ETag == "" {
		t.Error("expected non-empty ETag for AccountSource")
	}
	expURL := "/services/data/v63.0/ui-api/object-info/Account/picklist-values/012QP000000DsyiYAC/AccountSource"
	if pfv.URL != expURL {
		t.Errorf("unexpected URL: got %q, want %q", pfv.URL, expURL)
	}
	if len(pfv.Values) == 0 {
		t.Error("expected non-empty Values slice for AccountSource")
	}
	first := pfv.Values[0]
	if first.Value != "Advertisement" {
		t.Errorf("expected first picklist value Advertisement, got %q", first.Value)
	}
}
