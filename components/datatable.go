package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// DataTable renders an SLDS data table.
//
// headers is the ordered list of column labels.
// rows is a slice of maps from header to cell string.
func DataTable(headers []string, rows []map[string]string) masc.ComponentOrHTML {
	// If no headers, render nothing
	if len(headers) == 0 {
		return nil
	}
	// Table classes (striped for alternating row colors)
	tableClass := masc.Class("slds-table", "slds-table_cell-buffer", "slds-table_bordered", "slds-table_striped")
	// Build header row
	var headCells []masc.MarkupOrChild
	for _, h := range headers {
		headCells = append(headCells,
			elem.TableHeader(
				masc.Markup(
					masc.Property("scope", "col"),
				),
				elem.Div(
					masc.Markup(
						masc.Class("slds-truncate"),
						masc.Property("title", h),
					),
					masc.Text(h),
				),
			),
		)
	}
	// Combine header markup and header cells into arguments
	var headerRowArgs []masc.MarkupOrChild
	headerRowArgs = append(headerRowArgs, masc.Markup(masc.Class("slds-line-height_reset")))
	headerRowArgs = append(headerRowArgs, headCells...)
	headRow := elem.TableRow(headerRowArgs...)
	// Build body rows
	var bodyRows []masc.MarkupOrChild
	for _, row := range rows {
		var cells []masc.MarkupOrChild
		for i, h := range headers {
			cellText := row[h]
			content := elem.Div(
				masc.Markup(
					masc.Class("slds-truncate"),
					masc.Property("title", cellText),
				),
				masc.Text(cellText),
			)
			if i == 0 {
				cells = append(cells,
					elem.TableHeader(
						masc.Markup(
							masc.Property("scope", "row"),
							masc.Data("label", h),
						),
						content,
					),
				)
			} else {
				cells = append(cells,
					elem.TableData(
						masc.Markup(masc.Data("label", h)),
						content,
					),
				)
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

// TableWithActions renders an enhanced SLDS data table with action buttons per row.
// headers is the ordered list of column labels.
// rows is a slice of table rows containing cells and optional actions.
type TableRow struct {
	Cells   []TableCell
	Actions []masc.ComponentOrHTML
}

type TableCell struct {
	Content string
	Title   string // Optional tooltip text, defaults to Content if empty
}

// ComponentTableRow represents a table row with flexible content types
type ComponentTableRow struct {
	Cells   []ComponentTableCell
	Actions masc.ComponentOrHTML
}

// ComponentTableCell represents a table cell that can contain any component
type ComponentTableCell struct {
	Content masc.ComponentOrHTML
	Title   string // Optional tooltip text
}

func TableWithActions(headers []string, rows []TableRow) masc.ComponentOrHTML {
	if len(headers) == 0 {
		return nil
	}

	// Table classes (striped for alternating row colors)
	tableClass := masc.Class("slds-table", "slds-table_cell-buffer", "slds-table_bordered", "slds-table_striped")

	// Build header row
	var headCells []masc.MarkupOrChild
	for _, h := range headers {
		headCells = append(headCells,
			elem.TableHeader(
				masc.Markup(
					masc.Property("scope", "col"),
				),
				elem.Div(
					masc.Markup(
						masc.Class("slds-truncate"),
						masc.Property("title", h),
					),
					masc.Text(h),
				),
			),
		)
	}

	// Header row with styling
	var headerRowArgs []masc.MarkupOrChild
	headerRowArgs = append(headerRowArgs, masc.Markup(masc.Class("slds-line-height_reset")))
	headerRowArgs = append(headerRowArgs, headCells...)
	headRow := elem.TableRow(headerRowArgs...)

	// Build body rows
	var bodyRows []masc.MarkupOrChild
	for _, row := range rows {
		var cells []masc.MarkupOrChild

		// Add data cells
		for i, cell := range row.Cells {
			if i >= len(headers) {
				break // Don't exceed header count
			}

			title := cell.Title
			if title == "" {
				title = cell.Content
			}

			content := elem.Div(
				masc.Markup(
					masc.Class("slds-truncate"),
					masc.Property("title", title),
				),
				masc.Text(cell.Content),
			)

			if i == 0 {
				// First cell is a row header
				cells = append(cells,
					elem.TableHeader(
						masc.Markup(
							masc.Property("scope", "row"),
							masc.Data("label", headers[i]),
						),
						content,
					),
				)
			} else {
				// Regular data cell
				cells = append(cells,
					elem.TableData(
						masc.Markup(masc.Data("label", headers[i])),
						content,
					),
				)
			}
		}

		// Add action cell if actions exist
		if len(row.Actions) > 0 {
			actionArgs := []masc.MarkupOrChild{
				masc.Markup(masc.Class("slds-button-space")),
			}
			for _, action := range row.Actions {
				actionArgs = append(actionArgs, action)
			}
			actionCell := elem.TableData(
				masc.Markup(masc.Data("label", "Actions")),
				elem.Div(actionArgs...),
			)
			cells = append(cells, actionCell)
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

// TableWithComponents renders a data table with flexible component content
func TableWithComponents(headers []string, rows []ComponentTableRow) masc.ComponentOrHTML {
	if len(headers) == 0 {
		return nil
	}

	// Table classes (striped for alternating row colors)
	tableClass := masc.Class("slds-table", "slds-table_cell-buffer", "slds-table_bordered", "slds-table_striped")

	// Build header row
	var headCells []masc.MarkupOrChild
	for _, h := range headers {
		headCells = append(headCells,
			elem.TableHeader(
				masc.Markup(
					masc.Property("scope", "col"),
				),
				elem.Div(
					masc.Markup(
						masc.Class("slds-truncate"),
						masc.Property("title", h),
					),
					masc.Text(h),
				),
			),
		)
	}

	// Header row with styling
	var headerRowArgs []masc.MarkupOrChild
	headerRowArgs = append(headerRowArgs, masc.Markup(masc.Class("slds-line-height_reset")))
	headerRowArgs = append(headerRowArgs, headCells...)
	headRow := elem.TableRow(headerRowArgs...)

	// Build body rows
	var bodyRows []masc.MarkupOrChild
	for _, row := range rows {
		var cells []masc.MarkupOrChild

		// Add data cells
		for i, cell := range row.Cells {
			if i >= len(headers) {
				break // Don't exceed header count
			}

			content := elem.Div(
				masc.Markup(
					masc.Class("slds-truncate"),
					masc.Property("title", cell.Title),
				),
				cell.Content,
			)

			if i == 0 {
				// First cell is a row header
				cells = append(cells,
					elem.TableHeader(
						masc.Markup(
							masc.Property("scope", "row"),
							masc.Data("label", headers[i]),
						),
						content,
					),
				)
			} else {
				// Regular data cell
				cells = append(cells,
					elem.TableData(
						masc.Markup(masc.Data("label", headers[i])),
						content,
					),
				)
			}
		}

		// Add action cell if actions exist
		if row.Actions != nil {
			actionCell := elem.TableData(
				masc.Markup(masc.Data("label", "Actions")),
				row.Actions,
			)
			cells = append(cells, actionCell)
		}

		// Combine row markup and cells into arguments
		var rowArgs []masc.MarkupOrChild
		rowArgs = append(rowArgs, masc.Markup(masc.Class("slds-hint-parent")))
		rowArgs = append(rowArgs, cells...)
		bodyRows = append(bodyRows, elem.TableRow(rowArgs...))
	}

	// Assemble table in a section container
	return Section(SpaceLarge,
		elem.Table(
			masc.Markup(tableClass),
			elem.TableHead(headRow),
			elem.TableBody(bodyRows...),
		),
	)
}

// EmptyTable renders a centered message when no data is available.
func EmptyTable(message string) masc.ComponentOrHTML {
	return CenteredContainer(
		PaddingAround(SpaceMedium,
			masc.Text(message),
		),
	)
}

// LoadingTable renders a centered loading spinner for table loading states.
func LoadingTable() masc.ComponentOrHTML {
	return CenteredContainer(
		PaddingAround(SpaceMedium,
			Spinner("medium"),
		),
	)
}
