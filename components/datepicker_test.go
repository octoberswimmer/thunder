package components

import (
	"testing"
	"time"
)

// TestDatepickerBasic verifies that Datepicker returns non-nil for valid parameters.
func TestDatepickerBasic(t *testing.T) {
	comp := Datepicker("Date", time.Time{}, func(t time.Time) {})
	if comp == nil {
		t.Error("Datepicker returned nil for valid parameters")
	}
}

// TestDatepickerWithValue verifies that Datepicker works with a specific date value.
func TestDatepickerWithValue(t *testing.T) {
	date := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	comp := Datepicker("Birth Date", date, func(t time.Time) {})
	if comp == nil {
		t.Error("Datepicker returned nil with date value")
	}
}

// TestDatepickerWithZeroValue verifies that Datepicker works with zero time value.
func TestDatepickerWithZeroValue(t *testing.T) {
	comp := Datepicker("Optional Date", time.Time{}, func(t time.Time) {})
	if comp == nil {
		t.Error("Datepicker returned nil with zero time value")
	}
}

// TestDatepickerWithHandler verifies that Datepicker works with event handler.
func TestDatepickerWithHandler(t *testing.T) {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	handler := func(t time.Time) {}
	comp := Datepicker("Event Date", date, handler)
	if comp == nil {
		t.Error("Datepicker returned nil with event handler")
	}
}
