package main

import (
	"os"
	"path/filepath"
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
