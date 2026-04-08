BINARY     := skills
MODULE     := github.com/simonski/skills
VERSION    ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS    := -ldflags "-X $(MODULE)/cmd.Version=$(VERSION)"
BUILD_DIR  := dist
GOFLAGS    := -trimpath

.PHONY: help
help:
	@echo ""
	@echo "  skills — agentic skills manager"
	@echo ""
	@echo "  Usage: make <target>"
	@echo ""
	@echo "  Targets:"
	@echo "    build      Build the $(BINARY) binary into ./$(BUILD_DIR)/"
	@echo "    install    Install the $(BINARY) binary to GOPATH/bin"
	@echo "    test       Run all tests"
	@echo "    lint       Run go vet and staticcheck"
	@echo "    clean      Remove build artefacts"
	@echo "    release    Build release binaries for all supported platforms"
	@echo "    help       Show this help (default)"
	@echo ""

.DEFAULT_GOAL := help

.PHONY: build
build:
	@mkdir -p $(BUILD_DIR)
	go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) .
	@echo "Built $(BUILD_DIR)/$(BINARY) (version=$(VERSION))"

.PHONY: install
install:
	go install $(GOFLAGS) $(LDFLAGS) .
	@echo "Installed $(BINARY) (version=$(VERSION))"

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	go vet ./...
	@command -v staticcheck >/dev/null 2>&1 && staticcheck ./... || echo "(staticcheck not installed, skipping)"

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

.PHONY: release
release: clean
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin  GOARCH=amd64  go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-darwin-amd64  .
	GOOS=darwin  GOARCH=arm64  go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-darwin-arm64  .
	GOOS=linux   GOARCH=amd64  go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-amd64   .
	GOOS=linux   GOARCH=arm64  go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-arm64   .
	GOOS=windows GOARCH=amd64  go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe .
	@echo "Release binaries written to $(BUILD_DIR)/"
	@ls -lh $(BUILD_DIR)/
