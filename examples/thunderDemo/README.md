# Thunder Go WASM App

This directory contains a minimal Go application skeleton that demonstrates
using Thunder components and calling the Salesforce API.

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
   go get
   ```
2. Build the WASM binary:
   ```sh
   GOOS=js GOARCH=wasm go build -o bundle.wasm main.go
   ```
3. Create a zip archive for the static resource:
   ```sh
   zip thunderDemo.zip bundle.wasm
   ```
   and upload to your org.
4. Deploy the apex classes and LWC components in `../salesforce/`
5. Create an LWC component, import the resource and reference the WASM file within the zip:
```js
import Thunder from 'c/thunder';
import APP_URL from '@salesforce/resourceUrl/thunderDemo';

export default class ThunderDemo extends Thunder {
	connectedCallback() {
		this.app = APP_URL;
	}
}
```

No template is needed.  The template from the Thunder module extended by your
app is used.

Create a tab for the LWC or add it to a lightning page.
