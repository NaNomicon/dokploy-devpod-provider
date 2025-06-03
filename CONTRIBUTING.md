# Contributing to Dokploy DevPod Provider

Thank you for your interest in contributing to the Dokploy DevPod Provider! This guide will help you get started with development, testing, and submitting contributions.

## 🎯 Project Overview

This project provides a DevPod provider that integrates with Dokploy's deployment platform. The provider enables developers to create and manage development workspaces using Dokploy's infrastructure.

### Key Components

- **`provider.yaml`**: Main provider configuration file (Machine Provider pattern)
- **`Makefile`**: Comprehensive development tooling and automation
- **API Integration**: Handles communication with Dokploy's REST API
- **Lifecycle Management**: Manages workspace creation, deployment, and cleanup
- **Workspace Management**: Robust handling of active workspaces and cleanup
- **Cross-Platform Support**: Works on macOS, Linux, and Windows

### Architecture: Machine Provider Design

This provider implements DevPod's **Machine Provider** pattern with a two-layer architecture:

#### Layer 1: Machine Infrastructure (Provider Managed)

- **Purpose**: Provides the base environment where DevPod agent runs
- **Image**: Configured via `MACHINE_IMAGE` option (e.g., `ubuntu:22.04`)
- **Responsibility**: Basic OS, SSH server, Docker runtime
- **Managed By**: Dokploy provider

#### Layer 2: Development Environment (DevPod Managed)

- **Purpose**: Actual development container with tools and dependencies
- **Image**: Defined in `.devcontainer/devcontainer.json` or workspace configuration
- **Responsibility**: Development tools, language runtimes, project dependencies
- **Managed By**: DevPod agent

## 🚀 Getting Started

### Prerequisites

#### Required Tools

- [DevPod](https://devpod.sh/) installed locally (CLI or Desktop App)
- Access to a Dokploy instance for testing
- Basic knowledge of YAML, shell scripting, and REST APIs

#### Development Tools (Auto-installed by Makefile)

- `yq` - YAML processing
- `jq` - JSON processing
- `shellcheck` - Shell script linting
- `curl` - HTTP requests
- `docker` - Container management

### Quick Development Setup

```bash
# 1. Fork and clone the repository
git clone https://github.com/your-username/dokploy-devpod-provider.git
cd dokploy-devpod-provider

# 2. Install required development tools
make setup

# 3. Set up environment configuration
make setup-env
# Edit .env file with your Dokploy configuration

# 4. Install provider locally for development
make install-local

# 5. Configure provider from .env file
make configure-env

# 6. Test the provider
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
MACHINE_TYPE=small
MACHINE_IMAGE=ubuntu:22.04
AGENT_PATH=/opt/devpod/agent
```

**Important**: Never commit your `.env` file to version control. It's automatically gitignored.

## 🏗️ Development Workflow

### Understanding the Provider Structure

The `provider.yaml` file follows the DevPod Machine Provider format:

```yaml
name: dokploy # Provider name
version: v0.1.0 # Version
description: "Dokploy provider for DevPod - Create and manage development machines via Dokploy"
icon: "https://raw.githubusercontent.com/Dokploy/dokploy/refs/heads/canary/apps/dokploy/logo.png"

options: # Configuration options
  DOKPLOY_SERVER_URL: # Required
  DOKPLOY_API_TOKEN: # Required, password field
  DOKPLOY_PROJECT_NAME: # Optional
  DOKPLOY_SERVER_ID: # Optional
  MACHINE_TYPE: # Optional, enum: small/medium/large
  MACHINE_IMAGE: # Optional, default: ubuntu:22.04
  AGENT_PATH: # Optional

agent: # DevPod agent configuration
  path: ${AGENT_PATH}
  driver: docker
  inactivityTimeout: 10m

exec: # Provider commands
  init: |- # Initialize provider (validate connection)
    # Validate Dokploy API connection and authentication

  create: |- # Create machine
    # Create Dokploy application with machine image
    # Set up SSH access and return connection details

  start: |- # Start machine
    # Deploy/start the application

  stop: |- # Stop machine
    # Stop the application

  delete: |- # Delete machine
    # Clean up resources

  status: |- # Get machine status
    # Return machine status (Running/Stopped/Busy/NotFound)

  command: |- # Execute commands via SSH
    # Execute commands in the machine via SSH
```

### Image Management in Development

When developing and testing the provider, understand the two-layer architecture:

#### Machine Image (Provider Layer)

- **Purpose**: Base infrastructure where DevPod agent runs
- **Configuration**: `MACHINE_IMAGE` option (default: `ubuntu:22.04`)
- **Best Practices**: Use lightweight, well-supported images
- **Testing**: Test with different machine images to ensure compatibility

#### Development Image (DevPod Layer)

- **Purpose**: Actual development environment with tools
- **Configuration**: `.devcontainer/devcontainer.json` in test projects
- **Best Practices**: Use official DevContainer images when possible
- **Testing**: Test with various development scenarios

#### Testing Different Image Combinations

```bash
# Test with default machine image
make test-docker

# Test with Alpine (lightweight)
devpod provider set-options dokploy-dev --option MACHINE_IMAGE=alpine:latest
make test-docker

# Test with different development scenarios
devpod up https://github.com/microsoft/vscode-remote-try-node.git --provider dokploy-dev
devpod up https://github.com/microsoft/vscode-remote-try-python.git --provider dokploy-dev
```

### Makefile-Based Development

The project includes a comprehensive Makefile with 30+ commands for development:

#### Installation Management

```bash
make install-local     # Install provider locally for development
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
make configure-optional # Configure optional settings
make show-config       # Display current configuration
make clean-env         # Remove .env file
```

#### Testing Suite

```bash
make test-docker       # Test with Docker workspace
make test-git          # Test with Git repository
make test-lifecycle    # Complete lifecycle testing
make test-ssh          # SSH connection testing
make cleanup-test      # Clean up test workspaces
```

#### Validation and Quality

```bash
make validate          # YAML syntax and structure validation
make lint             # Shell script linting with shellcheck
make check-tools      # Verify required tools are installed
make check-devpod     # Check DevPod CLI availability
```

#### Workspace Management

```bash
make list-workspaces   # List all workspaces
make cleanup-workspaces # Clean up all provider workspaces
make list-providers    # List all DevPod providers
```

#### Tool Management

```bash
make setup            # Auto-install required tools (yq, jq, shellcheck)
make check-tools      # Check which tools are installed
make install-devpod-cli # Install DevPod CLI separately
```

#### Release Management

```bash
make version-check     # Check current version
make version-bump-patch # Bump patch version
make version-bump-minor # Bump minor version
make version-bump-major # Bump major version
make tag-release      # Create and push git tag
make release          # Full release process
```

#### Utilities

```bash
make debug-env        # Show debug information
make logs            # Show provider logs
make docs            # Generate documentation
make help            # Show all available commands
```

### API Authentication

The provider uses Bearer token authentication for all Dokploy API endpoints:

```bash
# All endpoints use the same authentication pattern
curl -H "Authorization: Bearer ${DOKPLOY_API_TOKEN}" \
  "${DOKPLOY_SERVER_URL}/api/endpoint"
```

### Error Handling Best Practices

1. **Validate inputs early**
2. **Provide clear error messages with actionable advice**
3. **Clean up resources on failure**
4. **Exit with appropriate codes**

Example:

```bash
if [ -z "${DOKPLOY_SERVER_URL}" ] || [ -z "${DOKPLOY_API_TOKEN}" ]; then
  echo "Error: DOKPLOY_SERVER_URL and DOKPLOY_API_TOKEN are required"
  exit 1
fi

# Test connection to Dokploy API
echo "Testing Dokploy API connection..."
if ! curl -s -f "${DOKPLOY_SERVER_URL}/api/settings.health" \
  -H "Authorization: Bearer ${DOKPLOY_API_TOKEN}" >/dev/null 2>&1; then
  echo "Error: Cannot connect to Dokploy server or invalid API token"
  echo "Please check:"
  echo "1. Server URL: ${DOKPLOY_SERVER_URL}"
  echo "2. API token is valid and generated from Settings > Profile > API/CLI"
  exit 1
fi
echo "✓ Dokploy API connection successful"
```

## 🧪 Testing

### Automated Testing with Makefile

The Makefile provides comprehensive testing capabilities:

```bash
# Run all tests
make test-lifecycle

# Test specific scenarios
make test-docker      # Test with Docker image workspace
make test-git         # Test with Git repository workspace
make test-ssh         # Test SSH connection

# Validate configuration
make validate         # YAML syntax and structure
make lint            # Shell script linting
```

### Manual Testing

#### 1. Test Provider Installation

```bash
make install-local
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

#### 2. Named Workspace

```bash
devpod up my-test-workspace --provider dokploy-dev
```

#### 3. Error Scenarios

Test with:

- Invalid API token
- Unreachable Dokploy server
- Invalid workspace names
- Network timeouts

### Workspace Management Testing

#### Active Workspace Scenarios

```bash
# Create a workspace
devpod up test-workspace --provider dokploy-dev

# Try to reinstall (should detect active workspace)
make reinstall

# Force reinstall (should handle active workspace)
make force-reinstall

# Clean up specific workspaces
make cleanup-workspaces
```

### Debugging

#### Enable Debug Mode

```bash
devpod up test-workspace --provider dokploy-dev --debug
```

#### Check Logs

```bash
# Provider logs
make logs

# Workspace logs
devpod logs test-workspace --follow

# Debug environment
make debug-env
```

#### Development Debugging

```bash
# Show current configuration
make show-config

# Validate provider
make validate

# Check tool availability
make check-tools
```

## 📝 Code Style and Standards

### Shell Script Guidelines

1. **Use strict error handling**

   ```bash
   set -e  # Exit on error (already included in provider scripts)
   ```

2. **Quote variables**

   ```bash
   echo "Server URL: ${DOKPLOY_SERVER_URL}"
   ```

3. **Check for required tools**

   ```bash
   if ! command -v curl >/dev/null 2>&1; then
     echo "Error: curl is required"
     exit 1
   fi
   ```

4. **Use meaningful variable names**

   ```bash
   APP_ID=$(echo "$RESPONSE" | jq -r '.applicationId // .id // empty')
   ```

5. **Provide helpful error messages**
   ```bash
   if [ -z "$APP_ID" ]; then
     echo "Error: Could not get application ID from response"
     echo "Response: $RESPONSE"
     exit 1
   fi
   ```

### YAML Guidelines

1. **Consistent indentation** (2 spaces)
2. **Clear descriptions** for all options
3. **Proper escaping** for shell scripts
4. **Logical organization** of sections
5. **Use password field** for sensitive options

### Makefile Guidelines

1. **Use colored output** for better UX
2. **Provide help text** for all targets
3. **Handle errors gracefully**
4. **Support cross-platform** operations
5. **Use .PHONY** for non-file targets

### Documentation

1. **Comment complex logic**
2. **Update README** for new features
3. **Include examples** for new options
4. **Document breaking changes**
5. **Update tutorial and journey** for significant changes

## 🔄 Contribution Process

### 1. Issue Discussion

Before starting work:

- Check existing issues
- Create an issue for new features
- Discuss approach with maintainers
- Consider backward compatibility

### 2. Development

1. **Create a feature branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Set up development environment**

   ```bash
   make setup
   make setup-env
   # Edit .env file
   make install-local
   ```

3. **Make your changes**

   - Follow code style guidelines
   - Add appropriate error handling
   - Update documentation
   - Add Makefile targets if needed

4. **Test thoroughly**

   ```bash
   make validate
   make lint
   make test-lifecycle
   make test-docker
   make test-git
   ```

### 3. Submission

1. **Commit your changes**

   ```bash
   git add .
   git commit -m "feat: add support for custom domains"
   ```

2. **Push to your fork**

   ```bash
   git push origin feature/your-feature-name
   ```

3. **Create a Pull Request**
   - Use descriptive title
   - Include detailed description
   - Reference related issues
   - Add testing instructions
   - Include Makefile command examples

### Commit Message Format

Use conventional commits:

```
type(scope): description

feat: add new feature
fix: fix bug
docs: update documentation
test: add tests
refactor: refactor code
chore: update build tools
```

## 🛠️ Workspace Management

### Understanding Active Workspace Issues

DevPod prevents provider deletion when workspaces are using it:

```bash
$ devpod provider delete dokploy
fatal cannot delete provider 'dokploy', because workspace 'my-project' is still using it
```

### Development Solutions

The Makefile provides several solutions:

```bash
# Check for active workspaces before operations
make reinstall         # Safe reinstall with workspace checking

# Handle active workspaces automatically
make force-reinstall   # Delete workspaces and reinstall
make force-uninstall   # Delete workspaces and remove provider

# Manual workspace management
make cleanup-workspaces # Delete all workspaces for this provider
make list-workspaces   # List all workspaces
```

### Testing Workspace Management

```bash
# Create test workspace
devpod up test-workspace --provider dokploy-dev

# Test workspace detection
make reinstall  # Should detect and warn about active workspace

# Test force operations
make force-reinstall  # Should handle workspace cleanup automatically

# Test cleanup
make cleanup-workspaces  # Should clean up all provider workspaces
```

## 🚀 Release Process

### Version Management

```bash
# Check current version
make version-check

# Bump version
make version-bump-patch  # For bug fixes
make version-bump-minor  # For new features
make version-bump-major  # For breaking changes

# Create release
make tag-release        # Create and push git tag
make release           # Full release process
```

### GitHub Releases

The project follows DevPod community provider patterns:

1. **Create a release** on GitHub with a version tag (e.g., `v0.1.0`)
2. **Attach the `provider.yaml`** file to the release
3. **DevPod will automatically** find and download the provider from the latest release

### Release Checklist

- [ ] Update version in `provider.yaml`
- [ ] Test all Makefile commands
- [ ] Run full test suite
- [ ] Update documentation
- [ ] Create GitHub release
- [ ] Test installation from GitHub

## 🐛 Reporting Issues

### Bug Reports

Include:

- DevPod version (`devpod version`)
- Dokploy version
- Provider configuration (sanitized, no tokens)
- Steps to reproduce
- Expected vs actual behavior
- Error messages and logs
- Output of `make debug-env`

### Feature Requests

Include:

- Use case description
- Proposed solution
- Alternative solutions considered
- Additional context
- Impact on existing functionality

### Workspace Management Issues

For workspace-related issues, include:

- Output of `devpod list`
- Output of `make list-workspaces`
- Workspace creation/deletion logs
- Provider logs (`make logs`)

## 📚 Development Resources

### Dokploy Resources

- [Dokploy Documentation](https://docs.dokploy.com/)
- [Dokploy GitHub](https://github.com/Dokploy/dokploy)
- [Dokploy API Reference](https://docs.dokploy.com/docs/api)

### DevPod Resources

- [DevPod Documentation](https://devpod.sh/docs)
- [Provider Development Guide](https://devpod.sh/docs/developing-providers/quickstart)
- [DevPod GitHub](https://github.com/loft-sh/devpod)
- [Community Providers](https://devpod.sh/docs/managing-providers/add-provider#community-providers)

### Community Provider Examples

- [Hetzner Provider](https://github.com/mrsimonemms/devpod-provider-hetzner)
- [Cloudbit Provider](https://github.com/cloudbit-ch/devpod-provider-cloudbit)
- [Scaleway Provider](https://github.com/dirien/devpod-provider-scaleway)
- [Flow Provider](https://github.com/flowswiss/devpod-provider-flow)

### Project Resources

- [Tutorial: Complete Guide](tutorial.md)
- [Development Journey](journey.md)

## 🤝 Community Guidelines

1. **Be respectful** and inclusive
2. **Help others** learn and contribute
3. **Share knowledge** and best practices
4. **Provide constructive feedback**
5. **Follow the code of conduct**
6. **Test thoroughly** before submitting
7. **Document your changes**

## 🏆 Recognition

Contributors will be recognized in:

- README acknowledgments
- Release notes
- Project documentation
- GitHub contributors list

## 🔧 Development Tips

### Cross-Platform Development

The Makefile supports multiple platforms:

- **macOS**: Uses Homebrew for tool installation
- **Ubuntu/Debian**: Uses APT package manager
- **RHEL/CentOS**: Uses YUM package manager
- **Manual**: Provides fallback instructions

### DevPod CLI Issues

If you encounter DevPod CLI issues:

```bash
# Check DevPod CLI availability
make check-devpod

# Install DevPod CLI separately
make install-devpod-cli

# Add DevPod app to PATH (macOS)
export PATH="/Applications/DevPod.app/Contents/MacOS:$PATH"
```

### Environment Management

```bash
# Create environment for different scenarios
cp .env .env.development
cp .env .env.staging

# Load specific environment
source .env.development
make configure-env
```

### Debugging Provider Issues

```bash
# Show comprehensive debug information
make debug-env

# Validate provider configuration
make validate

# Check all tools
make check-tools

# Test with verbose output
devpod up test-workspace --provider dokploy-dev --debug
```

Thank you for contributing to the Dokploy DevPod Provider! 🎉

Your contributions help make development environments more accessible and powerful for developers worldwide.
