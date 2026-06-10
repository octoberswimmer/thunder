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

func TestRecord_Children_Subquery(t *testing.T) {
	children := []forcequery.Record{
		{Fields: map[string]interface{}{"Name": "First"}},
		{Fields: map[string]interface{}{"Name": "Second"}},
	}
	raw := forcequery.Record{Fields: map[string]interface{}{"Clinic__r": children}}
	rec := Record{raw}

	got := rec.Children("Clinic__r")
	if len(got) != 2 {
		t.Fatalf("expected 2 children, got %d", len(got))
	}
	for i, want := range []string{"First", "Second"} {
		name, err := got[i].StringValue("Name")
		if err != nil {
			t.Fatalf("unexpected error reading child %d: %v", i, err)
		}
		if name != want {
			t.Errorf("child %d: expected %q, got %q", i, want, name)
		}
	}
}

func TestRecord_Children_Absent(t *testing.T) {
	raw := forcequery.Record{Fields: map[string]interface{}{"Name": "Test"}}
	rec := Record{raw}
	if got := rec.Children("Clinic__r"); got != nil {
		t.Errorf("expected nil for absent relationship, got %v", got)
	}
}

func TestRecord_Children_WrongType(t *testing.T) {
	raw := forcequery.Record{Fields: map[string]interface{}{"Clinic__r": "not-a-subquery"}}
	rec := Record{raw}
	if got := rec.Children("Clinic__r"); got != nil {
		t.Errorf("expected nil for non-subquery field, got %v", got)
	}
}

func TestRecord_Children_Empty(t *testing.T) {
	raw := forcequery.Record{Fields: map[string]interface{}{"Clinic__r": []forcequery.Record{}}}
	rec := Record{raw}
	if got := rec.Children("Clinic__r"); len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}
