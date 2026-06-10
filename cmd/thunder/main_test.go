package main

import (
	"archive/zip"
	"bytes"
	"io"
	"math/rand"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

// unzipChunkFile extracts a single named file from a static resource zip chunk.
func unzipChunkFile(t *testing.T, zipData []byte, name string) ([]byte, bool) {
	t.Helper()
	r, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		t.Fatalf("reading chunk zip: %v", err)
	}
	for _, f := range r.File {
		if f.Name != name {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("opening %s: %v", name, err)
		}
		defer rc.Close()
		b, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("reading %s: %v", name, err)
		}
		return b, true
	}
	return nil, false
}

// Test_splitAndZip_single_chunk_includes_manifest verifies that a bundle small
// enough for one resource is returned as a single chunk carrying a parts.json
// manifest recording a count of 1.
func Test_splitAndZip_single_chunk_includes_manifest(t *testing.T) {
	data := bytes.Repeat([]byte("thunder"), 1000)
	chunks, err := splitAndZipWithLimit(data, staticResourceLimit)
	if err != nil {
		t.Fatalf("splitAndZipWithLimit returned error: %v", err)
	}
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(chunks))
	}
	wasm, ok := unzipChunkFile(t, chunks[0], "bundle.wasm")
	if !ok {
		t.Fatal("chunk missing bundle.wasm")
	}
	if !bytes.Equal(wasm, data) {
		t.Error("single chunk bundle.wasm does not match original")
	}
	manifest, ok := unzipChunkFile(t, chunks[0], "parts.json")
	if !ok {
		t.Fatal("single chunk missing parts.json manifest")
	}
	if got := string(manifest); got != `{"parts":1}` {
		t.Errorf("manifest = %q; want %q", got, `{"parts":1}`)
	}
}

// Test_splitAndZip_splits_large_bundle verifies that an oversized bundle is split
// into multiple chunks, each under the limit, that reassemble to the original,
// with the first chunk's manifest recording the chunk count.
func Test_splitAndZip_splits_large_bundle(t *testing.T) {
	// Incompressible (pseudo-random) data so the compressed whole exceeds the
	// small limit and forces a split. Deterministic seed keeps the test stable.
	data := make([]byte, 256*1024)
	rng := rand.New(rand.NewSource(1))
	rng.Read(data)
	const limit = 64 * 1024

	chunks, err := splitAndZipWithLimit(data, limit)
	if err != nil {
		t.Fatalf("splitAndZipWithLimit returned error: %v", err)
	}
	if len(chunks) < 2 {
		t.Fatalf("expected multiple chunks, got %d", len(chunks))
	}

	// First chunk records the total count; trailing chunks carry no manifest.
	manifest, ok := unzipChunkFile(t, chunks[0], "parts.json")
	if !ok {
		t.Fatal("first chunk missing parts.json manifest")
	}
	wantManifest := `{"parts":` + strconv.Itoa(len(chunks)) + `}`
	if got := string(manifest); got != wantManifest {
		t.Errorf("manifest = %q; want %q", got, wantManifest)
	}
	if _, ok := unzipChunkFile(t, chunks[len(chunks)-1], "parts.json"); ok {
		t.Error("trailing chunk should not carry a parts.json manifest")
	}

	var reassembled []byte
	for i, c := range chunks {
		if len(c) > limit {
			t.Errorf("chunk %d size %d exceeds limit %d", i, len(c), limit)
		}
		wasm, ok := unzipChunkFile(t, c, "bundle.wasm")
		if !ok {
			t.Fatalf("chunk %d missing bundle.wasm", i)
		}
		reassembled = append(reassembled, wasm...)
	}
	if !bytes.Equal(reassembled, data) {
		t.Error("reassembled bundle does not match original")
	}
}

// Test_staticResourceNames verifies chunk 0 keeps the base name and the rest are
// suffixed Part1, Part2, ... in load order.
func Test_staticResourceNames(t *testing.T) {
	if got := staticResourceNames("MyApp", 1); len(got) != 1 || got[0] != "MyApp" {
		t.Errorf("staticResourceNames single = %v; want [MyApp]", got)
	}
	got := staticResourceNames("MyApp", 3)
	want := []string{"MyApp", "MyAppPart1", "MyAppPart2"}
	if len(got) != len(want) {
		t.Fatalf("staticResourceNames(3) = %v; want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("staticResourceNames(3)[%d] = %q; want %q", i, got[i], want[i])
		}
	}
}

// Test_generateAppJS_imports_base_resource verifies the wrapper always imports the
// base static resource and sets this.app, independent of chunk count.
func Test_generateAppJS_imports_base_resource(t *testing.T) {
	js := generateAppJS("osgo/thunder", "MyApp", "My App", "MyApp")
	for _, want := range []string{
		"import Thunder from 'osgo/thunder';",
		"import APP_URL from '@salesforce/resourceUrl/MyApp';",
		"this.app = APP_URL + '/bundle.wasm';",
		"this.appName = 'My App';",
	} {
		if !strings.Contains(js, want) {
			t.Errorf("generated wrapper missing %q\n%s", want, js)
		}
	}
}

func Test_generateVisualforcePage_wires_runtime_and_remoting(t *testing.T) {
	page := generateVisualforcePage("Clinic Scheduler", "ClinicSchedulerWasmExec", []string{"ClinicScheduler"})
	for _, want := range []string{
		`controller="GoBridge"`,
		`<apex:includeScript value="{!$Resource.ClinicSchedulerWasmExec}"/>`,
		`"{!$RemoteAction.GoBridge.remoteCallRest}"`,
		`globalThis.get = function`,
		`"{!URLFOR($Resource.ClinicScheduler, 'bundle.wasm')}"`,
		`startWithDiv(document.getElementById("thunder-app"))`,
		`new Go()`,
		`<title>Clinic Scheduler</title>`,
		// Absent ?id= must surface as undefined so api.RecordId() reports "no record".
		`var recordId = recordIdParam || undefined;`,
	} {
		if !strings.Contains(page, want) {
			t.Errorf("generated page missing %q\n%s", want, page)
		}
	}
}

func Test_generateVisualforcePage_lists_all_split_parts(t *testing.T) {
	page := generateVisualforcePage("App", "AppWasmExec", staticResourceNames("App", 3))
	for _, want := range []string{
		`"{!URLFOR($Resource.App, 'bundle.wasm')}"`,
		`"{!URLFOR($Resource.AppPart1, 'bundle.wasm')}"`,
		`"{!URLFOR($Resource.AppPart2, 'bundle.wasm')}"`,
	} {
		if !strings.Contains(page, want) {
			t.Errorf("split page missing %q", want)
		}
	}
}

func Test_buildVisualforcePackageXML_includes_all_types(t *testing.T) {
	xml := buildVisualforcePackageXML(staticResourceNames("App", 2), "AppWasmExec", "App", "APP", true)
	for _, want := range []string{
		"<members>App</members>",
		"<members>AppPart1</members>",
		"<members>AppWasmExec</members>",
		"<name>StaticResource</name>",
		"<members>GoBridge</members>",
		"<name>ApexClass</name>",
		"<members>Thunder_Settings__c</members>",
		"<name>ApexPage</name>",
		"<members>APP</members>",
		"<name>CustomTab</name>",
	} {
		if !strings.Contains(xml, want) {
			t.Errorf("package.xml missing %q\n%s", want, xml)
		}
	}
}

func Test_buildVisualforcePackageXML_omits_tab_without_flag(t *testing.T) {
	xml := buildVisualforcePackageXML(staticResourceNames("App", 1), "AppWasmExec", "App", "APP", false)
	if strings.Contains(xml, "<name>CustomTab</name>") {
		t.Errorf("package.xml should not contain a CustomTab when withTab is false\n%s", xml)
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
