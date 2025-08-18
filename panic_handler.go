package thunder

import (
	"fmt"
	"runtime/debug"
	"syscall/js"

	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/thunder/components"
)

type panicModel struct {
	masc.Core
	panicMessage string
	stackTrace   string
}

func (m *panicModel) Init() masc.Cmd {
	return nil
}

func (m *panicModel) Update(msg masc.Msg) (masc.Model, masc.Cmd) {
	return m, nil
}

func (m *panicModel) Render(send func(masc.Msg)) masc.ComponentOrHTML {
	// Wrap the modal in a div since Modal returns a List
	return elem.Div(
		components.Modal(
			"Application Error",
			elem.Div(
				elem.Div(
					masc.Markup(masc.Class("slds-text-color_error", "slds-m-bottom_small")),
					elem.Strong(masc.Text("An unexpected error occurred:")),
				),
				elem.Div(
					masc.Markup(masc.Class("slds-box", "slds-theme_shade", "slds-m-bottom_small")),
					elem.Code(
						masc.Markup(
							masc.Style("display", "block"),
							masc.Style("white-space", "pre-wrap"),
							masc.Style("word-break", "break-word"),
						),
						masc.Text(m.panicMessage),
					),
				),
				elem.Details(
					elem.Summary(
						masc.Markup(masc.Class("slds-text-link")),
						masc.Text("Stack Trace"),
					),
					elem.Code(
						masc.Markup(
							masc.Class("slds-box", "slds-theme_shade", "slds-m-top_x-small"),
							masc.Style("display", "block"),
							masc.Style("white-space", "pre-wrap"),
							masc.Style("word-break", "break-word"),
							masc.Style("font-size", "0.75rem"),
							masc.Style("max-height", "300px"),
							masc.Style("overflow-y", "auto"),
						),
						masc.Text(m.stackTrace),
					),
				),
			),
		),
	)
}

func handlePanic() {
	if r := recover(); r != nil {
		panicMessage := fmt.Sprintf("%v", r)
		stackTrace := string(debug.Stack())

		// Get the document
		doc := js.Global().Get("document")
		body := doc.Get("body")

		// Clear the existing content
		div := GetCurrentDiv()
		if !div.IsUndefined() && !div.IsNull() {
			// Clear existing content by removing all children
			for div.Get("firstChild").Truthy() {
				div.Call("removeChild", div.Get("firstChild"))
			}
		}

		// Always create a new div for the error modal to avoid conflicts
		errorDiv := doc.Call("createElement", "div")
		errorDiv.Set("id", "thunder-panic-modal")
		errorDiv.Set("style", "position: fixed; top: 0; left: 0; right: 0; bottom: 0; z-index: 9999;")
		body.Call("appendChild", errorDiv)

		model := &panicModel{
			panicMessage: panicMessage,
			stackTrace:   stackTrace,
		}

		// Run the panic display program
		go func() {
			defer func() {
				// Catch any panics in the panic handler itself
				if r2 := recover(); r2 != nil {
					// Last resort: just show an alert
					js.Global().Call("alert", fmt.Sprintf("Application Error: %v\n\nPanic handler also failed: %v", r, r2))
				}
			}()

			masc.NewProgram(
				model,
				masc.RenderTo(errorDiv),
			).Run()
		}()

		// Keep the program alive
		select {}
	}
}
