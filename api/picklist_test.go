package api

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUnmarshalPicklistFieldValues(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "picklistvalues.json"))
	if err != nil {
		t.Fatalf("failed to read picklistvalues.json: %v", err)
	}
	m, err := UnmarshalPicklistFieldValues(data)
	if err != nil {
		t.Fatalf("UnmarshalPicklistFieldValues returned error: %v", err)
	}
	if len(m) == 0 {
		t.Fatal("expected at least one picklist field, got none")
	}
	pfv, ok := m["AccountSource"]
	if !ok {
		t.Fatal("expected AccountSource in picklist field map")
	}
	if pfv.ETag == "" {
		t.Error("expected non-empty ETag for AccountSource")
	}
	expURL := "/services/data/v63.0/ui-api/object-info/Account/picklist-values/012QP000000DsyiYAC/AccountSource"
	if pfv.URL != expURL {
		t.Errorf("unexpected URL: got %q, want %q", pfv.URL, expURL)
	}
	if len(pfv.Values) == 0 {
		t.Error("expected non-empty Values slice for AccountSource")
	}
	first := pfv.Values[0]
	if first.Value != "Advertisement" {
		t.Errorf("expected first picklist value Advertisement, got %q", first.Value)
	}
}
