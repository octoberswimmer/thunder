//go:build dev
// +build dev

package thunder

import (
	"syscall/js"

	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/thunder/internal/panichandler"
	"github.com/octoberswimmer/thunder/internal/runtime"
)

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
	runtime.SetCurrentDiv(div)

	// Set up panic recovery
	defer panichandler.HandlePanic()

	// Create a panic handler that calls Thunder's panic handler
	thunderPanicHandler := func(panicValue interface{}) {
		js.Global().Get("console").Call("log", "Thunder panic handler called with:", panicValue)
		panichandler.ShowPanicModal(panicValue)
	}

	pgm := masc.NewProgram(model, masc.RenderTo(div), masc.WithPanicHandler(thunderPanicHandler))
	_, err := pgm.Run()
	if err != nil {
		panic(err)
	}
}

// GetCurrentDiv returns the current div element for the Thunder instance
func GetCurrentDiv() js.Value {
	return runtime.GetCurrentDiv()
}
