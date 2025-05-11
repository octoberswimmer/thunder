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

// ToggleModalMsg represents toggling the demo modal visibility.
type ToggleModalMsg struct{}

// AppModel holds application state (rows, modal visibility) and implements masc.Model.
type AppModel struct {
	masc.Core
	Rows      []map[string]string
	ShowModal bool
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
	case ToggleModalMsg:
		// Toggle modal visibility
		m.ShowModal = !m.ShowModal
		return m, nil
	default:
		return m, nil
	}
}

// Render renders the button or the data table based on state.
func (m *AppModel) Render(send func(masc.Msg)) masc.ComponentOrHTML {
	// Build action buttons
	elems := []masc.MarkupOrChild{
		components.Button("Fetch Accounts", components.VariantBrand, func(e *masc.Event) {
			send(FetchAccountsMsg{})
		}),
		components.Button("Show Modal", components.VariantNeutral, func(e *masc.Event) {
			send(ToggleModalMsg{})
		}),
	}
	// Include data table if rows fetched
	if len(m.Rows) > 0 {
		elems = append(elems,
			elem.Div(
				masc.Markup(masc.Class("slds-m-top_medium")),
				components.DataTable([]string{"Name"}, m.Rows),
			),
		)
	}
	// Compose main content: header and card
	children := []masc.MarkupOrChild{
		components.PageHeader(
			"Thunder Demo",
			"Go/WASM SLDS component demo",
		),
		components.Card("Accounts", elems...),
	}
	// Append modal overlay if toggled
	if m.ShowModal {
		// Show modal with close button inside
		children = append(children,
			components.Modal("Demo Modal",
				masc.Text("This is a demo modal"),
				components.Button("Close", components.VariantNeutral, func(e *masc.Event) {
					send(ToggleModalMsg{})
				}),
			),
		)
	}
	return elem.Div(children...)
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
