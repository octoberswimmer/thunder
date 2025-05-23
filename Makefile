 # Makefile for building Go WASM Thunder app into Salesforce static resource

 GOCMD = go
 GOWASMOS = js
 GOWASMARCH = wasm
 SRC_DIR = thunderDemo
 SRC_FILE = main.go
 OUTPUT = main/default/staticresources/thunderDemo.wasm

 .PHONY: all clean

 all: $(OUTPUT)

 $(OUTPUT): $(SRC_DIR)/$(SRC_FILE) $(SRC_DIR)/go.mod $(shell find components -type f)
	@echo "Building Go WASM binary -> $(OUTPUT)"
	cd $(SRC_DIR) && GOOS=$(GOWASMOS) GOARCH=$(GOWASMARCH) $(GOCMD) build -o ../$(OUTPUT) $(SRC_FILE)

 clean:
	@echo "Cleaning $(OUTPUT)"
	rm -f $(OUTPUT)
