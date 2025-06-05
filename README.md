# Dokploy DevPod Provider

<div align="center">

![Dokploy Logo](https://raw.githubusercontent.com/Dokploy/dokploy/refs/heads/canary/apps/dokploy/logo.png)

**A DevPod provider for [Dokploy](https://dokploy.com/) that creates development workspaces using Docker Compose services.**

[![Version](https://img.shields.io/github/v/release/NaNomicon/dokploy-devpod-provider?label=version)](https://github.com/NaNomicon/dokploy-devpod-provider/releases/latest)
[![Release](https://github.com/NaNomicon/dokploy-devpod-provider/actions/workflows/release.yml/badge.svg)](https://github.com/NaNomicon/dokploy-devpod-provider/actions/workflows/release.yml)

[![Go](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org)
[![DevPod](https://img.shields.io/badge/DevPod-v0.6.15-blue.svg)](https://github.com/loft-sh/devpod)
[![Dokploy](https://img.shields.io/badge/Dokploy-v0.22.7-green.svg)](https://github.com/Dokploy/dokploy)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

[Installation](#-installation) • [Contributing](#contributing)

[![Open in DevPod!](https://devpod.sh/assets/open-in-devpod.svg)](https://devpod.sh/open#https://github.com/NaNomicon/dokploy-devpod-provider)

</div>

---

## 🎉 Current Status

Core functionality is working and tested.

## ✨ What This Does

Creates development workspaces on your Dokploy server that you can connect to with DevPod. Think of it as your personal cloud development environment that automatically sets up everything you need.

### How It Works

```
DevPod → Dokploy Provider → Docker Compose Service → SSH Access → Your Code
```

1. **Creates infrastructure**: Spins up Docker Compose services in Dokploy
2. **Sets up SSH**: Configures secure root access to containers
3. **Installs DevPod agent**: DevPod connects and sets up your dev environment
4. **Ready to code**: Open in VS Code, clone repos, install dependencies automatically

## 🎯 Features

- ✅ **Docker Compose services** - Uses Dokploy's stable service infrastructure
- ✅ **SSH setup** - Automatic port mapping and authentication
- ✅ **DevPod integration** - Full workspace lifecycle support
- ✅ **Docker-in-Docker** - Complete container development capabilities
- ✅ **Git repositories** - Clone and develop any repository
- ✅ **Go binary helper** - Fast operations instead of slow shell scripts

## 📋 Requirements

- [DevPod](https://devpod.sh/) installed locally
- Dokploy server with API access
- API token with project creation permissions

## 🔄 Compatibility

| Component   | Supported Version | Status                |
| ----------- | ----------------- | --------------------- |
| **DevPod**  | v0.6.15           | ✅ Tested & Supported |
| **Dokploy** | v0.22.7           | ✅ Tested & Supported |

> **⚠️ Note**: We currently do not offer backward compatibility. Always use the exact versions listed above for the best experience. Newer versions may work but are not guaranteed to be compatible until tested and updated here.

## 🚀 Installation

### Quick Start

```bash
# Clone and build
git clone https://github.com/your-org/dokploy-devpod-provider
cd dokploy-devpod-provider
make build
make install-dev
```

### Configure Provider

```bash
devpod provider set-options dokploy-dev \
  DOKPLOY_SERVER_URL=https://your-dokploy-server.com \
  DOKPLOY_API_TOKEN=your-api-token
```

### Create Your First Workspace

```bash
# From a Git repository
devpod up https://github.com/microsoft/vscode-remote-try-node --provider dokploy-dev

# From a Docker image
devpod up ubuntu --provider dokploy-dev
```

## ⚙️ Configuration

| Option                 | Description                  | Default             | Required |
| ---------------------- | ---------------------------- | ------------------- | -------- |
| `DOKPLOY_SERVER_URL`   | Your Dokploy server URL      | -                   | ✅       |
| `DOKPLOY_API_TOKEN`    | API token for authentication | -                   | ✅       |
| `DOKPLOY_PROJECT_NAME` | Project name for workspaces  | `devpod-workspaces` | ❌       |

> **Note**: DevPod automatically manages agent installation, credentials injection, and auto-shutdown features.

## 🔧 Development

```bash
# Build binary
make build

# Install as development provider
make install-dev

# Test different scenarios
make test-git      # Test Git repository
make test-docker   # Test Docker image
make test-lifecycle # Full create/delete cycle

# Clean up stuck workspaces
make force-uninstall
```

## 🛠️ Architecture

### Container Setup Process

The provider creates workspaces with a 4-stage setup (~2-4 minutes):

1. **Docker daemon startup** (~30-60 seconds)
2. **Install SSH server + tools** (~30-60 seconds)
3. **Configure SSH for root user** (~10-20 seconds)
4. **Finalize SSH daemon** (~10-20 seconds)

### Technical Details

- **Base Image**: `cruizba/ubuntu-dind:latest` (Docker-in-Docker)
- **SSH Authentication**: Root access with key injection
- **Port Range**: 2222-2250 for SSH mappings
- **API Integration**: Dokploy REST API for service management

## 🐛 Troubleshooting

<details>
<summary><strong>SSH connection issues</strong></summary>

- Wait 2-4 minutes for full container setup
- Check if ports 2222-2250 are available
- Verify API token has correct permissions
</details>

<details>
<summary><strong>DevPod provider issues</strong></summary>

- Try `devpod provider delete dokploy-dev && make install-dev` to reinstall
- Check logs with `devpod up --debug`
</details>

<details>
<summary><strong>Container startup problems</strong></summary>

- Check Dokploy dashboard for service status
- Look at Docker Compose service logs in Dokploy
- Ensure Docker Swarm ports have propagated (can take 60+ seconds)
</details>

## ⚠️ Current Limitations

- Only tested on a few setups - might break in other environments
- SSH setup is slow (2-4 minutes)
- Port range is hardcoded (2222-2250)
- Error handling could be better
- No resource limits on containers
- Limited debugging tools

## 🤝 Contributing

This is a prototype that needs work! Help is welcome:

- 🧪 **Test on different environments** - Try it on your setup
- 📝 **Improve error messages** - Make failures clearer
- 🐛 **Fix edge cases and bugs** - There are definitely some
- 📚 **Write better documentation** - Always room for improvement

See [CONTRIBUTING.md](CONTRIBUTING.md) for development details.

## 📝 License

[MIT License](LICENSE) - feel free to use this however you want.

## 🙏 Acknowledgments

- [DevPod](https://devpod.sh/) for the amazing development environment platform
- [Dokploy](https://dokploy.com/) for the deployment infrastructure
- [Docker](https://docker.com/) for making containers awesome
- [cruizba/ubuntu-dind](https://hub.docker.com/r/cruizba/ubuntu-dind) for the excellent Ubuntu Docker-in-Docker base image
- [STACKIT DevPod Provider](https://github.com/stackitcloud/devpod-provider-stackit/) for the amazing provider codebase

---

<div align="center">
<sub>Built with ❤️ for developers who want better development environments</sub>
</div>
