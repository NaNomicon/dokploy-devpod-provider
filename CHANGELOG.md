# Dokploy DevPod Provider - Changelog

## v0.1.2 - Release Process Improvements

### Changes

- Added template system for generating `provider.yaml` with automatic checksums
- Fixed GitHub Actions build failures due to missing tools
- Fixed shell compatibility issues for Ubuntu/Linux builds
- Updated workflow to use Makefile targets instead of duplicating logic

### Bug Fixes

- Fixed missing `gox` and `gomplate` installation on Ubuntu systems
- Fixed bash-specific syntax that broke on `/bin/sh`
- Fixed template version variables to work with Git tags

### Development

- Added `make version-bump-*` commands for version management
- Added `make generate-provider` to create provider.yaml from template
- Updated GitHub Actions to use `make release-prepare`

This release focuses on fixing the build and release process. The provider functionality remains the same.

---

## v0.1.1 - Bug Fixes and Improvements

### 🐛 Bug Fixes

- **Binary Path Resolution**: Fixed provider not working when DevPod is run outside of development directory
- **Checksum Validation**: Improved checksums generation and validation process

### ✨ Features & Enhancements

- **Documentation**: Enhanced README with Ubuntu DIND details and SSH port checking TODO
- **Development Workflow**: Improved Makefile with better build and validation targets
- **Security Documentation**: Updated security guidelines and best practices
- **Monitoring**: Added GitHub workflows to monitor new versions from DevPod and Dokploy

### 🔧 Technical Changes

- Better binary path handling for cross-environment compatibility
- Enhanced checksum generation process
- Improved documentation structure
- Added automated version monitoring workflows

---

## v0.1.0 - First Release

A functional DevPod provider for Dokploy that creates development workspaces using Docker Compose services.

**Status**: Ready for v0.1.0 release - Core functionality is stable and well-tested. Suitable for development environments.

### What Works

- Creates Docker Compose services in Dokploy projects
- Sets up SSH access to containers with root user
- DevPod can connect and install its agent
- Basic workspace lifecycle (create, start, stop, delete)
- Git repository cloning works
- Docker-in-Docker container setup

### Technical Details

#### Docker Compose Implementation

- Uses Dokploy's Docker Compose services (not applications)
- Go CLI binary instead of shell scripts
- Ubuntu Docker-in-Docker base image (`cruizba/ubuntu-dind:latest`)
- Root SSH access for DevPod compatibility

#### Container Setup Process

Takes about 2-4 minutes to fully set up:

1. Docker daemon startup (~30-60 seconds)
2. Install SSH server and tools (~30-60 seconds)
3. Set up SSH keys for root user (~10-20 seconds)
4. Configure SSH daemon (~10-20 seconds)

#### API Integration

- Uses Dokploy REST API for Docker Compose management
- Basic error handling and logging
- API keys are redacted in logs

### CLI Commands

All commands implemented:

- `init` - Check API connectivity
- `create` - Create Docker Compose service with SSH setup
- `delete` - Delete Docker Compose service
- `start` - Start Docker Compose service
- `stop` - Stop Docker Compose service
- `status` - Get service status
- `command` - Execute commands via SSH

### Configuration Options

Simple configuration - only provider-specific options:

- `DOKPLOY_SERVER_URL` - Dokploy server URL (required)
- `DOKPLOY_API_TOKEN` - API token (required)
- `DOKPLOY_PROJECT_NAME` - Project name (default: `devpod-workspaces`)

**DevPod Agent**: DevPod automatically handles agent installation, credential injection, and auto-shutdown using its built-in defaults. No configuration needed.

### Current Limitations

- Only tested on a few setups
- SSH setup is slow (2-4 minutes)
- Port range is hardcoded (2222-2250)
- Error messages could be better
- No resource limits on containers
- Limited testing and edge case handling

### Known Issues

- Port mapping can take time due to Docker Swarm delays
- Some environments might not work properly
- Error handling needs improvement
- No graceful fallbacks for failures
- Limited debugging information

---

**Note**: This is a working prototype that demonstrates the basic concept. It needs more testing, better error handling, and polish before being ready for wider use.
