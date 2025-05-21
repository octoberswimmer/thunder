package api

import (
	"testing"

	forcequery "github.com/ForceCLI/force/lib/query"
)

func TestRecord_StringValue_Simple(t *testing.T) {
	raw := forcequery.Record{Fields: map[string]interface{}{"Name": "Test"}}
	rec := Record{raw}
	got, err := rec.StringValue("Name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Test" {
		t.Errorf("expected \"Test\", got %q", got)
	}
}

func TestRecord_StringValue_Nested(t *testing.T) {
	inner := forcequery.Record{Fields: map[string]interface{}{"Name": "Inner"}}
	raw := forcequery.Record{Fields: map[string]interface{}{"Inner__r": inner}}
	rec := Record{raw}
	got, err := rec.StringValue("Inner__r.Name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Inner" {
		t.Errorf("expected \"Inner\", got %q", got)
	}
}

func TestRecord_StringValue_Error(t *testing.T) {
	raw := forcequery.Record{Fields: map[string]interface{}{"Name": 123}}
	rec := Record{raw}
	_, err := rec.StringValue("Name")
	if err == nil {
		t.Fatalf("expected error for non-string field")
	}
}

func TestRecord_Value_InvalidPath(t *testing.T) {
	raw := forcequery.Record{Fields: map[string]interface{}{"Name": "Test"}}
	rec := Record{raw}
	_, err := rec.Value("Foo.Bar")
	if err == nil {
		t.Fatalf("expected error for invalid path")
	}
}
