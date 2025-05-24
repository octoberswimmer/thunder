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
│  └ lwc/                LWC wrappers (`go`, `thunder`)
├ components/            MASC components for Thunder apps
├ api/                   REST proxy for WASM apps, UI API metadata (GetObjectInfo, GetPicklistValuesByRecordType), and Record API for convenient field access (StringValue, Value)
└ thunderDemo/           sample Go MASC application
```

Key parts:
- **thunder LWC** (`c:thunder`):
  - Loads a Go WASM app as static resource, injects global API functions, and runs the app.
  - Exposes the `recordId` from Lightning record pages to Go WASM code via `globalThis.recordId`.
- **Thunder SLDS Components** (`components/`):
  - Go library offering SLDS-styled Masc components like `Button`, `DataTable`, `Grid`, `Textarea`, and `Stencil`.
- **Apex GoBridge** (`salesforce/classes/GoBridge.cls`):
  - `@AuraEnabled callRest(...)` to map REST API calls to Apex.
- **thunder CLI** (`cmd/thunder`):
  -  Command-line tool to serve a Thunder app while in development and, and to
	  build and deploy it to a Salesforce org.
**thunderDemo MASC app** (`thunderDemo/`):
  - Demonstrates SLDS-styled Masc components in a Go WASM app, organized into Actions, Data, and Layout tabs for interaction, data display, and grid layout examples.

Getting Started:
1. Install dependencies:
   - Go 1.24+ (with WASM support)
	- [Force CLI](https://github.com/forcecli/force)
2. Run the thunderDemo app locally:
   ```sh
   $ force login
   $ thunder serve ./thunderDemo
   ```
   This compiles `thunderDemo/main.go` and starts a web server to serve the app.
3. Deploy to Salesforce using `thunder deploy ./thunderDemo --tab`
4. Click **Fetch Accounts** to see a data table rendered from your Thunder app.

## Thunder Components

Thunder provides a comprehensive set of SLDS-styled components for building Lightning Experience applications:

### Form Components
- **`TextInput`**: Single-line text input with label and validation styling
- **`Textarea`**: Multi-line text input for longer content (e.g., addresses, descriptions)
- **`Select`**: Dropdown selection with picklist options
- **`Datepicker`**: Date input with SLDS calendar styling
- **`Checkbox`**: Boolean input with proper labeling
- **`RadioGroup`**: Multiple choice selection with radio buttons

### Layout Components  
- **`Grid`** & **`GridColumn`**: Responsive grid system for two-column layouts
  ```go
  t.Grid(
      t.GridColumn("1-of-2", /* first column content */),
      t.GridColumn("1-of-2", /* second column content */),
      t.GridColumn("1-of-1", /* full-width content */),
  )
  ```
- **`Card`**: Content containers with headers and proper spacing
- **`Page`** & **`PageHeader`**: Page-level layout with consistent styling
- **`Modal`**: Dialog overlays for secondary workflows

### Data Components
- **`DataTable`**: Feature-rich data tables with sorting and actions
- **`Lookup`**: Search and selection for related records

### UI Components
- **`Button`**: Action buttons with variant styling (Neutral, Brand, Destructive)
- **`Badge`**: Status indicators and labels
- **`Breadcrumb`**: Navigation hierarchy display
- **`Icon`**: SLDS icon integration
- **`ProgressBar`**: Progress indication for long-running operations
- **`Spinner`**: Loading indicators in multiple sizes
- **`Stencil`**: Skeleton placeholders for loading states
- **`Tabs`**: Tabbed content organization
- **`Toast`**: Notification messages

### Component Features
- **Consistent Spacing**: All components include proper SLDS margin classes
- **Responsive Design**: Grid system adapts to different screen sizes
- **Accessibility**: Full SLDS accessibility compliance
- **Event Handling**: Clean event binding with Go functions
- **Type Safety**: Strongly typed APIs for reliable development

### Example: Building a Form with Grid Layout

```go
func (m *AppModel) renderPatientForm(send func(masc.Msg)) masc.ComponentOrHTML {
    return t.Page(
        t.PageHeader("Patient Information", ""),
        t.Card("Patient Details",
            t.Grid(
                t.GridColumn("1-of-2",
                    t.TextInput("First Name", m.firstName, "", func(e *masc.Event) {
                        send(firstNameMsg(e.Target.Get("value").String()))
                    }),
                ),
                t.GridColumn("1-of-2",
                    t.TextInput("Last Name", m.lastName, "", func(e *masc.Event) {
                        send(lastNameMsg(e.Target.Get("value").String()))
                    }),
                ),
                t.GridColumn("1-of-1",
                    t.Textarea("Address", m.address, "Enter full address", 2, func(e *masc.Event) {
                        send(addressMsg(e.Target.Get("value").String()))
                    }),
                ),
                t.GridColumn("1-of-2",
                    // Show skeleton while loading picklist options
                    func() masc.ComponentOrHTML {
                        if len(m.stateOptions) == 0 {
                            return t.Stencil("State")
                        }
                        return t.Select("State", m.stateOptions, m.state, func(e *masc.Event) {
                            send(stateMsg(e.Target.Get("value").String()))
                        })
                    }(),
                ),
            ),
        ),
        t.Button("Save", t.VariantBrand, func(e *masc.Event) {
            send(saveMsg{})
        }),
    )
}
```

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

The CLI watches Go source files (`.go`, `go.mod`, `go.sum`) and automatically rebuilds the WASM bundle on changes. Refresh the browser to load the latest build.
API REST requests (via `/services/`) are automatically proxied through your active Salesforce CLI session. Be sure to run `force login` beforehand.

#### deploy
- `--tab, -t`: Also include a CustomTab in the deployment and open it for the app

`thunder deploy`:
- Builds a production WebAssembly bundle.
- Packages metadata (static resource, Apex classes, LWC wrappers, app LWC, and optional CustomTab) in-memory.
- Generates `package.xml` (includes CustomTab if requested) and deploys all metadata via your CLI session.
- With `--tab`, adds a CustomTab to the package, deploys it, and opens `/lightning/n/<app>` in your browser.
