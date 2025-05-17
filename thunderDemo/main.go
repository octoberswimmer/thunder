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

// TabChangeMsg represents selecting a different tab in the UI.
type TabChangeMsg struct{ Tab string }

// HideToastMsg represents closing the toast notification.
// HideToastMsg represents closing the toast notification.
type HideToastMsg struct{}

// QueryErrorMsg carries an error string when a query fails.
type QueryErrorMsg struct{ Err string }

// AppModel holds application state (input value, fetch limit, checkbox filter, rows, modal/toast visibility).
type AppModel struct {
	masc.Core
	// UI state
	InputValue  string
	Limit       string
	FilterAOnly bool
	FilterMode  string
	SelectedTab string
	// Data
	Rows []map[string]string
	// Overlays and toast state
	ShowModal    bool
	ShowToast    bool
	ToastVariant components.ToastVariant
	ToastHeader  string
	ToastMessage string
	// Loading flag for data fetch
	Loading bool
}

// Init returns no initial command.
// Init sets default limit on startup.
func (m *AppModel) Init() masc.Cmd {
	// Default fetch limit, filter mode, and selected tab
	m.Limit = "5"
	m.FilterMode = "contains"
	m.SelectedTab = "actions"
	m.Loading = false
	return nil
}

// Update handles messages and returns commands.
func (m *AppModel) Update(msg masc.Msg) (masc.Model, masc.Cmd) {
	switch msg := msg.(type) {
	case InputMsg:
		// Update input filter value
		m.InputValue = msg.Value
		return m, nil
	case LimitChangeMsg:
		// Update fetch limit and refetch
		m.Limit = msg.Limit
		cmd := m.fetchAccountsCmd(m.Limit)
		return m, cmd
	case CheckboxMsg:
		// Update checkbox filter
		m.FilterAOnly = msg.Checked
		return m, nil
	case FilterModeMsg:
		// Update filter mode
		m.FilterMode = msg.Mode
		return m, nil
	case FetchAccountsMsg:
		// Trigger asynchronous fetch command with selected limit
		m.SelectedTab = "data"
		return m, m.fetchAccountsCmd(m.Limit)
	case AccountsFetchedMsg:
		// Update model with fetched rows
		m.Rows = msg.Rows
		m.Loading = false
		return m, nil
	case QueryErrorMsg:
		// Display error toast on query failure
		m.Loading = false
		m.ShowToast = true
		m.ToastVariant = components.VariantError
		m.ToastHeader = "Error"
		m.ToastMessage = msg.Err
		return m, autoHideToastCmd()
	case ToggleModalMsg:
		// Toggle modal visibility
		m.ShowModal = !m.ShowModal
		return m, nil
	case ShowToastMsg:
		// Show success toast notification and schedule auto-hide
		m.ShowToast = true
		m.ToastVariant = components.VariantSuccess
		m.ToastHeader = "Success"
		m.ToastMessage = "This is a toast notification."
		return m, autoHideToastCmd()
	case HideToastMsg:
		// Hide toast notification
		m.ShowToast = false
		return m, nil
	case TabChangeMsg:
		m.SelectedTab = msg.Tab
		return m, nil
	default:
		return m, nil
	}
}

// Render renders the button or the data table based on state.
func (m *AppModel) Render(send func(masc.Msg)) masc.ComponentOrHTML {
	// Build two tab panes: actions and data
	// Actions pane: primary buttons
	actions := []masc.MarkupOrChild{
		components.Button("Fetch Accounts", components.VariantBrand, func(e *masc.Event) {
			send(FetchAccountsMsg{})
		}),
		components.Button("Show Modal", components.VariantNeutral, func(e *masc.Event) {
			send(ToggleModalMsg{})
		}),
		components.Button("Show Toast", components.VariantNeutral, func(e *masc.Event) {
			send(ShowToastMsg{})
		}),
		// Demo Badge and Pill components
		components.Badge("Demo Badge"),
		components.Pill("Tag1", func(e *masc.Event) {
			send(ShowToastMsg{})
		}),
		components.Pill("Tag2", nil),
		// Demo Icon component
		components.Icon(components.UtilityIcon, "settings", components.IconSmall),
	}
	// Data pane: spinner, filters, and table
	var data []masc.MarkupOrChild
	if m.Loading {
		// Show spinner while loading
		data = append(data,
			elem.Div(
				masc.Markup(masc.Class("slds-m-top_medium", "slds-align_absolute-center")),
				components.Spinner("medium"),
			),
		)
	} else if len(m.Rows) > 0 {
		// Limit select
		data = append(data,
			elem.Div(
				masc.Markup(masc.Class("slds-m-top_medium")),
				components.Select(
					"Limit",
					[]components.SelectOption{
						{Label: "5", Value: "5"},
						{Label: "10", Value: "10"},
						{Label: "20", Value: "20"},
						{Label: "5,000", Value: "5000"},
					},
					m.Limit,
					func(e *masc.Event) {
						send(LimitChangeMsg{Limit: e.Target.Get("value").String()})
					},
				),
			),
		)
		// Checkbox filter
		data = append(data,
			elem.Div(
				masc.Markup(masc.Class("slds-m-top_medium")),
				components.Checkbox(
					"Only names containing 'A'",
					m.FilterAOnly,
					func(e *masc.Event) {
						send(CheckboxMsg{Checked: !m.FilterAOnly})
					},
				),
			),
		)
		// Radio filter mode
		data = append(data,
			elem.Div(
				masc.Markup(masc.Class("slds-m-top_medium")),
				components.RadioGroup(
					"filtermode",
					"Filter Mode",
					[]components.RadioOption{{Label: "Contains", Value: "contains"}, {Label: "Starts With", Value: "startswith"}},
					m.FilterMode,
					func(mode string) { send(FilterModeMsg{Mode: mode}) },
				),
			),
		)
		// Text input filter
		data = append(data,
			elem.Div(
				masc.Markup(masc.Class("slds-m-top_medium")),
				components.TextInput("Filter by Name", m.InputValue, "Enter substring", func(e *masc.Event) {
					send(InputMsg{Value: e.Target.Get("value").String()})
				}),
			),
		)
		// Lookup filter for names
		// Build suggestions list from all account names
		var suggestions []components.LookupOption
		for _, r := range m.Rows {
			name := r["Name"]
			suggestions = append(suggestions, components.LookupOption{Label: name, Value: name})
		}
		data = append(data,
			elem.Div(
				masc.Markup(masc.Class("slds-m-top_medium")),
				components.Lookup(
					"Select by Name",
					suggestions,
					m.InputValue,
					func(val string) { send(InputMsg{Value: val}) },
					func(val string) { send(InputMsg{Value: val}) },
				),
			),
		)
		// Apply filters
		var filtered []map[string]string
		query := strings.ToLower(m.InputValue)
		for _, r := range m.Rows {
			name := r["Name"]
			lower := strings.ToLower(name)
			var match bool
			if query == "" || (m.FilterMode == "contains" && strings.Contains(lower, query)) || (m.FilterMode == "startswith" && strings.HasPrefix(lower, query)) {
				match = true
			}
			if match && (!m.FilterAOnly || strings.Contains(lower, "a")) {
				filtered = append(filtered, r)
			}
		}
		// Show filtering progress
		if len(m.Rows) > 0 {
			percent := len(filtered) * 100 / len(m.Rows)
			data = append(data,
				elem.Div(
					masc.Markup(masc.Class("slds-m-top_medium")),
					components.ProgressBar(percent),
				),
			)
		}
		// Data table
		data = append(data,
			elem.Div(
				masc.Markup(masc.Class("slds-m-top_medium")),
				components.DataTable([]string{"Name"}, filtered),
			),
		)
	}
	// Compose tabs inside the Accounts card
	tabs := components.Tabs(
		"accounts-tabs",
		[]components.TabOption{
			{Label: "Actions", Value: "actions", Content: elem.Div(actions...)},
			{Label: "Data", Value: "data", Content: elem.Div(data...)},
		},
		m.SelectedTab,
		func(tab string) { send(TabChangeMsg{Tab: tab}) },
	)
	// Build page layout with header and card
	header := components.PageHeader(
		"Thunder Demo",
		fmt.Sprintf("Mode: %s; Only A: %t", m.FilterMode, m.FilterAOnly),
	)
	card := components.Card("Accounts", tabs)
	// Page wraps header and card
	pageLayout := components.Page(header, card)
	// Include breadcrumb at top with padding and margin
	rawCrumbs := components.Breadcrumb([]components.BreadcrumbOption{
		{Label: "Home", Href: "#"},
		{Label: "Thunder Demo", Href: "#"},
	})
	crumbs := elem.Div(
		masc.Markup(masc.Class("slds-p-horizontal_medium", "slds-m-bottom_small")),
		rawCrumbs,
	)
	children := []masc.MarkupOrChild{crumbs, pageLayout}
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
			components.Toast(
				m.ToastVariant,
				m.ToastHeader,
				m.ToastMessage,
				func(e *masc.Event) { send(HideToastMsg{}) },
			),
		)
	}
	return elem.Div(children...)
}

// fetchAccountsCmd creates a Cmd that fetches accounts via JS and returns a Msg.
// It uses the provided limit value for the SOQL query.
func (m *AppModel) fetchAccountsCmd(limit string) masc.Cmd {
	m.Loading = true
	return func() masc.Msg {
		// Perform SOQL query via Query API
		soql := fmt.Sprintf("SELECT Name FROM Account LIMIT %s", limit)
		data, err := api.Query(soql)
		if err != nil {
			return QueryErrorMsg{Err: err.Error()}
		}
		var result map[string]any
		if err := json.Unmarshal(data, &result); err != nil {
			return QueryErrorMsg{Err: err.Error()}
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
