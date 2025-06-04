# Dokploy DevPod Provider

A high-performance DevPod provider for [Dokploy](https://dokploy.com/) that enables seamless development environment creation and management through Dokploy's container orchestration platform.

## 🚀 Features

- **⚡ Binary Helper**: Fast Go CLI binary instead of slow shell scripts
- **🔧 Automatic SSH Setup**: Intelligent port mapping and SSH configuration
- **🐳 Docker-in-Docker**: Native Docker support with `cruizba/ubuntu-dind:latest`
- **🛠️ Zero Configuration**: Automatic project and application management
- **📊 Comprehensive Debugging**: Detailed logging and error analysis
- **🔄 DevPod Compatible**: Full support for `.devcontainer.json` workflows
- **🌐 Cross-Platform**: Supports Linux, macOS, and Windows

## 📋 Prerequisites

- [DevPod CLI](https://devpod.sh/) installed
- Access to a Dokploy instance
- Dokploy API token (generate from Settings > Profile > API/CLI)

## 🛠️ Installation

### 1. Quick Install from GitHub

```bash
devpod provider add https://github.com/NaNomicon/dokploy-devpod-provider
```

### 2. Local Development Install

```bash
git clone https://github.com/NaNomicon/dokploy-devpod-provider
cd dokploy-devpod-provider
make install-dev
```

### 3. Configure Provider

```bash
# Option 1: Interactive configuration
make configure

# Option 2: Environment file configuration
make setup-env
# Edit .env file with your settings
make configure-env
```

## 🚀 Usage

### Create a Workspace from Git Repository

```bash
devpod up https://github.com/your-org/your-repo.git --provider dokploy
```

### Create a Workspace from Local Directory

```bash
devpod up ./my-project --provider dokploy
```

### Connect to Existing Workspace

```bash
devpod ssh my-workspace
```

## 🏗️ Architecture

### Binary Helper Implementation

This provider uses a **Go CLI binary helper** instead of shell scripts for superior performance and reliability:

```
DevPod → Binary Helper → Dokploy API → Docker Container → SSH Access
```

#### Key Components

- **`dokploy-provider` binary**: Cross-platform Go CLI handling all operations
- **Dokploy API client**: Comprehensive REST API integration
- **SSH client**: Automatic connection discovery and command execution
- **Configuration management**: Environment-based configuration loading

#### Commands Implemented

| Command   | Purpose                                      | Implementation                                      |
| --------- | -------------------------------------------- | --------------------------------------------------- |
| `init`    | Validate configuration and test connectivity | API health check + SSH validation                   |
| `create`  | Create new workspace with SSH setup          | Project/app creation + Docker config + port mapping |
| `delete`  | Remove workspace and cleanup resources       | Application deletion via API                        |
| `start`   | Start stopped workspace                      | Application deployment                              |
| `stop`    | Stop running workspace                       | Application stop                                    |
| `status`  | Get workspace status                         | API status mapping to DevPod states                 |
| `command` | Execute commands via SSH                     | Dynamic SSH discovery + command execution           |

### Container Setup Process

The provider creates workspaces using a 4-stage setup process:

1. **Stage 1**: Package update (1-2 minutes) - `apt-get update`
2. **Stage 2**: SSH installation (30-60 seconds) - install `openssh-server`
3. **Stage 3**: User setup (10-20 seconds) - create `devpod` user
4. **Stage 4**: SSH configuration (10-20 seconds) - configure SSH daemon

Total setup time: **2-4 minutes**

## 🔐 SSH Authentication & DevPod Integration

### How SSH Authentication Works

This provider follows DevPod's standard SSH connection pattern:

1. **Provider Creates Infrastructure**: Sets up Ubuntu container with SSH daemon
2. **DevPod Connects via SSH**: Uses standard SSH with password authentication
3. **DevPod Agent Installation**: DevPod installs its agent binary via SSH
4. **Workspace Management**: DevPod agent handles development environment setup

### Authentication Flow

```
DevPod → SSH Connection (password) → Container → DevPod Agent Installation → Workspace Setup
```

#### Container SSH Configuration

```bash
# SSH daemon configuration
PubkeyAuthentication yes          # For SSH key authentication
PasswordAuthentication yes        # For password authentication
AuthorizedKeysFile .ssh/authorized_keys
PermitRootLogin no
Port 22

# User setup
User: devpod
Password: devpod
Sudo: NOPASSWD:ALL
Docker group: yes
```

### Dynamic SSH Connection Discovery

The provider implements robust SSH connection handling:

1. **Create Phase**: Sets up container and configures SSH port mapping (2222-2230 range)
2. **Command Phase**: Dynamically retrieves SSH connection details using application ID
3. **API-Based Discovery**: Uses Dokploy API to find correct SSH port for each workspace
4. **Automatic Retry**: Handles Docker Swarm port propagation delays (60-120 seconds)

## ⏱️ Important: Docker Swarm Port Mapping Delay

**Expected Behavior**: New workspaces experience a 60+ second delay during SSH setup. This is **normal**.

### Why This Happens

Dokploy uses Docker Swarm for orchestration. Port mappings need time to propagate across the cluster (60-120 seconds).

### What You'll See

```
🎉 SSH port mapping configured successfully!
   Using port: 2222
⏳ Waiting for port to become available...
   This may take 60-120 seconds due to Docker Swarm propagation
```

## 📊 Configuration Options

| Option                 | Description                       | Default             | Required |
| ---------------------- | --------------------------------- | ------------------- | -------- |
| `DOKPLOY_SERVER_URL`   | Dokploy server URL                | -                   | ✅       |
| `DOKPLOY_API_TOKEN`    | API authentication token          | -                   | ✅       |
| `DOKPLOY_PROJECT_NAME` | Project name for workspaces       | `devpod-workspaces` | ❌       |
| `DOKPLOY_SERVER_ID`    | Server ID for multi-server setups | -                   | ❌       |
| `MACHINE_TYPE`         | Machine size (small/medium/large) | `small`             | ❌       |
| `AGENT_PATH`           | DevPod agent installation path    | `/opt/devpod/agent` | ❌       |

## 🐛 Troubleshooting

### Common Issues

#### 1. "Application not found" Error

- **Cause**: Application was deleted from Dokploy dashboard
- **Solution**: Delete the workspace and recreate it

#### 2. SSH Connection Timeout

- **Cause**: Port mapping still propagating (normal for 60-120 seconds)
- **Solution**: Wait 2-3 minutes and try again

#### 3. "Port already in use" Error

- **Cause**: Previous workspace using the same port
- **Solution**: Delete unused workspaces or wait for automatic cleanup

### Debug Mode

Enable detailed debugging:

```bash
devpod up <source> --provider dokploy --debug
```

### SSH Connection Issues

If you encounter SSH connection problems:

1. **Check Port Availability**: Ensure the SSH port (usually 2222) is accessible
2. **Wait for Propagation**: Docker Swarm port mapping can take 60-120 seconds
3. **DevPod Retry**: DevPod will automatically retry connections
4. **Manual Verification**: Test port accessibility with `nc -z <host> <port>`

### Provider Management

```bash
# Force reinstall (handles stuck workspaces)
make force-reinstall

# Clean up test workspaces
make cleanup-test

# Fix specific stuck workspace
make fix-stuck-workspace

# Nuclear option (delete everything)
make nuclear-cleanup
```

### Binary Issues

If the binary helper fails:

```bash
# Test binary directly
./dist/dokploy-provider --help

# Rebuild binary
make build

# Test specific commands
./dist/dokploy-provider init --verbose
```

## 🔧 Development

### Building from Source

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Test the binary
make test-build
```

### Development Workflow

```bash
# Setup development environment
make setup

# Install development provider
make install-dev

# Test with Docker workspace
make test-docker

# Test complete lifecycle
make test-lifecycle

# Clean up and reinstall
make force-reinstall
```

### Binary Helper Development

The provider is implemented as a Go CLI application:

```
├── cmd/                    # CLI commands
│   ├── root.go            # Root command and configuration
│   ├── init.go            # Initialize and validate provider
│   ├── create.go          # Create workspace
│   ├── delete.go          # Delete workspace
│   ├── start.go           # Start workspace
│   ├── stop.go            # Stop workspace
│   ├── status.go          # Get workspace status
│   └── command.go         # Execute commands via SSH
├── pkg/                   # Core packages
│   ├── options/           # Configuration management
│   ├── dokploy/           # Dokploy API client
│   ├── client/            # DevPod status types
│   └── ssh/               # SSH client for command execution
├── dist/                  # Built binaries
└── provider.yaml          # DevPod provider configuration
```

## 📊 Performance

### Deployment Speed

- **Binary Helper**: ~100ms command execution
- **Shell Scripts (previous)**: 1-3 seconds per operation
- **Improvement**: 10-30x faster operations

### Resource Usage

- **Binary Size**: ~9MB (statically linked)
- **Memory**: Minimal overhead (~10MB)
- **Network**: Efficient API calls with connection reuse

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test with `make test-lifecycle`
5. Submit a pull request

### Development Tools

```bash
# Install all development tools
make setup

# Check tool availability
make check-tools

# Validate provider configuration
make validate

# Run linting
make lint
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [DevPod](https://devpod.sh/) for the excellent development container platform
- [Dokploy](https://dokploy.com/) for the powerful container orchestration platform
- [Cobra](https://github.com/spf13/cobra) for the CLI framework
- [Logrus](https://github.com/sirupsen/logrus) for structured logging

---

**Note**: This provider is optimized for Dokploy's Docker Swarm architecture and includes built-in handling for platform-specific behaviors like port mapping propagation delays.
