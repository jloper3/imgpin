# -----------------------------------------
# imgpin – Makefile (With Test Support)
# -----------------------------------------

APP := imgpin
BIN := bin/$(APP)
PKG := ./cmd/$(APP)

GOPATH ?= $(shell go env GOPATH)
export GO111MODULE=on

# -----------------------------------------
# Default
# -----------------------------------------

.PHONY: all
all: build

# -----------------------------------------
# Build / Run
# -----------------------------------------

.PHONY: build
build:
	@echo "👉 Building $(APP)..."
	go build -o $(BIN) $(PKG)
	@echo "✔ Build complete: $(BIN)"

.PHONY: run
run: build
	@echo "👉 Running $(APP)..."
	./$(BIN)

# -----------------------------------------
# Test Targets
# -----------------------------------------

TESTPKGS := $(shell go list ./...)

.PHONY: test
test:
	@echo "👉 Running tests..."
	go test -v $(TESTPKGS)
	@echo "✔ Tests complete"

.PHONY: test-race
test-race:
	@echo "👉 Running tests with race detector..."
	go test -race -v $(TESTPKGS)
	@echo "✔ Race-safe"

.PHONY: test-cover
test-cover:
	@echo "👉 Running coverage..."
	go test -coverprofile=coverage.out $(TESTPKGS)
	@echo "✔ Coverage report created: coverage.out"

.PHONY: cover-html
cover-html: test-cover
	@echo "👉 Generating HTML coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "✔ Open coverage.html in your browser"

# -----------------------------------------
# Lint / Format / Hygiene
# -----------------------------------------

.PHONY: fmt
fmt:
	@echo "👉 Running go fmt..."
	go fmt ./...

.PHONY: vet
vet:
	@echo "👉 Running go vet..."
	go vet ./...

.PHONY: tidy
tidy:
	@echo "👉 Tidying module..."
	go mod tidy -v

.PHONY: lint
lint: fmt vet tidy

# -----------------------------------------
# Install / Remove
# -----------------------------------------

.PHONY: install
install: build
	@echo "👉 Installing binary into $$GOPATH/bin"
	mkdir -p "$(GOPATH)/bin"
	cp $(BIN) "$(GOPATH)/bin/$(APP)"
	@echo "✔ Installed to $(GOPATH)/bin/$(APP)"

.PHONY: uninstall
uninstall:
	@echo "👉 Removing $(APP) from $$GOPATH/bin"
	rm -f "$(GOPATH)/bin/$(APP)"
	@echo "✔ Uninstalled"

# -----------------------------------------
# Release
# -----------------------------------------

.PHONY: release
release: build
	@echo "👉 Packaging release..."
	rm -rf dist
	mkdir -p dist
	cp $(BIN) dist/
	cd dist && tar -czf $(APP).tar.gz $(APP)
	@echo "✔ Release created: dist/$(APP).tar.gz"

# -----------------------------------------
# Clean
# -----------------------------------------

.PHONY: clean
clean:
	@echo "👉 Cleaning..."
	rm -rf bin/
	rm -rf dist/
	rm -f coverage.out coverage.html
	@echo "✔ Cleaned"

# -----------------------------------------
# Debug Helpers
# -----------------------------------------

.PHONY: debug-env
debug-env:
	@echo "GOPATH = $(GOPATH)"
	@echo "BIN    = $(BIN)"
	@echo "PKG    = $(PKG)"

