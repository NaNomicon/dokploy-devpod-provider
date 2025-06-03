# Dokploy DevPod Provider Makefile
# Best practices for developing and managing a custom DevPod provider

# Variables
PROVIDER_NAME := dokploy
PROVIDER_FILE := provider.yaml
VERSION := $(shell grep '^version:' $(PROVIDER_FILE) | sed 's/version: *//')
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
	@workspaces=$$($(DEVPOD_CMD) list --output json 2>/dev/null | jq -r '.[] | select(.provider == "$(PROVIDER_NAME)" or .provider == "$(PROVIDER_NAME)-dev") | .id' 2>/dev/null || true); \
	if [ -n "$$workspaces" ]; then \
		echo "$(BLUE)Found workspaces to clean up:$(NC)"; \
		for ws in $$workspaces; do \
			echo "  - $$ws"; \
		done; \
		echo "$(YELLOW)Deleting workspaces with --force flag...$(NC)"; \
		for ws in $$workspaces; do \
			echo "$(BLUE)Deleting workspace: $$ws$(NC)"; \
			if $(DEVPOD_CMD) delete $$ws --force --debug; then \
				echo "$(GREEN)✓ Successfully deleted: $$ws$(NC)"; \
			else \
				echo "$(RED)✗ Failed to delete: $$ws$(NC)"; \
				echo "$(YELLOW)Trying alternative deletion methods...$(NC)"; \
				$(DEVPOD_CMD) stop $$ws --force 2>/dev/null || true; \
				$(DEVPOD_CMD) delete $$ws --force 2>/dev/null || true; \
			fi; \
		done; \
	else \
		echo "$(GREEN)No workspaces found for this provider$(NC)"; \
	fi

.PHONY: force-uninstall
force-uninstall: cleanup-workspaces ## Force remove provider and all its workspaces
	@echo "$(YELLOW)Force removing provider...$(NC)"
	@echo "$(BLUE)Current providers before deletion:$(NC)"
	@$(DEVPOD_CMD) provider list || true
	@echo "$(YELLOW)Attempting to delete providers...$(NC)"
	@if $(DEVPOD_CMD) provider list | grep -q "$(PROVIDER_NAME)-dev"; then \
		echo "$(BLUE)Deleting $(PROVIDER_NAME)-dev...$(NC)"; \
		$(DEVPOD_CMD) provider delete $(PROVIDER_NAME)-dev --debug || true; \
	fi
	@if $(DEVPOD_CMD) provider list | grep -q "$(PROVIDER_NAME)"; then \
		echo "$(BLUE)Deleting $(PROVIDER_NAME)...$(NC)"; \
		$(DEVPOD_CMD) provider delete $(PROVIDER_NAME) --debug || true; \
	fi
	@echo "$(BLUE)Current providers after deletion:$(NC)"
	@$(DEVPOD_CMD) provider list || true
	@echo "$(GREEN)Provider force removal completed$(NC)"

.PHONY: reinstall
reinstall: ## Reinstall the provider locally
	@echo "$(BLUE)Reinstalling provider...$(NC)"
	@echo "$(YELLOW)Checking for active workspaces...$(NC)"
	@if $(DEVPOD_CMD) list --output json 2>/dev/null | jq -r '.[].provider' 2>/dev/null | grep -q "$(PROVIDER_NAME)"; then \
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

.PHONY: force-reinstall
force-reinstall: force-uninstall install-local ## Force reinstall provider (deletes all workspaces)

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
	@for ws in $$($(DEVPOD_CMD) list --output json 2>/dev/null | jq -r '.[].id' 2>/dev/null | grep -E "test-workspace|vscode-remote-try" || true); do \
		echo "$(BLUE)Deleting test workspace: $$ws$(NC)"; \
		$(DEVPOD_CMD) delete $$ws --force --ignore-not-found 2>/dev/null || true; \
	done
	@echo "$(GREEN)Test workspaces cleaned up$(NC)"

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

.PHONY: version-bump-patch
version-bump-patch: ## Bump patch version (x.y.Z)
	@echo "$(BLUE)Bumping patch version...$(NC)"
	@current_version=$(VERSION); \
	new_version=$$(echo $$current_version | awk -F. '{$$3++; print $$1"."$$2"."$$3}'); \
	sed -i.bak "s/version: $$current_version/version: $$new_version/" $(PROVIDER_FILE); \
	rm -f $(PROVIDER_FILE).bak; \
	echo "$(GREEN)Version bumped from $$current_version to $$new_version$(NC)"

.PHONY: version-bump-minor
version-bump-minor: ## Bump minor version (x.Y.z)
	@echo "$(BLUE)Bumping minor version...$(NC)"
	@current_version=$(VERSION); \
	new_version=$$(echo $$current_version | awk -F. '{$$2++; $$3=0; print $$1"."$$2"."$$3}'); \
	sed -i.bak "s/version: $$current_version/version: $$new_version/" $(PROVIDER_FILE); \
	rm -f $(PROVIDER_FILE).bak; \
	echo "$(GREEN)Version bumped from $$current_version to $$new_version$(NC)"

.PHONY: version-bump-major
version-bump-major: ## Bump major version (X.y.z)
	@echo "$(BLUE)Bumping major version...$(NC)"
	@current_version=$(VERSION); \
	new_version=$$(echo $$current_version | awk -F. '{$$1++; $$2=0; $$3=0; print $$1"."$$2"."$$3}'); \
	sed -i.bak "s/version: $$current_version/version: $$new_version/" $(PROVIDER_FILE); \
	rm -f $(PROVIDER_FILE).bak; \
	echo "$(GREEN)Version bumped from $$current_version to $$new_version$(NC)"

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
	git push origin main
	git push origin "v$(VERSION)"
	@echo "$(GREEN)Release tag v$(VERSION) created and pushed$(NC)"

.PHONY: release
release: validate test-docker tag-release ## Create a full release (validate, test, tag)
	@echo "$(GREEN)Release v$(VERSION) completed!$(NC)"
	@echo "$(BLUE)Users can now install with:$(NC)"
	@echo "devpod provider add $(GITHUB_REPO)"

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
	@echo "  6. make release           # Create a release when ready"
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
	@echo "$(BLUE)GitHub Installation:$(NC)"
	@echo "  devpod provider add $(GITHUB_REPO)"
	@echo "" 