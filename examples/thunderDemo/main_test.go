package main

import (
	"testing"
	"time"

	"github.com/octoberswimmer/masc"
)

// TestLastModifiedDateChangeMsg verifies that the Update method sets LastModifiedDate.
func TestLastModifiedDateChangeMsg(t *testing.T) {
	m := &AppModel{}
	m.Init()
	date := time.Date(2023, 7, 20, 0, 0, 0, 0, time.UTC)
	_, _ = m.Update(LastModifiedDateChangeMsg{Value: date})
	if m.LastModifiedDate != date {
		t.Errorf("expected LastModifiedDate %v; got %v", date, m.LastModifiedDate)
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

// TestRenderLayoutContentBasic verifies that renderLayoutContent returns a non-nil component.
func TestRenderLayoutContentBasic(t *testing.T) {
	m := &AppModel{}
	m.Init()
	comp := m.renderLayoutContent(func(masc.Msg) {})
	if comp == nil {
		t.Error("expected renderLayoutContent to return a non-nil component")
	}
}
