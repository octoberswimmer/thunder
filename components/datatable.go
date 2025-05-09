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
	// Build header cells
	var headCells []masc.MarkupOrChild
	for _, h := range headers {
		headCells = append(headCells, elem.TableHeader(masc.Text(h)))
	}
	// Build body rows
	var bodyRows []masc.MarkupOrChild
	for _, row := range rows {
		var cells []masc.MarkupOrChild
		for _, h := range headers {
			cells = append(cells, elem.TableData(masc.Text(row[h])))
		}
		bodyRows = append(bodyRows, elem.TableRow(cells...))
	}
	// Assemble table
	return elem.Table(
		masc.Markup(masc.Class("slds-table", "slds-table_cell-buffer", "slds-table_bordered")),
		elem.TableHead(elem.TableRow(headCells...)),
		elem.TableBody(bodyRows...),
	)
}
