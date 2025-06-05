//go:build js && !dev
// +build js,!dev

package api

import (
	"errors"
	"log"
	"sync"
	"syscall/js"
	"time"
)

// auraMutex serializes all Aura-enabled calls to prevent Lightning XHR connection issues
var auraMutex sync.Mutex

// GetPicklistValuesByRecordType fetches picklist values for the given object and record type.
// It delegates to the global JavaScript getPicklistValuesByRecordType function.
// Returns a map of field names to PicklistFieldValue or an error.
// Aura calls are serialized to prevent Lightning XHR connection pool issues.
func GetPicklistValuesByRecordType(objectName, recordTypeId string) (map[string]PicklistFieldValue, error) {
	auraMutex.Lock()

	dataCh := make(chan []byte, 1) // Buffered channel to prevent blocking
	errCh := make(chan error, 1)   // Buffered channel to prevent blocking
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

	// Release mutex after setting up the callback but before waiting
	auraMutex.Unlock()

	for {
		select {
		case data := <-dataCh:
			return UnmarshalPicklistFieldValues(data)
		case err := <-errCh:
			return nil, err
		case <-time.After(2 * time.Second):
			log.Printf("Still waiting for GetPicklistValuesByRecordType %s %s response, continuing...", objectName, recordTypeId)
			continue
		}
	}
}

// GetObjectInfo fetches SObject metadata for the given object.
// It delegates to the global JavaScript getObjectInfo function.
// Returns an ObjectInfo struct or an error if the JS API returns an error.
// Aura calls are serialized to prevent Lightning XHR connection pool issues.
func GetObjectInfo(objectName string) (ObjectInfo, error) {
	auraMutex.Lock()

	dataCh := make(chan []byte, 1) // Buffered channel to prevent blocking
	errCh := make(chan error, 1)   // Buffered channel to prevent blocking
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

	// Release mutex after setting up the callback but before waiting
	auraMutex.Unlock()

	for {
		select {
		case data := <-dataCh:
			return UnmarshalObjectInfo(data)
		case err := <-errCh:
			return ObjectInfo{}, err
		case <-time.After(2 * time.Second):
			log.Printf("Still waiting for GetObjectInfo %s response, continuing...", objectName)
			continue
		}
	}
}
