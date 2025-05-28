package api

import (
	"testing"
)

func TestUnmarshalObjectInfo(t *testing.T) {
	// Sample Account object info JSON with essential fields for testing
	accountJSON := `{
		"apiName": "Account",
		"eTag": "550e8e2e-e09f-4c2c-8b93-4b7dc0f48b3b",
		"defaultRecordTypeId": "012000000000000AAA",
		"compactLayoutable": true,
		"createable": true,
		"custom": false,
		"deletable": true,
		"feedEnabled": true,
		"keyPrefix": "001",
		"label": "Account",
		"labelPlural": "Accounts",
		"layoutable": true,
		"mruEnabled": true,
		"queryable": true,
		"searchLayoutable": true,
		"searchable": true,
		"updateable": true,
		"nameFields": ["Name"],
		"childRelationships": [],
		"dependentFields": {},
		"themeInfo": {
			"color": "5867E8",
			"iconUrl": "/img/icon/t4v35/standard/account_120.png"
		},
		"recordTypeInfos": {
			"012000000000000AAA": {
				"available": true,
				"defaultRecordTypeMapping": true,
				"master": true,
				"name": "Master",
				"recordTypeId": "012000000000000AAA"
			}
		},
		"fields": {
			"Name": {
				"apiName": "Name",
				"calculated": false,
				"compound": false,
				"createable": true,
				"custom": false,
				"dataType": "String",
				"externalId": false,
				"filterable": true,
				"highScaleNumber": false,
				"htmlFormatted": false,
				"label": "Account Name",
				"length": 255,
				"nameField": true,
				"polymorphicForeignKey": false,
				"required": true,
				"searchPrefilterable": true,
				"sortable": true,
				"unique": false,
				"updateable": true,
				"controllingFields": [],
				"referenceToInfos": []
			},
			"Id": {
				"apiName": "Id",
				"calculated": false,
				"compound": false,
				"createable": false,
				"custom": false,
				"dataType": "Id",
				"externalId": false,
				"filterable": true,
				"highScaleNumber": false,
				"htmlFormatted": false,
				"label": "Account ID",
				"length": 18,
				"nameField": false,
				"polymorphicForeignKey": false,
				"required": false,
				"searchPrefilterable": false,
				"sortable": true,
				"unique": false,
				"updateable": false,
				"controllingFields": [],
				"referenceToInfos": []
			}
		}
	}`

	info, err := UnmarshalObjectInfo([]byte(accountJSON))
	if err != nil {
		t.Fatalf("UnmarshalObjectInfo returned error: %v", err)
	}
	if info.APIName != "Account" {
		t.Errorf("APIName: got %q, want %q", info.APIName, "Account")
	}
	if info.ETag == "" {
		t.Error("expected non-empty ETag")
	}
	if info.DefaultRecordTypeID == "" {
		t.Error("expected non-empty DefaultRecordTypeID")
	}
	if _, ok := info.RecordTypeInfos[info.DefaultRecordTypeID]; !ok {
		t.Errorf("RecordTypeInfos missing defaultRecordTypeId %q", info.DefaultRecordTypeID)
	}
	if len(info.Fields) == 0 {
		t.Error("expected non-empty Fields map")
	}
	if _, ok := info.Fields["Name"]; !ok {
		t.Error("expected field Name to be present")
	}
}
