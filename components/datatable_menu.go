package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
	"github.com/octoberswimmer/masc/event"
)

// ActionColumn represents a column configuration for actions
type ActionColumn struct {
	Actions []RowAction
}

// RowAction represents an action that can be performed on a row
type RowAction struct {
	Label string
	Name  string
}

// DataTableColumn represents a column configuration
type DataTableColumn struct {
	Label     string
	FieldName string
	Type      string
	Actions   *ActionColumn // Only set for action columns
}

// DataTableWithMenu renders an SLDS data table with a dropdown menu for actions (more Lightning-like)
func DataTableWithMenu(columns []DataTableColumn, rows []map[string]interface{}, onRowAction func(string, map[string]interface{})) masc.ComponentOrHTML {
	if len(columns) == 0 {
		return nil
	}

	// Table classes (striped for alternating row colors)
	tableClass := masc.Class("slds-table", "slds-table_cell-buffer", "slds-table_bordered", "slds-table_striped")

	// Build header row
	var headCells []masc.MarkupOrChild
	for _, col := range columns {
		if col.Type == "action" {
			// Action column header is empty
			headCells = append(headCells,
				elem.TableHeader(
					masc.Markup(
						masc.Property("scope", "col"),
						masc.Class("slds-text-align_right"),
					),
					elem.Div(
						masc.Markup(
							masc.Class("slds-truncate"),
						),
						elem.Span(
							masc.Markup(masc.Class("slds-assistive-text")),
							masc.Text("Actions"),
						),
					),
				),
			)
		} else {
			headCells = append(headCells,
				elem.TableHeader(
					masc.Markup(
						masc.Property("scope", "col"),
					),
					elem.Div(
						masc.Markup(
							masc.Class("slds-truncate"),
							masc.Property("title", col.Label),
						),
						masc.Text(col.Label),
					),
				),
			)
		}
	}

	// Combine header markup and header cells into arguments
	var headerRowArgs []masc.MarkupOrChild
	headerRowArgs = append(headerRowArgs, masc.Markup(masc.Class("slds-line-height_reset")))
	headerRowArgs = append(headerRowArgs, headCells...)
	headRow := elem.TableRow(headerRowArgs...)

	// Build body rows
	var bodyRows []masc.MarkupOrChild
	for rowIndex, row := range rows {
		var cells []masc.MarkupOrChild

		for colIndex, col := range columns {
			if col.Type == "action" {
				// Action column with menu button
				actionCell := renderMenuCell(col, row, onRowAction, rowIndex)
				cells = append(cells, actionCell)
			} else {
				// Regular data cell
				cellValue := ""
				if val, ok := row[col.FieldName]; ok {
					if strVal, isString := val.(string); isString {
						cellValue = strVal
					} else if boolVal, isBool := val.(bool); isBool && col.Type == "boolean" {
						if boolVal {
							cellValue = "Yes"
						} else {
							cellValue = "No"
						}
					}
				}

				content := elem.Div(
					masc.Markup(
						masc.Class("slds-truncate"),
						masc.Property("title", cellValue),
					),
					masc.Text(cellValue),
				)

				if colIndex == 0 {
					cells = append(cells,
						elem.TableHeader(
							masc.Markup(
								masc.Property("scope", "row"),
								masc.Data("label", col.Label),
							),
							content,
						),
					)
				} else {
					cells = append(cells,
						elem.TableData(
							masc.Markup(masc.Data("label", col.Label)),
							content,
						),
					)
				}
			}
		}

		// Combine row markup and cells into arguments
		var rowArgs []masc.MarkupOrChild
		rowArgs = append(rowArgs, masc.Markup(masc.Class("slds-hint-parent")))
		rowArgs = append(rowArgs, cells...)
		bodyRows = append(bodyRows, elem.TableRow(rowArgs...))
	}

	// Assemble table
	return elem.Table(
		masc.Markup(tableClass),
		elem.TableHead(headRow),
		elem.TableBody(bodyRows...),
	)
}

// renderMenuCell renders a table cell with a proper dropdown menu for actions
func renderMenuCell(col DataTableColumn, row map[string]interface{}, onRowAction func(string, map[string]interface{}), rowIndex int) masc.ComponentOrHTML {
	// Check if this row is loading
	isLoading := false
	if loadingVal, ok := row["isLoading"].(bool); ok {
		isLoading = loadingVal
	}

	// If loading, show spinner instead of menu
	if isLoading {
		return elem.TableData(
			masc.Markup(masc.Class("slds-text-align_right")),
			CenteredSpinner("small"),
		)
	}

	if col.Actions == nil || len(col.Actions.Actions) == 0 {
		return elem.TableData(
			masc.Markup(masc.Class("slds-text-align_right")),
		)
	}

	// Filter actions based on row status
	status := ""
	if statusVal, ok := row["Status"].(string); ok {
		status = statusVal
	}

	var validActions []RowAction
	for _, action := range col.Actions.Actions {
		// Show Edit and Delete for all rows
		if action.Name == "Edit" || action.Name == "Delete" {
			validActions = append(validActions, action)
		}
		// Show Activate only for Inactive rows
		if action.Name == "Activate" && status == "Inactive" {
			validActions = append(validActions, action)
		}
		// Show Deactivate only for Active rows
		if action.Name == "Deactivate" && status == "Active" {
			validActions = append(validActions, action)
		}
	}

	if len(validActions) == 0 {
		return elem.TableData(
			masc.Markup(masc.Class("slds-text-align_right")),
		)
	}

	// Create dropdown menu items
	var menuItems []masc.ComponentOrHTML
	for _, action := range validActions {
		// Capture the action name in the closure
		actionName := action.Name

		menuItems = append(menuItems,
			elem.ListItem(
				masc.Markup(masc.Class("slds-dropdown__item")),
				elem.Anchor(
					masc.Markup(
						masc.Attribute("href", "javascript:void(0);"),
						masc.Attribute("role", "menuitem"),
						masc.Attribute("tabindex", "-1"),
						masc.Class("slds-dropdown__item-link"),
						event.Click(func(e *masc.Event) {
							if onRowAction != nil {
								onRowAction(actionName, row)
							}
							// Close the dropdown after selecting an action via JavaScript
							e.Target.Call("closest", ".slds-dropdown").Get("previousElementSibling").Call("setAttribute", "aria-expanded", "false")
							e.Target.Call("closest", ".slds-dropdown").Get("style").Set("display", "none")
						}),
					),
					elem.Span(
						masc.Markup(masc.Class("slds-truncate")),
						masc.Text(action.Label),
					),
				),
			),
		)
	}

	// Create the dropdown structure
	dropdownId := "dropdown-menu-" + string(rune(rowIndex))

	return elem.TableData(
		masc.Markup(masc.Class("slds-text-align_right")),
		elem.Div(
			masc.Markup(
				masc.Class("slds-dropdown-trigger", "slds-dropdown-trigger_click"),
				masc.Data("dropdownid", dropdownId),
				masc.Style("position", "relative"),
			),
			// Trigger button
			elem.Button(
				masc.Markup(
					masc.Class("slds-button", "slds-button_icon", "slds-button_icon-border-filled", "slds-button_icon-x-small"),
					masc.Attribute("type", "button"),
					masc.Attribute("aria-haspopup", "true"),
					masc.Attribute("aria-expanded", "false"),
					masc.Attribute("title", "Show Actions"),
					// Simple dropdown toggle
					event.Click(func(e *masc.Event) {
						button := e.Target
						dropdown := button.Get("nextElementSibling")
						if dropdown.IsNull() {
							return
						}
						isOpen := button.Call("getAttribute", "aria-expanded").String() == "true"

						// Toggle this dropdown
						if isOpen {
							button.Call("setAttribute", "aria-expanded", "false")
							dropdown.Get("style").Set("display", "none")
						} else {
							button.Call("setAttribute", "aria-expanded", "true")
							dropdown.Get("style").Set("display", "block")
						}
					}),
				),
				masc.Text("⋯"),
				elem.Span(
					masc.Markup(masc.Class("slds-assistive-text")),
					masc.Text("Show Actions"),
				),
			),
			// Dropdown menu
			elem.Div(
				masc.Markup(
					masc.Class("slds-dropdown", "slds-dropdown_right"),
					masc.Property("role", "menu"),
					masc.Attribute("aria-labelledby", dropdownId),
					masc.Style("display", "none"), // Initially hidden
					masc.Style("position", "absolute"),
					masc.Style("z-index", "1000"),
					masc.Style("min-width", "6rem"),
					masc.Style("top", "100%"),
					masc.Style("right", "0"),
				),
				elem.UnorderedList(
					masc.Markup(
						masc.Class("slds-dropdown__list"),
						masc.Property("role", "presentation"),
					),
					masc.List(menuItems),
				),
			),
		),
	)
}
