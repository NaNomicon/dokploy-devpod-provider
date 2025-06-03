---
name: Bug report
about: Create a report to help us improve
title: ""
labels: bug
assignees: ""
---

**Describe the bug**
A clear and concise description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:

1. Run command '...'
2. See error

**Expected behavior**
A clear and concise description of what you expected to happen.

**Screenshots/Logs**
If applicable, add screenshots or logs to help explain your problem.

**Environment (please complete the following information):**

- OS: [e.g. macOS 14.5.0, Ubuntu 22.04, Windows 11]
- DevPod Version: [run `devpod version`]
- Provider Version: [e.g. v0.1.0]
- Dokploy Version: [e.g. v0.8.0]

**Provider Configuration:**

```yaml
# Please provide your provider configuration (remove sensitive data like API tokens)
DOKPLOY_SERVER_URL: https://your-dokploy.com
DOKPLOY_PROJECT_NAME: your-project
MACHINE_TYPE: small
MACHINE_IMAGE: ubuntu:22.04
AGENT_PATH: /opt/devpod/agent
```

**Debug Information:**

Please run the following commands and include the output (sanitize any sensitive information):

```bash
# Provider debug information
make debug-env

# Provider configuration
make show-config

# Workspace information (if applicable)
devpod list
make list-workspaces

# Provider logs
devpod provider logs dokploy
```

**Workspace Information (if applicable):**

- Workspace Name: [e.g. my-workspace]
- Workspace Status: [run `devpod status workspace-name`]
- Creation Method: [e.g. Git repository, Docker image, named workspace]

**Error Details:**

```
# Please include the full error message and stack trace
```

**Troubleshooting Attempted:**

Please check which troubleshooting steps you've already tried:

- [ ] Verified API token is valid and has correct permissions
- [ ] Checked Dokploy server is accessible
- [ ] Tested with `devpod up --debug` for verbose logging
- [ ] Ran `make validate` to check provider configuration
- [ ] Checked for active workspaces with `make list-workspaces`
- [ ] Tried recreating the workspace
- [ ] Checked Dokploy dashboard for application status

**Additional context**
Add any other context about the problem here, including:

- Recent changes to your setup
- Whether this worked before
- Any custom configurations or modifications
- Network or firewall considerations
- Related issues or error patterns
