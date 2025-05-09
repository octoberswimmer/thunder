package thunder

import (
	"syscall/js"

	"github.com/octoberswimmer/masc"
)

func Run(model masc.Model) {
	// Register startWithDiv: thunder host calls this
	js.Global().Set("startWithDiv", js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		div := args[0]
		// Launch Masc program rendering into this div
		go masc.NewProgram(
			model,
			masc.RenderTo(div),
		).Run()
		return nil
	}))
	// Keep Go runtime alive
	select {}
}
