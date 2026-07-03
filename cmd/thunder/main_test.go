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

	"golang.org/x/tools/go/packages"
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

// Test_splitAndZip_first_chunk_carries_extras verifies that extra files (e.g.
// wasm_exec.js for Visualforce apps) are packed into the first chunk only.
func Test_splitAndZip_first_chunk_carries_extras(t *testing.T) {
	// Incompressible data large enough to force a split so we can assert the
	// extra rides only in the first chunk.
	data := make([]byte, 256*1024)
	rng := rand.New(rand.NewSource(1))
	rng.Read(data)
	const limit = 64 * 1024
	exec := []byte("// wasm_exec.js runtime shim")

	chunks, err := splitAndZipWithLimit(data, limit, zipEntry{name: "wasm_exec.js", data: exec})
	if err != nil {
		t.Fatalf("splitAndZipWithLimit returned error: %v", err)
	}
	if len(chunks) < 2 {
		t.Fatalf("expected multiple chunks, got %d", len(chunks))
	}
	got, ok := unzipChunkFile(t, chunks[0], "wasm_exec.js")
	if !ok {
		t.Fatal("first chunk missing wasm_exec.js")
	}
	if !bytes.Equal(got, exec) {
		t.Errorf("wasm_exec.js content mismatch: got %q want %q", got, exec)
	}
	if _, ok := unzipChunkFile(t, chunks[1], "wasm_exec.js"); ok {
		t.Error("trailing chunk should not carry wasm_exec.js")
	}
}

// Test_runtimeExtras_ships_wasm_exec_without_visualforce_flag verifies the Go
// runtime is packed into the static resource even when --visualforce is not
// set, so --app-only and --watch redeploys of Visualforce apps keep working.
func Test_runtimeExtras_ships_wasm_exec_without_visualforce_flag(t *testing.T) {
	orig := deployVisualforce
	deployVisualforce = false
	defer func() { deployVisualforce = orig }()

	extras, err := runtimeExtras()
	if err != nil {
		t.Fatalf("runtimeExtras returned error: %v", err)
	}
	if len(extras) != 1 || extras[0].name != "wasm_exec.js" {
		t.Fatalf("expected a single wasm_exec.js entry, got %v", extras)
	}
	if len(extras[0].data) == 0 {
		t.Error("wasm_exec.js entry is empty")
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
	// Default (managed package) deployment references the namespaced osgo
	// GoBridge controller. The Go runtime loads from the app's own static
	// resource, where wasm_exec.js was packed next to bundle.wasm.
	page := generateVisualforcePage("Clinic Scheduler", "osgo.GoBridge", []string{"ClinicScheduler"})
	for _, want := range []string{
		`controller="osgo.GoBridge"`,
		`<apex:includeScript value="{!URLFOR($Resource.ClinicScheduler, 'wasm_exec.js')}"/>`,
		// $RemoteAction resolves the namespace automatically, so the merge field
		// uses the unqualified class name; a namespaced reference is rejected by
		// the Visualforce compiler.
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
	// The namespace must not leak into the RemoteAction merge field.
	if strings.Contains(page, `$RemoteAction.osgo.GoBridge`) {
		t.Errorf("RemoteAction reference must not include the namespace prefix\n%s", page)
	}
}

func Test_generateVisualforcePage_uses_unmanaged_gobridge_under_thunder_dev(t *testing.T) {
	// --thunder-dev deploys GoBridge unmanaged, so the page references it
	// without the osgo namespace.
	page := generateVisualforcePage("App", "GoBridge", []string{"App"})
	for _, want := range []string{
		`controller="GoBridge"`,
		`<apex:includeScript value="{!URLFOR($Resource.App, 'wasm_exec.js')}"/>`,
		`"{!$RemoteAction.GoBridge.remoteCallRest}"`,
	} {
		if !strings.Contains(page, want) {
			t.Errorf("thunder-dev page missing %q\n%s", want, page)
		}
	}
}

func Test_generateVisualforcePage_lists_all_split_parts(t *testing.T) {
	page := generateVisualforcePage("App", "osgo.GoBridge", staticResourceNames("App", 3))
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

func Test_buildVisualforcePackageXML_includes_all_types_under_thunder_dev(t *testing.T) {
	xml := buildVisualforcePackageXML(staticResourceNames("App", 2), "App", "APP", true, true)
	for _, want := range []string{
		"<members>App</members>",
		"<members>AppPart1</members>",
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

// Without --thunder-dev the GoBridge classes and Thunder Settings object come
// from the osgo managed package, so the manifest lists only the app's own
// metadata (the WASM static resources carry wasm_exec.js).
func Test_buildVisualforcePackageXML_omits_managed_package_members(t *testing.T) {
	xml := buildVisualforcePackageXML(staticResourceNames("App", 2), "App", "APP", true, false)
	for _, want := range []string{
		"<members>App</members>",
		"<members>AppPart1</members>",
		"<name>StaticResource</name>",
		"<name>ApexPage</name>",
		"<name>CustomTab</name>",
	} {
		if !strings.Contains(xml, want) {
			t.Errorf("package.xml missing %q\n%s", want, xml)
		}
	}
	for _, unwanted := range []string{
		"<members>GoBridge</members>",
		"<name>ApexClass</name>",
		"<members>Thunder_Settings__c</members>",
	} {
		if strings.Contains(xml, unwanted) {
			t.Errorf("package.xml should not contain %q without --thunder-dev\n%s", unwanted, xml)
		}
	}
}

func Test_buildVisualforcePackageXML_omits_tab_without_flag(t *testing.T) {
	xml := buildVisualforcePackageXML(staticResourceNames("App", 1), "App", "APP", false, false)
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

func Test_packageLoadError_returns_nil_for_clean_package(t *testing.T) {
	pkgs := []*packages.Package{{Name: "main"}}
	if err := packageLoadError(pkgs); err != nil {
		t.Errorf("expected nil for clean package, got: %v", err)
	}
}

func Test_packageLoadError_reports_no_package_found(t *testing.T) {
	err := packageLoadError(nil)
	if err == nil || !strings.Contains(err.Error(), "no Go package found") {
		t.Errorf("expected no-package error, got: %v", err)
	}
}

func Test_packageLoadError_surfaces_underlying_package_errors(t *testing.T) {
	pkgs := []*packages.Package{{
		Name: "",
		Errors: []packages.Error{
			{Msg: "go.work lists go 1.24.4 but module requires go >= 1.25.0"},
		},
	}}
	err := packageLoadError(pkgs)
	if err == nil || !strings.Contains(err.Error(), "go >= 1.25.0") {
		t.Errorf("expected underlying load error to be surfaced, got: %v", err)
	}
}
