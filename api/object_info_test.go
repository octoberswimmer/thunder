package api

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUnmarshalObjectInfo(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "Account.json"))
	if err != nil {
		t.Fatalf("failed to read Account.json: %v", err)
	}
	info, err := UnmarshalObjectInfo(data)
	if err != nil {
		t.Fatalf("UnmarshalObjectInfo returned error: %v", err)
	}
	if info.APIName != "Account" {
		t.Errorf("APIName: got %q, want %q", info.APIName, "Account")
	}
	if info.ETag == "" {
		t.Error("expected non-empty ETag")
	}
	if info.DefaultRecordTypeID == "" {
		t.Error("expected non-empty DefaultRecordTypeID")
	}
	if _, ok := info.RecordTypeInfos[info.DefaultRecordTypeID]; !ok {
		t.Errorf("RecordTypeInfos missing defaultRecordTypeId %q", info.DefaultRecordTypeID)
	}
	if len(info.Fields) == 0 {
		t.Error("expected non-empty Fields map")
	}
	if _, ok := info.Fields["Name"]; !ok {
		t.Error("expected field Name to be present")
	}
}
