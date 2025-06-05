package api

import (
	"fmt"

	"strings"

	forcequery "github.com/ForceCLI/force/lib/query"
)

// Query performs a SOQL query via the Force CLI query library,
// unifies pagination, and returns JSON with a 'records' array.
// It returns an error if the query fails.
// Query executes the SOQL query and returns wrapped Records for field access.
func Query(soql string) ([]Record, error) {
	if strings.Contains(soql, "PurchaserPlan") {
		// panic("querying for Purchaser Plan")
	}
	raw, err := forcequery.Eager(
		forcequery.InstanceUrl(""),
		forcequery.ApiVersion("v63.0"),
		forcequery.QS(soql),
		forcequery.HttpGet(Get),
	)
	if strings.Contains(soql, "PurchaserPlan") {
		// panic("got result Purchaser Plan")
	}
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	records := make([]Record, len(raw))
	for i, r := range raw {
		records[i] = Record{r}
	}
	return records, nil
}
