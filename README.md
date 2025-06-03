# Dokploy Provider for DevPod

A custom DevPod provider that enables seamless integration between [DevPod](https://devpod.sh/) and [Dokploy](https://dokploy.com/), allowing developers to create and manage development workspaces using Dokploy's deployment infrastructure.

## 🚀 Overview

This provider bridges DevPod's development environment management with Dokploy's deployment platform, enabling developers to:

- **Create development machines** backed by Dokploy infrastructure
- **SSH-based connection** for seamless development experience
- **Two-layer architecture** separating infrastructure from development concerns
- **Complete lifecycle management** (create, start, stop, delete, status)
- **Robust workspace management** with force operations and cleanup
- **Cross-platform development** with comprehensive tooling

## 🏗️ Architecture: Machine Provider Design

### Understanding the Two-Layer Architecture

This provider implements DevPod's **Machine Provider** pattern with a clear separation of concerns:

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

### Image Management: Two-Layer Approach

DevPod uses a **two-layer image approach** that separates infrastructure from development environments:

```
User Command: devpod up https://github.com/user/react-app.git --provider dokploy

Step 1: Dokploy Provider Creates Machine
┌─────────────────────────────────────┐
│ Dokploy Application                 │
│ Image: ubuntu:22.04 (MACHINE_IMAGE) │
│ + SSH Server                        │
│ + Docker Engine                     │
│ + Basic Unix Tools                  │
└─────────────────────────────────────┘

Step 2: DevPod Connects and Takes Over
┌─────────────────────────────────────┐
│ DevPod Agent (inside machine)       │
│ + Clones Git Repository             │
│ + Reads .devcontainer/devcontainer.json │
│ + Pulls Development Image           │
└─────────────────────────────────────┘

Step 3: Development Container Created
┌─────────────────────────────────────┐
│ Development Container               │
│ Image: node:18 (from devcontainer.json) │
│ + Your Source Code                  │
│ + Node.js Runtime                   │
│ + Development Tools                 │
│ + VS Code Extensions                │
└─────────────────────────────────────┘
```

#### Machine Image Best Practices

**Recommended Machine Images:**

- `ubuntu:22.04` - **Best overall choice** (good compatibility, well-tested)
- `ubuntu:20.04` - Stable, widely supported
- `debian:11` - Lightweight, stable
- `alpine:latest` - Minimal footprint (advanced users)

**Avoid:**

- Development-specific images (e.g., `node:18`) as machine images
- Images without SSH/Docker support
- Very old or deprecated images

#### Development Image Configuration

Use `.devcontainer/devcontainer.json` in your project:

```json
{
  "name": "React Development",
  "image": "mcr.microsoft.com/devcontainers/javascript-node:18",
  "features": {
    "ghcr.io/devcontainers/features/git:1": {}
  },
  "postCreateCommand": "npm install",
  "forwardPorts": [3000]
}
```

### Workflow Example

```bash
# 1. DevPod asks provider to create machine
devpod up my-project --provider dokploy

# 2. Provider creates Dokploy application with ubuntu:22.04
# 3. Provider sets up SSH access and returns connection details
# 4. DevPod connects via SSH and injects agent
# 5. Agent creates development container inside the machine
# 6. Developer gets full development environment
```

## 📋 Prerequisites

### DevPod Requirements

- [DevPod](https://devpod.sh/) installed and configured
- DevPod CLI or Desktop App

**Important Note**: The DevPod Desktop App doesn't automatically add the CLI to your PATH. If you need CLI access, either:

- Install CLI separately (recommended for development)
- Add app bundle to PATH: `export PATH="/Applications/DevPod.app/Contents/MacOS:$PATH"`

### Dokploy Requirements

- A running Dokploy instance (self-hosted or cloud)
- Dokploy API token with appropriate permissions
- Access to Dokploy server URL

### System Requirements

- `curl` and `jq` available on the system
- Docker access on the Dokploy server
- Network connectivity between DevPod and Dokploy server

## 🔧 Installation

### From GitHub (Recommended)

```bash
devpod provider add NaNomicon/dokploy-devpod-provider
devpod provider use dokploy
```

### From Local Development

```bash
git clone https://github.com/NaNomicon/dokploy-devpod-provider.git
cd dokploy-devpod-provider
make install-local
```

### From URL

```bash
devpod provider add https://raw.githubusercontent.com/NaNomicon/dokploy-devpod-provider/main/provider.yaml
```

## ⚙️ Configuration

### Required Options

| Option               | Description                                         | Example                       |
| -------------------- | --------------------------------------------------- | ----------------------------- |
| `DOKPLOY_SERVER_URL` | URL of your Dokploy server                          | `https://dokploy.example.com` |
| `DOKPLOY_API_TOKEN`  | API token from Dokploy Settings > Profile > API/CLI | `your_generated_token`        |

### Optional Options

| Option                 | Description                                                                        | Default             |
| ---------------------- | ---------------------------------------------------------------------------------- | ------------------- |
| `DOKPLOY_PROJECT_NAME` | Name of the project in Dokploy (will be automatically created if it doesn't exist) | `devpod-workspaces` |
| `DOKPLOY_SERVER_ID`    | ID of the server to deploy to                                                      | (uses default)      |
| `MACHINE_TYPE`         | Machine size (small/medium/large)                                                  | `small`             |
| `MACHINE_IMAGE`        | Docker image for the machine                                                       | `ubuntu:22.04`      |
| `AGENT_PATH`           | Path where DevPod agent is injected                                                | `/opt/devpod/agent` |

### Quick Configuration

```bash
devpod provider set-options dokploy \
  --option DOKPLOY_SERVER_URL=https://your-dokploy.com \
  --option DOKPLOY_API_TOKEN=your_api_token_here \
  --option DOKPLOY_PROJECT_NAME=devpod-workspaces
```

### Environment File Configuration (.env)

For easier configuration management, you can use a `.env` file:

#### Using Makefile (Development)

```bash
# 1. Create environment file from template
make setup-env

# 2. Edit .env file with your settings
# 3. Configure provider from .env file
make configure-env
```

#### Manual Configuration

```bash
# Create .env file
cat > .env << EOF
# Required Configuration
DOKPLOY_SERVER_URL=https://your-dokploy.com
DOKPLOY_API_TOKEN=your_api_token_here

# Optional Configuration
DOKPLOY_PROJECT_NAME=devpod-workspaces
DOKPLOY_SERVER_ID=
MACHINE_TYPE=small
MACHINE_IMAGE=ubuntu:22.04
AGENT_PATH=/opt/devpod/agent
EOF

# Configure provider
devpod provider set-options dokploy \
  --option DOKPLOY_SERVER_URL="$(grep DOKPLOY_SERVER_URL .env | cut -d= -f2)" \
  --option DOKPLOY_API_TOKEN="$(grep DOKPLOY_API_TOKEN .env | cut -d= -f2)"
```

## 🎯 Usage

### Creating a Workspace

#### From a Git Repository

```bash
devpod up https://github.com/your-username/your-repo.git --provider dokploy
```

#### With Specific Name

```bash
devpod up my-workspace --provider dokploy
```

### Managing Workspaces

```bash
# List workspaces
devpod list

# Start a workspace
devpod start my-workspace

# Stop a workspace
devpod stop my-workspace

# Delete a workspace
devpod delete my-workspace

# Force delete a workspace (no confirmation)
devpod delete my-workspace --force

# SSH into a workspace
devpod ssh my-workspace
```

### Workspace States

- **Running**: Workspace is active and accessible
- **Stopped**: Workspace is stopped but data is preserved
- **Busy**: Workspace is starting, stopping, or being modified
- **NotFound**: Workspace doesn't exist

## 🔧 Development Workflow

### Prerequisites for Development

```bash
# Install required tools
make setup

# Or manually:
brew install yq jq shellcheck  # macOS
sudo apt install jq shellcheck  # Ubuntu
```

### Quick Start Development

```bash
# 1. Clone the repository
git clone https://github.com/NaNomicon/dokploy-devpod-provider.git
cd dokploy-devpod-provider

# 2. Set up environment
make setup-env
# Edit .env file with your Dokploy configuration

# 3. Install provider locally
make install-local

# 4. Configure provider
make configure-env

# 5. Test the provider
make test-docker
```

### Development Commands

```bash
# Installation Management
make install-local     # Install provider locally
make reinstall         # Reinstall (checks for active workspaces)
make force-reinstall   # Force reinstall (handles active workspaces)
make uninstall         # Remove provider

# Configuration Management
make setup-env         # Create .env from template
make configure-env     # Configure from .env file
make configure         # Interactive configuration
make show-config       # Display current configuration

# Testing
make test-docker       # Test with Docker workspace
make test-git          # Test with Git repository
make test-lifecycle    # Complete lifecycle testing
make test-ssh          # SSH connection testing

# Validation
make validate          # YAML syntax and structure
make lint             # Shell script linting
make check-tools      # Verify required tools

# Workspace Management
make list-workspaces   # List all workspaces
make cleanup-test      # Clean up test workspaces
make cleanup-workspaces # Clean up all provider workspaces

# Utilities
make debug-env         # Show debug information
make help             # Show all available commands
```

## 🔐 Authentication

### Generating Dokploy API Token

1. Log into your Dokploy dashboard
2. Navigate to **Settings > Profile**
3. Go to the **API/CLI** section
4. Click **"Generate Token"**
5. Copy the generated token

### Authentication Methods

The provider uses Bearer token authentication for all Dokploy API endpoints:

```bash
# All endpoints use the same authentication pattern
curl -H "Authorization: Bearer ${TOKEN}" "${URL}/api/endpoint"
```

## 🏗️ How It Works

### Machine Provider Workflow

1. **Initialization**: Validates connection to Dokploy server and API token
2. **Project Management**: Checks if the specified project exists, creates it automatically if it doesn't
3. **Machine Creation**: Creates Dokploy application with machine image in the project
4. **SSH Setup**: Configures SSH server and creates devpod user
5. **Connection Details**: Returns SSH connection information to DevPod
6. **Agent Injection**: DevPod connects via SSH and injects development agent
7. **Development Environment**: Agent creates development container based on project configuration

### Architecture Flow

```
DevPod Client → Dokploy Provider → Dokploy API → Docker Container (Machine)
     ↓              ↓                ↓              ↓
  Commands    Machine Management   Application    SSH + Development Environment
                                   Deployment
```

## 🛠️ Workspace Management

### The Challenge: Active Workspaces

DevPod prevents provider deletion when workspaces are using it:

```bash
$ devpod provider delete dokploy
fatal cannot delete provider 'dokploy', because workspace 'my-project' is still using it
```

### Solutions

#### Option 1: Manual Cleanup

```bash
# List active workspaces
devpod list

# Delete specific workspaces
devpod delete workspace-name --force

# Then delete provider
devpod provider delete dokploy
```

#### Option 2: Using Makefile (Development)

```bash
# Clean up all workspaces for this provider
make cleanup-workspaces

# Force reinstall (deletes workspaces and reinstalls)
make force-reinstall

# Force uninstall (deletes workspaces and removes provider)
make force-uninstall
```

#### Option 3: Automated Cleanup Script

```bash
#!/bin/bash
# cleanup-dokploy-workspaces.sh

echo "Finding workspaces using dokploy provider..."
workspaces=$(devpod list --output json | jq -r '.[] | select(.provider == "dokploy" or .provider == "dokploy-dev") | .id')

if [ -n "$workspaces" ]; then
    echo "Found workspaces to clean up:"
    echo "$workspaces"

    for ws in $workspaces; do
        echo "Deleting workspace: $ws"
        devpod delete "$ws" --force
    done

    echo "Waiting for cleanup to complete..."
    sleep 5

    echo "Deleting provider..."
    devpod provider delete dokploy
else
    echo "No workspaces found using dokploy provider"
fi
```

## 🐛 Troubleshooting

### Common Issues

#### 1. Authentication Errors

```
Error: Cannot connect to Dokploy server or invalid API token
```

**Solutions:**

1. Verify your `DOKPLOY_SERVER_URL` is correct and accessible
2. Ensure your API token is valid and not expired
3. Check that the token has appropriate permissions
4. Regenerate the token if necessary

**Testing Connection:**

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "https://your-dokploy.com/api/settings.health"
```

#### 2. Provider Deletion Blocked by Active Workspaces

```
fatal cannot delete provider 'dokploy', because workspace 'my-project' is still using it
```

**Solutions:**

- Use `make cleanup-workspaces` to clean up all workspaces
- Use `make force-reinstall` for development
- Manually delete workspaces with `devpod delete workspace-name --force`

#### 3. Workspace Name Validation Errors

```
Error: Could not get application ID from response
Response: {"message":"Input validation failed"...}
```

**Root Cause:** The workspace name is empty or invalid.

**Solutions:**

1. Ensure `DEVPOD_WORKSPACE_ID` is properly set
2. Use valid workspace names (alphanumeric, hyphens allowed)
3. Check provider script for variable substitution issues

#### 4. SSH Connection Issues

**Problem:** Can't connect to workspace via SSH.

**Solutions:**

1. Verify SSH server is running in container
2. Check SSH port mapping
3. Ensure SSH keys are properly configured
4. Check firewall settings

**Debugging:**

```bash
# Test SSH connection manually
ssh -o StrictHostKeyChecking=no -p PORT user@host

# Check container SSH service
docker exec container-id service ssh status
```

### Debug Mode

#### Enable Verbose Logging

```bash
devpod up my-workspace --provider dokploy --debug
```

#### Check Provider Logs

```bash
devpod provider logs dokploy
```

#### Check Workspace Logs

```bash
devpod logs my-workspace --follow
```

#### Development Debugging

```bash
# Using Makefile for development
make debug-env

# Check provider configuration
make show-config

# Validate provider
make validate
```

## 🔧 Configuration Examples

### Development Environment

```bash
devpod provider set-options dokploy \
  --option MACHINE_IMAGE=node:18 \
  --option MACHINE_TYPE=small \
  --option DOKPLOY_PROJECT_NAME=dev-environments
```

### Production-like Environment

```bash
devpod provider set-options dokploy \
  --option MACHINE_IMAGE=ubuntu:22.04 \
  --option MACHINE_TYPE=medium \
  --option DOKPLOY_PROJECT_NAME=staging-workspaces
```

### Team Environment

```bash
devpod provider set-options dokploy \
  --option DOKPLOY_SERVER_URL=https://team-dokploy.company.com \
  --option DOKPLOY_PROJECT_NAME=team-development \
  --option MACHINE_TYPE=medium
```

## 📊 Monitoring

### Workspace Status

```bash
# Check workspace status
devpod status my-workspace

# View real-time logs
devpod logs my-workspace --follow

# List all workspaces
devpod list --output json
```

### Dokploy Dashboard

Monitor your workspaces directly in the Dokploy dashboard:

- Application status and health
- Resource usage (CPU, memory, disk)
- Deployment history
- Real-time logs

### Development Monitoring

```bash
# Show debug information
make debug-env

# List workspaces using this provider
make list-workspaces

# Check provider configuration
make show-config
```

## 🔒 Security Considerations

1. **API Token Security**: Store API tokens securely, never commit them to version control
2. **Environment Files**: Use `.env` files for local development, ensure they're gitignored
3. **Network Security**: Ensure secure communication between DevPod and Dokploy
4. **Access Control**: Use appropriate Dokploy user permissions
5. **Resource Limits**: Configure appropriate CPU and memory limits
6. **Image Security**: Use trusted Docker images and keep them updated
7. **SSH Security**: Provider sets up secure SSH access with proper user isolation

## 🚀 Best Practices

### Workspace Management

```bash
# Use descriptive names
devpod up frontend-dashboard --provider dokploy
devpod up api-user-service --provider dokploy

# Include environment or purpose
devpod up staging-frontend --provider dokploy
devpod up dev-microservice-auth --provider dokploy
```

### Resource Cleanup

```bash
# Regular cleanup of unused workspaces
devpod list | grep -E "(Stopped|NotFound)" | awk '{print $1}' | xargs -I {} devpod delete {} --force

# Using Makefile for development
make cleanup-test  # Clean up test workspaces
make cleanup-workspaces  # Clean up all provider workspaces
```

### Configuration Management

```bash
# Use environment files for team consistency
# .env.team
DOKPLOY_SERVER_URL=https://team-dokploy.company.com
DOKPLOY_PROJECT_NAME=team-development
MACHINE_TYPE=medium
MACHINE_IMAGE=node:18

# Load configuration
source .env.team
make configure-env
```

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on:

- Setting up the development environment
- Running tests and validation
- Code style and best practices
- Submitting pull requests

### Quick Development Setup

```bash
git clone https://github.com/NaNomicon/dokploy-devpod-provider.git
cd dokploy-devpod-provider
make setup
make setup-env
# Edit .env file
make install-local
make test-docker
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [DevPod](https://devpod.sh/) for the excellent development environment platform
- [Dokploy](https://dokploy.com/) for the powerful deployment infrastructure
- The open-source community for inspiration and support

## 📚 Additional Resources

- [DevPod Documentation](https://devpod.sh/docs)
- [Dokploy Documentation](https://docs.dokploy.com/)
- [DevPod Provider Development Guide](https://devpod.sh/docs/developing-providers/quickstart)
- [Community Providers](https://devpod.sh/docs/managing-providers/add-provider#community-providers)
- [Tutorial: Complete Guide](tutorial.md)
- [Development Journey](journey.md)

## 🔗 Quick Reference

### Essential Commands

```bash
# Provider management
devpod provider add NaNomicon/dokploy-devpod-provider
devpod provider use dokploy
devpod provider set-options dokploy --option KEY=VALUE

# Workspace management
devpod up workspace-name --provider dokploy
devpod ssh workspace-name
devpod stop workspace-name
devpod delete workspace-name --force

# Development (with Makefile)
make install-local
make configure-env
make test-docker
make force-reinstall
make cleanup-workspaces
```

### Configuration Template

```bash
# Required
DOKPLOY_SERVER_URL=https://your-dokploy.com
DOKPLOY_API_TOKEN=your_api_token

# Optional
DOKPLOY_PROJECT_NAME=devpod-workspaces
DOKPLOY_SERVER_ID=
MACHINE_TYPE=small
MACHINE_IMAGE=ubuntu:22.04
AGENT_PATH=/opt/devpod/agent
```

---

**Note**: This is a community provider for DevPod. It is not officially maintained by the DevPod or Dokploy teams.
