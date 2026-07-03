package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"compress/flate"
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

// Build-time variable that can be overridden with -ldflags
var osgoPackageVersionId = "04tKe0000008rAoIAI"

// osgoNamespace is the managed package namespace whose GoBridge proxy backs
// deployments that don't use --thunder-dev.
const osgoNamespace = "osgo"

// global state for serve command
var (
	servePort       int
	serveDir        string
	currentBuildDir string
	buildMutex      sync.RWMutex
	session         *forcecli.Force
	// deploy command flags
	deployDir         string
	deployTab         bool
	deployWatch       bool
	deployDebug       bool
	deployAppOnly     bool
	deployThunderDev  bool
	deployVisualforce bool
	deployName        string
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
	deployCmd.Flags().BoolVar(&deployAppOnly, "app-only", false, "Deploy only the static resource (WASM bundle)")
	deployCmd.Flags().BoolVar(&deployThunderDev, "thunder-dev", false, "Deploy unpackaged thunder dependencies instead of using the osgo package")
	deployCmd.Flags().BoolVar(&deployVisualforce, "visualforce", false, "Deploy the app as a Visualforce page (runs outside Lightning Web Security; needed for apps that use Web Workers)")
	deployCmd.Flags().StringVar(&deployName, "name", "", "Name for the app and tab (defaults to directory name)")
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

// packageLoadError returns a non-nil error describing why a Go package failed
// to load, or nil if the package loaded cleanly. It surfaces the underlying
// driver/build errors (e.g. workspace version mismatches) that packages.Load
// reports per-package rather than via its returned error.
func packageLoadError(pkgs []*packages.Package) error {
	if len(pkgs) == 0 {
		return fmt.Errorf("no Go package found")
	}
	var msgs []string
	for _, pkg := range pkgs {
		for _, e := range pkg.Errors {
			msgs = append(msgs, e.Error())
		}
	}
	if len(msgs) > 0 {
		return fmt.Errorf("%s", strings.Join(msgs, "\n"))
	}
	return nil
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
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return fmt.Errorf("failed to load Go package in %s: %w", serveDir, err)
	}
	if loadErr := packageLoadError(pkgs); loadErr != nil {
		return fmt.Errorf("failed to load Go package in %s: %w", serveDir, loadErr)
	}
	if pkgs[0].Name != "main" {
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

// toUpperSnakeCase converts a snake_case name to Upper_Snake_Case for tab names.
func toUpperSnakeCase(name string) string {
	parts := strings.Split(name, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.Title(p)
		}
	}
	return strings.Join(parts, "_")
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

// isOsgoPackageInstalled checks if the osgo package is installed in the org
func isOsgoPackageInstalled(force *forcecli.Force) (bool, error) {
	// Query using Tooling API to check for InstalledSubscriberPackage
	query := fmt.Sprintf("SELECT Id FROM InstalledSubscriberPackage WHERE SubscriberPackageVersionId = '%s' LIMIT 1", osgoPackageVersionId)
	result, err := force.Query(query, func(options *forcecli.QueryOptions) {
		options.IsTooling = true
	})
	if err != nil {
		// If the query fails, try checking for the namespace in ApexClass
		altQuery := "SELECT COUNT() FROM ApexClass WHERE NamespacePrefix = 'osgo'"
		altResult, altErr := force.Query(altQuery)
		if altErr != nil {
			return false, fmt.Errorf("failed to query for osgo package: %w", err)
		}
		if len(altResult.Records) > 0 && altResult.Records[0]["expr0"] != nil {
			count := int(altResult.Records[0]["expr0"].(float64))
			return count > 0, nil
		}
		return false, nil
	}
	return result.TotalSize > 0, nil
}

// installOsgoPackage installs the osgo package using the Tooling API
func installOsgoPackage(force *forcecli.Force) error {
	fmt.Printf("Installing osgo package %s...\n", osgoPackageVersionId)

	// Create PackageInstallRequest using Tooling API
	attrs := map[string]string{
		"SubscriberPackageVersionKey": osgoPackageVersionId,
		"EnableRss":                   "false",
		"NameConflictResolution":      "Block",
		"SecurityType":                "Full",
	}

	result, err := force.CreateToolingRecord("PackageInstallRequest", attrs)
	if err != nil {
		return fmt.Errorf("failed to create PackageInstallRequest: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("failed to create PackageInstallRequest: %v", result.Errors)
	}

	requestId := result.Id
	fmt.Printf("Package install request created: %s\n", requestId)

	// Poll for completion
	for i := 0; i < 120; i++ { // Poll for up to 10 minutes
		time.Sleep(5 * time.Second)

		query := fmt.Sprintf("SELECT Id, Status, Errors FROM PackageInstallRequest WHERE Id = '%s'", requestId)
		queryResult, err := force.Query(query, func(options *forcecli.QueryOptions) {
			options.IsTooling = true
		})
		if err != nil {
			return fmt.Errorf("failed to query PackageInstallRequest status: %w", err)
		}

		if len(queryResult.Records) == 0 {
			return fmt.Errorf("PackageInstallRequest not found: %s", requestId)
		}

		record := queryResult.Records[0]
		status := record["Status"].(string)

		if status == "SUCCESS" {
			fmt.Printf("Successfully installed osgo package\n")
			return nil
		} else if status == "ERROR" {
			errors := ""
			if record["Errors"] != nil {
				errors = fmt.Sprintf("%v", record["Errors"])
			}
			return fmt.Errorf("package installation failed: %s", errors)
		}

		if i%4 == 0 { // Print status every 20 seconds
			fmt.Printf("Installation status: %s\n", status)
		}
	}

	return fmt.Errorf("package installation timed out after 10 minutes")
}

// ensureOsgoPackageInstalled checks if osgo package is installed and installs it if not
func ensureOsgoPackageInstalled() error {
	creds, err := forcecli.ActiveCredentials(false)
	if err != nil {
		return fmt.Errorf("failed to get Salesforce credentials: %w", err)
	}

	force := forcecli.NewForce(&creds)

	installed, err := isOsgoPackageInstalled(force)
	if err != nil {
		// Log warning but continue - package check failed but deployment might still work
		fmt.Printf("Warning: Could not verify osgo package installation: %v\n", err)
		fmt.Printf("Continuing with deployment...\n")
		return nil
	}

	if !installed {
		fmt.Printf("osgo package not found, installing...\n")
		if err := installOsgoPackage(force); err != nil {
			return err
		}
	} else {
		if deployDebug {
			fmt.Printf("osgo package is already installed\n")
		}
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
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return fmt.Errorf("failed to load Go package in %s: %w", deployDir, err)
	}
	if loadErr := packageLoadError(pkgs); loadErr != nil {
		return fmt.Errorf("failed to load Go package in %s: %w", deployDir, loadErr)
	}
	if pkgs[0].Name != "main" {
		return fmt.Errorf("deploy directory %s is not package main", deployDir)
	}

	// Check for osgo package installation and install if needed (unless using --thunder-dev)
	if !deployThunderDev {
		if err := ensureOsgoPackageInstalled(); err != nil {
			return fmt.Errorf("failed to ensure osgo package is installed: %w", err)
		}
	}
	// Build production WASM bundle
	fmt.Printf("Building production WASM bundle in %s...\n", deployDir)
	absDir, _ := filepath.Abs(deployDir)

	// Always use directory name for file names to avoid conflicts
	baseName := filepath.Base(absDir)
	staticResourceName := sanitizeStaticResourceName(baseName)
	lwcName := sanitizeComponentName(baseName)
	appClass := toPascalCase(lwcName)
	buildDir, err := buildProdWASM(deployDir)
	if err != nil {
		return fmt.Errorf("Error building production WASM: %w", err)
	}
	fmt.Printf("Built production bundle at %s\n", buildDir)
	// Prepare metadata files in memory
	files := make(forcecli.ForceMetadataFiles)
	// Compress WASM bundle into one or more zip static resources, splitting it
	// into multiple pieces when it would exceed Salesforce's per-resource limit.
	wasmData, err := os.ReadFile(filepath.Join(buildDir, "bundle.wasm"))
	if err != nil {
		return err
	}
	firstChunkExtras, err := runtimeExtras()
	if err != nil {
		return err
	}
	zipChunks, err := splitAndZip(wasmData, firstChunkExtras...)
	if err != nil {
		return err
	}
	resourceNames := staticResourceNames(staticResourceName, len(zipChunks))
	if len(zipChunks) > 1 {
		fmt.Printf("WASM bundle exceeds the static resource limit; splitting into %d resources\n", len(zipChunks))
	}
	for i, chunk := range zipChunks {
		addStaticResource(files, resourceNames[i], chunk)
	}

	// If --app-only flag is set, only deploy the static resource(s)
	if deployAppOnly {
		// Generate minimal package.xml for just the static resource(s)
		pkg := `<?xml version="1.0" encoding="UTF-8"?>
<Package xmlns="http://soap.sforce.com/2006/04/metadata">
  <types>
` + staticResourceMembersXML(resourceNames) + `    <name>StaticResource</name>
  </types>
  <version>58.0</version>
</Package>`
		files["package.xml"] = []byte(pkg)

		// Perform deployment
		err = performDeployment(files, staticResourceName, "", false, nil)
		if err != nil {
			return err
		}

		// Cleanup temporary build directory
		if rmErr := os.RemoveAll(buildDir); rmErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to remove temp dir %s: %v\n", buildDir, rmErr)
		}

		return nil
	}

	// If --visualforce is set, deploy the app as a Visualforce page instead of an
	// LWC. Visualforce pages render in a plain iframe outside Lightning Web
	// Security, so apps that create Web Workers from blob URLs (which LWS rejects
	// with "Unsupported MIME type") run there. The REST proxy reaches the same
	// GoBridge logic through JavaScript Remoting rather than the @AuraEnabled path.
	if deployVisualforce {
		appName := deployName
		if appName == "" {
			appName = appClass
		}
		tabName := toUpperSnakeCase(lwcName)
		// Like the LWC deployment, the page reaches GoBridge from the osgo
		// managed package unless --thunder-dev deploys it unmanaged; the
		// controller and RemoteAction references are namespaced accordingly. The
		// Go runtime (wasm_exec.js) always ships inside the app's static resource
		// next to bundle.wasm, so neither path deploys a separate copy.
		controllerClass := osgoNamespace + ".GoBridge"
		if deployThunderDev {
			controllerClass = "GoBridge"
		}
		if err := addVisualforceMetadata(files, resourceNames, appClass, appName, tabName, controllerClass, deployThunderDev); err != nil {
			return err
		}
		files["package.xml"] = []byte(buildVisualforcePackageXML(resourceNames, appClass, tabName, deployTab, deployThunderDev))

		// Only the unmanaged GoBridge deployment includes GoBridgeTest to run.
		var runTests []string
		if deployThunderDev {
			runTests = []string{"GoBridgeTest"}
		}
		if err := performDeployment(files, staticResourceName, "", false, runTests); err != nil {
			return err
		}
		openVisualforceApp(appClass, tabName, deployTab)

		if rmErr := os.RemoveAll(buildDir); rmErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to remove temp dir %s: %v\n", buildDir, rmErr)
		}

		if deployWatch {
			fmt.Printf("Watching for changes in %s (WASM-only redeploys)...\n", deployDir)
			return watchAndRedeploy(deployDir, staticResourceName)
		}
		return nil
	}

	// Deploy thunder dependencies if using --thunder-dev flag
	if deployThunderDev {
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
	}

	// Generate LWC for the deployed app
	appComp := lwcName
	// JS wrapper for the app, importing the static resource
	// Import thunder from the appropriate namespace
	thunderImport := "osgo/thunder"
	if deployThunderDev {
		thunderImport = "c/thunder"
	}
	// Use deployName for appName if provided, otherwise use appClass
	appName := deployName
	if appName == "" {
		appName = appClass
	}

	js := generateAppJS(thunderImport, appClass, appName, staticResourceName)
	files[fmt.Sprintf("lwc/%s/%s.js", appComp, appComp)] = []byte(js)
	// JS meta
	files[fmt.Sprintf("lwc/%s/%s.js-meta.xml", appComp, appComp)] = []byte(lwcMetaXML(appClass))

	// Calculate tab name for both tab generation and package.xml
	tabName := toUpperSnakeCase(lwcName)

	// If requested, generate a CustomTab for the deployed app
	if deployTab {
		tabXml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<CustomTab xmlns="http://soap.sforce.com/2006/04/metadata">
    <label>%s</label>
    <lwcComponent>%s</lwcComponent>
    <motif>Custom75: Default</motif>
</CustomTab>`, appName, appComp)
		files[fmt.Sprintf("tabs/%s.tab-meta.xml", tabName)] = []byte(tabXml)
	}
	// Generate package.xml for the deployment
	files["package.xml"] = []byte(buildPackageXML(resourceNames, appComp, tabName, deployThunderDev, deployTab))
	// Perform initial deployment. When deploying the unpackaged thunder
	// dependencies (which include GoBridgeTest), run only that test.
	var runTests []string
	if deployThunderDev {
		runTests = []string{"GoBridgeTest"}
	}
	err = performDeployment(files, staticResourceName, appComp, deployTab, runTests)
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
func performDeployment(files forcecli.ForceMetadataFiles, staticResourceName, appComp string, openTab bool, runTests []string) error {
	creds, err := forcecli.ActiveCredentials(false)
	if err != nil {
		return fmt.Errorf("failed to load Salesforce credentials: %w", err)
	}
	fm := forcecli.NewForce(&creds)
	opts := forcecli.ForceDeployOptions{
		SinglePackage:   true,
		RollbackOnError: true,
	}
	// Run only the named Apex tests when any are deploying (e.g. GoBridgeTest),
	// rather than the org's entire test suite. With no tests, skip execution.
	if len(runTests) > 0 {
		opts.TestLevel = "RunSpecifiedTests"
		opts.RunTests = runTests
	} else {
		opts.RunTests = []string{}
	}
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

// redeployWASM builds and redeploys the WASM static resource(s). The LWC wrapper
// always imports the base resource and reads the chunk count from the base
// chunk's parts.json manifest at runtime, so it never needs regenerating when the
// chunk count changes between builds.
func redeployWASM(appDir, staticResourceName string) error {
	// Build production WASM bundle
	buildDir, err := buildProdWASM(appDir)
	if err != nil {
		return fmt.Errorf("error building production WASM: %w", err)
	}
	defer os.RemoveAll(buildDir)

	// Read and compress WASM bundle into one or more static resources
	wasmData, err := os.ReadFile(filepath.Join(buildDir, "bundle.wasm"))
	if err != nil {
		return err
	}
	firstChunkExtras, err := runtimeExtras()
	if err != nil {
		return err
	}
	zipChunks, err := splitAndZip(wasmData, firstChunkExtras...)
	if err != nil {
		return err
	}
	resourceNames := staticResourceNames(staticResourceName, len(zipChunks))

	files := make(forcecli.ForceMetadataFiles)
	for i, chunk := range zipChunks {
		addStaticResource(files, resourceNames[i], chunk)
	}

	// Generate minimal package.xml for just the static resource(s)
	pkg := `<?xml version="1.0" encoding="UTF-8"?>
<Package xmlns="http://soap.sforce.com/2006/04/metadata">
  <types>
` + staticResourceMembersXML(resourceNames) + `    <name>StaticResource</name>
  </types>
  <version>58.0</version>
</Package>`
	files["package.xml"] = []byte(pkg)

	// Deploy only the WASM static resource(s)
	return performDeployment(files, staticResourceName, "", false, nil)
}

// buildProdWASM compiles the Go app in appDir to WebAssembly for production.
// It strips debug symbols (-s -w) and trims source paths (-trimpath) to keep
// the bundle small enough to fit Salesforce's 5MB static-resource limit, and
// optionally post-processes the output with wasm-opt -Oz when available.
func buildProdWASM(appDir string) (string, error) {
	// create temporary build directory
	buildDir, err := os.MkdirTemp("", "thunder-deploy-*")
	if err != nil {
		return "", err
	}
	outWasm := filepath.Join(buildDir, "bundle.wasm")
	cmd := exec.Command("go", "build", "-trimpath", "-ldflags=-s -w", "-o", outWasm)

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

	optimizeWASM(outWasm)
	return buildDir, nil
}

// optimizeWASM runs wasm-opt -Oz on the given file when wasm-opt is found on
// PATH, replacing the file in place on success. Failures are logged but
// non-fatal — the un-optimized bundle remains usable.
func optimizeWASM(wasmPath string) {
	bin, err := exec.LookPath("wasm-opt")
	if err != nil {
		fmt.Println("Note: wasm-opt not found on PATH; skipping post-optimization (install Binaryen for further size reduction).")
		return
	}
	before, statErr := os.Stat(wasmPath)
	tmpOut := wasmPath + ".opt"
	// Enable the wasm features Go's compiler emits. --all-features is broad
	// but safe and avoids breakage as the Go toolchain adopts new features.
	cmd := exec.Command(bin, "-Oz", "--all-features", wasmPath, "-o", tmpOut)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("Running %s -Oz on bundle.wasm...\n", bin)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: wasm-opt failed (%v); using unoptimized bundle\n", err)
		os.Remove(tmpOut)
		return
	}
	if err := os.Rename(tmpOut, wasmPath); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not replace bundle with optimized output (%v); using unoptimized bundle\n", err)
		os.Remove(tmpOut)
		return
	}
	if statErr == nil {
		if after, err := os.Stat(wasmPath); err == nil {
			fmt.Printf("wasm-opt: %d -> %d bytes\n", before.Size(), after.Size())
		}
	}
}

// staticResourceLimit is Salesforce's maximum size for a single static resource (5 MB).
const staticResourceLimit = 5 * 1024 * 1024

// splitAndZip compresses the WebAssembly binary into one or more zip archives,
// each holding a contiguous slice of the bundle as bundle.wasm and each small
// enough to deploy as an individual Salesforce static resource. The first archive
// always carries a parts.json manifest recording the total chunk count (1 when
// the bundle fits in a single resource); the runtime loader treats a resource
// without that manifest as a legacy single-part app.
// runtimeExtras returns the extra files to pack into the first WASM static
// resource: the Go runtime (wasm_exec.js) from the SDK that built the bundle.
// Taking it from the same SDK keeps the runtime in lockstep with the compiler.
// Visualforce pages load it from the resource as a script; LWC deployments load
// the runtime through the go LWC and ignore it. It ships unconditionally so
// WASM-only redeploys (--app-only, --watch) don't need to know how the app was
// originally deployed.
func runtimeExtras() ([]zipEntry, error) {
	wasmExecPath := filepath.Join(runtime.GOROOT(), "lib", "wasm", "wasm_exec.js")
	wasmExec, err := os.ReadFile(wasmExecPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read wasm_exec.js from the Go SDK: %w", err)
	}
	return []zipEntry{{name: "wasm_exec.js", data: wasmExec}}, nil
}

func splitAndZip(wasmData []byte, firstChunkExtras ...zipEntry) ([][]byte, error) {
	return splitAndZipWithLimit(wasmData, staticResourceLimit, firstChunkExtras...)
}

// splitAndZipWithLimit is splitAndZip with an explicit per-resource byte limit
// (parameterized for testing). firstChunkExtras are additional files packed into
// the first archive only (e.g. wasm_exec.js for Visualforce deployments).
func splitAndZipWithLimit(wasmData []byte, limit int, firstChunkExtras ...zipEntry) ([][]byte, error) {
	// The single-resource case still carries a parts.json manifest (parts=1) so
	// newly deployed apps are described uniformly.
	whole, err := zipChunkWithManifest(wasmData, 1, firstChunkExtras...)
	if err != nil {
		return nil, err
	}
	if len(whole) <= limit {
		return [][]byte{whole}, nil
	}
	// Estimate the compression ratio from the whole-bundle archive so raw chunks
	// are sized to land comfortably under the limit once compressed. Target 90%
	// of the limit to leave headroom for variation in compressibility.
	ratio := float64(len(whole)) / float64(len(wasmData))
	rawChunk := int(float64(limit) * 0.9 / ratio)
	if rawChunk < 1 {
		rawChunk = limit
	}
	for {
		chunks, ok, err := chunkAndZip(wasmData, rawChunk, limit, firstChunkExtras...)
		if err != nil {
			return nil, err
		}
		if ok {
			return chunks, nil
		}
		// A chunk exceeded the limit (an atypically incompressible region);
		// shrink the raw chunk size and retry.
		rawChunk = rawChunk * 4 / 5
		if rawChunk < 1024 {
			return nil, fmt.Errorf("unable to split WASM bundle under the %d-byte static resource limit", limit)
		}
	}
}

// chunkAndZip slices wasmData into rawChunk-sized pieces and zips each. The first
// chunk also carries a parts.json manifest recording the total chunk count so the
// runtime loader knows how many sibling Part resources to fetch. It reports
// ok=false (without error) if any resulting archive exceeds limit so the caller
// can retry with a smaller chunk size.
func chunkAndZip(wasmData []byte, rawChunk, limit int, firstChunkExtras ...zipEntry) (chunks [][]byte, ok bool, err error) {
	count := (len(wasmData) + rawChunk - 1) / rawChunk
	for off, idx := 0, 0; off < len(wasmData); off, idx = off+rawChunk, idx+1 {
		end := off + rawChunk
		if end > len(wasmData) {
			end = len(wasmData)
		}
		var z []byte
		if idx == 0 {
			z, err = zipChunkWithManifest(wasmData[off:end], count, firstChunkExtras...)
		} else {
			z, err = zipBundle(wasmData[off:end])
		}
		if err != nil {
			return nil, false, err
		}
		if len(z) > limit {
			return nil, false, nil
		}
		chunks = append(chunks, z)
	}
	return chunks, true, nil
}

// staticResourceNames returns the static resource name for each WASM chunk in
// load order. Chunk 0 keeps the bare base name so the LWC wrapper's
// resourceUrl import is unchanged; additional chunks are suffixed Part1, Part2,
// ... and discovered at runtime from the base chunk's parts.json manifest.
func staticResourceNames(base string, n int) []string {
	names := make([]string, n)
	names[0] = base
	for i := 1; i < n; i++ {
		names[i] = fmt.Sprintf("%sPart%d", base, i)
	}
	return names
}

// staticResourceMetaXML is the metadata accompanying every WASM static resource.
const staticResourceMetaXML = `<?xml version="1.0" encoding="UTF-8"?>
<StaticResource xmlns="http://soap.sforce.com/2006/04/metadata">
	<cacheControl>Private</cacheControl>
	<contentType>application/zip</contentType>
</StaticResource>`

// addStaticResource registers a zipped WASM chunk and its metadata under name.
func addStaticResource(files forcecli.ForceMetadataFiles, name string, zipData []byte) {
	files["staticresources/"+name+".resource"] = zipData
	files["staticresources/"+name+".resource-meta.xml"] = []byte(staticResourceMetaXML)
}

// staticResourceMembersXML renders the <members> lines for a package.xml
// StaticResource type, one per chunk resource.
func staticResourceMembersXML(names []string) string {
	var b strings.Builder
	for _, n := range names {
		b.WriteString("    <members>" + n + "</members>\n")
	}
	return b.String()
}

// generateAppJS builds the LWC JavaScript wrapper for a deployed app. It always
// imports the base static resource and sets this.app; when the bundle is split
// across several resources Thunder discovers the additional chunks at runtime from
// the base chunk's parts.json manifest, so the wrapper itself does not change with
// chunk count.
func generateAppJS(thunderImport, appClass, appName, baseResourceName string) string {
	return fmt.Sprintf(`import Thunder from '%s';
import APP_URL from '@salesforce/resourceUrl/%s';

export default class %s extends Thunder {
	connectedCallback() {
		this.app = APP_URL + '/bundle.wasm';
		this.appName = '%s';
	}
}`, thunderImport, baseResourceName, appClass, appName)
}

// lwcMetaXML returns the LightningComponentBundle metadata for a generated app
// wrapper with the given master label.
func lwcMetaXML(masterLabel string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<LightningComponentBundle xmlns="http://soap.sforce.com/2006/04/metadata">
    <apiVersion>58.0</apiVersion>
    <isExposed>true</isExposed>
    <masterLabel>%s</masterLabel>
    <targets>
        <target>lightning__AppPage</target>
        <target>lightning__HomePage</target>
        <target>lightning__RecordAction</target>
        <target>lightning__RecordPage</target>
        <target>lightning__Tab</target>
    </targets>
    <targetConfigs>
        <targetConfig targets="lightning__RecordAction">
            <actionType>ScreenAction</actionType>
        </targetConfig>
    </targetConfigs>
</LightningComponentBundle>`, masterLabel)
}

// buildPackageXML assembles the deployment manifest covering the WASM static
// resources, the generated LWC (plus Thunder runtime components when deploying
// with --thunder-dev), and an optional CustomTab.
func buildPackageXML(resourceNames []string, appComp, tabName string, thunderDev, withTab bool) string {
	var b strings.Builder
	b.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	b.WriteString("<Package xmlns=\"http://soap.sforce.com/2006/04/metadata\">\n")
	b.WriteString("  <types>\n")
	b.WriteString(staticResourceMembersXML(resourceNames))
	b.WriteString("    <name>StaticResource</name>\n")
	b.WriteString("  </types>\n")
	if thunderDev {
		b.WriteString("  <types>\n    <members>GoBridge</members>\n    <members>GoBridgeTest</members>\n    <name>ApexClass</name>\n  </types>\n")
		b.WriteString("  <types>\n    <members>Thunder_Settings__c</members>\n    <name>CustomObject</name>\n  </types>\n")
		b.WriteString("  <types>\n    <members>go</members>\n    <members>thunder</members>\n    <members>" + appComp + "</members>\n    <name>LightningComponentBundle</name>\n  </types>\n")
	} else {
		b.WriteString("  <types>\n    <members>" + appComp + "</members>\n    <name>LightningComponentBundle</name>\n  </types>\n")
	}
	if withTab {
		b.WriteString("  <types>\n    <members>" + tabName + "</members>\n    <name>CustomTab</name>\n  </types>\n")
	}
	b.WriteString("  <version>58.0</version>\n")
	b.WriteString("</Package>")
	return b.String()
}

// addVisualforceMetadata registers everything a Visualforce-hosted app needs:
// the Visualforce page itself and an optional CustomTab pointing at it. When
// thunderDev is set it also bundles the GoBridge proxy classes and Thunder
// Settings object unmanaged, so the page's JavaScript Remoting can reach them.
// Otherwise the page references the osgo managed package's GoBridge and only the
// app metadata is deployed (like the LWC deployment). The Go runtime
// (wasm_exec.js) rides inside the app's WASM static resource in both cases.
func addVisualforceMetadata(files forcecli.ForceMetadataFiles, resourceNames []string, pageName, appName, tabName, controllerClass string, thunderDev bool) error {
	if thunderDev {
		// GoBridge proxy classes and the Thunder Settings object, deployed
		// unmanaged so the page can call GoBridge.remoteCallRest via JavaScript
		// Remoting.
		apexTemplates := []string{
			"classes/GoBridge.cls",
			"classes/GoBridge.cls-meta.xml",
			"classes/GoBridgeTest.cls",
			"classes/GoBridgeTest.cls-meta.xml",
			"objects/Thunder_Settings__c.object",
		}
		for _, src := range apexTemplates {
			data, err := salesforce.SalesforceMetadataFS.ReadFile(src)
			if err != nil {
				return fmt.Errorf("failed to read embedded %s: %w", src, err)
			}
			files[src] = data
		}
	}

	// The Visualforce page and its metadata.
	page := generateVisualforcePage(appName, controllerClass, resourceNames)
	files[fmt.Sprintf("pages/%s.page", pageName)] = []byte(page)
	files[fmt.Sprintf("pages/%s.page-meta.xml", pageName)] = []byte(visualforcePageMetaXML(appName))

	// Optional CustomTab pointing at the page.
	if deployTab {
		tabXml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<CustomTab xmlns="http://soap.sforce.com/2006/04/metadata">
    <label>%s</label>
    <page>%s</page>
    <motif>Custom75: Default</motif>
</CustomTab>`, appName, pageName)
		files[fmt.Sprintf("tabs/%s.tab-meta.xml", tabName)] = []byte(tabXml)
	}
	return nil
}

// wasmURLListJS renders the comma-separated list of Visualforce $Resource URLs
// for each WASM chunk, in load order, for embedding in the page's script. Unlike
// the LWC loader, the chunk count is known at deploy time, so the exact list is
// emitted rather than discovered from a runtime manifest.
func wasmURLListJS(resourceNames []string) string {
	var parts []string
	for _, n := range resourceNames {
		parts = append(parts, fmt.Sprintf("\"{!URLFOR($Resource.%s, 'bundle.wasm')}\"", n))
	}
	return strings.Join(parts, ",\n\t\t\t")
}

// generateVisualforcePage builds the Visualforce page that hosts the Go WASM app.
// The page loads the Go runtime and the WASM static resource(s), exposes the same
// global functions the app expects (get/post/patch/delete plus record/exit
// helpers), instantiates the module, and hands a container div to startWithDiv.
// REST calls reach <controllerClass>.remoteCallRest through JavaScript Remoting.
// The Go runtime (wasm_exec.js) is loaded from inside the first WASM static
// resource (resourceNames[0]), where it was packed next to bundle.wasm.
func generateVisualforcePage(appName, controllerClass string, resourceNames []string) string {
	wasmExecURL := fmt.Sprintf("{!URLFOR($Resource.%s, 'wasm_exec.js')}", resourceNames[0])
	// The $RemoteAction merge field takes only <controller>.<method> and appends
	// the controller's namespace automatically; including a namespace prefix here
	// (e.g. osgo.GoBridge) makes Visualforce reject the page. The controller
	// attribute, however, keeps the namespace so the page binds to the managed
	// class, so strip the namespace only for the RemoteAction reference.
	remoteActionClass := controllerClass
	if i := strings.LastIndex(remoteActionClass, "."); i >= 0 {
		remoteActionClass = remoteActionClass[i+1:]
	}
	return fmt.Sprintf(`<apex:page controller="%[4]s" showHeader="false" sidebar="false" standardStylesheets="false" applyHtmlTag="false" applyBodyTag="false" docType="html-5.0">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
	<meta charset="utf-8"/>
	<title>%[1]s</title>
	<apex:slds />
	<apex:includeScript value="%[2]s"/>
</head>
<body class="slds-scope">
	<div id="thunder-spinner" class="slds-spinner_container">
		<div role="status" class="slds-spinner slds-spinner_medium">
			<span class="slds-assistive-text">Loading</span>
			<div class="slds-spinner__dot-a"></div>
			<div class="slds-spinner__dot-b"></div>
		</div>
	</div>
	<div id="thunder-app" class="slds-scope"></div>
	<script type="text/javascript">
	(function () {
		"use strict";

		// Record context, if the page was opened with ?id=<recordId>. When absent
		// the merge field renders empty; expose undefined (not "") from
		// getRecordIdForDiv so api.RecordId() reports "no record" and the app
		// falls back to its own behavior instead of querying for an empty id.
		var recordIdParam = "{!$CurrentPage.parameters.id}";
		var recordId = recordIdParam || undefined;
		globalThis.recordId = recordIdParam;
		globalThis.getRecordIdForDiv = function () { return recordId; };

		// REST proxy: reach GoBridge through JavaScript Remoting. The Go runtime
		// expects a Promise that resolves to the JSON response string, mirroring
		// the LWC callRest proxy.
		function remoteCall(method, url, body) {
			return new Promise(function (resolve, reject) {
				Visualforce.remoting.Manager.invokeAction(
					"{!$RemoteAction.%[5]s.remoteCallRest}",
					method, url, body,
					function (result, event) {
						if (event.status) {
							resolve(result);
						} else {
							reject(event.message || "Remoting error");
						}
					},
					{ escape: false, timeout: 120000 }
				);
			});
		}

		globalThis.get = function (url) { return remoteCall("GET", url, null); };
		globalThis.post = function (url, body) { return remoteCall("POST", url, body); };
		globalThis.patch = function (url, body) { return remoteCall("PATCH", url, body); };
		globalThis.delete = function (url) { return remoteCall("DELETE", url, null); };

		// The UI API adapters (picklist/object-info) are Lightning-only; surface a
		// clear error instead of letting the Go runtime call an undefined function.
		function unsupported(name) {
			return function (config, cb) {
				cb({ error: name + " is not available on Visualforce pages" });
			};
		}
		globalThis.getPicklistValuesByRecordType = unsupported("getPicklistValuesByRecordType");
		globalThis.getObjectInfo = unsupported("getObjectInfo");

		// Exit/navigation helpers. On a standalone page, navigate the top window.
		globalThis.thunderExitToRecord = function (id) {
			if (id) { window.top.location.href = "/" + id; }
		};
		globalThis.thunderExit = function () {
			if (recordId) { window.top.location.href = "/" + recordId; }
			else { window.history.back(); }
		};
		globalThis.thunderCloseModal = function () { window.history.back(); };

		// A WASM bundle larger than Salesforce's 5MB static resource limit is split
		// across several resources; concatenate the chunks before instantiating.
		var wasmUrls = [
			%[3]s
		];

		function concatBuffers(buffers) {
			if (buffers.length === 1) { return buffers[0]; }
			var total = buffers.reduce(function (sum, buf) { return sum + buf.byteLength; }, 0);
			var combined = new Uint8Array(total);
			var offset = 0;
			buffers.forEach(function (buf) {
				combined.set(new Uint8Array(buf), offset);
				offset += buf.byteLength;
			});
			return combined.buffer;
		}

		function hideSpinner() {
			var spinner = document.getElementById("thunder-spinner");
			if (spinner) { spinner.style.display = "none"; }
		}

		function start() {
			Promise.all(wasmUrls.map(function (u) { return fetch(u); }))
				.then(function (responses) {
					var failed = responses.find(function (r) { return !r.ok; });
					if (failed) {
						return failed.text().then(function (t) { throw new Error(t); });
					}
					return Promise.all(responses.map(function (r) { return r.arrayBuffer(); }));
				})
				.then(function (buffers) {
					var src = concatBuffers(buffers);
					var go = new Go();
					return WebAssembly.instantiate(src, go.importObject).then(function (result) {
						go.run(result.instance);
						return new Promise(function (resolve) { setTimeout(resolve, 1000); });
					});
				})
				.then(function () {
					hideSpinner();
					startWithDiv(document.getElementById("thunder-app"));
				})
				.catch(function (err) {
					hideSpinner();
					var pre = document.createElement("pre");
					pre.innerText = (err && err.message) || String(err);
					document.getElementById("thunder-app").appendChild(pre);
				});
		}

		if (document.readyState === "loading") {
			document.addEventListener("DOMContentLoaded", start);
		} else {
			start();
		}
	})();
	</script>
</body>
</html>
</apex:page>`, appName, wasmExecURL, wasmURLListJS(resourceNames), controllerClass, remoteActionClass)
}

// visualforcePageMetaXML returns the ApexPage metadata for a generated page.
func visualforcePageMetaXML(label string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<ApexPage xmlns="http://soap.sforce.com/2006/04/metadata">
    <apiVersion>58.0</apiVersion>
    <availableInTouch>true</availableInTouch>
    <label>%s</label>
</ApexPage>`, label)
}

// buildVisualforcePackageXML assembles the deployment manifest for a Visualforce
// app: the WASM static resources (which carry wasm_exec.js), the page, and an
// optional CustomTab. Under thunderDev it also lists the GoBridge proxy classes
// and the Thunder Settings object that are deployed unmanaged; otherwise those
// come from the osgo managed package and are omitted.
func buildVisualforcePackageXML(resourceNames []string, pageName, tabName string, withTab, thunderDev bool) string {
	var b strings.Builder
	b.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	b.WriteString("<Package xmlns=\"http://soap.sforce.com/2006/04/metadata\">\n")
	b.WriteString("  <types>\n")
	b.WriteString(staticResourceMembersXML(resourceNames))
	b.WriteString("    <name>StaticResource</name>\n")
	b.WriteString("  </types>\n")
	if thunderDev {
		b.WriteString("  <types>\n    <members>GoBridge</members>\n    <members>GoBridgeTest</members>\n    <name>ApexClass</name>\n  </types>\n")
		b.WriteString("  <types>\n    <members>Thunder_Settings__c</members>\n    <name>CustomObject</name>\n  </types>\n")
	}
	b.WriteString("  <types>\n    <members>" + pageName + "</members>\n    <name>ApexPage</name>\n  </types>\n")
	if withTab {
		b.WriteString("  <types>\n    <members>" + tabName + "</members>\n    <name>CustomTab</name>\n  </types>\n")
	}
	b.WriteString("  <version>58.0</version>\n")
	b.WriteString("</Package>")
	return b.String()
}

// openVisualforceApp opens the deployed Visualforce app in the browser: the
// CustomTab when one was deployed, otherwise the page directly.
func openVisualforceApp(pageName, tabName string, withTab bool) {
	creds, err := forcecli.ActiveCredentials(false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not determine instance URL to open app: %v\n", err)
		return
	}
	var url string
	if withTab {
		url = fmt.Sprintf("%s/lightning/n/%s", creds.InstanceUrl, tabName)
	} else {
		url = fmt.Sprintf("%s/apex/%s", creds.InstanceUrl, pageName)
	}
	if deployDebug {
		fmt.Printf("Debug: Opening URL: %s\n", url)
	}
	if err := desktop.Open(url); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open app URL: %v\n", err)
	}
}

// zipEntry is a single file to place in a static resource zip archive.
type zipEntry struct {
	name string
	data []byte
}

// zipFiles compresses the given entries into a zip archive for StaticResource
// deployment. It uses flate.BestCompression to squeeze the result as close to
// Salesforce's 5MB per-static-resource limit as possible.
func zipFiles(entries []zipEntry) ([]byte, error) {
	buf := &bytes.Buffer{}
	zw := zip.NewWriter(buf)
	zw.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})
	for _, e := range entries {
		w, err := zw.Create(e.name)
		if err != nil {
			return nil, err
		}
		if _, err := w.Write(e.data); err != nil {
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// zipBundle compresses the WebAssembly binary into a single-file zip archive
// containing only bundle.wasm. It is used for the trailing chunks of a split
// bundle, which carry no manifest. (A resource with neither additional chunks nor
// a parts.json manifest is what legacy, pre-split apps deployed, and the runtime
// loader treats such a resource as a single-part bundle.)
func zipBundle(wasmData []byte) ([]byte, error) {
	return zipFiles([]zipEntry{{name: "bundle.wasm", data: wasmData}})
}

// zipChunkWithManifest compresses the first chunk of a split bundle, adding a
// parts.json manifest that records how many static resources the full bundle was
// split across. The runtime loader reads it to fetch and concatenate the
// remaining Part resources before instantiating the WASM module. extras are any
// additional files to pack alongside (e.g. wasm_exec.js for Visualforce apps).
func zipChunkWithManifest(wasmData []byte, parts int, extras ...zipEntry) ([]byte, error) {
	manifest := fmt.Sprintf(`{"parts":%d}`, parts)
	entries := []zipEntry{
		{name: "bundle.wasm", data: wasmData},
		{name: "parts.json", data: []byte(manifest)},
	}
	entries = append(entries, extras...)
	return zipFiles(entries)
}
