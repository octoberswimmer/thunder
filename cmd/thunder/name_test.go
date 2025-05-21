package main

import "testing"

func TestSanitizeStaticResourceName(t *testing.T) {
	tests := []struct{ input, want string }{
		{"My App!", "MyApp"},
		{"123Name", "A123Name"},
		{"-hello-", "hello"},
		{"", "App"},
	}
	for _, tc := range tests {
		if got := sanitizeStaticResourceName(tc.input); got != tc.want {
			t.Errorf("sanitizeStaticResourceName(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

func TestSanitizeComponentName(t *testing.T) {
	tests := []struct{ input, want string }{
		{"myApp", "myApp"},
		{"MyApp", "MyApp"},
		{"My App!", "My_App"},
		{"--Test--Name--", "Test_Name"},
		{"123name", "a123name"},
		{"", "app"},
		{"Foo__Bar", "Foo_Bar"},
	}
	for _, tc := range tests {
		if got := sanitizeComponentName(tc.input); got != tc.want {
			t.Errorf("sanitizeComponentName(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct{ input, want string }{
		{"my_app", "MyApp"},
		{"test", "Test"},
		{"foo_bar_baz", "FooBarBaz"},
		{"", ""},
	}
	for _, tc := range tests {
		if got := toPascalCase(tc.input); got != tc.want {
			t.Errorf("toPascalCase(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}
