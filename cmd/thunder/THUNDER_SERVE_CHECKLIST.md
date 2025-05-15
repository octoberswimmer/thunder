# Thunder Serve CLI Checklist

This checklist tracks the implementation of a `thunder serve` CLI tool to build and serve Thunder WASM apps locally.

- [ ] Scaffold CLI directory `cmd/thunder` and `main.go` with basic flag parsing (port, app dir)
- [ ] Accept `--port` (default 8000) and `--dir` (path to Thunder app, default `.`)
- [ ] Implement `buildWASM()`:
  - Use `GOOS=js GOARCH=wasm go build -o bundle.wasm` in a temp workspace
  - Apply source overlay to:
    - Replace `thunder.Run(...)` with a dev-mode mount (e.g., `masc.RenderIntoNode(...)`)
  - [x] Replace `api.Get` calls to proxy via the current Salesforce session (using Force CLI library)
- [ ] Generate `index.html` that loads:
  - `wasm_exec.js` from Go SDK
  - `bundle.wasm`
  - Injects/links SLDS CSS into the page (via static resource or CDN)
  - Initializes the WASM module and bootstraps the app
- [ ] Start HTTP server:
  - Serve `index.html`, `wasm_exec.js`, `bundle.wasm`, and static resources
  - Auto-reload on rebuild
- [ ] Watch source files (Go code) and rebuild/reload on changes
- [ ] Provide CLI help and usage examples
- [ ] Test against `thunderDemo` sample app
- [ ] Document `thunder serve` usage in `README.md`