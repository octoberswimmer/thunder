//go:build js && !dev
// +build js,!dev

package api

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"syscall/js"
	"time"
)

// getMutex serializes all GET requests to prevent Lightning XHR connection issues
var getMutex sync.Mutex

// Get performs a GET via JS proxy, automatically following Salesforce cursor pagination.
// It returns an error if the underlying promise is rejected.
// GET requests are serialized to prevent Lightning XHR connection pool issues.
func Get(url string) ([]byte, error) {
	if len(url) > 100 && url[50:60] == "PurchaserP" {
		panic(fmt.Sprintf("DEBUG: Get() called for PurchaserPlan URL: %s", url))
	}
	getMutex.Lock()

	dataCh := make(chan []byte, 1) // Buffered channel to prevent blocking
	errCh := make(chan error, 1)   // Buffered channel to prevent blocking

	go func() {

		// Check what the get call returns
		result := js.Global().Call("get", url)
		resultType := result.Type()
		if strings.Contains(url, "PurchaserPlan") {
			// panic("called get for Purchaser Plan query: " + url)
		}

		// Panic if our assumption about get() returning an object is wrong
		if resultType == js.TypeUndefined || resultType == js.TypeNull {
			panic(fmt.Sprintf("ASSUMPTION FAILED: get() returned %s for URL: %s", resultType.String(), url))
		}

		// Panic if our assumption about it being a thenable object is wrong
		thenMethod := result.Get("then")
		if thenMethod.Type() == js.TypeUndefined {
			panic(fmt.Sprintf("ASSUMPTION FAILED: promise object missing .then() method for URL: %s, result type: %s", url, resultType.String()))
		}
		if thenMethod.Type() != js.TypeFunction {
			panic(fmt.Sprintf("ASSUMPTION FAILED: .then() is not a function, got %s for URL: %s", thenMethod.Type().String(), url))
		}

		// Panic if our assumption about catch method is wrong
		catchMethod := result.Get("catch")
		if catchMethod.Type() != js.TypeFunction {
			panic(fmt.Sprintf("ASSUMPTION FAILED: .catch() is not a function, got %s for URL: %s", catchMethod.Type().String(), url))
		}

		then := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if strings.Contains(url, "PurchaserPlan") {
				// panic("then callback for Purchaser Plan called")
			}
			// Panic if our assumption about success callback args is wrong
			if len(args) != 1 {
				panic(fmt.Sprintf("ASSUMPTION FAILED: success callback expected 1 arg, got %d for URL: %s", len(args), url))
			}
			if args[0].Type() != js.TypeString {
				panic(fmt.Sprintf("ASSUMPTION FAILED: success callback arg not string, got %s for URL: %s", args[0].Type().String(), url))
			}
			dataCh <- []byte(args[0].String())
			return nil
		})

		catch := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if strings.Contains(url, "PurchaserPlan") {
				panic("catch callback for Purchaser Plan called")
			}
			// Panic if our assumption about error callback args is wrong
			if len(args) != 1 {
				panic(fmt.Sprintf("ASSUMPTION FAILED: error callback expected 1 arg, got %d for URL: %s", len(args), url))
			}
			errCh <- errors.New(args[0].String())
			return nil
		})

		result.Call("then", then)
		result.Call("catch", catch)

	}()

	// Release mutex after setting up the promise chain but before waiting
	getMutex.Unlock()

	for {
		select {
		case data := <-dataCh:
			return data, nil
		case err := <-errCh:
			return nil, err
		case <-time.After(2 * time.Second):
			// Continue waiting without logging to avoid deadlocks
			continue
		}
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
	for {
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
		case <-time.After(2 * time.Second):
			log.Printf("Still waiting for POST %s response, continuing...", url)
			continue
		}
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
	for {
		select {
		case data := <-dataCh:
			return data, nil
		case err := <-errCh:
			return nil, err
		case <-time.After(2 * time.Second):
			log.Printf("Still waiting for PATCH %s response, continuing...", url)
			continue
		}
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
	for {
		select {
		case data := <-dataCh:
			return data, nil
		case err := <-errCh:
			return nil, err
		case <-time.After(2 * time.Second):
			log.Printf("Still waiting for DELETE %s response, continuing...", url)
			continue
		}
	}
}
