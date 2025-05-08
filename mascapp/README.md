# MASC Go WASM App

This directory contains a minimal Go application skeleton that demonstrates calling the Salesforce REST proxy injected by the Vecty LWC.

Steps to build and deploy:
1. Initialize the Go module:
   ```sh
   cd mascapp
   go mod init mascapp
   ```
2. (Optionally) add the `masc` library:
   ```sh
   go get github.com/octoberswimmer/masc
   ```
3. Build the WASM binary:
   ```sh
   GOOS=js GOARCH=wasm go build -o ../main/default/staticresources/masc.wasm main.go
   ```
4. Create a static resource in `main/default/staticresources`:
   - Place `masc.wasm` in `main/default/staticresources/`
   - Add a `masc.resource-meta.xml` alongside (copy `hello.resource-meta.xml` and adjust `contentType` if needed).
5. In your LWC or Vecty Tester component, import the resource:
   ```js
   import MASC_APP from '@salesforce/resourceUrl/masc';
   ```
   and pass it to `<c-vecty app={MASC_APP}></c-vecty>`.

Now you can click the button in the Vecty Tester to trigger your Go/WASM code, which calls the `get` method proxied by your Apex `GoBridge`.