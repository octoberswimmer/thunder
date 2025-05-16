package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopyFile(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "subdir", "dst.txt")
	data := []byte("hello thunder")
	if err := os.WriteFile(src, data, 0644); err != nil {
		t.Fatalf("writing source file: %v", err)
	}
	if err := copyFile(src, dst); err != nil {
		t.Fatalf("copyFile returned error: %v", err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("reading destination file: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("copied data = %q; want %q", got, data)
	}
}

// Tests for serve command validations
func TestRunServe_InvalidDir(t *testing.T) {
	// Non-existent directory
	serveDir = "/does/not/exist"
	servePort = 0
	err := runServe(nil, nil)
	if err == nil || !strings.Contains(err.Error(), "Invalid app directory") {
		t.Fatalf("Expected invalid app directory error, got: %v", err)
	}
}

func TestRunServe_NotMainPackage(t *testing.T) {
	tmp := t.TempDir()
	// Create a Go file with non-main package
	mainFile := filepath.Join(tmp, "app.go")
	src := []byte("package foo\nfunc main() {}")
	if err := os.WriteFile(mainFile, src, 0644); err != nil {
		t.Fatalf("writing file: %v", err)
	}
	// Initialize as Go module so go list works
	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module test"), 0644); err != nil {
		t.Fatalf("writing go.mod: %v", err)
	}
	serveDir = tmp
	servePort = 0
	err := runServe(nil, nil)
	if err == nil || !strings.Contains(err.Error(), "is not package main") {
		t.Fatalf("Expected package main error, got: %v", err)
	}
}

// Tests for deploy command validations
func TestRunDeploy_InvalidDir(t *testing.T) {
	// Non-existent directory
	deployDir = "/absent"
	deployTab = false
	err := runDeploy(nil, nil)
	if err == nil || !strings.Contains(err.Error(), "Invalid app directory") {
		t.Fatalf("Expected invalid app directory error, got: %v", err)
	}
}

func TestRunDeploy_NotMainPackage(t *testing.T) {
	tmp := t.TempDir()
	// Create a Go file with non-main package
	mainFile := filepath.Join(tmp, "app.go")
	src := []byte("package bar\nfunc main() {}")
	if err := os.WriteFile(mainFile, src, 0644); err != nil {
		t.Fatalf("writing file: %v", err)
	}
	// Initialize as Go module so go list works
	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module test"), 0644); err != nil {
		t.Fatalf("writing go.mod: %v", err)
	}
	deployDir = tmp
	deployTab = false
	err := runDeploy(nil, nil)
	if err == nil || !strings.Contains(err.Error(), "is not package main") {
		t.Fatalf("Expected package main error for deploy, got: %v", err)
	}
}
