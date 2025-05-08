package main

import (
    "syscall/js"
    // Uncomment and use masc APIs instead of direct JS calls:
    // "github.com/octoberswimmer/masc"
)

func main() {
    // Minimal Go WASM using syscall/js to invoke the vecty LWC proxy methods

    // Export a startWithDiv function for the vecty LWC to call
    js.Global().Set("startWithDiv", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        // args[0] is the container div element passed in by vecty.js
        div := args[0]
        // Create a heading
        h1 := js.Global().Get("document").Call("createElement", "h1")
        h1.Set("innerText", "MASC Example")
        div.Call("appendChild", h1)
        // Create a button that triggers a SOQL query via the vecty LWC proxy
        btn := js.Global().Get("document").Call("createElement", "button")
        btn.Set("innerText", "Query Accounts")
        btn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
            // Use the global 'get' function injected by the vecty LWC
            promise := js.Global().Call("get", "/services/data/v58.0/query?q=SELECT+Name+FROM+Account+LIMIT+5")
            promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
                pre := js.Global().Get("document").Call("createElement", "pre")
                pre.Set("innerText", args[0].String())
                div.Call("appendChild", pre)
                return nil
            }))
            return nil
        }))
        div.Call("appendChild", btn)
        return nil
    }))

    // Prevent the Go program from exiting
    select {}
}