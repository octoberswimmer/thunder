# Thunder Go WASM App

This directory contains a minimal Go application skeleton that demonstrates calling the Salesforce REST proxy injected by the Thunder LWC.

## UI Demonstration

The demo is organized into four tabs: **Actions**, **Data**, **Object Info**, and **Layout** (grid demonstration).

This demo includes the following SLDS components:
- Badge: displays a simple label
- Pill: displays a pill with an optional remove button (clicking remove shows a toast)
- Icon: renders an SLDS icon (utility/action/standard)
- Datepicker: select a date to filter Accounts by LastModifiedDate
- Breadcrumbs: displays navigation path (Home > Demo)
- Tabs: group related content sections (Actions, Data, Object Info, Layout)
- Card: wrap content within SLDS cards
- Button: trigger actions like fetching data, fetching object info, showing modal and toast
- Modal: display overlay content
- Toast: show notifications
- Spinner: indicate loading state
- Select, Checkbox, RadioGroup, TextInput, Lookup: filter and select options for data
- ProgressBar: show the percentage of filtered results
- DataTable: display tabular data for Accounts
- Grid: demonstrate SLDS grid layout with cards

### Object Info Tab

The **Object Info** tab demonstrates the `GetObjectInfo` API functionality:
- Click "Get Account Info" to fetch Account object metadata
- Displays a spinner while loading
- Shows comprehensive object information including:
  - Basic info (API name, label, key prefix, custom status)
  - Object capabilities (createable, updateable, deletable, queryable, searchable)
  - Additional metadata (feed enabled, MRU enabled, theme info, field/relationship counts)

## Deploy

Deploy the app to your Salesforce org using the `thunder` CLI:

```
$ thunder deploy ./thunderDemo/
```

Manual steps to build and deploy:

1. Initialize the Go module:
   ```sh
   cd thunderDemo
   go mod init thunderDemo
   ```
2. Add the `masc` and `thunder` libraries:
   ```sh
   go get github.com/octoberswimmer/masc
   go get github.com/octoberswimmer/thunder
   ```
3. Build the WASM binary:
   ```sh
   GOOS=js GOARCH=wasm go build -o ../main/default/staticresources/thunderDemo.wasm main.go
   ```
4. Create a static resource in `main/default/staticresources`:
   - Place `thunderDemo.wasm` in `main/default/staticresources/`
   - Add a `thunderDemo.resource-meta.xml` alongside
5. In your LWC component, import the resource:

```js
import { LightningElement } from 'lwc';
// Load the compiled Go WASM Thunder app as a static resource
import THUNDER_APP from '@salesforce/resourceUrl/thunderDemo';

export default class ThunderDemo extends LightningElement {
        // Pass the Thunder WASM app URL into the thunder component
        thunderApp = THUNDER_APP;
}
```
   and pass it to `<c-thunder app={THUNDER_APP}></c-thunder>`.

Now you can click the button in the Thunder Demo to trigger your Go/WASM code, which calls the `get` method proxied by your Apex `GoBridge`.
