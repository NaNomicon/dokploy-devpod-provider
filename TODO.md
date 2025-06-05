# TODO: Dokploy DevPod Provider

## Security

- [ ] Container user hardening
      Run containers with non-root user when possible, better privilege separation

- [ ] Complete analyzation of the project to ensure security

## Tests

- [ ] Unit tests
      Add unit tests for core functions (options loading, client methods, etc.)

- [ ] Integration tests
      Test API integration with mock Dokploy server

- [ ] End-to-end tests
      Full workspace lifecycle testing (create, start, stop, delete)

- [ ] Error handling tests
      Test various failure scenarios and error recovery

- [ ] Performance tests
      Test with multiple workspaces and concurrent operations

## Core Features

- [ ] Multi-server support for Dokploy
      Add support for DOKPLOY_SERVER_ID to handle multiple Dokploy servers

- [ ] SSH connection pooling
      Reuse SSH connections for better performance in command execution (if possible)

- [ ] Non-root SSH option
      Optional devpod user setup for environments that don't allow root access

- [ ] Properly implement all features that DevPod offers

## Improvements

- [ ] Port range configuration
      Allow custom SSH port ranges instead of fixed 2222-2250

- [ ] Workspace backup/restore
      Add ability to backup and restore workspace configurations (maybe)

- [ ] Enhanced error messages
      More specific error messages with actionable troubleshooting steps
