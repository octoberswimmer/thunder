//go:build (!js || !wasm) && !dev
// +build !js !wasm
// +build !dev

package thunder

import "github.com/octoberswimmer/masc"

// Run is a stub for non-WASM builds and will panic if called.
func Run(model masc.Model) {
	panic("thunder.Run is not supported outside the WASM environment")
}
