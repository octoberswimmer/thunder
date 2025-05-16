# Thunder CLI Checklist

This checklist tracks the implementation of the `thunder` CLI with `serve` and `deploy` subcommands.

- [x] Scaffold `cmd/thunder` directory with `main.go` and basic flag parsing
- [x] Migrate CLI to use Cobra with root command and subcommands

## Serve Subcommand
- [x] Implement `thunder serve` subcommand
  - [x] Validate that the serve directory contains a `main` package (using `go list`)
  - [x] Fetch Salesforce session (instance URL, access token) via Force CLI library
  - [x] Build WASM bundle in dev mode (`GOOS=js GOARCH=wasm`, `-tags dev`)
  - [x] Automatically open default web browser at `http://localhost:<port>`
  - [x] Start HTTP server:
    - Serve `bundle.wasm`, `wasm_exec.js`, `index.html`, and static assets
    - Proxy `/services/...` REST calls to Salesforce org
    - Watch source files and auto-rebuild on changes

## Deploy Subcommand
- [x] Implement `thunder deploy` subcommand
- [x] Embed metadata templates using `//go:embed` (Apex classes, LWC, static resource metadata)
  - [x] Build WASM bundle for production (omit `dev` tag)
  - [x] Generate metadata folder structure in-memory:
    - Static resource for WASM bundle
    - `GoBridge.cls` and `GoBridgeTest.cls`
    - LWC components (`go`, `thunder`, and user app)
  - [x] Deploy metadata to org via Force CLI library
  - [x] Add CustomTab metadata when `--tab` flag is set
  - [x] Open browser to `/lightning/n/<app>` after deploy with `--tab`

## Common Tasks
- [x] Provide CLI help and usage examples
- [x] Write tests for Cobra commands (`serve` and `deploy`)
- [x] Document new commands in `README.md`
