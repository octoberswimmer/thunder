package api

import (
	"encoding/json"
	"fmt"

	forcequery "github.com/ForceCLI/force/lib/query"
)

// Query performs a SOQL query via the Force CLI query library,
// unifies pagination, and returns JSON with a 'records' array.
// It returns an error if the query fails.
func Query(soql string) ([]byte, error) {
	records, err := forcequery.Eager(
		forcequery.InstanceUrl(""),
		forcequery.ApiVersion("v63.0"),
		forcequery.QS(soql),
		forcequery.HttpGet(Get),
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	raw := make([]map[string]interface{}, len(records))
	for i, r := range records {
		raw[i] = r.Raw
	}
	payload := map[string]interface{}{"records": raw}
	out, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query results: %w", err)
	}
	return out, nil
}
