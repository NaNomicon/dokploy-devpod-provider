# Dokploy DevPod Provider Makefile
# Best practices for developing and managing a custom DevPod provider

# Variables
PROVIDER_NAME := dokploy
PROVIDER_FILE := provider.yaml
PROVIDER_DEV_FILE := provider-dev.yaml
VERSION := $(shell grep '^version:' $(PROVIDER_FILE) | sed 's/version: *v*//')
GITHUB_REPO := NaNomicon/dokploy-devpod-provider
TEST_WORKSPACE := test-workspace-$(shell date +%s)
TEST_REPO := https://github.com/microsoft/vscode-remote-try-node.git

# Load .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

# DevPod command detection
DEVPOD_CMD := $(shell \
	if command -v devpod >/dev/null 2>&1; then \
		echo "devpod"; \
	elif [ -f "/Applications/DevPod.app/Contents/MacOS/devpod" ]; then \
		echo "/Applications/DevPod.app/Contents/MacOS/devpod"; \
	elif [ -f "/Applications/DevPod.app/Contents/Resources/devpod" ]; then \
		echo "/Applications/DevPod.app/Contents/Resources/devpod"; \
	else \
		echo "devpod-not-found"; \
	fi)

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

# Default target
.DEFAULT_GOAL := help

##@ Build

# Go build variables
BINARY_NAME := dokploy-provider
BUILD_DIR := dist
LDFLAGS := -ldflags="-s -w"
VERSIONED_BINARY_NAME := $(BINARY_NAME)-$(VERSION)

.PHONY: build
build: ## Build binary for current platform
	@echo "$(BLUE)Building $(BINARY_NAME) v$(VERSION) for current platform...$(NC)"
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "$(GREEN)✓ Binary built: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

.PHONY: build-versioned
build-versioned: ## Build versioned binary for current platform
	@echo "$(BLUE)Building $(VERSIONED_BINARY_NAME) for current platform...$(NC)"
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(VERSIONED_BINARY_NAME) .
	@echo "$(GREEN)✓ Versioned binary built: $(BUILD_DIR)/$(VERSIONED_BINARY_NAME)$(NC)"

.PHONY: build-all
build-all: ## Build binaries for all supported platforms
	@echo "$(BLUE)Building $(BINARY_NAME) v$(VERSION) for all platforms...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@echo "$(YELLOW)Building for Linux AMD64...$(NC)"
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	@echo "$(YELLOW)Building for Linux ARM64...$(NC)"
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	@echo "$(YELLOW)Building for macOS AMD64...$(NC)"
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@echo "$(YELLOW)Building for macOS ARM64...$(NC)"
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@echo "$(YELLOW)Building for Windows AMD64...$(NC)"
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "$(GREEN)✓ All binaries built in $(BUILD_DIR)/$(NC)"

.PHONY: build-all-versioned
build-all-versioned: ## Build versioned binaries for all supported platforms
	@echo "$(BLUE)Building $(BINARY_NAME) v$(VERSION) for all platforms with version in filename...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@echo "$(YELLOW)Building for Linux AMD64...$(NC)"
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-amd64 .
	@echo "$(YELLOW)Building for Linux ARM64...$(NC)"
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-arm64 .
	@echo "$(YELLOW)Building for macOS AMD64...$(NC)"
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-darwin-amd64 .
	@echo "$(YELLOW)Building for macOS ARM64...$(NC)"
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-darwin-arm64 .
	@echo "$(YELLOW)Building for Windows AMD64...$(NC)"
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe .
	@echo "$(GREEN)✓ All versioned binaries built in $(BUILD_DIR)/$(NC)"

.PHONY: checksums
checksums: build-all ## Generate SHA256 checksums for all binaries
	@echo "$(BLUE)Generating SHA256 checksums...$(NC)"
	@cd $(BUILD_DIR) && \
	for file in $(BINARY_NAME)-*; do \
		if [ -f "$$file" ] && [[ "$$file" != *.sha256* ]]; then \
			echo "$(YELLOW)Generating checksum for $$file...$(NC)"; \
			if command -v sha256sum >/dev/null 2>&1; then \
				sha256sum "$$file" > "$$file.sha256"; \
			elif command -v shasum >/dev/null 2>&1; then \
				shasum -a 256 "$$file" > "$$file.sha256"; \
			else \
				echo "$(RED)Error: No SHA256 tool available$(NC)"; \
				exit 1; \
			fi; \
		fi; \
	done
	@echo "$(GREEN)✓ Checksums generated in $(BUILD_DIR)/$(NC)"

.PHONY: checksums-versioned
checksums-versioned: build-all-versioned ## Generate SHA256 checksums for all versioned binaries
	@echo "$(BLUE)Generating SHA256 checksums for versioned binaries...$(NC)"
	@cd $(BUILD_DIR) && \
	for file in $(BINARY_NAME)-$(VERSION)-*; do \
		if [ -f "$$file" ]; then \
			echo "$(YELLOW)Generating checksum for $$file...$(NC)"; \
			if command -v sha256sum >/dev/null 2>&1; then \
				sha256sum "$$file" > "$$file.sha256"; \
			elif command -v shasum >/dev/null 2>&1; then \
				shasum -a 256 "$$file" > "$$file.sha256"; \
			else \
				echo "$(RED)Error: No SHA256 tool available$(NC)"; \
				exit 1; \
			fi; \
		fi; \
	done
	@echo "$(GREEN)✓ Checksums generated for versioned binaries in $(BUILD_DIR)/$(NC)"

.PHONY: update-provider-checksums
update-provider-checksums: checksums ## Update provider.yaml with actual checksums
	@echo "$(BLUE)Updating provider.yaml with generated checksums...$(NC)"
	@if [ ! -f $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64.sha256 ]; then \
		echo "$(RED)Error: Checksums not found. Run 'make checksums' first$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Backing up provider.yaml...$(NC)"
	@cp $(PROVIDER_FILE) $(PROVIDER_FILE).backup
	@echo "$(YELLOW)Updating checksums in provider.yaml using line-by-line approach...$(NC)"
	@# Process each binary checksum
	@temp_file=$(BUILD_DIR)/provider_temp.yaml; \
	cp $(PROVIDER_FILE) $$temp_file; \
	echo "$(BLUE)Processing checksums...$(NC)"; \
	if [ -f "$(BUILD_DIR)/$(BINARY_NAME)-linux-amd64.sha256" ]; then \
		checksum=$$(cut -d' ' -f1 "$(BUILD_DIR)/$(BINARY_NAME)-linux-amd64.sha256"); \
		echo "$(BLUE)  Updating linux-amd64: $$checksum$(NC)"; \
		if [[ "$$OSTYPE" == "darwin"* ]]; then \
			sed -i '' '/dokploy-provider-linux-amd64$$/{ n; s/checksum: "[^"]*"/checksum: "'$$checksum'"/; }' $$temp_file; \
		else \
			sed -i '/dokploy-provider-linux-amd64$$/ { n; s/checksum: "[^"]*"/checksum: "'$$checksum'"/; }' $$temp_file; \
		fi; \
	fi; \
	if [ -f "$(BUILD_DIR)/$(BINARY_NAME)-linux-arm64.sha256" ]; then \
		checksum=$$(cut -d' ' -f1 "$(BUILD_DIR)/$(BINARY_NAME)-linux-arm64.sha256"); \
		echo "$(BLUE)  Updating linux-arm64: $$checksum$(NC)"; \
		if [[ "$$OSTYPE" == "darwin"* ]]; then \
			sed -i '' '/dokploy-provider-linux-arm64$$/{ n; s/checksum: "[^"]*"/checksum: "'$$checksum'"/; }' $$temp_file; \
		else \
			sed -i '/dokploy-provider-linux-arm64$$/ { n; s/checksum: "[^"]*"/checksum: "'$$checksum'"/; }' $$temp_file; \
		fi; \
	fi; \
	if [ -f "$(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64.sha256" ]; then \
		checksum=$$(cut -d' ' -f1 "$(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64.sha256"); \
		echo "$(BLUE)  Updating darwin-amd64: $$checksum$(NC)"; \
		if [[ "$$OSTYPE" == "darwin"* ]]; then \
			sed -i '' '/dokploy-provider-darwin-amd64$$/{ n; s/checksum: "[^"]*"/checksum: "'$$checksum'"/; }' $$temp_file; \
		else \
			sed -i '/dokploy-provider-darwin-amd64$$/ { n; s/checksum: "[^"]*"/checksum: "'$$checksum'"/; }' $$temp_file; \
		fi; \
	fi; \
	if [ -f "$(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64.sha256" ]; then \
		checksum=$$(cut -d' ' -f1 "$(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64.sha256"); \
		echo "$(BLUE)  Updating darwin-arm64: $$checksum$(NC)"; \
		if [[ "$$OSTYPE" == "darwin"* ]]; then \
			sed -i '' '/dokploy-provider-darwin-arm64$$/{ n; s/checksum: "[^"]*"/checksum: "'$$checksum'"/; }' $$temp_file; \
		else \
			sed -i '/dokploy-provider-darwin-arm64$$/ { n; s/checksum: "[^"]*"/checksum: "'$$checksum'"/; }' $$temp_file; \
		fi; \
	fi; \
	if [ -f "$(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe.sha256" ]; then \
		checksum=$$(cut -d' ' -f1 "$(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe.sha256"); \
		echo "$(BLUE)  Updating windows-amd64: $$checksum$(NC)"; \
		if [[ "$$OSTYPE" == "darwin"* ]]; then \
			sed -i '' '/dokploy-provider-windows-amd64\.exe$$/{ n; s/checksum: "[^"]*"/checksum: "'$$checksum'"/; }' $$temp_file; \
		else \
			sed -i '/dokploy-provider-windows-amd64\.exe$$/ { n; s/checksum: "[^"]*"/checksum: "'$$checksum'"/; }' $$temp_file; \
		fi; \
	fi; \
	mv $$temp_file $(PROVIDER_FILE)
	@echo "$(GREEN)✓ provider.yaml updated with checksums$(NC)"
	@echo "$(YELLOW)Backup saved as $(PROVIDER_FILE).backup$(NC)"

.PHONY: show-checksums
show-checksums: ## Display generated checksums with uniqueness validation
	@echo "$(BLUE)Generated checksums for each platform/architecture:$(NC)"
	@if [ -d $(BUILD_DIR) ]; then \
		cd $(BUILD_DIR) && \
		echo "$(YELLOW)Platform-specific checksums (each should be unique):$(NC)"; \
		echo ""; \
		for file in *.sha256; do \
			if [ -f "$$file" ]; then \
				platform=$$(echo "$$file" | sed 's/$(BINARY_NAME)-//; s/\.sha256//'); \
				checksum=$$(cut -d' ' -f1 "$$file"); \
				printf "$(GREEN)%-20s$(NC) %s\n" "$$platform:" "$$checksum"; \
			fi; \
		done; \
		echo ""; \
		echo "$(BLUE)Uniqueness validation:$(NC)"; \
		unique_count=$$(for file in *.sha256; do [ -f "$$file" ] && cut -d' ' -f1 "$$file"; done | sort -u | wc -l | tr -d ' '); \
		total_count=$$(for file in *.sha256; do [ -f "$$file" ] && echo "1"; done | wc -l | tr -d ' '); \
		if [ "$$unique_count" = "$$total_count" ]; then \
			echo "$(GREEN)✓ All $$total_count checksums are unique (correct)$(NC)"; \
		else \
			echo "$(RED)✗ Only $$unique_count unique checksums out of $$total_count total (ERROR!)$(NC)"; \
			echo "$(YELLOW)Duplicate checksums detected - this indicates a build problem$(NC)"; \
		fi; \
	else \
		echo "$(RED)No checksums found. Run 'make checksums' first$(NC)"; \
	fi

.PHONY: verify-checksums
verify-checksums: ## Verify existing checksums against binaries
	@echo "$(BLUE)Verifying checksums...$(NC)"
	@cd $(BUILD_DIR) && \
	failed=0; \
	for file in *.sha256; do \
		if [ -f "$$file" ] && [[ "$$file" != *.sha256.sha256* ]]; then \
			echo "$(YELLOW)Verifying $$file...$(NC)"; \
			if command -v sha256sum >/dev/null 2>&1; then \
				if sha256sum -c "$$file"; then \
					echo "$(GREEN)✓ $$file verified$(NC)"; \
				else \
					echo "$(RED)✗ $$file verification failed$(NC)"; \
					failed=1; \
				fi; \
			elif command -v shasum >/dev/null 2>&1; then \
				if shasum -a 256 -c "$$file"; then \
					echo "$(GREEN)✓ $$file verified$(NC)"; \
				else \
					echo "$(RED)✗ $$file verification failed$(NC)"; \
					failed=1; \
				fi; \
			fi; \
		fi; \
	done; \
	if [ $$failed -eq 0 ]; then \
		echo "$(GREEN)✓ All checksums verified$(NC)"; \
	else \
		echo "$(RED)Some checksums failed verification$(NC)"; \
		exit 1; \
	fi

.PHONY: validate-provider-checksums
validate-provider-checksums: checksums ## Validate that provider.yaml has correct and unique checksums
	@echo "$(BLUE)Validating provider.yaml checksums against generated checksums...$(NC)"
	@if [ ! -d $(BUILD_DIR) ] || [ ! -f $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64.sha256 ]; then \
		echo "$(RED)Error: Generated checksums not found. Run 'make checksums' first$(NC)"; \
		exit 1; \
	fi
	@failed=0; \
	echo "$(YELLOW)Checking each platform checksum:$(NC)"; \
	for platform in "linux-amd64" "linux-arm64" "darwin-amd64" "darwin-arm64" "windows-amd64.exe"; do \
		if [ -f "$(BUILD_DIR)/$(BINARY_NAME)-$$platform.sha256" ]; then \
			generated_checksum=$$(cut -d' ' -f1 "$(BUILD_DIR)/$(BINARY_NAME)-$$platform.sha256"); \
			provider_checksum=$$(grep -A1 "dokploy-provider-$$platform" $(PROVIDER_FILE) | grep "checksum:" | sed 's/.*checksum: *"\([^"]*\)".*/\1/'); \
			if [ "$$generated_checksum" = "$$provider_checksum" ]; then \
				echo "$(GREEN)✓ $$platform: checksums match$(NC)"; \
			else \
				echo "$(RED)✗ $$platform: checksum mismatch$(NC)"; \
				echo "  Generated: $$generated_checksum"; \
				echo "  Provider:  $$provider_checksum"; \
				failed=1; \
			fi; \
		fi; \
	done; \
	echo ""; \
	echo "$(BLUE)Checking checksum uniqueness in provider.yaml:$(NC)"; \
	provider_checksums=$$(grep "checksum:" $(PROVIDER_FILE) | sed 's/.*checksum: *"\([^"]*\)".*/\1/' | sort); \
	unique_provider_checksums=$$(echo "$$provider_checksums" | sort -u); \
	provider_count=$$(echo "$$provider_checksums" | wc -l | tr -d ' '); \
	unique_count=$$(echo "$$unique_provider_checksums" | wc -l | tr -d ' '); \
	if [ "$$provider_count" = "$$unique_count" ]; then \
		echo "$(GREEN)✓ All $$provider_count checksums in provider.yaml are unique$(NC)"; \
	else \
		echo "$(RED)✗ Only $$unique_count unique checksums out of $$provider_count total in provider.yaml$(NC)"; \
		echo "$(YELLOW)Duplicate checksums found:$(NC)"; \
		echo "$$provider_checksums" | sort | uniq -d; \
		failed=1; \
	fi; \
	if [ $$failed -eq 0 ]; then \
		echo ""; \
		echo "$(GREEN)✓ All checksums are valid and unique$(NC)"; \
	else \
		echo ""; \
		echo "$(RED)✗ Checksum validation failed$(NC)"; \
		echo "$(YELLOW)Run 'make update-provider-checksums' to fix$(NC)"; \
		exit 1; \
	fi

.PHONY: release-prepare
release-prepare: validate version-check build-all checksums update-provider-checksums validate-provider-checksums ## Prepare everything for release
	@echo "$(BLUE)Preparing release v$(VERSION)...$(NC)"
	@echo "$(GREEN)✓ Release v$(VERSION) prepared successfully!$(NC)"
	@echo ""
	@echo "$(BLUE)Release artifacts:$(NC)"
	@ls -la $(BUILD_DIR)/
	@echo ""
	@echo "$(BLUE)Next steps:$(NC)"
	@echo "  1. Review the updated provider.yaml (URLs and checksums updated)"
	@echo "  2. Test the provider: make test-docker"
	@echo "  3. Create release: make tag-release"
	@echo "  4. Upload binaries to GitHub releases"
	@echo "  5. Update CHANGELOG.md"

.PHONY: release-clean
release-clean: ## Clean release artifacts but keep backups
	@echo "$(YELLOW)Cleaning release artifacts...$(NC)"
	@if [ -d $(BUILD_DIR) ]; then \
		rm -f $(BUILD_DIR)/*.sha256; \
		echo "$(GREEN)✓ Checksum files removed$(NC)"; \
	fi
	@if [ -f $(PROVIDER_FILE).backup ]; then \
		echo "$(YELLOW)provider.yaml backup preserved: $(PROVIDER_FILE).backup$(NC)"; \
	fi

.PHONY: restore-provider
restore-provider: ## Restore provider.yaml from backup
	@echo "$(YELLOW)Restoring provider.yaml from backup...$(NC)"
	@if [ -f $(PROVIDER_FILE).backup ]; then \
		cp $(PROVIDER_FILE).backup $(PROVIDER_FILE); \
		echo "$(GREEN)✓ provider.yaml restored from backup$(NC)"; \
	else \
		echo "$(RED)No backup found: $(PROVIDER_FILE).backup$(NC)"; \
		exit 1; \
	fi

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	@echo "$(GREEN)✓ Build artifacts cleaned$(NC)"

.PHONY: deps
deps: ## Download and tidy Go dependencies
	@echo "$(BLUE)Managing Go dependencies...$(NC)"
	go mod download
	go mod tidy
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

.PHONY: test-build
test-build: build ## Test the built binary
	@echo "$(BLUE)Testing built binary...$(NC)"
	@if [ -f "$(BUILD_DIR)/$(BINARY_NAME)" ]; then \
		echo "$(YELLOW)Testing binary help command...$(NC)"; \
		$(BUILD_DIR)/$(BINARY_NAME) --help; \
		echo "$(GREEN)✓ Binary test completed$(NC)"; \
	else \
		echo "$(RED)✗ Binary not found: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"; \
		exit 1; \
	fi

##@ Development

.PHONY: setup
setup: ## Install required development tools
	@echo "$(BLUE)Installing required development tools...$(NC)"
	@echo "$(YELLOW)Checking for package managers...$(NC)"
	@if command -v brew >/dev/null 2>&1; then \
		echo "$(GREEN)✓ Homebrew found$(NC)"; \
		echo "$(YELLOW)Installing tools via Homebrew...$(NC)"; \
		brew install yq jq shellcheck 2>/dev/null || true; \
		echo "$(YELLOW)Checking DevPod installation...$(NC)"; \
		if [ -d "/Applications/DevPod.app" ] || command -v devpod >/dev/null 2>&1; then \
			echo "$(GREEN)✓ DevPod already installed$(NC)"; \
		else \
			echo "$(YELLOW)Installing DevPod Desktop App...$(NC)"; \
			brew install --cask devpod 2>/dev/null || true; \
		fi; \
	elif command -v apt-get >/dev/null 2>&1; then \
		echo "$(GREEN)✓ APT found$(NC)"; \
		echo "$(YELLOW)Installing tools via APT...$(NC)"; \
		sudo apt-get update && sudo apt-get install -y jq shellcheck curl; \
		echo "$(YELLOW)Installing yq...$(NC)"; \
		sudo wget -qO /usr/local/bin/yq https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64; \
		sudo chmod +x /usr/local/bin/yq; \
		echo "$(YELLOW)Checking DevPod installation...$(NC)"; \
		if command -v devpod >/dev/null 2>&1; then \
			echo "$(GREEN)✓ DevPod already installed$(NC)"; \
		else \
			echo "$(YELLOW)Installing DevPod CLI...$(NC)"; \
			curl -L -o devpod "https://github.com/loft-sh/devpod/releases/latest/download/devpod-linux-amd64"; \
			sudo install -c -m 0755 devpod /usr/local/bin && rm -f devpod; \
		fi; \
	elif command -v yum >/dev/null 2>&1; then \
		echo "$(GREEN)✓ YUM found$(NC)"; \
		echo "$(YELLOW)Installing tools via YUM...$(NC)"; \
		sudo yum install -y jq ShellCheck curl; \
		echo "$(YELLOW)Installing yq...$(NC)"; \
		sudo wget -qO /usr/local/bin/yq https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64; \
		sudo chmod +x /usr/local/bin/yq; \
		echo "$(YELLOW)Checking DevPod installation...$(NC)"; \
		if command -v devpod >/dev/null 2>&1; then \
			echo "$(GREEN)✓ DevPod already installed$(NC)"; \
		else \
			echo "$(YELLOW)Installing DevPod CLI...$(NC)"; \
			curl -L -o devpod "https://github.com/loft-sh/devpod/releases/latest/download/devpod-linux-amd64"; \
			sudo install -c -m 0755 devpod /usr/local/bin && rm -f devpod; \
		fi; \
	else \
		echo "$(YELLOW)⚠ No supported package manager found$(NC)"; \
		echo "$(BLUE)Please install manually:$(NC)"; \
		echo "  - yq: https://github.com/mikefarah/yq#install"; \
		echo "  - jq: https://stedolan.github.io/jq/download/"; \
		echo "  - shellcheck: https://github.com/koalaman/shellcheck#installing"; \
		echo "  - devpod: https://devpod.sh/docs/getting-started/install"; \
	fi
	@echo "$(GREEN)Tool installation completed$(NC)"

.PHONY: check-tools
check-tools: ## Check if required tools are installed
	@echo "$(BLUE)Checking required tools...$(NC)"
	@echo -n "yq: "; \
	if command -v yq >/dev/null 2>&1; then \
		echo "$(GREEN)✓ installed ($(shell yq --version))$(NC)"; \
	else \
		echo "$(RED)✗ missing$(NC)"; \
		echo "  Install: brew install yq  # or make setup"; \
	fi
	@echo -n "jq: "; \
	if command -v jq >/dev/null 2>&1; then \
		echo "$(GREEN)✓ installed ($(shell jq --version))$(NC)"; \
	else \
		echo "$(RED)✗ missing$(NC)"; \
		echo "  Install: brew install jq  # or make setup"; \
	fi
	@echo -n "shellcheck: "; \
	if command -v shellcheck >/dev/null 2>&1; then \
		echo "$(GREEN)✓ installed ($(shell shellcheck --version | head -n2 | tail -n1))$(NC)"; \
	else \
		echo "$(RED)✗ missing$(NC)"; \
		echo "  Install: brew install shellcheck  # or make setup"; \
	fi
	@echo -n "devpod: "; \
	if command -v devpod >/dev/null 2>&1; then \
		echo "$(GREEN)✓ installed ($(shell devpod version))$(NC)"; \
	else \
		echo "$(RED)✗ missing$(NC)"; \
		echo "  Install: brew install --cask devpod  # macOS Desktop App"; \
		echo "  Or CLI: curl -L -o devpod https://github.com/loft-sh/devpod/releases/latest/download/devpod-darwin-amd64 && sudo install -c -m 0755 devpod /usr/local/bin"; \
		echo "  Or visit: https://devpod.sh/docs/getting-started/install"; \
	fi

.PHONY: check-devpod
check-devpod: ## Check DevPod CLI availability and provide setup instructions
	@echo "$(BLUE)Checking DevPod CLI availability...$(NC)"
	@if [ "$(DEVPOD_CMD)" = "devpod-not-found" ]; then \
		echo "$(RED)✗ DevPod CLI not found$(NC)"; \
		echo "$(YELLOW)DevPod Desktop App may be installed but CLI is not in PATH$(NC)"; \
		echo "$(BLUE)To fix this, choose one option:$(NC)"; \
		echo ""; \
		echo "$(YELLOW)Option 1: Install DevPod CLI separately$(NC)"; \
		echo "  curl -L -o devpod https://github.com/loft-sh/devpod/releases/latest/download/devpod-darwin-amd64"; \
		echo "  sudo install -c -m 0755 devpod /usr/local/bin && rm -f devpod"; \
		echo ""; \
		echo "$(YELLOW)Option 2: Add DevPod.app CLI to PATH$(NC)"; \
		echo "  echo 'export PATH=\"/Applications/DevPod.app/Contents/MacOS:\$$PATH\"' >> ~/.zshrc"; \
		echo "  source ~/.zshrc"; \
		echo ""; \
		echo "$(YELLOW)Option 3: Use make install-devpod-cli$(NC)"; \
		echo "  make install-devpod-cli"; \
		echo ""; \
		exit 1; \
	elif [ "$(DEVPOD_CMD)" != "devpod" ]; then \
		echo "$(YELLOW)⚠ Using DevPod from app bundle: $(DEVPOD_CMD)$(NC)"; \
		echo "$(BLUE)Consider adding to PATH for easier access$(NC)"; \
	else \
		echo "$(GREEN)✓ DevPod CLI available: $(DEVPOD_CMD)$(NC)"; \
	fi

.PHONY: install-devpod-cli
install-devpod-cli: ## Install DevPod CLI to /usr/local/bin
	@echo "$(BLUE)Installing DevPod CLI...$(NC)"
	@if command -v devpod >/dev/null 2>&1; then \
		echo "$(GREEN)✓ DevPod CLI already in PATH$(NC)"; \
	else \
		echo "$(YELLOW)Downloading DevPod CLI...$(NC)"; \
		curl -L -o devpod "https://github.com/loft-sh/devpod/releases/latest/download/devpod-darwin-amd64"; \
		sudo install -c -m 0755 devpod /usr/local/bin; \
		rm -f devpod; \
		echo "$(GREEN)✓ DevPod CLI installed to /usr/local/bin/devpod$(NC)"; \
	fi

.PHONY: install-local
install-local: check-devpod check-tools validate ## Install provider locally for development
	@echo "$(BLUE)Installing provider locally...$(NC)"
	@if [ -f .env ] && [ -n "$(DOKPLOY_SERVER_URL)" ] && [ -n "$(DOKPLOY_API_TOKEN)" ]; then \
		echo "$(YELLOW)Installing with configuration from .env file...$(NC)"; \
		$(DEVPOD_CMD) provider add ./$(PROVIDER_FILE) --name $(PROVIDER_NAME)-dev \
			--option DOKPLOY_SERVER_URL="$(DOKPLOY_SERVER_URL)" \
			--option DOKPLOY_API_TOKEN="$(DOKPLOY_API_TOKEN)" \
			$$([ -n "$(DOKPLOY_PROJECT_NAME)" ] && echo "--option DOKPLOY_PROJECT_NAME=$(DOKPLOY_PROJECT_NAME)") \
			$$([ -n "$(DOKPLOY_SERVER_ID)" ] && echo "--option DOKPLOY_SERVER_ID=$(DOKPLOY_SERVER_ID)") \
			$$([ -n "$(MACHINE_TYPE)" ] && echo "--option MACHINE_TYPE=$(MACHINE_TYPE)") \
			$$([ -n "$(MACHINE_IMAGE)" ] && echo "--option MACHINE_IMAGE=$(MACHINE_IMAGE)") \
			$$([ -n "$(AGENT_PATH)" ] && echo "--option AGENT_PATH=$(AGENT_PATH)"); \
	else \
		echo "$(YELLOW)Installing without configuration (use make configure-env or make configure)...$(NC)"; \
		$(DEVPOD_CMD) provider add ./$(PROVIDER_FILE) --name $(PROVIDER_NAME)-dev --use=false; \
	fi
	@echo "$(GREEN)Provider installed as '$(PROVIDER_NAME)-dev'$(NC)"
	@if [ ! -f .env ] || [ -z "$(DOKPLOY_SERVER_URL)" ] || [ -z "$(DOKPLOY_API_TOKEN)" ]; then \
		echo "$(YELLOW)To configure and use the provider, run:$(NC)"; \
		echo "  make configure-env  # (if you have .env file)"; \
		echo "  make configure      # (for interactive setup)"; \
	fi

.PHONY: install-dev
install-dev: check-devpod build ## Install provider locally for development using local binary
	@echo "$(BLUE)Installing development provider locally...$(NC)"
	@if [ -f .env ] && [ -n "$(DOKPLOY_SERVER_URL)" ] && [ -n "$(DOKPLOY_API_TOKEN)" ]; then \
		echo "$(YELLOW)Installing with configuration from .env file...$(NC)"; \
		$(DEVPOD_CMD) provider add ./$(PROVIDER_DEV_FILE) --name $(PROVIDER_NAME)-dev \
			--option DOKPLOY_SERVER_URL="$(DOKPLOY_SERVER_URL)" \
			--option DOKPLOY_API_TOKEN="$(DOKPLOY_API_TOKEN)" \
			--option DOKPLOY_PROVIDER_PATH="$(PWD)/$(BUILD_DIR)/$(BINARY_NAME)" \
			$$([ -n "$(DOKPLOY_PROJECT_NAME)" ] && echo "--option DOKPLOY_PROJECT_NAME=$(DOKPLOY_PROJECT_NAME)") \
			$$([ -n "$(DOKPLOY_SERVER_ID)" ] && echo "--option DOKPLOY_SERVER_ID=$(DOKPLOY_SERVER_ID)") \
			$$([ -n "$(MACHINE_TYPE)" ] && echo "--option MACHINE_TYPE=$(MACHINE_TYPE)") \
			$$([ -n "$(AGENT_PATH)" ] && echo "--option AGENT_PATH=$(AGENT_PATH)"); \
	else \
		echo "$(YELLOW)Installing without configuration (use make configure-env or make configure)...$(NC)"; \
		$(DEVPOD_CMD) provider add ./$(PROVIDER_DEV_FILE) --name $(PROVIDER_NAME)-dev --use=false; \
	fi
	@echo "$(GREEN)Development provider installed as '$(PROVIDER_NAME)-dev'$(NC)"
	@if [ ! -f .env ] || [ -z "$(DOKPLOY_SERVER_URL)" ] || [ -z "$(DOKPLOY_API_TOKEN)" ]; then \
		echo "$(YELLOW)To configure and use the provider, run:$(NC)"; \
		echo "  make configure-env  # (if you have .env file)"; \
		echo "  make configure      # (for interactive setup)"; \
	fi

.PHONY: install-github
install-github: ## Install provider from GitHub repository
	@echo "$(BLUE)Installing provider from GitHub...$(NC)"
	$(DEVPOD_CMD) provider add $(GITHUB_REPO) --name $(PROVIDER_NAME)
	@echo "$(GREEN)Provider installed from GitHub$(NC)"
	@echo "$(YELLOW)To configure and use the provider, run:$(NC)"
	@echo "  $(DEVPOD_CMD) provider use $(PROVIDER_NAME)"

.PHONY: uninstall
uninstall: ## Remove the provider from DevPod
	@echo "$(YELLOW)Removing provider...$(NC)"
	@echo "$(BLUE)Checking for active workspaces...$(NC)"
	@if $(DEVPOD_CMD) list --output json 2>/dev/null | jq -r '.[].provider' 2>/dev/null | grep -q "$(PROVIDER_NAME)"; then \
		echo "$(RED)⚠ Active workspaces found using this provider:$(NC)"; \
		$(DEVPOD_CMD) list --output json 2>/dev/null | jq -r '.[] | select(.provider == "$(PROVIDER_NAME)" or .provider == "$(PROVIDER_NAME)-dev") | "  - " + .id + " (status: " + .status + ")"' 2>/dev/null || true; \
		echo "$(YELLOW)Please stop/delete these workspaces first:$(NC)"; \
		echo "  make cleanup-workspaces  # Delete all workspaces for this provider"; \
		echo "  devpod delete <workspace-name>  # Delete specific workspace"; \
		echo "  devpod stop <workspace-name>   # Stop specific workspace"; \
		exit 1; \
	fi
	-$(DEVPOD_CMD) provider delete $(PROVIDER_NAME)-dev --ignore-not-found 2>/dev/null || true
	-$(DEVPOD_CMD) provider delete $(PROVIDER_NAME) --ignore-not-found 2>/dev/null || true
	@echo "$(GREEN)Provider removed$(NC)"

.PHONY: cleanup-workspaces
cleanup-workspaces: ## Delete all workspaces using this provider
	@echo "$(YELLOW)Cleaning up workspaces for provider $(PROVIDER_NAME)...$(NC)"
	@echo "$(BLUE)Checking for workspaces using this provider...$(NC)"
	@workspaces=$$($(DEVPOD_CMD) list --output json 2>/dev/null | jq -r '.[] | select(.provider.name == "$(PROVIDER_NAME)" or .provider.name == "$(PROVIDER_NAME)-dev") | .id' 2>/dev/null || true); \
	if [ -n "$$workspaces" ]; then \
		echo "$(BLUE)Found workspaces to clean up:$(NC)"; \
		for ws in $$workspaces; do \
			echo "  - $$ws"; \
		done; \
		echo "$(YELLOW)Deleting workspaces with --force flag...$(NC)"; \
		for ws in $$workspaces; do \
			echo "$(BLUE)Deleting workspace: $$ws$(NC)"; \
			$(DEVPOD_CMD) delete $$ws --force --ignore-not-found 2>/dev/null || true; \
		done; \
		echo "$(GREEN)✓ Workspaces cleaned up$(NC)"; \
	else \
		echo "$(GREEN)No workspaces found for this provider$(NC)"; \
	fi

.PHONY: force-reinstall
force-reinstall: ## Force reinstall provider (handles active workspaces and stubborn providers)
	@echo "$(YELLOW)Force reinstalling provider $(PROVIDER_NAME)...$(NC)"
	@$(MAKE) cleanup-workspaces
	@echo "$(YELLOW)Force removing provider...$(NC)"
	@echo "$(BLUE)Current providers before deletion:$(NC)"
	@$(DEVPOD_CMD) provider list || true
	@echo ""
	@echo "$(YELLOW)Attempting to delete providers (with retries)...$(NC)"
	@for attempt in 1 2 3; do \
		echo "$(BLUE)Deletion attempt $$attempt/3...$(NC)"; \
		for provider in $(PROVIDER_NAME) $(PROVIDER_NAME)-dev; do \
			if $(DEVPOD_CMD) provider list --output json 2>/dev/null | jq -e "has(\"$$provider\")" >/dev/null 2>&1; then \
				echo "$(BLUE)Deleting $$provider...$(NC)"; \
				if $(DEVPOD_CMD) provider delete $$provider 2>/dev/null; then \
					echo "$(GREEN)✓ Successfully deleted $$provider$(NC)"; \
				else \
					echo "$(RED)Failed to delete $$provider on attempt $$attempt$(NC)"; \
				fi; \
			fi; \
		done; \
		if ! $(DEVPOD_CMD) provider list --output json 2>/dev/null | jq -e "has(\"$(PROVIDER_NAME)\") or has(\"$(PROVIDER_NAME)-dev\")" >/dev/null 2>&1; then \
			echo "$(GREEN)✓ All providers deleted successfully$(NC)"; \
			break; \
		fi; \
		if [ $$attempt -lt 3 ]; then \
			echo "$(YELLOW)Retrying in 3 seconds...$(NC)"; \
			sleep 3; \
			echo "$(BLUE)Re-checking for remaining workspaces...$(NC)"; \
			$(MAKE) cleanup-workspaces; \
		fi; \
	done
	@if $(DEVPOD_CMD) provider list --output json 2>/dev/null | jq -e "has(\"$(PROVIDER_NAME)\") or has(\"$(PROVIDER_NAME)-dev\")" >/dev/null 2>&1; then \
		echo "$(RED)⚠ Provider deletion failed after 3 attempts$(NC)"; \
		echo "$(YELLOW)Attempting nuclear cleanup...$(NC)"; \
		$(MAKE) nuclear-cleanup-providers; \
	fi
	@echo "$(BLUE)Current providers after deletion:$(NC)"
	@$(DEVPOD_CMD) provider list || true
	@echo ""
	@echo "$(GREEN)Provider force removal completed$(NC)"
	@$(MAKE) build
	@$(MAKE) install-dev

##@ Configuration

.PHONY: setup-env
setup-env: ## Create .env file from .env.example
	@echo "$(BLUE)Setting up environment configuration...$(NC)"
	@if [ ! -f .env ]; then \
		if [ -f .env.example ]; then \
			cp .env.example .env; \
			echo "$(GREEN)✓ Created .env file from .env.example$(NC)"; \
			echo "$(YELLOW)Please edit .env file with your Dokploy configuration$(NC)"; \
		else \
			echo "$(RED)✗ .env.example file not found$(NC)"; \
			exit 1; \
		fi; \
	else \
		echo "$(YELLOW)⚠ .env file already exists$(NC)"; \
	fi

.PHONY: configure-env
configure-env: setup-env ## Configure provider using .env file
	@echo "$(BLUE)Configuring provider from .env file...$(NC)"
	@if [ ! -f .env ]; then \
		echo "$(RED)✗ .env file not found. Run 'make setup-env' first$(NC)"; \
		exit 1; \
	fi
	@if [ -z "$(DOKPLOY_SERVER_URL)" ] || [ -z "$(DOKPLOY_API_TOKEN)" ]; then \
		echo "$(RED)✗ Required variables not set in .env file$(NC)"; \
		echo "$(YELLOW)Please ensure DOKPLOY_SERVER_URL and DOKPLOY_API_TOKEN are set$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Configuring provider with environment variables...$(NC)"
	$(DEVPOD_CMD) provider set-options $(PROVIDER_NAME)-dev \
		--option DOKPLOY_SERVER_URL="$(DOKPLOY_SERVER_URL)" \
		--option DOKPLOY_API_TOKEN="$(DOKPLOY_API_TOKEN)" \
		$$([ -n "$(DOKPLOY_PROJECT_NAME)" ] && echo "--option DOKPLOY_PROJECT_NAME=$(DOKPLOY_PROJECT_NAME)") \
		$$([ -n "$(DOKPLOY_SERVER_ID)" ] && echo "--option DOKPLOY_SERVER_ID=$(DOKPLOY_SERVER_ID)") \
		$$([ -n "$(MACHINE_TYPE)" ] && echo "--option MACHINE_TYPE=$(MACHINE_TYPE)") \
		$$([ -n "$(MACHINE_IMAGE)" ] && echo "--option MACHINE_IMAGE=$(MACHINE_IMAGE)") \
		$$([ -n "$(AGENT_PATH)" ] && echo "--option AGENT_PATH=$(AGENT_PATH)")
	@echo "$(GREEN)Provider configured successfully$(NC)"

.PHONY: configure
configure: ## Configure provider with required options (interactive)
	@echo "$(BLUE)Configuring provider...$(NC)"
	@read -p "Enter Dokploy Server URL: " server_url; \
	read -s -p "Enter Dokploy API Token: " api_token; \
	echo; \
	$(DEVPOD_CMD) provider set-options $(PROVIDER_NAME)-dev \
		--option DOKPLOY_SERVER_URL=$$server_url \
		--option DOKPLOY_API_TOKEN=$$api_token
	@echo "$(GREEN)Provider configured$(NC)"

.PHONY: configure-optional
configure-optional: ## Configure optional provider settings
	@echo "$(BLUE)Configuring optional settings...$(NC)"
	@read -p "Enter project name [devpod-workspaces]: " project_name; \
	project_name=$${project_name:-devpod-workspaces}; \
	read -p "Enter server ID (optional): " server_id; \
	read -p "Enter machine type [small]: " machine_type; \
	machine_type=$${machine_type:-small}; \
	read -p "Enter machine image [ubuntu:22.04]: " machine_image; \
	machine_image=$${machine_image:-ubuntu:22.04}; \
	$(DEVPOD_CMD) provider set-options $(PROVIDER_NAME)-dev \
		--option DOKPLOY_PROJECT_NAME=$$project_name \
		--option MACHINE_TYPE=$$machine_type \
		--option MACHINE_IMAGE=$$machine_image \
		$$([ -n "$$server_id" ] && echo "--option DOKPLOY_SERVER_ID=$$server_id")
	@echo "$(GREEN)Optional settings configured$(NC)"

.PHONY: show-config
show-config: ## Show current provider configuration
	@echo "$(BLUE)Current provider configuration:$(NC)"
	$(DEVPOD_CMD) provider options $(PROVIDER_NAME)-dev || $(DEVPOD_CMD) provider options $(PROVIDER_NAME)

##@ Testing

.PHONY: test-docker
test-docker: ## Test provider with Docker image workspace
	@echo "$(BLUE)Testing provider with Docker image...$(NC)"
	$(DEVPOD_CMD) up $(TEST_WORKSPACE)-docker \
		--provider $(PROVIDER_NAME)-dev \
		--workspace-image ubuntu:22.04 \
		--debug
	@echo "$(GREEN)Docker workspace test completed$(NC)"

.PHONY: test-git
test-git: ## Test provider with Git repository workspace
	@echo "$(BLUE)Testing provider with Git repository...$(NC)"
	$(DEVPOD_CMD) up $(TEST_REPO) \
		--provider $(PROVIDER_NAME)-dev \
		--debug
	@echo "$(GREEN)Git workspace test completed$(NC)"

.PHONY: test-lifecycle
test-lifecycle: ## Test complete workspace lifecycle
	@echo "$(BLUE)Testing complete workspace lifecycle...$(NC)"
	@workspace_name=$(TEST_WORKSPACE)-lifecycle; \
	echo "Creating workspace: $$workspace_name"; \
	$(DEVPOD_CMD) up $$workspace_name --provider $(PROVIDER_NAME)-dev --workspace-image ubuntu:22.04; \
	echo "$(YELLOW)Checking workspace status...$(NC)"; \
	$(DEVPOD_CMD) status $$workspace_name; \
	echo "$(YELLOW)Stopping workspace...$(NC)"; \
	$(DEVPOD_CMD) stop $$workspace_name; \
	echo "$(YELLOW)Starting workspace...$(NC)"; \
	$(DEVPOD_CMD) start $$workspace_name; \
	echo "$(YELLOW)Deleting workspace...$(NC)"; \
	$(DEVPOD_CMD) delete $$workspace_name --force
	@echo "$(GREEN)Lifecycle test completed$(NC)"

.PHONY: test-ssh
test-ssh: ## Test SSH connection to workspace
	@echo "$(BLUE)Testing SSH connection...$(NC)"
	@workspace_name=$(TEST_WORKSPACE)-ssh; \
	$(DEVPOD_CMD) up $$workspace_name --provider $(PROVIDER_NAME)-dev --workspace-image ubuntu:22.04; \
	echo "$(YELLOW)Testing SSH connection...$(NC)"; \
	$(DEVPOD_CMD) ssh $$workspace_name -- echo "SSH connection successful"; \
	$(DEVPOD_CMD) delete $$workspace_name --force
	@echo "$(GREEN)SSH test completed$(NC)"

.PHONY: cleanup-test
cleanup-test: ## Clean up test workspaces
	@echo "$(YELLOW)Cleaning up test workspaces...$(NC)"
	@for ws in $$($(DEVPOD_CMD) list --output json 2>/dev/null | jq -r '.[].id' 2>/dev/null | grep -E "test-workspace|vscode-remote-try|test-debug|test-" || true); do \
		echo "$(BLUE)Deleting test workspace: $$ws$(NC)"; \
		$(DEVPOD_CMD) delete $$ws --force --ignore-not-found 2>/dev/null || true; \
	done
	@echo "$(GREEN)Test workspaces cleaned up$(NC)"

.PHONY: nuclear-cleanup
nuclear-cleanup: ## Nuclear option: delete ALL workspaces and providers
	@echo "$(RED)⚠ WARNING: This will delete ALL DevPod workspaces and providers!$(NC)"
	@echo "$(YELLOW)This is a nuclear option for when things are completely stuck.$(NC)"
	@read -p "Are you sure? Type 'yes' to continue: " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		echo "$(RED)Performing nuclear cleanup...$(NC)"; \
		echo "$(BLUE)Deleting all workspaces...$(NC)"; \
		for ws in $$($(DEVPOD_CMD) list --output json 2>/dev/null | jq -r '.[].id' 2>/dev/null || true); do \
			echo "$(YELLOW)Deleting workspace: $$ws$(NC)"; \
			$(DEVPOD_CMD) delete $$ws --force --ignore-not-found 2>/dev/null || true; \
		done; \
		sleep 5; \
		echo "$(BLUE)Deleting all providers...$(NC)"; \
		for provider in $$($(DEVPOD_CMD) provider list --output json 2>/dev/null | jq -r '.[].name' 2>/dev/null || true); do \
			echo "$(YELLOW)Deleting provider: $$provider$(NC)"; \
			$(DEVPOD_CMD) provider delete $$provider --force --ignore-not-found 2>/dev/null || true; \
		done; \
		echo "$(GREEN)Nuclear cleanup completed$(NC)"; \
	else \
		echo "$(GREEN)Nuclear cleanup cancelled$(NC)"; \
	fi

.PHONY: fix-stuck-workspace
fix-stuck-workspace: ## Fix a specific stuck workspace (interactive)
	@echo "$(BLUE)Fix stuck workspace$(NC)"
	@echo "$(YELLOW)Current workspaces:$(NC)"
	@$(DEVPOD_CMD) list || true
	@read -p "Enter workspace name to force delete: " ws_name; \
	if [ -n "$$ws_name" ]; then \
		echo "$(BLUE)Attempting to force delete workspace: $$ws_name$(NC)"; \
		echo "$(YELLOW)Method 1: Standard force delete...$(NC)"; \
		$(DEVPOD_CMD) delete $$ws_name --force 2>/dev/null || true; \
		echo "$(YELLOW)Method 2: Stop then delete...$(NC)"; \
		$(DEVPOD_CMD) stop $$ws_name --force 2>/dev/null || true; \
		$(DEVPOD_CMD) delete $$ws_name --force 2>/dev/null || true; \
		echo "$(YELLOW)Method 3: Ignore not found...$(NC)"; \
		$(DEVPOD_CMD) delete $$ws_name --force --ignore-not-found 2>/dev/null || true; \
		echo "$(GREEN)Attempted all deletion methods for: $$ws_name$(NC)"; \
	else \
		echo "$(RED)No workspace name provided$(NC)"; \
	fi

##@ Validation

.PHONY: validate
validate: ## Validate provider.yaml syntax and structure
	@echo "$(BLUE)Validating provider configuration...$(NC)"
	@if [ ! -f $(PROVIDER_FILE) ]; then \
		echo "$(RED)Error: $(PROVIDER_FILE) not found$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Checking YAML syntax...$(NC)"
	@if command -v yq >/dev/null 2>&1; then \
		yq eval '.' $(PROVIDER_FILE) >/dev/null && echo "$(GREEN)✓ YAML syntax valid$(NC)"; \
	elif command -v python3 >/dev/null 2>&1; then \
		python3 -c "import yaml; yaml.safe_load(open('$(PROVIDER_FILE)'))" && echo "$(GREEN)✓ YAML syntax valid$(NC)"; \
	else \
		echo "$(YELLOW)⚠ Cannot validate YAML syntax$(NC)"; \
		echo "$(BLUE)Install validation tools with:$(NC)"; \
		echo "  brew install yq           # macOS with Homebrew"; \
		echo "  sudo apt install python3  # Ubuntu/Debian"; \
		echo "  make setup                # Auto-install"; \
	fi
	@echo "$(YELLOW)Checking required fields...$(NC)"
	@grep -q "^name:" $(PROVIDER_FILE) && echo "$(GREEN)✓ Name field present$(NC)" || (echo "$(RED)✗ Name field missing$(NC)" && exit 1)
	@grep -q "^version:" $(PROVIDER_FILE) && echo "$(GREEN)✓ Version field present$(NC)" || (echo "$(RED)✗ Version field missing$(NC)" && exit 1)
	@grep -q "^exec:" $(PROVIDER_FILE) && echo "$(GREEN)✓ Exec section present$(NC)" || (echo "$(RED)✗ Exec section missing$(NC)" && exit 1)
	@echo "$(GREEN)Provider validation completed$(NC)"

.PHONY: lint
lint: ## Lint provider configuration and scripts
	@echo "$(BLUE)Linting provider...$(NC)"
	@if command -v shellcheck >/dev/null 2>&1; then \
		echo "$(YELLOW)Running shellcheck on embedded scripts...$(NC)"; \
		grep -A 1000 "exec:" $(PROVIDER_FILE) | grep -E "^\s+[a-z]+:\s*\|-" -A 1000 | \
		sed 's/^[[:space:]]*//' | shellcheck -s bash - || true; \
	else \
		echo "$(YELLOW)⚠ Shellcheck not available$(NC)"; \
		echo "$(BLUE)Install shellcheck with:$(NC)"; \
		echo "  brew install shellcheck    # macOS with Homebrew"; \
		echo "  sudo apt install shellcheck # Ubuntu/Debian"; \
		echo "  make setup                 # Auto-install"; \
	fi
	@echo "$(GREEN)Linting completed$(NC)"

##@ Documentation

.PHONY: docs
docs: ## Generate documentation
	@echo "$(BLUE)Generating documentation...$(NC)"
	@echo "# Dokploy DevPod Provider" > USAGE.md
	@echo "" >> USAGE.md
	@echo "## Installation" >> USAGE.md
	@echo "" >> USAGE.md
	@echo "\`\`\`bash" >> USAGE.md
	@echo "$(DEVPOD_CMD) provider add $(GITHUB_REPO)" >> USAGE.md
	@echo "\`\`\`" >> USAGE.md
	@echo "" >> USAGE.md
	@echo "## Configuration" >> USAGE.md
	@echo "" >> USAGE.md
	@grep -A 10 "options:" $(PROVIDER_FILE) | grep -E "^\s+[A-Z_]+:" | sed 's/^[[:space:]]*/- /' >> USAGE.md
	@echo "$(GREEN)Documentation generated in USAGE.md$(NC)"

##@ Release Management

.PHONY: version-check
version-check: ## Check if version is properly set
	@echo "$(BLUE)Current version: $(VERSION)$(NC)"
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Error: Version not found in $(PROVIDER_FILE)$(NC)"; \
		exit 1; \
	fi

# Helper function to update version and URLs
define update-version-and-urls
	@current_version=$(VERSION); \
	new_version=$(1); \
	echo "$(YELLOW)Updating version from $$current_version to $$new_version$(NC)"; \
	if [[ "$$OSTYPE" == "darwin"* ]]; then \
		sed -i.bak "s/version: v*$$current_version/version: v$$new_version/" $(PROVIDER_FILE); \
		sed -i.bak2 "s|releases/download/[^/]*/|releases/download/v$$new_version/|g" $(PROVIDER_FILE); \
	else \
		sed -i.bak "s/version: v*$$current_version/version: v$$new_version/" $(PROVIDER_FILE); \
		sed -i.bak2 "s|releases/download/[^/]*/|releases/download/v$$new_version/|g" $(PROVIDER_FILE); \
	fi; \
	rm -f $(PROVIDER_FILE).bak $(PROVIDER_FILE).bak2; \
	echo "$(GREEN)✓ Version bumped from $$current_version to $$new_version$(NC)"; \
	echo "$(GREEN)✓ URLs updated to use v$$new_version$(NC)"
endef

.PHONY: version-bump-patch
version-bump-patch: ## Bump patch version (x.y.Z) and update URLs
	@echo "$(BLUE)Bumping patch version...$(NC)"
	$(call update-version-and-urls,$$(echo $(VERSION) | awk -F. '{$$3++; print $$1"."$$2"."$$3}'))

.PHONY: version-bump-minor
version-bump-minor: ## Bump minor version (x.Y.z) and update URLs
	@echo "$(BLUE)Bumping minor version...$(NC)"
	$(call update-version-and-urls,$$(echo $(VERSION) | awk -F. '{$$2++; $$3=0; print $$1"."$$2"."$$3}'))

.PHONY: version-bump-major
version-bump-major: ## Bump major version (X.y.z) and update URLs
	@echo "$(BLUE)Bumping major version...$(NC)"
	$(call update-version-and-urls,$$(echo $(VERSION) | awk -F. '{$$1++; $$2=0; $$3=0; print $$1"."$$2"."$$3}'))

.PHONY: tag-release
tag-release: validate version-check ## Create and push git tag for release
	@echo "$(BLUE)Creating release tag v$(VERSION)...$(NC)"
	@if git tag | grep -q "^v$(VERSION)$$"; then \
		echo "$(RED)Error: Tag v$(VERSION) already exists$(NC)"; \
		exit 1; \
	fi
	git add $(PROVIDER_FILE)
	git commit -m "Release v$(VERSION)" || true
	git tag -a "v$(VERSION)" -m "Release v$(VERSION)"
	@echo "$(GREEN)Release tag v$(VERSION) created$(NC)"

.PHONY: release
release: release-prepare test-docker tag-release ## Create a full release (prepare, test, tag)
	@echo "$(GREEN)Release v$(VERSION) completed!$(NC)"
	@echo "$(BLUE)Users can now install with:$(NC)"
	@echo "devpod provider add $(GITHUB_REPO)"
	@echo ""
	@echo "$(YELLOW)Don't forget to:$(NC)"
	@echo "  1. Upload binaries from $(BUILD_DIR)/ to GitHub releases"
	@echo "  2. Update the release description"
	@echo "  3. Mark as latest release"
	@echo "  4. Push the tag to GitHub: git push origin v$(VERSION)"

##@ Utilities

.PHONY: list-workspaces
list-workspaces: ## List all DevPod workspaces
	@echo "$(BLUE)DevPod workspaces:$(NC)"
	$(DEVPOD_CMD) list

.PHONY: list-providers
list-providers: ## List all DevPod providers
	@echo "$(BLUE)DevPod providers:$(NC)"
	$(DEVPOD_CMD) provider list

.PHONY: logs
logs: ## Show provider logs
	@echo "$(BLUE)Provider logs:$(NC)"
	$(DEVPOD_CMD) provider logs $(PROVIDER_NAME)-dev || $(DEVPOD_CMD) provider logs $(PROVIDER_NAME)

.PHONY: debug-env
debug-env: ## Show debug environment information
	@echo "$(BLUE)Debug Environment Information:$(NC)"
	@echo "DevPod version: $$($(DEVPOD_CMD) version)"
	@echo "DevPod command: $(DEVPOD_CMD)"
	@echo "Provider file: $(PROVIDER_FILE)"
	@echo "Provider version: $(VERSION)"
	@echo "GitHub repo: $(GITHUB_REPO)"
	@echo "Test workspace: $(TEST_WORKSPACE)"

.PHONY: clean-env
clean-env: ## Remove .env file
	@echo "$(YELLOW)Removing .env file...$(NC)"
	@if [ -f .env ]; then \
		rm .env; \
		echo "$(GREEN)✓ .env file removed$(NC)"; \
	else \
		echo "$(YELLOW)⚠ .env file not found$(NC)"; \
	fi

.PHONY: clean
clean: cleanup-test uninstall ## Clean up everything (workspaces, providers)
	@echo "$(GREEN)Cleanup completed$(NC)"

##@ Help

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\n$(BLUE)Dokploy DevPod Provider Makefile$(NC)\n\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n$(YELLOW)%s$(NC)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(BLUE)Quick Start:$(NC)"
	@echo "  0. make setup             # Install required tools (yq, jq, shellcheck)"
	@echo "  1. make install-local     # Install provider locally"
	@echo "  2. make setup-env         # Create .env file from template"
	@echo "  3. Edit .env file with your Dokploy settings"
	@echo "  4. make configure-env     # Configure provider from .env file"
	@echo "  5. make test-docker       # Test with a Docker workspace"
	@echo "  6. make release-prepare   # Prepare release (builds + checksums)"
	@echo "  7. make release           # Create a release when ready"
	@echo ""
	@echo "$(BLUE)Alternative Configuration:$(NC)"
	@echo "  make configure            # Interactive configuration (instead of steps 2-4)"
	@echo ""
	@echo "$(BLUE)Workspace Management:$(NC)"
	@echo "  make cleanup-workspaces   # Delete all workspaces using this provider"
	@echo "  make force-reinstall      # Delete workspaces and reinstall provider"
	@echo "  make force-uninstall      # Delete workspaces and remove provider"
	@echo ""
	@echo "$(BLUE)Tool Installation:$(NC)"
	@echo "  make setup                # Auto-install all required tools"
	@echo "  make check-tools          # Check which tools are installed"
	@echo "  brew install yq jq shellcheck  # Manual install (macOS)"
	@echo ""
	@echo "$(BLUE)Release Management:$(NC)"
	@echo "  make build-all            # Build binaries for all platforms"
	@echo "  make checksums            # Generate SHA256 checksums"
	@echo "  make show-checksums       # Display generated checksums"
	@echo "  make verify-checksums     # Verify checksums against binaries"
	@echo "  make release-prepare      # Complete release preparation (URLs + checksums)"
	@echo "  make version-bump-patch   # Bump patch version (0.1.0 -> 0.1.1)"
	@echo "  make restore-provider     # Restore provider.yaml from backup"
	@echo ""
	@echo "$(BLUE)GitHub Installation:$(NC)"
	@echo "  devpod provider add $(GITHUB_REPO)"
	@echo ""

.PHONY: nuclear-cleanup-providers
nuclear-cleanup-providers: ## Nuclear option: force delete providers using all methods
	@echo "$(RED)⚠ Nuclear cleanup: Attempting to force delete providers$(NC)"
	@echo "$(YELLOW)This will try multiple deletion methods...$(NC)"
	@for provider in $(PROVIDER_NAME) $(PROVIDER_NAME)-dev; do \
		if $(DEVPOD_CMD) provider list --output json 2>/dev/null | jq -e "has(\"$$provider\")" >/dev/null 2>&1; then \
			echo "$(BLUE)Nuclear deletion of $$provider...$(NC)"; \
			$(DEVPOD_CMD) provider delete $$provider --ignore-not-found 2>/dev/null || true; \
			$(DEVPOD_CMD) provider delete $$provider --ignore-not-found 2>/dev/null || true; \
			$(DEVPOD_CMD) provider delete $$provider 2>/dev/null || true; \
		fi; \
	done
	@echo "$(GREEN)Nuclear cleanup completed$(NC)"

.PHONY: force-uninstall
force-uninstall: cleanup-workspaces ## Force remove provider and all its workspaces
	@echo "$(YELLOW)Force removing provider...$(NC)"
	@echo "$(BLUE)Current providers before deletion:$(NC)"
	@$(DEVPOD_CMD) provider list || true
	@echo ""
	@echo "$(YELLOW)Attempting to delete providers (with retries)...$(NC)"
	@for attempt in 1 2 3; do \
		echo "$(BLUE)Deletion attempt $$attempt/3...$(NC)"; \
		for provider in $(PROVIDER_NAME) $(PROVIDER_NAME)-dev; do \
			if $(DEVPOD_CMD) provider list --output json 2>/dev/null | jq -e "has(\"$$provider\")" >/dev/null 2>&1; then \
				echo "$(BLUE)Deleting $$provider...$(NC)"; \
				if $(DEVPOD_CMD) provider delete $$provider 2>/dev/null; then \
					echo "$(GREEN)✓ Successfully deleted $$provider$(NC)"; \
				else \
					echo "$(RED)Failed to delete $$provider on attempt $$attempt$(NC)"; \
				fi; \
			fi; \
		done; \
		if ! $(DEVPOD_CMD) provider list --output json 2>/dev/null | jq -e "has(\"$(PROVIDER_NAME)\") or has(\"$(PROVIDER_NAME)-dev\")" >/dev/null 2>&1; then \
			echo "$(GREEN)✓ All providers deleted successfully$(NC)"; \
			break; \
		fi; \
		if [ $$attempt -lt 3 ]; then \
			echo "$(YELLOW)Retrying in 3 seconds...$(NC)"; \
			sleep 3; \
			echo "$(BLUE)Re-checking for remaining workspaces...$(NC)"; \
			$(MAKE) cleanup-workspaces; \
		fi; \
	done
	@if $(DEVPOD_CMD) provider list --output json 2>/dev/null | jq -e "has(\"$(PROVIDER_NAME)\") or has(\"$(PROVIDER_NAME)-dev\")" >/dev/null 2>&1; then \
		echo "$(RED)⚠ Provider deletion failed after 3 attempts$(NC)"; \
		echo "$(YELLOW)Attempting nuclear cleanup...$(NC)"; \
		$(MAKE) nuclear-cleanup-providers; \
	fi
	@echo "$(BLUE)Current providers after deletion:$(NC)"
	@$(DEVPOD_CMD) provider list || true
	@echo ""
	@echo "$(GREEN)Provider force removal completed$(NC)"

.PHONY: reinstall
reinstall: ## Reinstall the provider locally
	@echo "$(BLUE)Reinstalling provider...$(NC)"
	@echo "$(YELLOW)Checking for active workspaces...$(NC)"
	@if $(DEVPOD_CMD) list --output json 2>/dev/null | jq -r '.[] | select(.provider.name == "$(PROVIDER_NAME)" or .provider.name == "$(PROVIDER_NAME)-dev") | .id' 2>/dev/null | grep -q .; then \
		echo "$(RED)⚠ Active workspaces found. Use one of these options:$(NC)"; \
		echo "  make force-reinstall     # Delete workspaces and reinstall"; \
		echo "  make cleanup-workspaces  # Just delete workspaces"; \
		echo "  devpod delete <name>     # Delete specific workspace"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Removing existing provider...$(NC)"
	-$(DEVPOD_CMD) provider delete $(PROVIDER_NAME)-dev --ignore-not-found 2>/dev/null || true
	-$(DEVPOD_CMD) provider delete $(PROVIDER_NAME) --ignore-not-found 2>/dev/null || true
	@echo "$(YELLOW)Installing provider...$(NC)"
	@$(MAKE) install-local 