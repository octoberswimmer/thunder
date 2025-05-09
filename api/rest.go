package api

import (
	"syscall/js"
)

func Get(url string) []byte {
	ch := make(chan []byte)
	// Call global get() proxy to SOQL endpoint
	js.Global().Call("get", url).
		Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			ch <- []byte(args[0].String())
			return nil
		}))
	// Wait for JS promise callback to send rows
	v := <-ch
	return v
}
