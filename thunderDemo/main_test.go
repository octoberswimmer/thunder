package main

import (
   "testing"

   "github.com/octoberswimmer/masc"
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

// TestRenderBasic verifies that Render returns a non-nil component.
func TestRenderBasic(t *testing.T) {
	m := &AppModel{}
	m.Init()
	comp := m.Render(func(masc.Msg) {})
	if comp == nil {
		t.Error("expected Render to return a non-nil component")
	}
}
