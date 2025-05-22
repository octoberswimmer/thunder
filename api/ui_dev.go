//go:build dev
// +build dev

package api

import (
	"fmt"
)

// GetPicklistValuesByRecordType fetches picklist values for the given object and record type.
// It delegates to the global JavaScript getPicklistValuesByRecordType function.
// Returns a map of field names to PicklistFieldValue or an error.
func GetPicklistValuesByRecordType(objectName, recordTypeId string) (map[string]PicklistFieldValue, error) {
	data, err := Get(fmt.Sprintf("/services/data/v63.0/ui-api/object-info/%s/picklist-values/%s", objectName, recordTypeId))
	if err != nil {
		return nil, err
	}
	return UnmarshalPicklistFieldValues(data)
}

// GetObjectInfo fetches SObject metadata for the given object.
// It delegates to the global JavaScript getObjectInfo function.
// Returns an ObjectInfo struct or an error.
func GetObjectInfo(objectName string) (ObjectInfo, error) {
	data, err := Get(fmt.Sprintf("/services/data/v63.0/ui-api/object-info/%s", objectName))
	if err != nil {
		return ObjectInfo{}, err
	}
	return UnmarshalObjectInfo(data)
}
