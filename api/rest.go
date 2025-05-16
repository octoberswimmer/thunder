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
