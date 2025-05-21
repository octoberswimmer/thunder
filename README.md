# Thunder
## Masc (Go WASM Apps) + SLDS on Salesforce

This repository provides libraries for building applications for Salesforce's
Lightning Experience in Go using SLDS-styled [masc](https://github.com/octoberswimmer/masc) components, which follow the [Elm Architecture](https://guide.elm-lang.org/architecture/).

Thunder is made up of these masc components, a `thunder` LWC for running the
compiled wasm applications, a lightweight API designed to mirror the REST
API, and a CLI to make development a joy.

Repository Structure:
```
. (repo root)
├ cmd/thunder/           CLI source (Cobra commands: serve, deploy)
│  ├ main.go             CLI implementation
│  ├ main_test.go        CLI command tests
├ salesforce/            embedded metadata templates
│  ├ classes/            Apex classes
│  ├ staticresources/    StaticResource metadata
│  └ lwc/                LWC wrappers (`go`, `thunder`)
├ components/            MASC components for Thunder apps
├ api/                   REST proxy for WASM apps and Record API for convenient field access (StringValue, Value)
└ thunderDemo/           sample Go MASC application
```

Key parts:
- **thunder LWC** (`c:thunder`):
  - Loads a Go WASM app as static resource, injects global `get`/`post`/`put`/`delete` functions, and runs the app.
  - Exposes the `recordId` from Lightning record pages to Go WASM code via `globalThis.recordId`.
- **Thunder SLDS Components** (`components/`):
  - Go library offering SLDS-styled Masc components like `Button` and `DataTable`.
- **Apex GoBridge** (`salesforce/classes/GoBridge.cls`):
  - `@AuraEnabled callRest(...)` to map REST API calls to Apex.
- **thunder CLI** (`cmd/thunder`):
  -  Command-line tool to serve a Thunder app while in development and, and to
	  build and deploy it to a Salesforce org.
- **thunderDemo MASC app** (`thunderDemo/`):
  - Implements a Masc Model with a button and data table, using Thunder SLDS components.

Getting Started:
1. Install dependencies:
   - Go 1.24+ (with WASM support)
   - Force CLI
2. Run the thunderDemo app locally:
   ```sh
   $ force login
   $ thunder ./thunderDemo
   ```
   This compiles `thunderDemo/main.go` and starts a web server to serve the app.
3. Deploy to Salesforce using `thunder deploy -d ./thunderDemo --tab`
4. Click **Fetch Accounts** to see a data table rendered from your Thunder app.

## Thunder CLI
Thunder provides a CLI with two subcommands, `serve` and `deploy`, for local development and deployment of Go WASM apps on Salesforce.

### Installation
```sh
go install github.com/octoberswimmer/thunder/cmd/thunder@latest
```


### Usage
```sh
thunder serve [dir] --port PORT   # build & serve locally (defaults to current dir)
thunder deploy [dir] [--tab]      # deploy app to Salesforce org (defaults to current dir)
```

#### serve
 - `--port, -p`: Port to serve on (default `8000`)

`thunder serve`:
- Builds the app in dev mode (`GOOS=js GOARCH=wasm -tags dev`).
- Serves on `http://localhost:PORT`, auto-rebuilds on file changes.
- Proxies `/services/...` REST calls to your Salesforce org via CLI auth.
- Opens your default browser to the served app URL.

#### deploy
- `--tab, -t`: Also include a CustomTab in the deployment and open it for the app

`thunder deploy`:
- Builds a production WebAssembly bundle.
- Packages metadata (static resource, Apex classes, LWC wrappers, app LWC, and optional CustomTab) in-memory.
- Generates `package.xml` (includes CustomTab if requested) and deploys all metadata via your CLI session.
- With `--tab`, adds a CustomTab to the package, deploys it, and opens `/lightning/n/<app>` in your browser.

The CLI watches Go source files (`.go`, `go.mod`, `go.sum`) and automatically rebuilds the WASM bundle on changes. Refresh the browser to load the latest build.
API REST requests (via `/services/`) are automatically proxied through your active Salesforce CLI session. Be sure to run `force login` beforehand.

For details on implementing additional SLDS components, see `THUNDER_CHECKLIST.md`.
