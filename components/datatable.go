// Package components provides SLDS-styled Masc components under the thunder namespace.
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
