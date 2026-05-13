GO         ?= go
GOPATH     := $(shell $(GO) env GOPATH)
BIN        := bin
DEFS_REPO  := https://github.com/holidays/definitions.git
DEFS_TAG   ?= v6.1.1

.PHONY: build holidays gen-holidays vet staticcheck test generate update-definitions clean

build: holidays gen-holidays

holidays:
	$(GO) build -o $(BIN)/holidays ./cmd/holidays

gen-holidays:
	$(GO) build -o $(BIN)/gen-holidays ./cmd/gen-holidays

vet:
	$(GO) vet ./...

staticcheck:
	@command -v staticcheck >/dev/null 2>&1 || $(GO) install honnef.co/go/tools/cmd/staticcheck@latest
	$(GOPATH)/bin/staticcheck ./...

test: vet
	$(GO) test ./...

generate: gen-holidays
	$(BIN)/gen-holidays -in definitions -out pkg/definitions

update-definitions:
	cd definitions && git fetch --tags && git checkout $(DEFS_TAG)
	@echo "Submodule now at $(DEFS_TAG). Don't forget: git add definitions && git commit"

clean:
	rm -rf $(BIN)
