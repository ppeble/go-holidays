GO         ?= go
GOPATH     := $(shell $(GO) env GOPATH)
BIN        := bin
DEFS_REPO  := https://github.com/holidays/definitions.git
DEFS_TAG   ?= v8.1.0

.PHONY: build holidays gen-holidays vet staticcheck test parity generate update-definitions clean

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

# parity runs the Ruby<->Go comparison suite (build-tagged, excluded from `test`).
# Requires Ruby and the `holidays` gem installed, plus the `definitions/` submodule
# checked out (the oracle loads our v8.1.0 YAML into the gem via load_custom).
# See parity/README.md for the design and prerequisites.
parity:
	@ls definitions/*.yaml >/dev/null 2>&1 || { \
		echo "error: definitions/ submodule is empty; the oracle would load no v8.1.0 YAML and every region would mismatch." >&2; \
		echo "Run: git submodule update --init definitions" >&2; \
		exit 1; \
	}
	$(GO) test -tags parity -timeout=20m ./parity/...

generate: gen-holidays
	$(BIN)/gen-holidays -in definitions -out pkg/definitions

update-definitions:
	cd definitions && git fetch --tags && git checkout $(DEFS_TAG)
	@echo "Submodule now at $(DEFS_TAG). Don't forget: git add definitions && git commit"

clean:
	rm -rf $(BIN)
