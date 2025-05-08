package main

import (
   "syscall/js"

   "github.com/octoberswimmer/masc"
   "github.com/octoberswimmer/masc/elem"
   "github.com/octoberswimmer/thunder/components"
)

// AppModel is a Masc model that renders a single SLDS button.
type AppModel struct {
	masc.Core
}

// Init implements masc.Model.
func (m *AppModel) Init() masc.Cmd { return nil }

// Update implements masc.Model.
func (m *AppModel) Update(msg masc.Msg) (masc.Model, masc.Cmd) {
	// No state updates required
	return m, nil
}

// Render returns the SLDS-styled button wrapped in a <div>.
func (m *AppModel) Render(send func(masc.Msg)) masc.ComponentOrHTML {
   return elem.Div(
       components.Button("Query Accounts", components.VariantBrand, func(e *masc.Event) {
           // Use vecty-provided global 'get' for SOQL
           promise := js.Global().Call("get", "/services/data/v58.0/query?q=SELECT+Name+FROM+Account+LIMIT+5")
           promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
               pre := js.Global().Get("document").Call("createElement", "pre")
               pre.Set("innerText", args[0].String())
               // Append to body
               js.Global().Get("document").Get("body").Call("appendChild", pre)
               return nil
           }))
       }),
   )
}

func main() {
	// Register startWithDiv callback for vecty host
	js.Global().Set("startWithDiv", js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		div := args[0]
		// Launch Masc program, rendering into the provided div
		go func() {
			program := masc.NewProgram(&AppModel{}, masc.RenderTo(div))
			program.Run()
		}()
		return nil
	}))
	// Keep Go runtime alive
	select {}
}
