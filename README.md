# Dokploy DevPod Provider

A high-performance DevPod provider for [Dokploy](https://dokploy.com/) that enables seamless development environment creation and management through Dokploy's container orchestration platform.

## 🚀 Features

- **⚡ Fast Deployment**: Uses Alpine Linux for 6x faster package installation
- **🔧 Automatic SSH Setup**: Intelligent port mapping and SSH configuration
- **🐳 Docker Swarm Integration**: Native Dokploy/Docker Swarm compatibility
- **🛠️ Zero Configuration**: Automatic project and application management
- **📊 Comprehensive Debugging**: Detailed logging and error analysis
- **🔄 DevPod Compatible**: Full support for `.devcontainer.json` workflows

## 📋 Prerequisites

- [DevPod CLI](https://devpod.sh/) installed
- Access to a Dokploy instance
- Dokploy API token (generate from Settings > Profile > API/CLI)

## 🛠️ Installation

### 1. Clone the Provider

```bash
git clone <repository-url>
cd dokploy-devpod-provider
```

### 2. Configure Environment

Create a `.env` file with your Dokploy credentials:

```bash
# Required
DOKPLOY_SERVER_URL=https://your-dokploy-instance.com
DOKPLOY_API_TOKEN=your-api-token-here

# Optional
DOKPLOY_PROJECT_NAME=devpod-workspaces
DOKPLOY_SERVER_ID=your-server-id  # For multi-server setups
MACHINE_TYPE=small
```

### 3. Install Provider

```bash
make install-local
```

## 🚀 Usage

### Create a Workspace from Git Repository

```bash
devpod up https://github.com/your-org/your-repo.git --provider dokploy-dev
```

### Create a Workspace from Local Directory

```bash
devpod up ./my-project --provider dokploy-dev
```

### Connect to Existing Workspace

```bash
devpod ssh my-workspace
```

## 🔐 SSH Authentication & DevPod Integration

### How SSH Authentication Works

This provider is designed to work seamlessly with DevPod's built-in SSH authentication system:

1. **Provider Creates Infrastructure**: Sets up Alpine container with SSH daemon and hybrid authentication
2. **DevPod Connects**: Uses its built-in SSH key management to connect to the container
3. **Agent Injection**: DevPod automatically injects its agent and SSH keys
4. **Workspace Management**: DevPod agent handles development environment setup

### Authentication Flow

```
DevPod → SSH Connection → Container (hybrid auth) → Agent Injection → SSH Keys → Secure Connection
```

#### Container SSH Configuration

The provider configures the container with:

```bash
# Hybrid authentication support
PubkeyAuthentication yes          # For DevPod's SSH keys
PasswordAuthentication yes        # For initial connection
AuthorizedKeysFile .ssh/authorized_keys
PermitRootLogin no
```

#### DevPod's Role

- **SSH Key Generation**: DevPod automatically generates or uses existing SSH keys
- **Key Injection**: Injects public keys into `/home/devpod/.ssh/authorized_keys`
- **Agent Deployment**: Installs DevPod agent for workspace management
- **Credential Forwarding**: Handles Git credentials, Docker credentials, etc.

### Why This Approach Works

- **No Manual SSH Setup**: DevPod handles all SSH key management automatically
- **Secure by Default**: Uses SSH keys for ongoing connections
- **Fallback Support**: Password authentication available if needed
- **Platform Agnostic**: Works across different development environments

## ⏱️ Important: Docker Swarm Port Mapping Delay

**Expected Behavior**: When creating a new workspace, you'll see a 60+ second delay during SSH setup. This is **normal and expected** behavior.

### Why This Happens

Dokploy uses Docker Swarm for container orchestration. When the provider creates SSH port mappings, Docker Swarm needs time to propagate these mappings across the cluster. This process typically takes 60-120 seconds.

### What You'll See

```
🎉 SSH port mapping configured successfully!
   Using port: 2222

ℹ️  NOTICE: Docker Swarm Port Mapping Delay
   Dokploy uses Docker Swarm for container orchestration, which requires
   time for port mappings to propagate across the cluster. This 60+ second
   delay is normal and expected behavior, not a provider issue.

   • Port mapping API: ✅ Completed successfully
   • Port propagation: ⏳ In progress (60-120 seconds typical)
   • SSH accessibility: ⏳ Will be available after propagation

DEBUG: Sleeping for 60 seconds to allow Dokploy port mapping to propagate...
```

### This is NOT a Bug

- ✅ The provider is working correctly
- ✅ Port mapping was created successfully
- ⏳ Docker Swarm is propagating the mapping
- 🎯 SSH will be accessible once propagation completes

## 🏗️ Architecture

### DevPod Two-Layer Architecture

The provider works with DevPod's two-layer architecture:

1. **Layer 1 (Infrastructure)**: Alpine Linux container with SSH access
   - Managed by this Dokploy provider
   - Provides the base environment and SSH connectivity
2. **Layer 2 (Development Environment)**: Your actual development tools
   - Managed by DevPod agent
   - Installs Node.js, Python, Docker, etc. based on your `.devcontainer.json`

### Dokploy Integration

```
DevPod CLI → Dokploy Provider → Dokploy API → Docker Swarm → Alpine Container
```

## 🔧 Configuration Options

| Option                 | Description                       | Default             | Required |
| ---------------------- | --------------------------------- | ------------------- | -------- |
| `DOKPLOY_SERVER_URL`   | Your Dokploy server URL           | -                   | ✅       |
| `DOKPLOY_API_TOKEN`    | API token from Dokploy            | -                   | ✅       |
| `DOKPLOY_PROJECT_NAME` | Project name for workspaces       | `devpod-workspaces` | ❌       |
| `DOKPLOY_SERVER_ID`    | Server ID for multi-server        | -                   | ❌       |
| `MACHINE_TYPE`         | Machine size (small/medium/large) | `small`             | ❌       |
| `AGENT_PATH`           | DevPod agent installation path    | `/opt/devpod/agent` | ❌       |

## 🐛 Troubleshooting

### Common Issues

#### 1. "Application not found" Error

- **Cause**: Application was deleted from Dokploy dashboard
- **Solution**: Delete the workspace and recreate it

#### 2. SSH Connection Timeout

- **Cause**: Port mapping still propagating
- **Solution**: Wait 2-3 minutes and try again

#### 3. "Port already in use" Error

- **Cause**: Previous workspace using the same port
- **Solution**: Delete unused workspaces or wait for automatic cleanup

### Debug Mode

Enable detailed debugging:

```bash
devpod up <source> --provider dokploy-dev --debug
```

### SSH Connection Issues

If you encounter SSH connection problems:

1. **Check Port Availability**: Ensure the SSH port (usually 2222) is accessible
2. **Wait for Propagation**: Docker Swarm port mapping can take 60-120 seconds
3. **DevPod Retry**: DevPod will automatically retry connections
4. **Manual Verification**: Test port accessibility with `nc -z <host> <port>`

### DevPod Agent Issues

If DevPod agent injection fails:

```bash
# Check if container is running
devpod ssh <workspace-name>

# If SSH works but agent fails, check logs
devpod logs <workspace-name>
```

### Container Access

For direct container access (debugging only):

```bash
# SSH with password (if needed for debugging)
ssh -p 2222 devpod@your-dokploy-host.com
# Password: devpod (only for initial setup/debugging)
```

**Note**: DevPod will automatically set up SSH key authentication, so password access is mainly for debugging purposes.

## 📊 Performance

### Deployment Speed

- **Alpine Linux**: ~5 seconds for package installation
- **Ubuntu (previous)**: 30+ seconds for package installation
- **Improvement**: 6x faster deployment

### Resource Usage

- **Base Image**: Alpine Linux (~5MB)
- **Memory**: Minimal overhead
- **CPU**: Efficient container startup

## 🔄 Development

### Testing

```bash
# Test with Git repository
make test-git

# Test with local directory
make test-local

# Validate provider configuration
make validate
```

### Local Development

```bash
# Install development dependencies
make install-dev

# Run linting
make lint

# Clean up test workspaces
make clean
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📝 License

[MIT License](LICENSE)

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/your-org/dokploy-devpod-provider/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-org/dokploy-devpod-provider/discussions)
- **Dokploy**: [Dokploy Documentation](https://docs.dokploy.com/)
- **DevPod**: [DevPod Documentation](https://devpod.sh/docs)

## 🙏 Acknowledgments

- [Dokploy](https://dokploy.com/) for the excellent container platform
- [DevPod](https://devpod.sh/) for the development environment framework
- [Alpine Linux](https://alpinelinux.org/) for the lightweight base image

---

**Note**: This provider is optimized for Dokploy's Docker Swarm architecture and includes built-in handling for platform-specific behaviors like port mapping propagation delays.
