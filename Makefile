BINARY       := skills
MODULE       := github.com/simonski/skills
VERSION_FILE := VERSION
VERSION      := $(shell cat $(VERSION_FILE) 2>/dev/null | tr -d '[:space:]' || echo "0.1.0")
BUILD_DIR    := dist
GOFLAGS      := -trimpath
TAP_REPO     := https://github.com/simonski/homebrew-tap.git

.PHONY: help
help:
	@echo ""
	@echo "  skills — agentic skills manager"
	@echo ""
	@echo "  Usage: make <target>"
	@echo ""
	@echo "  Targets:"
	@echo "    build      Build the $(BINARY) binary into ./$(BUILD_DIR)/ (bumps patch version)"
	@echo "    install    Install the $(BINARY) binary to GOPATH/bin"
	@echo "    test       Run all tests"
	@echo "    lint       Run go vet and staticcheck"
	@echo "    clean      Remove build artefacts"
	@echo "    release    Build release archives for all supported platforms"
	@echo "    publish    Tag, release to GitHub, and update the Homebrew tap"
	@echo "    help       Show this help (default)"
	@echo ""

.DEFAULT_GOAL := help

# _bump_version increments the patch component of the VERSION file.
define _bump_version
	@CUR=$$(cat $(VERSION_FILE) | tr -d '[:space:]'); \
	MAJOR=$$(echo $$CUR | cut -d. -f1); \
	MINOR=$$(echo $$CUR | cut -d. -f2); \
	PATCH=$$(echo $$CUR | cut -d. -f3); \
	NEW="$$MAJOR.$$MINOR.$$((PATCH + 1))"; \
	echo "$$NEW" > $(VERSION_FILE); \
	echo "Version bumped: $$CUR -> $$NEW"
endef

.PHONY: build
build:
	$(call _bump_version)
	@mkdir -p $(BUILD_DIR)
	@NEW_VER=$$(cat $(VERSION_FILE) | tr -d '[:space:]'); \
	go build $(GOFLAGS) -ldflags "-X $(MODULE)/cmd.Version=$$NEW_VER" -o $(BUILD_DIR)/$(BINARY) .; \
	echo "Built $(BUILD_DIR)/$(BINARY) (version=$$NEW_VER)"

.PHONY: install
install:
	@NEW_VER=$$(cat $(VERSION_FILE) | tr -d '[:space:]'); \
	go install $(GOFLAGS) -ldflags "-X $(MODULE)/cmd.Version=$$NEW_VER" .; \
	echo "Installed $(BINARY) (version=$$NEW_VER)"

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

# _build_archive compiles and packages a single platform archive.
# Usage: $(call _build_archive,<GOOS>,<GOARCH>,<VERSION>)
define _build_archive
	@VER=$(3); \
	GOOS=$(1) GOARCH=$(2) go build $(GOFLAGS) -ldflags "-X $(MODULE)/cmd.Version=$$VER" \
		-o $(BUILD_DIR)/$(BINARY)-$(1)-$(2)$(if $(filter windows,$(1)),.exe,) .; \
	if [ "$(1)" = "windows" ]; then \
		zip -j $(BUILD_DIR)/$(BINARY)-$(1)-$(2).zip $(BUILD_DIR)/$(BINARY)-$(1)-$(2).exe && \
		rm $(BUILD_DIR)/$(BINARY)-$(1)-$(2).exe; \
	else \
		tar -czf $(BUILD_DIR)/$(BINARY)-$(1)-$(2).tar.gz -C $(BUILD_DIR) $(BINARY)-$(1)-$(2) && \
		rm $(BUILD_DIR)/$(BINARY)-$(1)-$(2); \
	fi
endef

.PHONY: release
release: clean
	@mkdir -p $(BUILD_DIR)
	@VER=$$(cat $(VERSION_FILE) | tr -d '[:space:]'); \
	echo "Building release archives for v$$VER..."; \
	GOOS=darwin  GOARCH=amd64 go build $(GOFLAGS) -ldflags "-X $(MODULE)/cmd.Version=$$VER" -o $(BUILD_DIR)/$(BINARY)-darwin-amd64  .; \
	tar -czf $(BUILD_DIR)/$(BINARY)-darwin-amd64.tar.gz -C $(BUILD_DIR) $(BINARY)-darwin-amd64 && rm $(BUILD_DIR)/$(BINARY)-darwin-amd64; \
	GOOS=darwin  GOARCH=arm64 go build $(GOFLAGS) -ldflags "-X $(MODULE)/cmd.Version=$$VER" -o $(BUILD_DIR)/$(BINARY)-darwin-arm64  .; \
	tar -czf $(BUILD_DIR)/$(BINARY)-darwin-arm64.tar.gz -C $(BUILD_DIR) $(BINARY)-darwin-arm64 && rm $(BUILD_DIR)/$(BINARY)-darwin-arm64; \
	GOOS=linux   GOARCH=amd64 go build $(GOFLAGS) -ldflags "-X $(MODULE)/cmd.Version=$$VER" -o $(BUILD_DIR)/$(BINARY)-linux-amd64   .; \
	tar -czf $(BUILD_DIR)/$(BINARY)-linux-amd64.tar.gz  -C $(BUILD_DIR) $(BINARY)-linux-amd64  && rm $(BUILD_DIR)/$(BINARY)-linux-amd64; \
	GOOS=linux   GOARCH=arm64 go build $(GOFLAGS) -ldflags "-X $(MODULE)/cmd.Version=$$VER" -o $(BUILD_DIR)/$(BINARY)-linux-arm64   .; \
	tar -czf $(BUILD_DIR)/$(BINARY)-linux-arm64.tar.gz  -C $(BUILD_DIR) $(BINARY)-linux-arm64  && rm $(BUILD_DIR)/$(BINARY)-linux-arm64; \
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -ldflags "-X $(MODULE)/cmd.Version=$$VER" -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe .; \
	zip -j $(BUILD_DIR)/$(BINARY)-windows-amd64.zip $(BUILD_DIR)/$(BINARY)-windows-amd64.exe && rm $(BUILD_DIR)/$(BINARY)-windows-amd64.exe; \
	echo "Release archives written to $(BUILD_DIR)/"; \
	ls -lh $(BUILD_DIR)/

# publish: tag the current version, create a GitHub release with archives,
# and update the simonski/homebrew-tap formula.
#
# Prerequisites:
#   - gh (GitHub CLI) authenticated
#   - git configured with push access to this repo and the tap repo
.PHONY: publish
publish: release
	@command -v gh >/dev/null 2>&1 || { echo "Error: gh CLI is required for publish"; exit 1; }
	@VER=$$(cat $(VERSION_FILE) | tr -d '[:space:]'); \
	TAG="v$$VER"; \
	echo "Publishing $$TAG..."; \
	\
	SHA_DARWIN_AMD64=$$(shasum -a 256 $(BUILD_DIR)/$(BINARY)-darwin-amd64.tar.gz  | awk '{print $$1}'); \
	SHA_DARWIN_ARM64=$$(shasum -a 256 $(BUILD_DIR)/$(BINARY)-darwin-arm64.tar.gz  | awk '{print $$1}'); \
	SHA_LINUX_AMD64=$$(shasum -a 256  $(BUILD_DIR)/$(BINARY)-linux-amd64.tar.gz   | awk '{print $$1}'); \
	SHA_LINUX_ARM64=$$(shasum -a 256  $(BUILD_DIR)/$(BINARY)-linux-arm64.tar.gz   | awk '{print $$1}'); \
	\
	git add $(VERSION_FILE); \
	git diff --cached --quiet || git commit -m "chore: release $$TAG"; \
	git tag -a "$$TAG" -m "Release $$TAG"; \
	git push origin HEAD; \
	git push origin "$$TAG"; \
	\
	gh release create "$$TAG" \
		$(BUILD_DIR)/$(BINARY)-darwin-amd64.tar.gz \
		$(BUILD_DIR)/$(BINARY)-darwin-arm64.tar.gz \
		$(BUILD_DIR)/$(BINARY)-linux-amd64.tar.gz \
		$(BUILD_DIR)/$(BINARY)-linux-arm64.tar.gz \
		$(BUILD_DIR)/$(BINARY)-windows-amd64.zip \
		--title "skills $$TAG" \
		--notes "Release $$TAG"; \
	\
	TAPDIR=$$(mktemp -d); \
	git clone $(TAP_REPO) "$$TAPDIR"; \
	mkdir -p "$$TAPDIR/Formula"; \
	sed \
		-e "s/SKILLS_VERSION/$$VER/g" \
		-e "s/SKILLS_SHA_DARWIN_AMD64/$$SHA_DARWIN_AMD64/g" \
		-e "s/SKILLS_SHA_DARWIN_ARM64/$$SHA_DARWIN_ARM64/g" \
		-e "s/SKILLS_SHA_LINUX_AMD64/$$SHA_LINUX_AMD64/g" \
		-e "s/SKILLS_SHA_LINUX_ARM64/$$SHA_LINUX_ARM64/g" \
		packaging/homebrew/skills.rb.template > "$$TAPDIR/Formula/skills.rb"; \
	cd "$$TAPDIR" && \
	git add Formula/skills.rb && \
	git commit -m "skills: update to $$TAG" && \
	git push origin HEAD; \
	rm -rf "$$TAPDIR"; \
	echo "Published $$TAG and updated homebrew tap."
