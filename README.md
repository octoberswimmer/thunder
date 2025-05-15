# Thunder
## Masc (Go WASM Apps) + SLDS on Salesforce

This repository provides libraries for building applications for Salesforce's
Lightning Experience in Go using SLDS-styled [masc](https://github.com/octoberswimmer/masc) components, which follow the [Elm Architecture](https://guide.elm-lang.org/architecture/).

Thunder is made up of these masc components, a `thunder` LWC for running the
compiled wasm applications, and a lightweight API designed to mirror the REST
API.

Repository Structure:
```
.                   # root
├── components/      # Thunder SLDS components (Go masc components)
│   ├── button.go    # SLDS button
│   └── datatable.go # SLDS data table
├── main/            # Salesforce metadata (sfdx source format)
│   ├── default/
│   │   ├── classes/         # Apex classes (GoBridge + tests)
│   │   ├── lwc/             # Lightning Web Components (thunder, thunderDemo)
│   │   └── staticresources/ # WASM binaries (hello.wasm, masc.wasm)
├── thunderDemo/     # Go MASC application source (uses Thunder components)
│   ├── main.go       # Masc program rendering button and data table
│   └── go.mod        # Go module for thunderDemo
├── Makefile         # Builds thunderDemo/main.go to thunderDemo.wasm
└── README.md        # This file
```

Key parts:
- **thunder LWC** (`c:thunder`):
  - Loads a Go WASM app as static resource, injects global `get`/`post`/`put`/`delete` functions, and runs the app.
- **Thunder SLDS Components** (`components/`):
  - Go library offering SLDS-styled Masc components like `Button` and `DataTable`.
- **thunderDemo MASC app** (`thunderDemo/`):
  - Implements a Masc Model with a button and data table, using Thunder SLDS components.
- **thunderDemo LWC** (`c:thunder-demo`):
  - Passes `thunderDemo.wasm` into `thunder` and mounts it in Lightning pages.
- **Apex GoBridge** (`main/default/classes/GoBridge.cls`):
  - `@AuraEnabled callRest(...)` for REST and native SOQL calls.

Getting Started:
1. Install dependencies:
   - Go 1.24+ (with WASM support)
   - Force CLI or Salesforce CLI (sfdx)
2. Build the sample thunderDemo Go WASM app:
   ```sh
   make
   ```
   This compiles `thunderDemo/main.go` to `main/default/staticresources/thunderDemo.wasm`.
3. Deploy to Salesforce using `force` or `sfdx`.
5. Grant yourself access to **Thunder Demo** Tab.
5. Open **Thunder Demo** Tab in your org.
6. Click **Fetch Accounts** to see a data table rendered from your Go WASM app.

## Thunder Serve CLI

Thunder provides a local development CLI to build and serve Thunder apps with automatic rebuilds.

**Install the CLI:**
```sh
go install github.com/octoberswimmer/thunder/cmd/thunder@latest
```

**Usage:**
```sh
thunder --dir path/to/app --port 8000
```

**Flags:**
- `--dir`: Path to the Thunder app directory (default `.`)
- `--port`: Port to serve on (default `8000`)

The CLI watches Go source files (`.go`, `go.mod`, `go.sum`) and automatically rebuilds the WASM bundle on changes. Refresh the browser to load the latest build.
API REST requests (via `/services/`) are automatically proxied through your active Salesforce CLI session. Be sure to run `force login` beforehand.

For details on implementing additional SLDS components, see `THUNDER_CHECKLIST.md`.
