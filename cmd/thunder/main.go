package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

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
	// instanceURL and accessToken for Salesforce REST proxy
	instanceURL string
	accessToken string
	// deploy command flags
	deployDir string
	deployTab bool
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

func init() {
	// serve flags (port only; app dir is optional positional arg)
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8000, "Port to serve on")
	// deploy flags (app dir is optional positional arg)
	deployCmd.Flags().BoolVarP(&deployTab, "tab", "t", false, "Deploy and open a CustomTab for the app")
	// add subcommands
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(deployCmd)
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
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
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

// watchAndRebuild watches Go source files and rebuilds the WASM bundle on change.
func watchAndRebuild(appDir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up file watcher: %v\n", err)
		return
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
	listCmd.Env = os.Environ()
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

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
				ext := filepath.Ext(event.Name)
				if ext == ".go" || ext == ".mod" || ext == ".sum" {
					fmt.Printf("File changed (%s), rebuilding...\n", event.Name)
					newBuildDir, err := buildWASM(appDir)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error rebuilding WASM: %v\n", err)
						continue
					}
					buildMutex.Lock()
					old := currentBuildDir
					currentBuildDir = newBuildDir
					buildMutex.Unlock()
					os.RemoveAll(old)
					fmt.Println("Rebuild complete")
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Fprintf(os.Stderr, "Watcher error: %v\n", err)
		}
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
func fetchAuthInfo() (string, string, error) {
	creds, err := forcecli.ActiveCredentials(false)
	if err != nil {
		return "", "", err
	}
	return creds.InstanceUrl, creds.AccessToken, nil
}

// proxyHandler forwards requests under /services/ to the Salesforce instance
// using the stored session credentials.
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Construct target URL
	target := instanceURL + r.RequestURI
	// Create new request
	req, err := http.NewRequest(r.Method, target, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Copy headers
	for k, vv := range r.Header {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}
	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)
	// Forward request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	// Copy response headers
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
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
	cfg := &packages.Config{
		Mode: packages.NeedName,
		Dir:  serveDir,
	}
	pkgs, _ := packages.Load(cfg, ".")
	if len(pkgs) == 0 || pkgs[0].Name != "main" {
		return fmt.Errorf("serve directory %s is not package main", serveDir)
	}
	// Fetch Salesforce auth info
	instanceURL, accessToken, err = fetchAuthInfo()
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

	fmt.Printf("Serving Thunder app on port %d (watching %s)...\n", servePort, serveDir)
	// Open default browser to the served app
	urlStr := fmt.Sprintf("http://localhost:%d", servePort)
	go func() {
		if err := desktop.Open(urlStr); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open browser: %v\n", err)
		}
	}()
	// Proxy Salesforce REST API requests
	http.HandleFunc("/services/", proxyHandler)
	// Serve static assets
	http.HandleFunc("/bundle.wasm", wasmHandler)
	http.HandleFunc("/wasm_exec.js", wasmExecHandler)
	// Serve index HTML
	http.HandleFunc("/", indexHandler)
	address := fmt.Sprintf(":%d", servePort)
	if err := http.ListenAndServe(address, nil); err != nil {
		return fmt.Errorf("Error starting server: %w", err)
	}
	return nil
}

// runDeploy is a stub for the deploy subcommand.
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
	cfg := &packages.Config{
		Mode: packages.NeedName,
		Dir:  deployDir,
	}
	pkgs, _ := packages.Load(cfg, ".")
	if len(pkgs) == 0 || pkgs[0].Name != "main" {
		return fmt.Errorf("serve directory %s is not package main", deployDir)
	}
	// Build production WASM bundle
	fmt.Printf("Building production WASM bundle in %s...\n", deployDir)
	absDir, _ := filepath.Abs(deployDir)
	resourceName := filepath.Base(absDir)
	buildDir, err := buildProdWASM(deployDir)
	if err != nil {
		return fmt.Errorf("Error building production WASM: %w", err)
	}
	fmt.Printf("Built production bundle at %s\n", buildDir)
	// Prepare metadata files in memory
	files := make(forcecli.ForceMetadataFiles)
	// WASM static resource
	wasmData, err := os.ReadFile(filepath.Join(buildDir, "bundle.wasm"))
	if err != nil {
		return err
	}
	// Add WASM bundle as a StaticResource with .resource extension
	files["staticresources/"+resourceName+".resource"] = wasmData
	staticResourceMetadata := `<?xml version="1.0" encoding="UTF-8"?>
<StaticResource xmlns="http://soap.sforce.com/2006/04/metadata">
	<cacheControl>Private</cacheControl>
	<contentType>application/wasm</contentType>
</StaticResource>`
	files["staticresources/"+resourceName+".resource-meta.xml"] = []byte(staticResourceMetadata)
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
	appComp := resourceName
	appClass := strings.Title(appComp)
	// JS
	js := fmt.Sprintf(`import Thunder from 'c/thunder';
import APP_URL from '@salesforce/resourceUrl/%s';

export default class %s extends Thunder {
	connectedCallback() {
		this.app = APP_URL;
	}
}`, appComp, appClass)
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
	pkg := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
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
    <members>go</members>
    <members>thunder</members>
    <members>%s</members>
    <name>LightningComponentBundle</name>
  </types>
  <version>58.0</version>
</Package>`, resourceName, resourceName)
	files["package.xml"] = []byte(pkg)
	// Perform deployment
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
	// Cleanup temporary build directory
	if rmErr := os.RemoveAll(buildDir); rmErr != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to remove temp dir %s: %v\n", buildDir, rmErr)
	}
	fmt.Printf("Deployment complete: %+v\n", result)
	// Open new tab in Salesforce if requested
	if deployTab {
		tabUrl := fmt.Sprintf("%s/lightning/n/%s", creds.InstanceUrl, resourceName)
		if err := desktop.Open(tabUrl); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open tab URL: %v\n", err)
		}
	}
	return nil
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
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
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
