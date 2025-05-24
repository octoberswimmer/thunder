package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/octoberswimmer/masc"
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

// LastModifiedDateChangeMsg represents selecting a date to filter Accounts by LastModifiedDate.
type LastModifiedDateChangeMsg struct{ Value string }

// ShowToastMsg represents clicking the show toast button.
type ShowToastMsg struct{}

// FetchObjectInfoMsg represents the user clicking the fetch object info button.
type FetchObjectInfoMsg struct{}

// ObjectInfoFetchedMsg carries object info retrieved from the UI API.
type ObjectInfoFetchedMsg struct{ Info api.ObjectInfo }

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
	// LastModifiedDate is the date filter for querying Accounts.
	LastModifiedDate string
	// Data
	Rows       []map[string]string
	ObjectInfo *api.ObjectInfo
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
	m.LastModifiedDate = ""
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
	case LastModifiedDateChangeMsg:
		// Update LastModifiedDate filter
		m.LastModifiedDate = msg.Value
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
	case FetchObjectInfoMsg:
		m.SelectedTab = "objectinfo"
		return m, m.fetchObjectInfoCmd()
	case ObjectInfoFetchedMsg:
		m.ObjectInfo = &msg.Info
		m.Loading = false
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
	children := []masc.MarkupOrChild{
		m.renderBreadcrumb(),
		m.renderPageLayout(send),
	}
	if m.ShowModal {
		children = append(children, m.renderModal(send))
	}
	if m.ShowToast {
		children = append(children, m.renderToast(send))
	}
	return components.Container(children...)
}

// renderActionsContent builds the Actions tab content.
func (m *AppModel) renderActionsContent(send func(masc.Msg)) masc.ComponentOrHTML {
	return components.Container(
		components.Button("Fetch Accounts", components.VariantBrand, func(e *masc.Event) {
			send(FetchAccountsMsg{})
		}),
		components.Button("Get Account Info", components.VariantBrand, func(e *masc.Event) {
			send(FetchObjectInfoMsg{})
		}),
		components.Button("Show Modal", components.VariantNeutral, func(e *masc.Event) {
			send(ToggleModalMsg{})
		}),
		components.Button("Show Toast", components.VariantNeutral, func(e *masc.Event) {
			send(ShowToastMsg{})
		}),
		components.Datepicker("Modified Since", m.LastModifiedDate, func(e *masc.Event) {
			send(LastModifiedDateChangeMsg{Value: e.Target.Get("value").String()})
		}),
		components.Badge("Demo Badge"),
		components.Pill("Tag1", func(e *masc.Event) {
			send(ShowToastMsg{})
		}),
		components.Pill("Tag2", nil),
		components.Icon(components.UtilityIcon, "settings", components.IconSmall),
	)
}

// renderDataContent builds the Data tab content.
func (m *AppModel) renderDataContent(send func(masc.Msg)) masc.ComponentOrHTML {
	var data []masc.MarkupOrChild
	if m.Loading {
		data = append(data,
			components.MarginTop(components.SpaceMedium,
				components.LoadingSpinner("medium"),
			),
		)
	} else if len(m.Rows) > 0 {
		data = append(data,
			components.MarginTop(components.SpaceMedium,
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
		data = append(data,
			components.MarginTop(components.SpaceMedium,
				components.Checkbox(
					"Only names containing 'A'",
					m.FilterAOnly,
					func(e *masc.Event) {
						send(CheckboxMsg{Checked: !m.FilterAOnly})
					},
				),
			),
		)
		data = append(data,
			components.MarginTop(components.SpaceMedium,
				components.RadioGroup(
					"filtermode",
					"Filter Mode",
					[]components.RadioOption{
						{Label: "Contains", Value: "contains"},
						{Label: "Starts With", Value: "startswith"},
					},
					m.FilterMode,
					func(mode string) { send(FilterModeMsg{Mode: mode}) },
				),
			),
		)
		data = append(data,
			components.MarginTop(components.SpaceMedium,
				components.TextInput("Filter by Name", m.InputValue, "Enter substring", func(e *masc.Event) {
					send(InputMsg{Value: e.Target.Get("value").String()})
				}),
			),
		)
		var suggestions []components.LookupOption
		for _, r := range m.Rows {
			name := r["Name"]
			suggestions = append(suggestions, components.LookupOption{Label: name, Value: name})
		}
		data = append(data,
			components.MarginTop(components.SpaceMedium,
				components.Lookup(
					"Select by Name",
					suggestions,
					m.InputValue,
					func(val string) { send(InputMsg{Value: val}) },
					func(val string) { send(InputMsg{Value: val}) },
				),
			),
		)
		var filtered []map[string]string
		query := strings.ToLower(m.InputValue)
		for _, r := range m.Rows {
			name := r["Name"]
			lower := strings.ToLower(name)
			var match bool
			if query == "" ||
				(m.FilterMode == "contains" && strings.Contains(lower, query)) ||
				(m.FilterMode == "startswith" && strings.HasPrefix(lower, query)) {
				match = true
			}
			if match && (!m.FilterAOnly || strings.Contains(lower, "a")) {
				filtered = append(filtered, r)
			}
		}
		if len(m.Rows) > 0 {
			percent := len(filtered) * 100 / len(m.Rows)
			data = append(data,
				components.MarginTop(components.SpaceMedium,
					components.ProgressBar(percent),
				),
			)
		}
		data = append(data,
			components.MarginTop(components.SpaceMedium,
				components.DataTable([]string{"Name", "First Contact"}, filtered),
			),
		)
	} else {
		// Show default message when no data has been fetched yet
		data = append(data,
			components.Spacer(components.SpaceOptions{
				PaddingHorizontal: components.SpaceMedium,
				MarginTop:         components.SpaceMedium,
			},
				components.Text("Click 'Fetch Accounts' to load account data and explore the data table features."),
			),
		)
	}
	return components.Container(data...)
}

// renderPageLayout builds the main page header, tabs, and card.
func (m *AppModel) renderPageLayout(send func(masc.Msg)) masc.ComponentOrHTML {
	actions := m.renderActionsContent(send)
	data := m.renderDataContent(send)
	tabs := components.Tabs(
		"accounts-tabs",
		[]components.TabOption{
			{Label: "Actions", Value: "actions", Content: actions},
			{Label: "Data", Value: "data", Content: data},
			{Label: "Object Info", Value: "objectinfo", Content: m.renderObjectInfoContent(send)},
			{Label: "Layout", Value: "layout", Content: m.renderLayoutContent(send)},
		},
		m.SelectedTab,
		func(tab string) { send(TabChangeMsg{Tab: tab}) },
	)
	header := components.PageHeader(
		"Thunder Demo",
		fmt.Sprintf("Mode: %s; Only A: %t", m.FilterMode, m.FilterAOnly),
	)
	card := components.Card("Accounts", tabs)
	return components.Page(header, card)
}

// renderBreadcrumb builds the page breadcrumb.
func (m *AppModel) renderBreadcrumb() masc.ComponentOrHTML {
	raw := components.Breadcrumb([]components.BreadcrumbOption{
		{Label: "Home", Href: "#"},
		{Label: "Thunder Demo", Href: "#"},
	})
	return components.Spacer(components.SpaceOptions{
		PaddingHorizontal: components.SpaceMedium,
		MarginBottom:      components.SpaceSmall,
	}, raw)
}

// renderObjectInfoContent builds the Object Info tab content.
func (m *AppModel) renderObjectInfoContent(send func(masc.Msg)) masc.ComponentOrHTML {
	if m.Loading && m.ObjectInfo == nil {
		return components.Spacer(components.SpaceOptions{
			PaddingHorizontal: components.SpaceMedium,
			MarginTop:         components.SpaceMedium,
		},
			components.LoadingSpinner("medium"),
		)
	}

	if m.ObjectInfo == nil {
		return components.Spacer(components.SpaceOptions{
			PaddingHorizontal: components.SpaceMedium,
			MarginTop:         components.SpaceMedium,
		},
			components.Text("Click 'Get Account Info' to fetch Account object metadata."),
		)
	}

	info := m.ObjectInfo
	return components.Spacer(components.SpaceOptions{
		PaddingHorizontal: components.SpaceMedium,
		MarginTop:         components.SpaceMedium,
	},
		components.MarginBottom(components.SpaceMedium,
			components.Heading("Account Object Information", components.HeadingMedium),
		),
		components.Grid(
			components.GridColumn("1-of-2", components.Card("Basic Info", components.Container(
				components.Paragraph(fmt.Sprintf("API Name: %s", info.APIName)),
				components.Paragraph(fmt.Sprintf("Label: %s", info.Label)),
				components.Paragraph(fmt.Sprintf("Label Plural: %s", info.LabelPlural)),
				components.Paragraph(fmt.Sprintf("Key Prefix: %s", info.KeyPrefix)),
				components.Paragraph(fmt.Sprintf("Custom: %t", info.Custom)),
			))),
			components.GridColumn("1-of-2", components.Card("Capabilities", components.Container(
				components.Paragraph(fmt.Sprintf("Createable: %t", info.Createable)),
				components.Paragraph(fmt.Sprintf("Updateable: %t", info.Updateable)),
				components.Paragraph(fmt.Sprintf("Deletable: %t", info.Deletable)),
				components.Paragraph(fmt.Sprintf("Queryable: %t", info.Queryable)),
				components.Paragraph(fmt.Sprintf("Searchable: %t", info.Searchable)),
			))),
		),
		components.MarginTop(components.SpaceMedium,
			components.Card("Additional Info", components.Container(
				components.Paragraph(fmt.Sprintf("Feed Enabled: %t", info.FeedEnabled)),
				components.Paragraph(fmt.Sprintf("MRU Enabled: %t", info.MRUEnabled)),
				components.Paragraph(fmt.Sprintf("Layoutable: %t", info.Layoutable)),
				components.Paragraph(fmt.Sprintf("Theme Color: %s", info.ThemeInfo.Color)),
				components.Paragraph(fmt.Sprintf("Number of Fields: %d", len(info.Fields))),
				components.Paragraph(fmt.Sprintf("Number of Child Relationships: %d", len(info.ChildRelationships))),
			)),
		),
	)
}

// renderLayoutContent builds the Layout tab content with a grid demonstration.
func (m *AppModel) renderLayoutContent(send func(masc.Msg)) masc.ComponentOrHTML {
	return components.Spacer(components.SpaceOptions{
		PaddingHorizontal: components.SpaceMedium,
		MarginTop:         components.SpaceLarge,
	},
		components.Text("Grid Demonstration:"),
		components.Grid(
			components.GridColumn("1-of-3", components.Card("Column 1", components.Text("This is column 1"))),
			components.GridColumn("1-of-3", components.Card("Column 2", components.Text("This is column 2"))),
			components.GridColumn("1-of-3", components.Card("Column 3", components.Text("This is column 3"))),
		),
	)
}

// renderModal builds the demo modal overlay.
func (m *AppModel) renderModal(send func(masc.Msg)) masc.ComponentOrHTML {
	return components.Modal(
		"Demo Modal",
		masc.Text("This is a demo modal"),
		components.Button("Close", components.VariantNeutral, func(e *masc.Event) {
			send(ToggleModalMsg{})
		}),
	)
}

// renderToast builds the toast notification.
func (m *AppModel) renderToast(send func(masc.Msg)) masc.ComponentOrHTML {
	return components.Toast(
		m.ToastVariant,
		m.ToastHeader,
		m.ToastMessage,
		func(e *masc.Event) { send(HideToastMsg{}) },
	)
}

// fetchAccountsCmd creates a Cmd that fetches accounts via JS and returns a Msg.
// It uses the provided limit value for the SOQL query.
func (m *AppModel) fetchAccountsCmd(limit string) masc.Cmd {
	m.Loading = true
	return func() masc.Msg {
		// Build SOQL query with optional LastModifiedDate filter
		var soql string
		if m.LastModifiedDate != "" {
			t, err := time.Parse("2006-01-02", m.LastModifiedDate)
			if err != nil {
				return QueryErrorMsg{Err: err.Error()}
			}
			dt := t.UTC().Format("2006-01-02T15:04:05Z")
			soql = fmt.Sprintf("SELECT Name, (SELECT Name FROM Contacts ORDER BY CreatedDate Desc LIMIT 1) FROM Account WHERE LastModifiedDate >= %s LIMIT %s", dt, limit)
		} else {
			soql = fmt.Sprintf("SELECT Name, (SELECT Name FROM Contacts ORDER BY CreatedDate Desc LIMIT 1) FROM Account LIMIT %s", limit)
		}
		data, err := api.Query(soql)
		if err != nil {
			return QueryErrorMsg{Err: err.Error()}
		}
		rows := make([]map[string]string, len(data))
		for i, r := range data {
			name, err := r.StringValue("Name")
			if err != nil {
				return QueryErrorMsg{Err: err.Error()}
			}
			contactName, _ := r.StringValue("let c = Contacts | first(); c?.Name")
			rows[i] = map[string]string{"Name": name, "First Contact": contactName}
		}
		return AccountsFetchedMsg{Rows: rows}
	}
}

// fetchObjectInfoCmd creates a Cmd that fetches Account object info via the UI API.
func (m *AppModel) fetchObjectInfoCmd() masc.Cmd {
	m.Loading = true
	return func() masc.Msg {
		info, err := api.GetObjectInfo("Account")
		if err != nil {
			return QueryErrorMsg{Err: err.Error()}
		}
		return ObjectInfoFetchedMsg{Info: info}
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
