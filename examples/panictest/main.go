package main

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/thunder"
	"github.com/octoberswimmer/thunder/components"
)

type PanicTestMsg struct{}

type PanicTestModel struct {
	masc.Core
}

func (m *PanicTestModel) Init() masc.Cmd {
	return nil
}

func (m *PanicTestModel) Update(msg masc.Msg) (masc.Model, masc.Cmd) {
	switch msg.(type) {
	case PanicTestMsg:
		panic("Test panic: This is a deliberate panic to test error handling!")
	}
	return m, nil
}

func (m *PanicTestModel) Render(send func(masc.Msg)) masc.ComponentOrHTML {
	return components.Page(
		components.PageHeader("Panic Test Application", "Test the panic handler"),
		elem.Div(
			masc.Markup(masc.Class("slds-m-around_medium")),
			elem.Div(
				masc.Text("Click the button below to trigger a panic and see the error modal:"),
			),
			elem.Div(
				masc.Markup(masc.Class("slds-m-top_medium")),
				components.Button(
					"Trigger Panic",
					components.VariantDestructive,
					func(e *masc.Event) {
						send(PanicTestMsg{})
					},
				),
			),
		),
	)
}

func main() {
	thunder.Run(&PanicTestModel{})
}
