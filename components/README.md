# Thunder SLDS Component Library

`github.com/octoberswimmer/thunder/components` provides a set of Salesforce Lightning Design System (SLDS)–styled UI components for Go WASM applications built with [Masc](https://github.com/octoberswimmer/masc) (which is a crossbreed between the Bubble Tea and Vecty rendering models).

## Features
- Button (neutral, brand, destructive)
- DataTable (basic tabular display)
- Card (container with header and body)
- Page Header (page-level heading with optional subtitle and actions)
- Modal (dialog overlay)
- Toast (notification overlay)
- TextInput (labeled text input)
- Select (dropdown)
- Checkbox (boolean input)
  
## Installation
Add Thunder as a dependency in your Go WASM module:
```sh
go get github.com/octoberswimmer/thunder@latest
```

## Usage
Import the components package in your Go code:
```go
import (
    "github.com/octoberswimmer/masc"
    "github.com/octoberswimmer/thunder/components"
)
```

### Button
Render an SLDS button:
```go
btn := components.Button(
    "Save Record",              // label
    components.VariantBrand,     // SLDS style variant
    func(e *masc.Event) {        // click handler
        // your logic here
    },
)
```

### DataTable
Display tabular data with headers and rows:
```go
headers := []string{"Name", "Industry"}
rows := []map[string]string{
    {"Name": "Acme Corp", "Industry": "Manufacturing"},
    {"Name": "Foo Ltd",  "Industry": "Technology"},
}
table := components.DataTable(headers, rows)
```

## Integration with Masc
In your Masc model, render components just like any other Masc component:

```go
import (
"github.com/octoberswimmer/masc"
"github.com/octoberswimmer/thunder"
"github.com/octoberswimmer/thunder/components"
)

type AppModel struct { masc.Core }

func (m *AppModel) Init() masc.Cmd { return nil }
func (m *AppModel) Update(msg masc.Msg) (masc.Model, masc.Cmd) { ... }
func (m *AppModel) Render(send func(masc.Msg)) masc.ComponentOrHTML {
    return components.Button("Fetch Data", components.VariantBrand, func(e *masc.Event) {
        send(FetchDataMsg{})
    })
}

func main() {
    thunder.Run(&AppModel{})
}
```

## Roadmap
See `THUNDER_COMPONENTS_CHECKLIST.md` for upcoming SLDS components:
- Card, Modal, Toast, Inputs, Tabs, Breadcrumbs, Progress Indicators, and more.
