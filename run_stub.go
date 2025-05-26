//go:build !js && !dev
// +build !js,!dev

package thunder

import "github.com/octoberswimmer/masc"

// Run initializes a Thunder application for deployment in Salesforce Lightning.
//
// In production builds (GOOS=js GOARCH=wasm), this function registers the global
// "startWithDiv" JavaScript function that Lightning Web Components call to launch
// the Go WASM application within a specific DOM element.
//
// In development builds (GOOS=js GOARCH=wasm -tags dev), this function directly
// renders the application into the "app" div element provided by thunder serve.
//
// The model parameter should implement masc.Model with Init(), Update(), and Render()
// methods following the Elm Architecture pattern.
//
// This function will panic if called outside a WebAssembly environment, as Thunder
// applications are designed to run only in the browser via Lightning Web Components.
//
// Example usage:
//
//	func main() {
//	    thunder.Run(&MyAppModel{})
//	}
func Run(model masc.Model) {
	panic("thunder.Run is not supported outside the WASM environment")
}
