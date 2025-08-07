//go:build dev
// +build dev

package thunder

import (
	"syscall/js"

	"github.com/octoberswimmer/masc"
)

// currentDiv stores the div element for the current Thunder instance
var currentDiv js.Value

// Run initializes a Thunder application for development mode with thunder serve.
// In development builds (with -tags dev), this function directly renders the application
// into the "app" div element that thunder serve provides in its HTML template.
//
// The model parameter should implement masc.Model with Init(), Update(), and Render() methods
// following the Elm Architecture pattern. Unlike production builds, this version runs
// synchronously and returns when the application exits.
//
// This function is only compiled when building with the "dev" build tag, typically
// used during local development with the Thunder CLI serve command.
func Run(model masc.Model) {
	doc := js.Global().Get("document")
	div := doc.Call("getElementById", "app")

	// Store the div element for this instance
	currentDiv = div

	pgm := masc.NewProgram(model, masc.RenderTo(div))
	_, err := pgm.Run()
	if err != nil {
		panic(err)
	}
}

// GetCurrentDiv returns the current div element for the Thunder instance
func GetCurrentDiv() js.Value {
	return currentDiv
}
