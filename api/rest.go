//go:build js && !dev
// +build js,!dev

package api

import (
	"errors"
	"syscall/js"
)

// Get performs a GET via JS proxy, automatically following Salesforce cursor pagination.
// It returns an error if the underlying promise is rejected.
func Get(url string) ([]byte, error) {
	dataCh := make(chan []byte)
	errCh := make(chan error)
	js.Global().Call("get", url).
		Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			dataCh <- []byte(args[0].String())
			return nil
		})).
		Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// args[0] is the error from JS
			errCh <- errors.New(args[0].String())
			return nil
		}))
	select {
	case data := <-dataCh:
		return data, nil
	case err := <-errCh:
		return nil, err
	}
}

// Post performs a POST via JS proxy, automatically following Salesforce cursor pagination.
// It returns an error if the underlying promise is rejected.
// For composite requests, it returns CompositeErrors if any sub-requests fail.
func Post(url string, body []byte) ([]byte, error) {
	dataCh := make(chan []byte)
	errCh := make(chan error)
	js.Global().Call("post", url, string(body)).
		Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			dataCh <- []byte(args[0].String())
			return nil
		})).
		Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			errCh <- errors.New(args[0].String())
			return nil
		}))
	select {
	case data := <-dataCh:
		// Check if this is a composite request response
		if isCompositeRequest(url, body) {
			if compositeErrs, err := parseCompositeResponse(data); err == nil && compositeErrs.HasErrors() {
				return data, compositeErrs
			}
		}
		return data, nil
	case err := <-errCh:
		return nil, err
	}
}

// Patch performs a PATCH via JS proxy, automatically following Salesforce cursor pagination.
// It returns an error if the underlying promise is rejected.
func Patch(url string, body []byte) ([]byte, error) {
	dataCh := make(chan []byte)
	errCh := make(chan error)
	js.Global().Call("patch", url, string(body)).
		Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			dataCh <- []byte(args[0].String())
			return nil
		})).
		Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			errCh <- errors.New(args[0].String())
			return nil
		}))
	select {
	case data := <-dataCh:
		return data, nil
	case err := <-errCh:
		return nil, err
	}
}

// Delete performs a DELETE via JS proxy.
// It returns an error if the underlying promise is rejected.
func Delete(url string) ([]byte, error) {
	dataCh := make(chan []byte)
	errCh := make(chan error)
	js.Global().Call("delete", url).
		Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			dataCh <- []byte(args[0].String())
			return nil
		})).
		Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			errCh <- errors.New(args[0].String())
			return nil
		}))
	select {
	case data := <-dataCh:
		return data, nil
	case err := <-errCh:
		return nil, err
	}
}
