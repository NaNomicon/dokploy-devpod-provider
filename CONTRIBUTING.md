# Contributing to Dokploy DevPod Provider

Thank you for your interest in contributing to the Dokploy DevPod Provider! This guide will help you get started with development, testing, and submitting contributions.

## 🎯 Project Overview

This project provides a DevPod provider that integrates with Dokploy's deployment platform using **Docker Compose services**. The provider enables developers to create and manage development workspaces using Dokploy's infrastructure with a high-performance Go binary helper.

### Key Components

- **`provider.yaml`**: Main provider configuration file (Machine Provider pattern)
- **`dokploy-provider` binary**: Go CLI binary for all operations
- **Docker Compose Integration**: Creates and manages Docker Compose services in Dokploy
- **SSH Infrastructure**: Root-based SSH access for maximum DevPod compatibility
- **Makefile**: Comprehensive development tooling and automation
- **Cross-Platform Support**: Works on macOS, Linux, and Windows

### Architecture: Docker Compose Services

This provider implements DevPod's **Machine Provider** pattern using Dokploy's **Docker Compose services**:

#### Layer 1: Docker Compose Infrastructure (Provider Managed)

- **Purpose**: Provides the base Docker Compose container where DevPod agent runs
- **Base Image**: `cruizba/ubuntu-dind:latest` (Docker-in-Docker with Ubuntu)
- **SSH Setup**: Root-based SSH access with key injection
- **Responsibility**: Docker daemon, SSH server, basic OS tools
- **Managed By**: Dokploy provider via Docker Compose API

#### Layer 2: Development Environment (DevPod Managed)

- **Purpose**: Actual development container with tools and dependencies
- **Image**: Defined in `.devcontainer/devcontainer.json` or workspace configuration
- **Responsibility**: Development tools, language runtimes, project dependencies
- **Managed By**: DevPod agent inside the Docker Compose container

#### Binary Helper Architecture

The provider uses a **Go CLI binary** instead of shell scripts:

```
DevPod → Binary Helper → Dokploy API → Docker Compose Service → SSH Access
```

**Commands**:

- `dokploy-provider init` - Validate configuration and connectivity
- `dokploy-provider create` - Create Docker Compose service with SSH setup
- `dokploy-provider delete` - Remove Docker Compose service
- `dokploy-provider start` - Start Docker Compose service
- `dokploy-provider stop` - Stop Docker Compose service
- `dokploy-provider status` - Get service status
- `dokploy-provider command` - Execute commands via SSH

## 🚀 Getting Started

### Prerequisites

#### Required Tools

- [DevPod](https://devpod.sh/) installed locally (CLI or Desktop App)
- Access to a Dokploy instance for testing
- Go 1.22+ for binary development
- Docker for testing
- Basic knowledge of Go, Docker Compose, and REST APIs

#### Development Tools (Auto-installed by Makefile)

- `yq` - YAML processing
- `jq` - JSON processing
- `curl` - HTTP requests
- `docker` - Container management

### Quick Development Setup

```bash
# 1. Fork and clone the repository
git clone https://github.com/your-username/dokploy-devpod-provider.git
cd dokploy-devpod-provider

# 2. Install required development tools
make setup

# 3. Build the binary helper
make build

# 4. Set up environment configuration
make setup-env
# Edit .env file with your Dokploy configuration

# 5. Install provider locally for development
make install-dev

# 6. Configure provider from .env file
make configure-env

# 7. Test the provider
make test-docker
```

### Environment Configuration

Create a `.env` file for your development configuration:

```bash
# Required Configuration
DOKPLOY_SERVER_URL=https://your-test-dokploy.com
DOKPLOY_API_TOKEN=your_test_token

# Optional Configuration
DOKPLOY_PROJECT_NAME=devpod-development
DOKPLOY_SERVER_ID=
AGENT_PATH=/opt/devpod/agent
AGENT_DATA_PATH=/opt/devpod/agent-data
INACTIVITY_TIMEOUT=10m
INJECT_GIT_CREDENTIALS=true
INJECT_DOCKER_CREDENTIALS=true
```

**Important**: Never commit your `.env` file to version control. It's automatically gitignored.

## 🏗️ Development Workflow

### Understanding the Provider Structure

The `provider.yaml` file follows the DevPod Machine Provider format with Docker Compose support:

```yaml
name: dokploy # Provider name
version: v0.1.0 # Version
description: "DevPod on Dokploy - Docker Compose services"
icon: "https://raw.githubusercontent.com/Dokploy/dokploy/refs/heads/canary/apps/dokploy/logo.png"

options: # Configuration options
  DOKPLOY_SERVER_URL: # Required
  DOKPLOY_API_TOKEN: # Required, password field
  DOKPLOY_PROJECT_NAME: # Optional
  AGENT_PATH: # Optional, default: /opt/devpod/agent
  AGENT_DATA_PATH: # Optional
  INACTIVITY_TIMEOUT: # Optional
  INJECT_GIT_CREDENTIALS: # Optional
  INJECT_DOCKER_CREDENTIALS: # Optional

agent: # DevPod agent configuration
  path: ${AGENT_PATH}
  dataPath: ${AGENT_DATA_PATH}
  inactivityTimeout: ${INACTIVITY_TIMEOUT}
  injectGitCredentials: ${INJECT_GIT_CREDENTIALS}
  injectDockerCredentials: ${INJECT_DOCKER_CREDENTIALS}
  docker:
    install: true

binaries: # Binary helper distribution
  DOKPLOY_PROVIDER_BINARY:
    - os: linux
      arch: amd64
      path: https://github.com/.../dokploy-provider-linux-amd64

exec: # Provider commands (using binary helper)
  init: ${DOKPLOY_PROVIDER_BINARY} init
  create: ${DOKPLOY_PROVIDER_BINARY} create
  start: ${DOKPLOY_PROVIDER_BINARY} start
  stop: ${DOKPLOY_PROVIDER_BINARY} stop
  delete: ${DOKPLOY_PROVIDER_BINARY} delete
  status: ${DOKPLOY_PROVIDER_BINARY} status
  command: ${DOKPLOY_PROVIDER_BINARY} command
```

### Docker Compose Service Management

The provider creates Docker Compose services in Dokploy with the following characteristics:

#### Container Configuration

- **Base Image**: `cruizba/ubuntu-dind:latest`
- **Privileged Mode**: Enabled for Docker-in-Docker
- **Port Mapping**: SSH port from range 2222-2250
- **SSH Access**: Root user with SSH key injection
- **Docker Daemon**: Full Docker-in-Docker capabilities

#### Setup Process (4 Stages)

1. **Docker Daemon Startup** (30-60s): Start Docker-in-Docker
2. **SSH Installation** (30-60s): Install SSH server and tools
3. **SSH Key Setup** (10-20s): Inject SSH keys for root user
4. **SSH Configuration** (10-20s): Configure SSH daemon for root access

### Binary Helper Development

The provider is implemented as a Go CLI application:

```
├── cmd/                    # CLI commands
│   ├── root.go            # Root command and configuration
│   ├── init.go            # Initialize and validate provider
│   ├── create.go          # Create Docker Compose service
│   ├── delete.go          # Delete Docker Compose service
│   ├── start.go           # Start Docker Compose service
│   ├── stop.go            # Stop Docker Compose service
│   ├── status.go          # Get Docker Compose service status
│   └── command.go         # Execute commands via SSH
├── pkg/                   # Core packages
│   ├── options/           # Configuration management
│   ├── dokploy/           # Dokploy API client (Docker Compose support)
│   ├── client/            # DevPod status types
│   └── ssh/               # SSH client for command execution
├── templates/             # Docker Compose and setup templates
│   ├── docker-compose.yml # Docker Compose template
│   └── setup-root.sh     # Container setup script
├── dist/                  # Built binaries
└── provider.yaml          # DevPod provider configuration
```

### Makefile-Based Development

The project includes a comprehensive Makefile with 30+ commands for development:

#### Binary Management

```bash
make build             # Build binary for current platform
make build-all         # Build for all platforms
make test-build        # Test binary functionality
make install-dev       # Install development provider with local binary
```

#### Installation Management

```bash
make install-dev       # Install development provider (uses local binary)
make install-github    # Install from GitHub repository
make reinstall         # Reinstall (checks for active workspaces)
make force-reinstall   # Force reinstall (handles active workspaces)
make uninstall         # Remove provider
make force-uninstall   # Force remove provider and all workspaces
```

#### Configuration Management

```bash
make setup-env         # Create .env from template
make configure-env     # Configure from .env file
make configure         # Interactive configuration
make show-config       # Display current configuration
make clean-env         # Remove .env file
```

#### Testing Suite

```bash
make test-docker       # Test with Docker workspace
make test-git          # Test with Git repository
make test-compose      # Test Docker Compose functionality
make test-ssh          # SSH connection testing
make test-lifecycle    # Complete lifecycle testing
make cleanup-test      # Clean up test workspaces
```

#### Validation and Quality

```bash
make validate          # YAML syntax and structure validation
make lint             # Go code linting
make check-tools      # Verify required tools are installed
make check-devpod     # Check DevPod CLI availability
```

#### Workspace Management

```bash
make list-workspaces   # List all workspaces
make cleanup-workspaces # Clean up all provider workspaces
make list-providers    # List all DevPod providers
make fix-stuck-workspace # Handle stuck workspaces
```

#### Tool Management

```bash
make setup            # Auto-install required tools
make check-tools      # Check which tools are installed
make install-devpod-cli # Install DevPod CLI separately
```

### API Integration

The provider uses `x-api-key` header authentication for all Dokploy API endpoints:

```go
// Example API call
req.Header.Set("x-api-key", c.apiToken)
req.Header.Set("Content-Type", "application/json")

// Docker Compose API endpoints
POST /api/compose.create     # Create Docker Compose service
POST /api/compose.update     # Update Docker Compose configuration
POST /api/compose.deploy     # Deploy Docker Compose service
POST /api/compose.start      # Start Docker Compose service
POST /api/compose.stop       # Stop Docker Compose service
POST /api/compose.delete     # Delete Docker Compose service
GET  /api/compose.one        # Get Docker Compose service details
```

### Error Handling Best Practices

1. **Validate inputs early in Go functions**
2. **Provide clear error messages with actionable advice**
3. **Clean up Docker Compose resources on failure**
4. **Use structured logging with logrus**

Example:

```go
func (c *Client) CreateCompose(req CreateComposeRequest) (*Compose, error) {
    if req.Name == "" {
        return nil, fmt.Errorf("compose service name is required")
    }

    if req.ProjectID == "" {
        return nil, fmt.Errorf("project ID is required")
    }

    c.logger.Infof("Creating Docker Compose service: %s", req.Name)

    resp, err := c.makeRequest("POST", "/api/compose.create", req)
    if err != nil {
        return nil, fmt.Errorf("failed to create compose service: %w", err)
    }
    defer resp.Body.Close()

    // Handle response...
}
```

## 🧪 Testing

### Automated Testing with Makefile

The Makefile provides comprehensive testing capabilities:

```bash
# Build and test binary
make build
make test-build

# Run all tests
make test-lifecycle

# Test specific scenarios
make test-docker      # Test with Docker image workspace
make test-git         # Test with Git repository workspace
make test-compose     # Test Docker Compose functionality
make test-ssh         # Test SSH connection

# Validate configuration
make validate         # YAML syntax and structure
make lint            # Go code linting
```

### Manual Testing

#### 1. Test Provider Installation

```bash
make build
make install-dev
make show-config
```

#### 2. Test Configuration

```bash
make setup-env
# Edit .env file with test configuration
make configure-env
```

#### 3. Test Workspace Lifecycle

```bash
# Create workspace
devpod up test-ws --provider dokploy-dev --debug

# Check status
devpod status test-ws

# SSH into workspace
devpod ssh test-ws

# Stop workspace
devpod stop test-ws

# Delete workspace
devpod delete test-ws --force
```

### Testing Different Scenarios

#### 1. Git Repository Deployment

```bash
devpod up https://github.com/microsoft/vscode-remote-try-node.git --provider dokploy-dev
```

#### 2. Docker Image Workspace

```bash
devpod up ubuntu:22.04 --provider dokploy-dev
```

#### 3. Error Scenarios

Test with:

- Invalid API token
- Unreachable Dokploy server
- Invalid workspace names
- Network timeouts
- Docker Compose deployment failures

### Docker Compose Testing

#### Test Docker Compose Service Creation

```bash
# Test with verbose logging
DOKPLOY_PROVIDER_DEV=true ./dist/dokploy-provider create --verbose

# Test service status
DOKPLOY_PROVIDER_DEV=true ./dist/dokploy-provider status --verbose

# Test SSH connectivity
DOKPLOY_PROVIDER_DEV=true ./dist/dokploy-provider command --verbose
```

### Debugging

#### Enable Debug Mode

```bash
devpod up test-workspace --provider dokploy-dev --debug
```

#### Binary Debug Mode

```bash
# Enable development mode
export DOKPLOY_PROVIDER_DEV=true

# Run commands with verbose logging
./dist/dokploy-provider create --verbose
./dist/dokploy-provider status --verbose
```

#### Debug Log Locations

- **Status Command**: `/tmp/dokploy-provider-status.log` (development mode)
- **Other Commands**: stderr output
- **API Debugging**: All requests/responses logged with security redaction

## 🔄 Release Process

### Version Management

```bash
make version-check     # Check current version
make version-bump-patch # Bump patch version
make version-bump-minor # Bump minor version
make version-bump-major # Bump major version
```

### Building for Release

```bash
# Build for all platforms
make build-all

# Create release artifacts
make release-artifacts

# Tag and push release
make tag-release
```

### Binary Distribution

The provider distributes binaries for multiple platforms:

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

Binaries are automatically downloaded by DevPod based on the `binaries` section in `provider.yaml`.

## 🤝 Contributing Guidelines

### Code Style

- **Go Code**: Follow standard Go conventions, use `gofmt`
- **YAML**: Use 2-space indentation
- **Shell Scripts**: Use shellcheck for linting
- **Documentation**: Keep README and CONTRIBUTING up to date

### Pull Request Process

1. **Fork** the repository
2. **Create feature branch** from main
3. **Build and test** with `make test-lifecycle`
4. **Update documentation** if needed
5. **Submit pull request** with clear description

### Testing Requirements

- All new features must include tests
- Binary must build successfully for all platforms
- Provider must pass `make test-lifecycle`
- Docker Compose functionality must be tested

### Documentation

- Update README.md for user-facing changes
- Update CONTRIBUTING.md for development changes
- Document new configuration options
- Include examples for new features

## 🐛 Troubleshooting Development

### Common Issues

1. **Binary Build Failures**: Check Go version and dependencies
2. **Docker Compose API Errors**: Verify Dokploy server connectivity
3. **SSH Connection Issues**: Check port mapping and Docker Swarm propagation
4. **Provider Installation Failures**: Use `make force-reinstall`

### Debug Commands

```bash
# Check provider status
make show-config

# Debug binary directly
./dist/dokploy-provider --help
./dist/dokploy-provider init --verbose

# Check Docker Compose services in Dokploy dashboard
# Look for services in the configured project

# Test SSH connectivity manually
nc -z <host> <port>  # Test port availability
```

Remember: The provider uses **Docker Compose services** in Dokploy with **root SSH access** for maximum DevPod compatibility.
