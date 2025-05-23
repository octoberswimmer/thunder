package main

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func Test_indexHandler_serves_index_html(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	indexHandler(w, req)
	res := w.Result()
	defer res.Body.Close()
	if got := res.Header.Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q; want %q", got, "text/html; charset=utf-8")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}
	if got := string(body); got != indexHTML {
		t.Errorf("Body = %q; want %q", got, indexHTML)
	}
}

func Test_wasmHandler_serves_wasm_file(t *testing.T) {
	dir, err := os.MkdirTemp("", "test-build-*")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}
	defer os.RemoveAll(dir)
	data := []byte("wasm-content")
	path := filepath.Join(dir, "bundle.wasm")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("writing wasm file: %v", err)
	}
	buildMutex.Lock()
	currentBuildDir = dir
	buildMutex.Unlock()
	req := httptest.NewRequest("GET", "/bundle.wasm", nil)
	w := httptest.NewRecorder()
	wasmHandler(w, req)
	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}
	if got := string(body); got != string(data) {
		t.Errorf("Body = %q; want %q", got, string(data))
	}
}

func Test_wasmExecHandler_serves_wasm_exec_js(t *testing.T) {
	dir, err := os.MkdirTemp("", "test-build-*")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}
	defer os.RemoveAll(dir)
	data := []byte("exec-js-content")
	path := filepath.Join(dir, "wasm_exec.js")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("writing wasm_exec.js file: %v", err)
	}
	buildMutex.Lock()
	currentBuildDir = dir
	buildMutex.Unlock()
	req := httptest.NewRequest("GET", "/wasm_exec.js", nil)
	w := httptest.NewRecorder()
	wasmExecHandler(w, req)
	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}
	if got := string(body); got != string(data) {
		t.Errorf("Body = %q; want %q", got, string(data))
	}
}

// Test_zipBundle_creates_zip_with_bundleWasm verifies that zipBundle correctly creates
// a zip archive containing a single bundle.wasm file with the original data.
func Test_zipBundle_creates_zip_with_bundleWasm(t *testing.T) {
	data := []byte("wasm-content")
	zipData, err := zipBundle(data)
	if err != nil {
		t.Fatalf("zipBundle returned error: %v", err)
	}
	r, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		t.Fatalf("failed to read zip: %v", err)
	}
	if len(r.File) != 1 {
		t.Fatalf("expected 1 file in zip, got %d", len(r.File))
	}
	f := r.File[0]
	if f.Name != "bundle.wasm" {
		t.Errorf("expected file name bundle.wasm, got %s", f.Name)
	}
	rc, err := f.Open()
	if err != nil {
		t.Fatalf("failed to open zipped file: %v", err)
	}
	defer rc.Close()
	content, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("failed to read zipped file content: %v", err)
	}
	if !bytes.Equal(content, data) {
		t.Errorf("zipped content mismatch: expected %q, got %q", data, content)
	}
}
