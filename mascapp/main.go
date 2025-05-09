package main

import (
	"syscall/js"

	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/thunder/components"
)

// Msg types
// FetchAccountsMsg represents the user clicking the fetch button.
type FetchAccountsMsg struct{}

// AccountsFetchedMsg carries rows retrieved from the REST proxy.
type AccountsFetchedMsg struct{ Rows []map[string]string }

// AppModel holds application state (rows) and implements masc.Model.
type AppModel struct {
	masc.Core
	Rows []map[string]string
}

// Init returns no initial command.
func (m *AppModel) Init() masc.Cmd { return nil }

// Update handles messages and returns commands.
func (m *AppModel) Update(msg masc.Msg) (masc.Model, masc.Cmd) {
	switch msg.(type) {
	case FetchAccountsMsg:
		// Trigger asynchronous fetch command
		return m, fetchAccountsCmd()
	case AccountsFetchedMsg:
		// Update model with fetched rows
		m.Rows = msg.(AccountsFetchedMsg).Rows
		return m, nil
	default:
		return m, nil
	}
}

// Render renders the button or the data table based on state.
func (m *AppModel) Render(send func(masc.Msg)) masc.ComponentOrHTML {
	elems := []masc.MarkupOrChild{}
	elems = append(elems, components.Button("Fetch Accounts", components.VariantBrand, func(e *masc.Event) {
		send(FetchAccountsMsg{})
	}))
	if len(m.Rows) > 0 {
		elems = append(elems, components.DataTable([]string{"Name"}, m.Rows))
	}
	return elem.Div(
		elems...,
	)
}

// fetchAccountsCmd creates a Cmd that fetches accounts via JS and returns a Msg.
func fetchAccountsCmd() masc.Cmd {
	return func() masc.Msg {
		ch := make(chan []map[string]string)
		// Call global get() proxy to SOQL endpoint
		js.Global().Call(
			"get",
			"/services/data/v58.0/query?q=SELECT+Name+FROM+Account+LIMIT+5",
		).Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			parsed := js.Global().Get("JSON").Call("parse", args[0].String())
			recs := parsed.Get("records")
			n := recs.Length()
			rows := make([]map[string]string, n)
			for i := 0; i < n; i++ {
				r := recs.Index(i)
				rows[i] = map[string]string{"Name": r.Get("Name").String()}
			}
			ch <- rows
			return nil
		}))
		// Wait for JS promise callback to send rows
		rows := <-ch
		return AccountsFetchedMsg{Rows: rows}
	}
}

func main() {
	// Register startWithDiv: Vecty host calls this
	js.Global().Set("startWithDiv", js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		div := args[0]
		// Launch Masc program rendering into this div
		go masc.NewProgram(
			&AppModel{},
			masc.RenderTo(div),
		).Run()
		return nil
	}))
	// Keep Go runtime alive
	select {}
}
