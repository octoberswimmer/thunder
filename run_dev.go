//go:build dev
// +build dev

package thunder

import (
	"syscall/js"

	"github.com/octoberswimmer/masc"
)

func Run(model masc.Model) {
	doc := js.Global().Get("document")
	div := doc.Call("getElementById", "app")

	pgm := masc.NewProgram(model, masc.RenderTo(div))
	_, err := pgm.Run()
	if err != nil {
		panic(err)
	}
}
