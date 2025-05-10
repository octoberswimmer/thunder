package main

import (
	"encoding/json"

	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/thunder"
	"github.com/octoberswimmer/thunder/api"
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
		// Add margin between the fetch button and the data table
		elems = append(elems,
			elem.Div(
				masc.Markup(masc.Class("slds-m-top_medium")),
				components.DataTable([]string{"Name"}, m.Rows),
			),
		)
	}
	// Assemble header and content inside a container div
	return elem.Div(
		// Page header with title and optional subtitle
		components.PageHeader(
			"Thunder Demo",
			"Go/WASM SLDS component demo",
		),
		// Card containing button and data table
		components.Card("Accounts", elems...),
	)
}

// fetchAccountsCmd creates a Cmd that fetches accounts via JS and returns a Msg.
func fetchAccountsCmd() masc.Cmd {
	return func() masc.Msg {
		data := api.Get("/services/data/v58.0/query?q=SELECT+Name+FROM+Account+LIMIT+5")
		var result map[string]any
		err := json.Unmarshal(data, &result)
		if err != nil {
			panic(err.Error())
		}
		recs := result["records"].([]any)
		rows := make([]map[string]string, len(recs))
		for i, r := range recs {
			v := r.(map[string]any)
			name := v["Name"].(string)
			rows[i] = map[string]string{"Name": name}
		}
		return AccountsFetchedMsg{Rows: rows}
	}
}

func main() {
	thunder.Run(&AppModel{})
}
