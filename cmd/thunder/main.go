package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
	"unicode"

	salesforce "github.com/octoberswimmer/thunder/salesforce"
	"golang.org/x/tools/go/packages"

	desktop "github.com/ForceCLI/force/desktop"
	forcecli "github.com/ForceCLI/force/lib"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

// global state for serve command
var (
	servePort       int
	serveDir        string
	currentBuildDir string
	buildMutex      sync.RWMutex
	session         *forcecli.Force
	// deploy command flags
	deployDir   string
	deployTab   bool
	deployWatch bool
	deployDebug bool
	// build command flags
	buildDev    bool
	buildOutput string
)

// indexHTML is the HTML template served for the Thunder app root.
const indexHTML = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Thunder App</title>
    <link rel="stylesheet" href="https://unpkg.com/@salesforce-ux/design-system@latest/assets/styles/salesforce-lightning-design-system.min.css">
    <script src="wasm_exec.js"></script>
    <script>
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("bundle.wasm"), go.importObject).then((result) => {
            go.run(result.instance);
        });
    </script>
</head>
<body>
    <div id="app"></div>
</body>
</html>`

// root command
var rootCmd = &cobra.Command{Use: "thunder"}

// serve command
var serveCmd = &cobra.Command{
	Use:   "serve [dir]",
	Short: "Build and serve the Thunder app locally",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runServe,
}

// deploy command stub
var deployCmd = &cobra.Command{
	Use:   "deploy [dir]",
	Short: "Deploy the Thunder app to a Salesforce org",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runDeploy,
}

// build command
var buildCmd = &cobra.Command{
	Use:   "build [dir]",
	Short: "Build the Thunder app to WebAssembly",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runBuild,
}

func init() {
	// serve flags (port only; app dir is optional positional arg)
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8000, "Port to serve on")
	// deploy flags (app dir is optional positional arg)
	deployCmd.Flags().BoolVarP(&deployTab, "tab", "t", false, "Deploy and open a CustomTab for the app")
	deployCmd.Flags().BoolVarP(&deployWatch, "watch", "w", false, "Watch for changes and automatically redeploy WASM bundle")
	deployCmd.Flags().BoolVar(&deployDebug, "debug", false, "Enable debug output")
	// build flags
	buildCmd.Flags().BoolVarP(&buildDev, "dev", "d", false, "Build with development tags")
	buildCmd.Flags().StringVarP(&buildOutput, "output", "o", "./build", "Output directory for build artifacts")
	// add subcommands
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(buildCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// buildWASM compiles the Go app in appDir to WebAssembly and prepares assets.
func buildWASM(appDir string) (string, error) {
	// create temporary build directory
	buildDir, err := os.MkdirTemp("", "thunder-build-*")
	if err != nil {
		return "", err
	}
	// build WASM binary
	outWasm := filepath.Join(buildDir, "bundle.wasm")
	cmd := exec.Command("go", "build", "-o", outWasm, "-tags", "dev")

	// Set up environment with smart GOWORK handling
	env := append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	if shouldDisableWorkspace(appDir) {
		env = append(env, "GOWORK=off")
		fmt.Printf("Note: Disabling go.work for standalone module build\n")
	}
	cmd.Env = env

	absPath, err := filepath.Abs(appDir)
	if err != nil {
		return "", fmt.Errorf("failed to set app dir: %w", err)
	}
	cmd.Dir = absPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	// copy wasm_exec.js from Go SDK
	wasmExecSrc := filepath.Join(runtime.GOROOT(), "lib", "wasm", "wasm_exec.js")
	wasmExecDst := filepath.Join(buildDir, "wasm_exec.js")
	if err := copyFile(wasmExecSrc, wasmExecDst); err != nil {
		return "", err
	}
	return buildDir, nil
}

// findGoWork searches for go.work file starting from dir and walking up
func findGoWork(dir string) string {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return ""
	}

	for {
		workFile := filepath.Join(absDir, "go.work")
		if _, err := os.Stat(workFile); err == nil {
			return workFile
		}

		parent := filepath.Dir(absDir)
		if parent == absDir {
			break // reached root
		}
		absDir = parent
	}
	return ""
}

// parseWorkspaceModules parses go.work file and returns the list of module directories
func parseWorkspaceModules(workFile string) ([]string, error) {
	file, err := os.Open(workFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var modules []string
	scanner := bufio.NewScanner(file)
	inUseBlock := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "use (") {
			inUseBlock = true
			continue
		}
		if strings.HasPrefix(line, "use ") && !strings.Contains(line, "(") {
			// Single line use directive
			module := strings.TrimSpace(strings.TrimPrefix(line, "use"))
			if !strings.HasPrefix(module, "//") && module != "" {
				modules = append(modules, module)
			}
			continue
		}
		if inUseBlock {
			if strings.HasPrefix(line, ")") {
				inUseBlock = false
				continue
			}
			if strings.HasPrefix(line, "//") || line == "" {
				continue
			}
			modules = append(modules, line)
		}
	}

	// Convert relative paths to absolute paths relative to workspace file
	workDir := filepath.Dir(workFile)
	for i, module := range modules {
		if !filepath.IsAbs(module) {
			modules[i] = filepath.Join(workDir, module)
		}
	}

	return modules, scanner.Err()
}

// shouldDisableWorkspace determines if GOWORK should be disabled for the target directory
func shouldDisableWorkspace(targetDir string) bool {
	workFile := findGoWork(targetDir)
	if workFile == "" {
		return false // No workspace to disable
	}

	modules, err := parseWorkspaceModules(workFile)
	if err != nil {
		// If we can't parse the workspace, be conservative and don't disable it
		return false
	}

	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return false
	}

	// Check if target directory matches any workspace module
	for _, module := range modules {
		absModule, err := filepath.Abs(module)
		if err != nil {
			continue
		}
		if absTarget == absModule {
			return false // Target is in workspace, keep GOWORK enabled
		}
	}

	return true // Target not in workspace, disable GOWORK
}

// findFreePort finds and reserves a free port, returning both the port and listener
func findFreePort(preferredPort int) (int, net.Listener, error) {
	// First try the preferred port
	if preferredPort > 0 {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", preferredPort))
		if err == nil {
			return preferredPort, ln, nil
		}
		fmt.Printf("Port %d is in use, finding alternative...\n", preferredPort)
	}

	// Let OS assign a free port
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, nil, err
	}

	port := ln.Addr().(*net.TCPAddr).Port
	return port, ln, nil
}

// copyFile copies a file from src to dst, creating parent directories.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

// watchFiles watches Go source files and calls the provided callback on changes.
func watchFiles(appDir string, onRebuild func() error) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error setting up file watcher: %w", err)
	}
	defer watcher.Close()

	// watch app and local module directories recursively
	// Determine module directories via `go list`
	gomodcacheBytes, err := exec.Command("go", "env", "GOMODCACHE").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting GOMODCACHE: %v\n", err)
	}
	gomodcache := strings.TrimSpace(string(gomodcacheBytes))

	listCmd := exec.Command("go", "list", "-C", appDir, "-m", "-mod=readonly", "-f", "{{.Dir}}", "all")

	// Use same environment setup as build commands for consistency
	env := os.Environ()
	if shouldDisableWorkspace(appDir) {
		env = append(env, "GOWORK=off")
	}
	listCmd.Env = env

	out, err := listCmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing modules: %v\n", err)
	}
	roots := make(map[string]struct{})
	// always include the app directory
	roots[appDir] = struct{}{}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		// skip modules in GOMODCACHE
		if strings.HasPrefix(line, gomodcache) {
			continue
		}
		roots[line] = struct{}{}
	}
	// Walk and watch each root directory
	for root := range roots {
		err = filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if info.IsDir() {
				if watchErr := watcher.Add(path); watchErr != nil {
					fmt.Fprintf(os.Stderr, "Error watching %s: %v\n", path, watchErr)
				}
			}
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking %s for file watching: %v\n", root, err)
		}
	}

	// Debounce mechanism for rebuilds
	rebuildTimer := time.NewTimer(0)
	if !rebuildTimer.Stop() {
		<-rebuildTimer.C
	}
	rebuildPending := false
	debounceDelay := 500 * time.Millisecond

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
				ext := filepath.Ext(event.Name)
				if ext == ".go" || ext == ".mod" || ext == ".sum" {
					fmt.Printf("File changed (%s), scheduling rebuild...\n", event.Name)
					// Reset the timer to debounce multiple rapid changes
					if !rebuildTimer.Stop() && rebuildPending {
						<-rebuildTimer.C
					}
					rebuildTimer.Reset(debounceDelay)
					rebuildPending = true
				}
			}
		case <-rebuildTimer.C:
			if rebuildPending {
				err := onRebuild()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error during rebuild: %v\n", err)
				}
				rebuildPending = false
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Fprintf(os.Stderr, "Watcher error: %v\n", err)
		}
	}
}

// watchAndRebuild watches Go source files and rebuilds the WASM bundle on change.
func watchAndRebuild(appDir string) {
	err := watchFiles(appDir, func() error {
		fmt.Println("Rebuilding...")
		newBuildDir, err := buildWASM(appDir)
		if err != nil {
			return fmt.Errorf("error rebuilding WASM: %w", err)
		}
		buildMutex.Lock()
		old := currentBuildDir
		currentBuildDir = newBuildDir
		buildMutex.Unlock()
		os.RemoveAll(old)
		fmt.Println("Rebuild complete")
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Watch error: %v\n", err)
	}
}

// serve starts an HTTP server on the given port, serving files from dir.
func serve(port int, dir string) error {
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)
	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(addr, nil)
}

// fetchAuthInfo retrieves the Salesforce instance URL and access token
// from the active Force CLI session.
func fetchAuthInfo() (*forcecli.Force, error) {
	creds, err := forcecli.ActiveCredentials(false)
	if err != nil {
		return nil, err
	}
	f := forcecli.NewForce(&creds)
	return f, nil
}

// proxyHandler forwards requests under /services/ to the Salesforce instance
// using the stored session credentials.
// proxyHandler forwards requests under /services/ to the Salesforce instance,
// renewing the session automatically if the access token expires.
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	var resp *http.Response
	for attempt := 0; attempt < 2; attempt++ {
		target := session.Credentials.InstanceUrl + r.RequestURI
		req, err := http.NewRequest(r.Method, target, bytes.NewReader(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for k, vv := range r.Header {
			for _, v := range vv {
				req.Header.Add(k, v)
			}
		}
		req.Header.Set("Authorization", "Bearer "+session.Credentials.AccessToken)

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		if resp.StatusCode == http.StatusUnauthorized && attempt == 0 {
			resp.Body.Close()
			fmt.Fprintf(os.Stderr, "Salesforce session expired, refreshing credentials and retrying\n")
			err := session.RefreshSession()
			if err != nil {
				http.Error(w, fmt.Sprintf("Error renewing session: %v", err), http.StatusBadGateway)
				return
			}
			continue
		}
		break
	}
	if resp == nil {
		http.Error(w, "no response from Salesforce", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// wasmHandler serves the bundle.wasm file from the current build directory.
func wasmHandler(w http.ResponseWriter, r *http.Request) {
	buildMutex.RLock()
	dirPath := currentBuildDir
	buildMutex.RUnlock()
	http.ServeFile(w, r, filepath.Join(dirPath, "bundle.wasm"))
}

// wasmExecHandler serves the wasm_exec.js file from the current build directory.
func wasmExecHandler(w http.ResponseWriter, r *http.Request) {
	buildMutex.RLock()
	dirPath := currentBuildDir
	buildMutex.RUnlock()
	http.ServeFile(w, r, filepath.Join(dirPath, "wasm_exec.js"))
}

// indexHandler serves the indexHTML template directly.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Only serve index for root path and paths that don't match other handlers
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
}

// settingsHandler serves Thunder Settings from environment variables for dev mode.
func settingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create settings response from environment variables
	settings := map[string]interface{}{
		"Google_Maps_API_Key__c": os.Getenv("GOOGLE_MAPS_API_KEY"),
		"error":                  false,
		"message":                "",
	}

	// If no API key is set, return an error
	if settings["Google_Maps_API_Key__c"] == "" {
		settings["error"] = true
		settings["message"] = "GOOGLE_MAPS_API_KEY environment variable not set"
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(settings); err != nil {
		http.Error(w, "Failed to encode settings", http.StatusInternalServerError)
	}
}

// runServe builds the WASM bundle and serves the app with auto-rebuild.
func runServe(cmd *cobra.Command, args []string) error {
	// Determine app directory (optional positional argument)
	if len(args) > 0 {
		serveDir = args[0]
	} else {
		serveDir = "."
	}
	// Validate app directory
	info, err := os.Stat(serveDir)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("Invalid app directory: %s", serveDir)
	}

	// Set up environment for package validation
	env := os.Environ()
	if shouldDisableWorkspace(serveDir) {
		env = append(env, "GOWORK=off")
	}

	cfg := &packages.Config{
		Mode: packages.NeedName,
		Dir:  serveDir,
		Env:  env,
	}
	pkgs, _ := packages.Load(cfg, ".")
	if len(pkgs) == 0 || pkgs[0].Name != "main" {
		return fmt.Errorf("serve directory %s is not package main", serveDir)
	}
	// Fetch Salesforce auth info
	session, err = fetchAuthInfo()
	if err != nil {
		return fmt.Errorf("Error fetching Salesforce auth info: %w", err)
	}
	fmt.Printf("Building WASM bundle in %s...\n", serveDir)
	buildDir, err := buildWASM(serveDir)
	if err != nil {
		return fmt.Errorf("Error building WASM: %w", err)
	}
	buildMutex.Lock()
	currentBuildDir = buildDir
	buildMutex.Unlock()
	go watchAndRebuild(serveDir)

	// Find and reserve a free port
	actualPort, listener, err := findFreePort(servePort)
	if err != nil {
		return fmt.Errorf("Error finding free port: %w", err)
	}
	defer listener.Close()

	fmt.Printf("Serving Thunder app on port %d (watching %s)...\n", actualPort, serveDir)

	// Set up HTTP handlers
	http.HandleFunc("/services/", proxyHandler)
	http.HandleFunc("/cometd/", proxyHandler)
	http.HandleFunc("/api/settings", settingsHandler)
	http.HandleFunc("/bundle.wasm", wasmHandler)
	http.HandleFunc("/wasm_exec.js", wasmExecHandler)
	http.HandleFunc("/", indexHandler)

	// Start the server in a goroutine so we can open browser after it starts
	server := &http.Server{}
	serverStarted := make(chan bool, 1)
	serverErr := make(chan error, 1)

	go func() {
		// Signal that we're about to start the server
		serverStarted <- true
		serverErr <- server.Serve(listener)
	}()

	// Wait for server to start, then open browser
	urlStr := fmt.Sprintf("http://localhost:%d", actualPort)
	go func() {
		select {
		case <-serverStarted:
			// Give server a moment to fully initialize
			time.Sleep(100 * time.Millisecond)
			// Try to open browser
			if err := desktop.Open(urlStr); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to open browser: %v\n", err)
			}
		case err := <-serverErr:
			// Server failed to start immediately, don't open browser
			if err != nil && err != http.ErrServerClosed {
				fmt.Fprintf(os.Stderr, "Server failed to start: %v\n", err)
			}
		}
	}()

	// Wait for server to finish or fail
	return <-serverErr
}

// sanitizeStaticResourceName converts a name to a valid static resource API name (alphanumeric, begins with letter).
func sanitizeStaticResourceName(name string) string {
	re := regexp.MustCompile(`[^A-Za-z0-9]+`)
	name = re.ReplaceAllString(name, "")
	if len(name) == 0 {
		name = "App"
	}
	if !unicode.IsLetter(rune(name[0])) {
		name = "A" + name
	}
	return name
}

// sanitizeComponentName converts an arbitrary name to a valid LWC component name (snake_case, lowercase, begins with letter, no consecutive underscores).
func sanitizeComponentName(name string) string {
	re := regexp.MustCompile(`[^A-Za-z0-9]+`)
	name = re.ReplaceAllString(name, "_")
	// collapse multiple underscores
	name = regexp.MustCompile(`_+`).ReplaceAllString(name, "_")
	name = strings.Trim(name, "_")
	if name == "" {
		name = "app"
	}
	if !unicode.IsLetter(rune(name[0])) {
		name = "a" + name
	}
	return name
}

// toPascalCase converts a snake_case name to PascalCase for component class names.
func toPascalCase(name string) string {
	parts := strings.Split(name, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.Title(p)
		}
	}
	return strings.Join(parts, "")
}

// runBuild handles the build subcommand to compile the app to WebAssembly.
func runBuild(cmd *cobra.Command, args []string) error {
	// Determine app directory (optional positional argument)
	buildDir := "."
	if len(args) > 0 {
		buildDir = args[0]
	}

	// Validate app directory
	info, err := os.Stat(buildDir)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("Invalid app directory: %s", buildDir)
	}

	// Build the WASM
	fmt.Printf("Building Thunder app from %s...\n", buildDir)
	var tempBuildDir string
	if buildDev {
		tempBuildDir, err = buildWASM(buildDir)
	} else {
		tempBuildDir, err = buildProdWASM(buildDir)
	}
	if err != nil {
		return fmt.Errorf("Error building WASM: %w", err)
	}
	defer os.RemoveAll(tempBuildDir)

	// Create output directory if it doesn't exist
	outputDir := buildOutput
	if !filepath.IsAbs(outputDir) {
		outputDir = filepath.Join(buildDir, outputDir)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("Failed to create output directory: %w", err)
	}

	// Copy build artifacts to output directory
	wasmSrc := filepath.Join(tempBuildDir, "bundle.wasm")
	wasmDst := filepath.Join(outputDir, "bundle.wasm")
	if err := copyFile(wasmSrc, wasmDst); err != nil {
		return fmt.Errorf("Failed to copy WASM bundle: %w", err)
	}

	// Copy wasm_exec.js for dev builds
	if buildDev {
		wasmExecSrc := filepath.Join(tempBuildDir, "wasm_exec.js")
		wasmExecDst := filepath.Join(outputDir, "wasm_exec.js")
		if err := copyFile(wasmExecSrc, wasmExecDst); err != nil {
			return fmt.Errorf("Failed to copy wasm_exec.js: %w", err)
		}
	}

	// Report success
	mode := "production"
	if buildDev {
		mode = "development"
	}
	fmt.Printf("\nSuccessfully built Thunder app (%s mode)\n", mode)
	fmt.Printf("Output directory: %s\n", outputDir)
	fmt.Printf("Files generated:\n")
	fmt.Printf("  - bundle.wasm\n")
	if buildDev {
		fmt.Printf("  - wasm_exec.js\n")
	}

	return nil
}

// runDeploy handles the deploy subcommand with optional watch functionality.
func runDeploy(cmd *cobra.Command, args []string) error {
	// Determine app directory (optional positional argument)
	if len(args) > 0 {
		deployDir = args[0]
	} else {
		deployDir = "."
	}
	// Validate app directory
	info, err := os.Stat(deployDir)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("Invalid app directory: %s", deployDir)
	}

	// Set up environment for package validation
	env := os.Environ()
	if shouldDisableWorkspace(deployDir) {
		env = append(env, "GOWORK=off")
	}

	cfg := &packages.Config{
		Mode: packages.NeedName,
		Dir:  deployDir,
		Env:  env,
	}
	pkgs, _ := packages.Load(cfg, ".")
	if len(pkgs) == 0 || pkgs[0].Name != "main" {
		return fmt.Errorf("serve directory %s is not package main", deployDir)
	}
	// Build production WASM bundle
	fmt.Printf("Building production WASM bundle in %s...\n", deployDir)
	absDir, _ := filepath.Abs(deployDir)
	rawName := filepath.Base(absDir)
	staticResourceName := sanitizeStaticResourceName(rawName)
	lwcName := sanitizeComponentName(rawName)
	appClass := toPascalCase(lwcName)
	buildDir, err := buildProdWASM(deployDir)
	if err != nil {
		return fmt.Errorf("Error building production WASM: %w", err)
	}
	fmt.Printf("Built production bundle at %s\n", buildDir)
	// Prepare metadata files in memory
	files := make(forcecli.ForceMetadataFiles)
	// Compress WASM bundle into a zip static resource
	wasmData, err := os.ReadFile(filepath.Join(buildDir, "bundle.wasm"))
	if err != nil {
		return err
	}
	zipData, err := zipBundle(wasmData)
	if err != nil {
		return err
	}
	files["staticresources/"+staticResourceName+".resource"] = zipData
	staticResourceMetadata := `<?xml version="1.0" encoding="UTF-8"?>
<StaticResource xmlns="http://soap.sforce.com/2006/04/metadata">
	<cacheControl>Private</cacheControl>
	<contentType>application/zip</contentType>
</StaticResource>`
	files["staticresources/"+staticResourceName+".resource-meta.xml"] = []byte(staticResourceMetadata)
	// Apex classes
	apexTemplates := []struct{ src, dst string }{
		{"classes/GoBridge.cls", "classes/GoBridge.cls"},
		{"classes/GoBridge.cls-meta.xml", "classes/GoBridge.cls-meta.xml"},
		{"classes/GoBridgeTest.cls", "classes/GoBridgeTest.cls"},
		{"classes/GoBridgeTest.cls-meta.xml", "classes/GoBridgeTest.cls-meta.xml"},
	}
	for _, t := range apexTemplates {
		if data, err := salesforce.SalesforceMetadataFS.ReadFile(t.src); err == nil {
			files[t.dst] = data
		}
	}

	// Custom Objects (Thunder Settings)
	objectTemplates := []struct{ src, dst string }{
		{"objects/Thunder_Settings__c.object", "objects/Thunder_Settings__c.object"},
	}
	for _, t := range objectTemplates {
		if data, err := salesforce.SalesforceMetadataFS.ReadFile(t.src); err == nil {
			files[t.dst] = data
		}
	}
	// LWC components (runtime and wrapper)
	for _, comp := range []string{"go", "thunder"} {
		base := "lwc/" + comp
		entries, _ := salesforce.SalesforceMetadataFS.ReadDir(base)
		for _, e := range entries {
			if data, err := salesforce.SalesforceMetadataFS.ReadFile(base + "/" + e.Name()); err == nil {
				files["lwc/"+comp+"/"+e.Name()] = data
			}
		}
	}
	// Generate LWC for the deployed app
	appComp := lwcName
	// JS wrapper for the app, importing the static resource
	js := fmt.Sprintf(`import Thunder from 'c/thunder';
import APP_URL from '@salesforce/resourceUrl/%s';

export default class %s extends Thunder {
	connectedCallback() {
		this.app = APP_URL + '/bundle.wasm';
	}
}`, staticResourceName, appClass)
	files[fmt.Sprintf("lwc/%s/%s.js", appComp, appComp)] = []byte(js)
	// JS meta
	meta := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<LightningComponentBundle xmlns="http://soap.sforce.com/2006/04/metadata">
    <apiVersion>58.0</apiVersion>
    <isExposed>true</isExposed>
    <masterLabel>%s</masterLabel>
    <targets>
        <target>lightning__AppPage</target>
        <target>lightning__RecordAction</target>
        <target>lightning__RecordPage</target>
        <target>lightning__Tab</target>
    </targets>
    <targetConfigs>
        <targetConfig targets="lightning__RecordAction">
            <actionType>ScreenAction</actionType>
        </targetConfig>
    </targetConfigs>
</LightningComponentBundle>`, appClass)
	files[fmt.Sprintf("lwc/%s/%s.js-meta.xml", appComp, appComp)] = []byte(meta)
	// If requested, generate a CustomTab for the deployed app
	if deployTab {
		tabXml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<CustomTab xmlns="http://soap.sforce.com/2006/04/metadata">
    <label>%s</label>
    <lwcComponent>%s</lwcComponent>
    <motif>Custom75: Default</motif>
</CustomTab>`, appClass, appComp)
		files[fmt.Sprintf("tabs/%s.tab-meta.xml", appComp)] = []byte(tabXml)
	}
	// Generate package.xml for the deployment
	packageTpl := `<?xml version="1.0" encoding="UTF-8"?>
<Package xmlns="http://soap.sforce.com/2006/04/metadata">
  <types>
    <members>%s</members>
    <name>StaticResource</name>
  </types>
  <types>
    <members>GoBridge</members>
    <members>GoBridgeTest</members>
    <name>ApexClass</name>
  </types>
  <types>
    <members>Thunder_Settings__c</members>
    <name>CustomObject</name>
  </types>
  <types>
    <members>go</members>
    <members>thunder</members>
    <members>%s</members>
    <name>LightningComponentBundle</name>
  </types>`
	if deployTab {
		packageTpl += `
  <types>
    <members>%s</members>
    <name>CustomTab</name>
  </types>`
	}
	packageTpl += `
  <version>58.0</version>
</Package>`
	var pkg string
	if deployTab {
		pkg = fmt.Sprintf(packageTpl, staticResourceName, appComp, appComp)
	} else {
		pkg = fmt.Sprintf(packageTpl, staticResourceName, appComp)
	}
	files["package.xml"] = []byte(pkg)
	// Perform initial deployment
	err = performDeployment(files, staticResourceName, appComp, deployTab)
	if err != nil {
		return err
	}

	// Cleanup temporary build directory
	if rmErr := os.RemoveAll(buildDir); rmErr != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to remove temp dir %s: %v\n", buildDir, rmErr)
	}

	// If watch flag is set, start watching for changes
	if deployWatch {
		fmt.Printf("Watching for changes in %s (WASM-only redeploys)...\n", deployDir)
		return watchAndRedeploy(deployDir, staticResourceName)
	}

	return nil
}

// performDeployment deploys the given metadata files to Salesforce
func performDeployment(files forcecli.ForceMetadataFiles, staticResourceName, appComp string, openTab bool) error {
	creds, err := forcecli.ActiveCredentials(false)
	if err != nil {
		return fmt.Errorf("failed to load Salesforce credentials: %w", err)
	}
	fm := forcecli.NewForce(&creds)
	opts := forcecli.ForceDeployOptions{SinglePackage: true}
	fmt.Printf("Deploying metadata to %s...\n", creds.InstanceUrl)
	result, err := fm.Metadata.Deploy(files, opts)
	if err != nil {
		return fmt.Errorf("deployment failed: %w", err)
	}
	if !result.Success {
		fmt.Fprintf(os.Stderr, "Deployment errors:\n")
		for _, failure := range result.Details.ComponentFailures {
			fmt.Fprintf(os.Stderr, "- %s:%d %s: %s\n", failure.FileName, failure.LineNumber, failure.ProblemType, failure.Problem)
		}
		for _, tf := range result.Details.RunTestResult.TestFailures {
			fmt.Fprintf(os.Stderr, "- Test %s.%s: %s\n", tf.Name, tf.MethodName, tf.Message)
		}
		return fmt.Errorf("metadata deployment completed with errors")
	}
	fmt.Printf("Deployment complete: %+v\n", result)

	// Open new tab in Salesforce if requested
	if openTab {
		tabUrl := fmt.Sprintf("%s/lightning/n/%s", creds.InstanceUrl, appComp)
		if deployDebug {
			fmt.Printf("Debug: Opening URL: %s\n", tabUrl)
		}
		if err := desktop.Open(tabUrl); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open tab URL: %v\n", err)
		}
	}
	return nil
}

// watchAndRedeploy watches for Go source file changes and redeploys only the WASM bundle
func watchAndRedeploy(appDir, staticResourceName string) error {
	return watchFiles(appDir, func() error {
		fmt.Println("Rebuilding and redeploying WASM...")
		err := redeployWASM(appDir, staticResourceName)
		if err != nil {
			return fmt.Errorf("error redeploying WASM: %w", err)
		}
		fmt.Println("WASM redeploy complete")
		return nil
	})
}

// redeployWASM builds and redeploys only the WASM static resource
func redeployWASM(appDir, staticResourceName string) error {
	// Build production WASM bundle
	buildDir, err := buildProdWASM(appDir)
	if err != nil {
		return fmt.Errorf("error building production WASM: %w", err)
	}
	defer os.RemoveAll(buildDir)

	// Read and compress WASM bundle
	wasmData, err := os.ReadFile(filepath.Join(buildDir, "bundle.wasm"))
	if err != nil {
		return err
	}
	zipData, err := zipBundle(wasmData)
	if err != nil {
		return err
	}

	// Create minimal deployment with just the static resource
	files := make(forcecli.ForceMetadataFiles)
	files["staticresources/"+staticResourceName+".resource"] = zipData
	staticResourceMetadata := `<?xml version="1.0" encoding="UTF-8"?>
<StaticResource xmlns="http://soap.sforce.com/2006/04/metadata">
	<cacheControl>Private</cacheControl>
	<contentType>application/zip</contentType>
</StaticResource>`
	files["staticresources/"+staticResourceName+".resource-meta.xml"] = []byte(staticResourceMetadata)

	// Generate minimal package.xml for just the static resource
	pkg := `<?xml version="1.0" encoding="UTF-8"?>
<Package xmlns="http://soap.sforce.com/2006/04/metadata">
  <types>
    <members>` + staticResourceName + `</members>
    <name>StaticResource</name>
  </types>
  <version>58.0</version>
</Package>`
	files["package.xml"] = []byte(pkg)

	// Deploy only the WASM static resource
	return performDeployment(files, staticResourceName, "", false)
}

// buildProdWASM compiles the Go app in appDir to WebAssembly for production.
func buildProdWASM(appDir string) (string, error) {
	// create temporary build directory
	buildDir, err := os.MkdirTemp("", "thunder-deploy-*")
	if err != nil {
		return "", err
	}
	outWasm := filepath.Join(buildDir, "bundle.wasm")
	cmd := exec.Command("go", "build", "-o", outWasm)

	// Set up environment with smart GOWORK handling
	env := append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	if shouldDisableWorkspace(appDir) {
		env = append(env, "GOWORK=off")
		fmt.Printf("Note: Disabling go.work for standalone module deployment\n")
	}
	cmd.Env = env

	abs, err := filepath.Abs(appDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}
	cmd.Dir = abs
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return buildDir, nil
}

// zipBundle compresses the WebAssembly binary into a zip archive for StaticResource deployment.
func zipBundle(wasmData []byte) ([]byte, error) {
	buf := &bytes.Buffer{}
	zw := zip.NewWriter(buf)
	w, err := zw.Create("bundle.wasm")
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(wasmData); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
