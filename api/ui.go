//go:build js && !dev
// +build js,!dev

package api

import (
	"errors"
	"syscall/js"
)

// GetPicklistValuesByRecordType fetches picklist values for the given object and record type.
// It delegates to the global JavaScript getPicklistValuesByRecordType function.
// Returns a map of field names to PicklistFieldValue or an error.
func GetPicklistValuesByRecordType(objectName, recordTypeId string) (map[string]PicklistFieldValue, error) {
	dataCh := make(chan []byte)
	errCh := make(chan error)
	config := map[string]interface{}{
		"objectApiName": objectName,
		"recordTypeId":  recordTypeId,
	}
	js.Global().Call("getPicklistValuesByRecordType", js.ValueOf(config), js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		result := args[0]
		if errVal := result.Get("error"); errVal.Type() != js.TypeUndefined && errVal.Type() != js.TypeNull {
			errCh <- errors.New(errVal.String())
			return nil
		}
		if dataVal := result.Get("data"); dataVal.Type() != js.TypeUndefined && dataVal.Type() != js.TypeNull {
			json := js.Global().Get("JSON").Call("stringify", dataVal).String()
			dataCh <- []byte(json)
			return nil
		}
		errCh <- errors.New("getPicklistValuesByRecordType returned no data")
		return nil
	}))
	select {
	case data := <-dataCh:
		return UnmarshalPicklistFieldValues(data)
	case err := <-errCh:
		return nil, err
	}
}

// GetObjectInfo fetches SObject metadata for the given object.
// It delegates to the global JavaScript getObjectInfo function.
// Returns an ObjectInfo struct or an error if the JS API returns an error.
func GetObjectInfo(objectName string) (ObjectInfo, error) {
	dataCh := make(chan []byte)
	errCh := make(chan error)
	config := map[string]interface{}{
		"objectApiName": objectName,
	}
	js.Global().Call("getObjectInfo", js.ValueOf(config), js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		result := args[0]
		if errVal := result.Get("error"); errVal.Type() != js.TypeUndefined && errVal.Type() != js.TypeNull {
			errCh <- errors.New(errVal.String())
			return nil
		}
		if dataVal := result.Get("data"); dataVal.Type() != js.TypeUndefined && dataVal.Type() != js.TypeNull {
			json := js.Global().Get("JSON").Call("stringify", dataVal).String()
			dataCh <- []byte(json)
			return nil
		}
		errCh <- errors.New("getObjectInfo returned no data")
		return nil
	}))
	select {
	case data := <-dataCh:
		return UnmarshalObjectInfo(data)
	case err := <-errCh:
		return ObjectInfo{}, err
	}
}
