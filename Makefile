 # Makefile for building Go WASM MASC app into Salesforce static resource

 GOCMD = go
 GOWASMOS = js
 GOWASMARCH = wasm
 SRC_DIR = mascapp
 SRC_FILE = main.go
 OUTPUT = main/default/staticresources/masc.wasm

 .PHONY: all clean

 all: $(OUTPUT)

 $(OUTPUT): $(SRC_DIR)/$(SRC_FILE) $(SRC_DIR)/go.mod
	@echo "Building Go WASM binary -> $(OUTPUT)"
	cd $(SRC_DIR) && GOOS=$(GOWASMOS) GOARCH=$(GOWASMARCH) $(GOCMD) build -o ../$(OUTPUT) $(SRC_FILE)

 clean:
	@echo "Cleaning $(OUTPUT)"
	rm -f $(OUTPUT)