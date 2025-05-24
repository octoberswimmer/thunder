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
└ examples/              example Thunder applications
   ├ thunderDemo/        main demo app showcasing all components
   └ validation/         comprehensive form validation example
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
**Example Applications** (`examples/`):
  - **thunderDemo**: Demonstrates all Thunder components in a Go WASM app, organized into Actions, Data, ObjectInfo, and Layout tabs. Showcases component-only development without direct SLDS classes or elem usage.
  - **validation**: Comprehensive form validation example with real-time error handling, demonstrating ValidationState and validated components.

Getting Started:
1. Install dependencies:
   - Go 1.24+ (with WASM support)
	- [Force CLI](https://github.com/forcecli/force)
2. Run the thunderDemo app locally:
   ```sh
   $ force login
   $ thunder serve ./examples/thunderDemo
   ```
   This compiles `examples/thunderDemo/main.go` and starts a web server to serve the app.
3. Deploy to Salesforce using `thunder deploy ./examples/thunderDemo --tab`
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
- **`ValidatedTextInput`**, **`ValidatedTextarea`**, etc.: Form components with built-in validation state management

### Layout Components  
- **`Grid`** & **`GridColumn`**: Responsive grid system for flexible layouts
  ```go
  components.Grid(
      components.GridColumn("1-of-2", /* first column content */),
      components.GridColumn("1-of-2", /* second column content */),
      components.GridColumn("1-of-1", /* full-width content */),
  )
  ```
- **`CenteredGrid`**: Grid with center alignment for loading states and centered content
- **`Card`**: Content containers with headers and proper spacing
- **`Page`** & **`PageHeader`**: Page-level layout with consistent styling
- **`Modal`**: Dialog overlays for secondary workflows
- **`Container`**: Basic layout wrapper to avoid direct element usage
- **`Spacer`**: Flexible spacing container with margin/padding options
- **`MarginTop`**, **`MarginBottom`**, **`PaddingHorizontal`**, etc.: Semantic spacing components

### Data Components
- **`DataTable`**: Feature-rich data tables with sorting and actions
- **`Lookup`**: Search and selection for related records

### UI Components
- **`Button`**: Action buttons with variant styling (Neutral, Brand, Destructive)
- **`LoadingButton`**: Button with built-in spinner and disabled state
- **`Badge`**: Status indicators and labels
- **`Breadcrumb`**: Navigation hierarchy display
- **`Icon`**: SLDS icon integration
- **`ProgressBar`**: Progress indication for long-running operations
- **`Spinner`**: Loading indicators in multiple sizes
- **`LoadingSpinner`**: Centered loading spinner for containers
- **`Stencil`**: Skeleton placeholders for loading states
- **`Tabs`**: Tabbed content organization
- **`Toast`**: Notification messages

### Text Components
- **`Text`**: Styled text with size variants (Small, Regular, Large)
- **`Paragraph`**: Paragraph elements with text styling
- **`Heading`**: Heading elements (H1/H2/H3) with semantic sizing (Small, Medium, Large)

### Component Features
- **Complete Abstraction**: No direct SLDS classes or masc elements required in application code
- **Semantic APIs**: Type-safe spacing, sizing, and styling options
- **Consistent Spacing**: Semantic spacing components (Spacer, MarginTop, etc.)
- **Responsive Design**: Grid system adapts to different screen sizes
- **Accessibility**: Full SLDS accessibility compliance
- **Event Handling**: Clean event binding with Go functions
- **Type Safety**: Strongly typed APIs for reliable development
- **Loading States**: Built-in support for loading spinners and disabled states

### Example: Component-Only Development

Thunder components provide complete abstraction from SLDS classes and DOM elements:

```go
// Instead of using elem.Div with SLDS classes
elem.Div(
    masc.Markup(masc.Class("slds-m-top_medium", "slds-align_absolute-center")),
    components.Spinner("medium"),
)

// Use semantic layout components
components.MarginTop(components.SpaceMedium,
    components.LoadingSpinner("medium"),
)

// Complex layouts with semantic spacing
func (m *AppModel) renderPatientForm(send func(masc.Msg)) masc.ComponentOrHTML {
    return components.Page(
        components.PageHeader("Patient Information", "Enter patient details"),
        components.Card("Patient Details",
            components.Grid(
                components.GridColumn("1-of-2",
                    components.ValidatedTextInput("First Name", m.firstName, "",
                        components.ValidationState{
                            Required: true,
                            HasError: m.hasError("firstName"),
                            ErrorMessage: m.errors["firstName"],
                        },
                        func(e *masc.Event) {
                            send(firstNameMsg(e.Target.Get("value").String()))
                        },
                    ),
                ),
                components.GridColumn("1-of-2",
                    components.ValidatedTextInput("Last Name", m.lastName, "",
                        components.ValidationState{Required: true},
                        func(e *masc.Event) {
                            send(lastNameMsg(e.Target.Get("value").String()))
                        },
                    ),
                ),
            ),
            // Loading button with built-in spinner - conditional rendering with masc.If
            masc.If(m.isSubmitting,
                components.LoadingButton("Saving...", components.VariantBrand),
            ),
            masc.If(!m.isSubmitting,
                components.Button("Save", components.VariantBrand, func(e *masc.Event) {
                    send(saveMsg{})
                }),
            ),
        ),
    )
}
```

## Form Validation

Thunder provides validated form components that handle error states, required field validation, and help text. Each validated component includes:

- **Error State Management**: Red styling for validation errors
- **Required Field Indicators**: Asterisk (*) for required fields  
- **Help Text**: Descriptive text below form fields
- **Real-time Validation**: Immediate feedback on user input

### Validated Components
- **`ValidatedTextInput`**: Text input with validation state
- **`ValidatedTextarea`**: Multi-line text with validation
- **`ValidatedSelect`**: Dropdown selection with validation
- **`ValidatedDatepicker`**: Date input with validation

### ValidationState
All validated components use the `ValidationState` struct:
```go
type ValidationState struct {
    HasError     bool   // Show error styling
    Required     bool   // Show asterisk indicator
    ErrorMessage string // Error text to display
    HelpText     string // Help text below field
}
```

### Example: Validated Form
```go
validationState := components.ValidationState{
    HasError:     len(m.email) > 0 && !isValidEmail(m.email),
    Required:     true,
    ErrorMessage: "Please enter a valid email address",
    HelpText:     "We'll use this to send important updates",
}

components.ValidatedTextInput("Email", m.email, "user@example.com", validationState, func(e *masc.Event) {
    send(emailChangedMsg(e.Target.Get("value").String()))
})
```

## Examples

The `examples/` directory contains complete Thunder applications demonstrating different patterns:

### thunderDemo
The main demonstration app showcasing all Thunder components across four tabs:
- **Actions**: Buttons, badges, icons, and date pickers
- **Data**: Interactive data table with filtering, pagination, and controls
- **ObjectInfo**: Salesforce metadata display using the UI API
- **Layout**: Grid system demonstration

Features component-only development with no direct SLDS classes or masc elements.

```sh
thunder serve ./examples/thunderDemo
```

### validation
Comprehensive form validation example demonstrating:
- Real-time field validation with `ValidationState`
- Required field indicators and error messages
- Loading states with `LoadingButton`
- Semantic spacing with `MarginTop` and `Spacer`

Shows how to build robust forms using only Thunder components.

```sh
thunder serve ./examples/validation
```

Both examples are complete Go modules that can be run independently with `thunder serve` or deployed with `thunder deploy`.

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
