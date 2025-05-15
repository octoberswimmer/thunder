package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	forcecli "github.com/ForceCLI/force/lib"
	"github.com/fsnotify/fsnotify"
)

// currentBuildDir holds the latest build output directory; buildMutex protects access.
var (
	currentBuildDir string
	buildMutex      sync.RWMutex
	// instanceURL and accessToken for Salesforce REST proxy
	instanceURL string
	accessToken string
)

// Usage: thunder --dir path/to/app --port 8000
func main() {
	// CLI flags
	port := flag.Int("port", 8000, "Port to serve on")
	dir := flag.String("dir", ".", "Path to Thunder app directory")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: thunder --dir PATH --port PORT\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Validate app directory
	info, err := os.Stat(*dir)
	if err != nil || !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Invalid app directory: %s\n", *dir)
		os.Exit(1)
	}

	// Fetch Salesforce auth info for proxying API requests
	instanceURL, accessToken, err = fetchAuthInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching Salesforce auth info: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Building WASM bundle in %s...\n", *dir)
	buildDir, err := buildWASM(*dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building WASM: %v\n", err)
		os.Exit(1)
	}
	// prepare for serving and start watching for changes
	buildMutex.Lock()
	currentBuildDir = buildDir
	buildMutex.Unlock()
	go watchAndRebuild(*dir)

	fmt.Printf("Serving Thunder app on port %d (watching %s)...\n", *port, *dir)
	// serve files from the latest build directory, rebuilding on changes
	// Proxy Salesforce REST API requests
	http.HandleFunc("/services/", proxyHandler)
	// Serve files from the latest build directory
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buildMutex.RLock()
		dirPath := currentBuildDir
		buildMutex.RUnlock()
		fs := http.FileServer(http.Dir(dirPath))
		fs.ServeHTTP(w, r)
	})
	address := fmt.Sprintf(":%d", *port)
	if err := http.ListenAndServe(address, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
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
	// generate index.html to load WASM module
	indexHTML := `<!DOCTYPE html>
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
	if err := os.WriteFile(filepath.Join(buildDir, "index.html"), []byte(indexHTML), 0644); err != nil {
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

	// watch directories recursively
	err = filepath.Walk(appDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking project for file watching: %v\n", err)
		return
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
