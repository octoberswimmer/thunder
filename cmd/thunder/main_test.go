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
	if f.Method != zip.Deflate {
		t.Errorf("expected compression method Deflate, got %d", f.Method)
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

// Test_zipBundle_compresses_highly_compressible_data verifies that zipBundle
// actually compresses content (regression test for using flate.BestCompression).
func Test_zipBundle_compresses_highly_compressible_data(t *testing.T) {
	// 64 KiB of zeros — should compress to a tiny fraction of original size.
	data := make([]byte, 64*1024)
	zipData, err := zipBundle(data)
	if err != nil {
		t.Fatalf("zipBundle returned error: %v", err)
	}
	// Even with zip overhead, the result should be far smaller than the input
	// if compression is actually engaged.
	if len(zipData) >= len(data)/4 {
		t.Errorf("zipBundle did not compress highly-compressible data: input %d bytes, output %d bytes", len(data), len(zipData))
	}
}

// Test_optimizeWASM_missing_binary_is_noop ensures optimizeWASM does not panic
// or modify the input file when wasm-opt is unavailable.
func Test_optimizeWASM_missing_binary_is_noop(t *testing.T) {
	// Point PATH at an empty directory so wasm-opt cannot be found.
	emptyDir := t.TempDir()
	t.Setenv("PATH", emptyDir)

	dir := t.TempDir()
	wasm := filepath.Join(dir, "bundle.wasm")
	data := []byte("\x00asm\x01\x00\x00\x00") // minimal wasm header bytes
	if err := os.WriteFile(wasm, data, 0644); err != nil {
		t.Fatalf("writing wasm: %v", err)
	}
	optimizeWASM(wasm) // must not panic
	got, err := os.ReadFile(wasm)
	if err != nil {
		t.Fatalf("reading wasm after optimizeWASM: %v", err)
	}
	if !bytes.Equal(got, data) {
		t.Errorf("optimizeWASM mutated bundle when wasm-opt was missing")
	}
}
