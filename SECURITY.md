# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |

## Reporting a Vulnerability

We take the security of the DevPod Dokploy Provider seriously. If you discover a security vulnerability, please follow the responsible disclosure process outlined below.

### How to Report

**Please do NOT create a public GitHub issue for security vulnerabilities.**

Instead, please report security issues via one of the following methods:

1. **Private Security Advisory** (Preferred): Use GitHub's private vulnerability reporting feature by clicking the "Security" tab and then "Report a vulnerability"
2. **Email**: Send an email to the project maintainers through GitHub

### What to Include

When reporting a vulnerability, please include:

- A clear description of the vulnerability
- Steps to reproduce the issue
- Affected versions (if known)
- Potential impact assessment
- Any suggested fixes or mitigations
- Your contact information for follow-up

### Response Process

- **Acknowledgment**: We will acknowledge receipt of your report within **48 hours**
- **Assessment**: We will assess the vulnerability and determine its severity within **5 business days**
- **Updates**: We will provide regular updates on our progress
- **Resolution**: We aim to release a fix within **30 days** for critical vulnerabilities, **90 days** for others

### Disclosure Timeline

We follow a coordinated disclosure process:

1. **Day 0**: Vulnerability reported and acknowledged
2. **Day 1-5**: Initial assessment and validation
3. **Day 6-30/90**: Development and testing of fix
4. **Day 30/90+**: Public disclosure after fix is released

We may publicly disclose the vulnerability earlier if:

- A fix is available and deployed
- The vulnerability is already publicly known
- We mutually agree with the reporter

## Security Considerations for DevPod Dokploy Provider

### What We Secure

This security policy covers:

- The DevPod Dokploy Provider CLI binary
- Docker Compose configurations and templates
- SSH key handling and authentication
- API interactions with Dokploy servers
- Container isolation and security contexts

### Known Security Model

The DevPod Dokploy Provider operates with the following security model:

- **SSH Access**: Uses root-based SSH authentication for maximum DevPod compatibility
- **API Communication**: All Dokploy API calls use HTTPS with API key authentication
- **Container Isolation**: Workspaces run in isolated Docker containers with appropriate security contexts
- **Key Management**: SSH keys are securely injected into containers and properly configured

### Best Practices for Users

To maintain security when using this provider:

1. **Keep Updated**: Always use the latest version of the provider
2. **Secure Dokploy**: Ensure your Dokploy server is properly secured and updated
3. **Network Security**: Use private networks or VPNs when possible
4. **API Keys**: Rotate Dokploy API keys regularly and store them securely
5. **SSH Keys**: Use strong SSH keys and rotate them periodically

## Out of Scope

The following are outside the scope of this security policy:

- Vulnerabilities in DevPod core (report to the DevPod project)
- Vulnerabilities in Dokploy server (report to the Dokploy project)
- Issues with Docker or container runtime (report to respective projects)
- Social engineering attacks
- Physical security issues

## Security Updates

Security updates will be:

- Released as patch versions (e.g., 0.1.1, 0.1.2)
- Documented in the CHANGELOG.md with severity assessment
- Announced via GitHub releases with security advisory
- Tagged with appropriate CVE identifiers when applicable

## Contact

For security-related questions or concerns:

- **Vulnerabilities**: Use GitHub's private vulnerability reporting
- **General Security Questions**: Create a GitHub discussion with the "security" tag
- **Documentation**: Check existing security advisories and documentation

Thank you for helping keep the DevPod Dokploy Provider secure! 🔒
