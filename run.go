//go:build js && !dev
// +build js,!dev

package thunder

import (
	"syscall/js"

	"github.com/octoberswimmer/masc"
)

// currentDiv stores the div element for the current Thunder instance
var currentDiv js.Value

// Run initializes a Thunder application for production deployment in Salesforce Lightning.
// It registers the global "startWithDiv" JavaScript function that Lightning Web Components
// call to launch the Go WASM application within a specific DOM element.
//
// The model parameter should implement masc.Model with Init(), Update(), and Render() methods
// following the Elm Architecture pattern. When deployed to Salesforce, the Thunder LWC will
// invoke startWithDiv with a target div element to render the application.
//
// This function blocks indefinitely to keep the Go runtime alive for the duration of the
// Lightning page session. It should only be called from main() in production builds.
func Run(model masc.Model) {
	// Register startWithDiv: thunder host calls this
	js.Global().Set("startWithDiv", js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		div := args[0]
		// Store the div element for this instance
		currentDiv = div
		// Launch Masc program rendering into this div with panic handler
		go func() {
			defer handlePanic()
			masc.NewProgram(
				model,
				masc.RenderTo(div),
				masc.WithoutCatchPanics(),
			).Run()
		}()
		return nil
	}))
	// Keep Go runtime alive
	select {}
}

// GetCurrentDiv returns the current div element for the Thunder instance
func GetCurrentDiv() js.Value {
	return currentDiv
}
