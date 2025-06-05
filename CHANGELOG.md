# Dokploy DevPod Provider - Changelog

## v0.1.2 - Template-based Release System

### 🚀 Major Features

- **Template-based Release System**: Implemented `gomplate`-based template system for generating `provider.yaml`

  - Automatic checksum calculation from built binaries
  - Environment variable-driven template rendering
  - Eliminates manual checksum updates and human errors

- **Enhanced Version Management**: Switched to Git tag-based versioning
  - Version automatically detected from Git tags (`git describe --tags`)
  - Proper separation of version number (`0.1.2`) and tag (`v0.1.2`) in templates
  - Eliminates circular dependency between Makefile and provider.yaml

### ✨ Release Workflow Improvements

- **Simplified CI/CD**: GitHub Actions workflow now uses Makefile targets

  - DRY principle: no duplication between local and CI builds
  - Consistent build process via `make release-prepare`
  - Automatic tool installation via `make setup`

- **Cross-compilation**: Added `gox` for parallel binary building

  - Faster builds across all supported platforms
  - Consistent binary naming and output structure

- **Version Bumping**: Added semantic version bump targets
  - `make version-bump-patch` - Bump patch version and create Git tag
  - `make version-bump-minor` - Bump minor version and create Git tag
  - `make version-bump-major` - Bump major version and create Git tag

### 🔧 Development Experience

- **Enhanced Makefile**: Comprehensive release management targets

  - `make generate-provider` - Generate provider.yaml from template
  - `make release-prepare` - Complete release preparation pipeline
  - `make validate-provider-checksums` - Verify checksum accuracy
  - `make restore-provider` - Restore provider.yaml from backup

- **Template Architecture**: Clean separation of concerns
  - Template file: `hack/provider.yaml.tpl`
  - Version variables: `VERSION_NUMBER` and `VERSION_TAG`
  - Automatic binary checksum injection via `file.Read | crypto.SHA256`

### 🐛 Bug Fixes

- **Checksum Consistency**: Eliminated checksum mismatches between local and CI builds
- **Version Handling**: Fixed template version variable usage for proper Git tag integration
- **Release Artifacts**: Proper binary and checksum file organization

### 📚 Documentation

- **Release Process**: Updated documentation for new template-based workflow
- **Acknowledgments**: Added credit to STACKIT DevPod Provider for template inspiration
- **Version Management**: Clarified Git tag-based versioning approach

### 🔄 Breaking Changes

- **Version Detection**: Now requires Git tags for version detection (no longer uses provider.yaml)
- **Build Process**: `make generate-provider` must be run to update provider.yaml after version changes

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
