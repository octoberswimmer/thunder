package main

import (
	"testing"
)

// TestLastModifiedDateChangeMsg verifies that the Update method sets LastModifiedDate.
func TestLastModifiedDateChangeMsg(t *testing.T) {
	m := &AppModel{}
	m.Init()
	date := "2023-07-20"
	_, _ = m.Update(LastModifiedDateChangeMsg{Value: date})
	if m.LastModifiedDate != date {
		t.Errorf("expected LastModifiedDate %q; got %q", date, m.LastModifiedDate)
	}
}
