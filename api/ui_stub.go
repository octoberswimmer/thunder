//go:build !js
// +build !js

package api

// GetPicklistValuesByRecordType is a stub implementation for non-WASM builds and will panic if called.
func GetPicklistValuesByRecordType(objectName, recordTypeId string) (map[string]PicklistFieldValue, error) {
	panic("api.GetPicklistValuesByRecordType is not supported outside the WASM environment")
}

// GetObjectInfo is a stub implementation for non-WASM builds and will panic if called.
func GetObjectInfo(objectName string) (ObjectInfo, error) {
	panic("api.GetObjectInfo is not supported outside the WASM environment")
}
