package api

import (
	"fmt"

	forcequery "github.com/ForceCLI/force/lib/query"
)

// Query performs a SOQL query via the Force CLI query library,
// unifies pagination, and returns JSON with a 'records' array.
// It returns an error if the query fails.
func Query(soql string) ([]forcequery.Record, error) {
	records, err := forcequery.Eager(
		forcequery.InstanceUrl(""),
		forcequery.ApiVersion("v63.0"),
		forcequery.QS(soql),
		forcequery.HttpGet(Get),
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return records, nil
}
