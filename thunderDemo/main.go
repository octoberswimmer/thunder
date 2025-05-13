package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

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

// InputMsg represents user input in the text field.
type InputMsg struct{ Value string }

// LimitChangeMsg represents changing the fetch limit via select dropdown.
type LimitChangeMsg struct{ Limit string }

// CheckboxMsg represents toggling the demo checkbox filter.
type CheckboxMsg struct{ Checked bool }

// FilterModeMsg represents changing the filter mode (contains vs startsWith).
type FilterModeMsg struct{ Mode string }

// ShowToastMsg represents clicking the show toast button.
type ShowToastMsg struct{}

// HideToastMsg represents closing the toast notification.
type HideToastMsg struct{}

// AppModel holds application state (input value, fetch limit, checkbox filter, rows, modal/toast visibility).
type AppModel struct {
	masc.Core
	InputValue  string
	Limit       string
	FilterAOnly bool
	FilterMode  string
	Rows        []map[string]string
	ShowModal   bool
	ShowToast   bool
}

// Init returns no initial command.
// Init sets default limit on startup.
func (m *AppModel) Init() masc.Cmd {
	// Default fetch limit and filter mode
	m.Limit = "5"
	m.FilterMode = "contains"
	return nil
}

// Update handles messages and returns commands.
func (m *AppModel) Update(msg masc.Msg) (masc.Model, masc.Cmd) {
	switch msg.(type) {
	case InputMsg:
		// Update input filter value
		m.InputValue = msg.(InputMsg).Value
		return m, nil
	case LimitChangeMsg:
		// Update fetch limit and refetch
		m.Limit = msg.(LimitChangeMsg).Limit
		return m, fetchAccountsCmd(m.Limit)
	case CheckboxMsg:
		// Update checkbox filter
		m.FilterAOnly = msg.(CheckboxMsg).Checked
		return m, nil
	case FilterModeMsg:
		// Update filter mode
		m.FilterMode = msg.(FilterModeMsg).Mode
		return m, nil
	case FetchAccountsMsg:
		// Trigger asynchronous fetch command with selected limit
		return m, fetchAccountsCmd(m.Limit)
	case AccountsFetchedMsg:
		// Update model with fetched rows
		m.Rows = msg.(AccountsFetchedMsg).Rows
		return m, nil
	case ToggleModalMsg:
		// Toggle modal visibility
		m.ShowModal = !m.ShowModal
		return m, nil
	case ShowToastMsg:
		// Show toast notification and schedule auto-hide
		m.ShowToast = true
		return m, autoHideToastCmd()
	case HideToastMsg:
		// Hide toast notification
		m.ShowToast = false
		return m, nil
	default:
		return m, nil
	}
}

// Render renders the button or the data table based on state.
func (m *AppModel) Render(send func(masc.Msg)) masc.ComponentOrHTML {
	// Build primary action buttons
	elems := []masc.MarkupOrChild{
		components.Button("Fetch Accounts", components.VariantBrand, func(e *masc.Event) {
			send(FetchAccountsMsg{})
		}),
		components.Button("Show Modal", components.VariantNeutral, func(e *masc.Event) {
			send(ToggleModalMsg{})
		}),
		components.Button("Show Toast", components.VariantNeutral, func(e *masc.Event) {
			send(ShowToastMsg{})
		}),
	}
	if len(m.Rows) > 0 {
		// Render limit select with spacing above
		elems = append(elems,
			elem.Div(
				masc.Markup(masc.Class("slds-m-top_medium")),
				components.Select(
					"Limit",
					[]components.SelectOption{{Label: "5", Value: "5"}, {Label: "10", Value: "10"}, {Label: "20", Value: "20"}},
					m.Limit,
					func(e *masc.Event) {
						send(LimitChangeMsg{Limit: e.Target.Get("value").String()})
					},
				),
			),
		)
		// Render checkbox filter with spacing above
		elems = append(elems,
			elem.Div(
				masc.Markup(masc.Class("slds-m-top_medium")),
				components.Checkbox(
					"Only names containing 'A'",
					m.FilterAOnly,
					func(e *masc.Event) {
						// Toggle checkbox filter state
						send(CheckboxMsg{Checked: !m.FilterAOnly})
					},
				),
			),
		)
		// Render filter mode radio group with spacing above
		elems = append(elems,
			elem.Div(
				masc.Markup(masc.Class("slds-m-top_medium")),
				components.RadioGroup(
					"filtermode",
					"Filter Mode",
					[]components.RadioOption{
						{Label: "Contains", Value: "contains"},
						{Label: "Starts With", Value: "startswith"},
					},
					m.FilterMode,
					func(mode string) {
						send(FilterModeMsg{Mode: mode})
					},
				),
			),
		)
		// Render filter input with spacing above
		elems = append(elems,
			elem.Div(
				masc.Markup(masc.Class("slds-m-top_medium")),
				components.TextInput("Filter by Name", m.InputValue, "Enter substring", func(e *masc.Event) {
					send(InputMsg{Value: e.Target.Get("value").String()})
				}),
			),
		)
		// Filter rows by input substring and checkbox
		var filtered []map[string]string
		query := strings.ToLower(m.InputValue)
		for _, r := range m.Rows {
			name := r["Name"]
			lower := strings.ToLower(name)
			// Determine if query matches based on selected mode
			var matchesQuery bool
			if query == "" {
				matchesQuery = true
			} else if m.FilterMode == "contains" {
				matchesQuery = strings.Contains(lower, query)
			} else if m.FilterMode == "startswith" {
				matchesQuery = strings.HasPrefix(lower, query)
			}
			// Apply checkbox filter for letter 'a'
			if matchesQuery && (!m.FilterAOnly || strings.Contains(lower, "a")) {
				filtered = append(filtered, r)
			}
		}
		// Render data table with spacing above
		elems = append(elems,
			elem.Div(
				masc.Markup(masc.Class("slds-m-top_medium")),
				components.DataTable([]string{"Name"}, filtered),
			),
		)
	}
	// Compose main content: header and card
	children := []masc.MarkupOrChild{
		components.PageHeader(
			"Thunder Demo",
			fmt.Sprintf("Mode: %s; Only A: %t", m.FilterMode, m.FilterAOnly),
		),
		components.Card("Accounts", elems...),
	}
	// Append modal overlay if toggled
	if m.ShowModal {
		children = append(children,
			components.Modal("Demo Modal",
				masc.Text("This is a demo modal"),
				components.Button("Close", components.VariantNeutral, func(e *masc.Event) {
					send(ToggleModalMsg{})
				}),
			),
		)
	}
	// Append toast notification if toggled
	if m.ShowToast {
		children = append(children,
			components.Toast(components.VariantSuccess,
				"Success",
				"This is a toast notification.",
				func(e *masc.Event) { send(HideToastMsg{}) },
			),
		)
	}
	// Optionally append modal and toast overlays
	if m.ShowModal {
		children = append(children,
			components.Modal("Demo Modal",
				masc.Text("This is a demo modal"),
				components.Button("Close", components.VariantNeutral, func(e *masc.Event) {
					send(ToggleModalMsg{})
				}),
			),
		)
	}
	if m.ShowToast {
		children = append(children,
			components.Toast(components.VariantSuccess,
				"Success",
				"This is a toast notification.",
				func(e *masc.Event) { send(HideToastMsg{}) },
			),
		)
	}
	return elem.Div(children...)
}

// fetchAccountsCmd creates a Cmd that fetches accounts via JS and returns a Msg.
// It uses the provided limit value for the SOQL query.
func fetchAccountsCmd(limit string) masc.Cmd {
	return func() masc.Msg {
		// Build query with dynamic LIMIT
		url := "/services/data/v58.0/query?q=SELECT+Name+FROM+Account+LIMIT+" + limit
		data := api.Get(url)
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

// autoHideToastCmd returns a Cmd that waits then hides the toast
func autoHideToastCmd() masc.Cmd {
	return func() masc.Msg {
		time.Sleep(3 * time.Second)
		return HideToastMsg{}
	}
}

func main() {
	thunder.Run(&AppModel{})
}
