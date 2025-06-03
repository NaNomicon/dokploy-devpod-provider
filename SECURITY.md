# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Which versions are eligible for receiving such patches depends on the CVSS v3.0 Rating:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |

## Reporting a Vulnerability

If you discover a security vulnerability, please report it to us responsibly:

### Private Disclosure

For security vulnerabilities, please **DO NOT** create a public GitHub issue. Instead:

1. **Email**: Send details to the maintainers via GitHub's private vulnerability reporting feature
2. **Include**:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)
   - Provider version and configuration (sanitized)
   - DevPod and Dokploy versions

### What to Expect

- **Acknowledgment**: We'll acknowledge receipt within 48 hours
- **Assessment**: We'll assess the vulnerability within 5 business days
- **Updates**: We'll provide regular updates on our progress
- **Resolution**: We'll work to resolve critical issues within 30 days

### Disclosure Timeline

- **Day 0**: Vulnerability reported
- **Day 1-2**: Acknowledgment and initial assessment
- **Day 3-7**: Detailed analysis and fix development
- **Day 8-30**: Testing, review, and release preparation
- **Day 30+**: Public disclosure (coordinated with reporter)

## Security Architecture

### Machine Provider Security Model

This provider implements DevPod's **Machine Provider** pattern with security considerations:

#### Layer 1: Machine Infrastructure Security

- **SSH Access Control**: Secure SSH setup with dedicated `devpod` user
- **User Isolation**: Non-root user with sudo access for development tasks
- **Network Security**: Controlled port exposure and firewall configuration
- **Container Security**: Proper resource limits and security contexts

#### Layer 2: Development Environment Security

- **Agent Security**: DevPod agent runs with minimal required permissions
- **Container Isolation**: Development containers isolated from host system
- **Resource Limits**: CPU and memory limits prevent resource exhaustion
- **Image Security**: Base images from trusted sources with security updates

### SSH Security Implementation

The provider implements secure SSH access:

```bash
# Secure SSH setup in containers
- Creates dedicated `devpod` user with limited privileges
- Configures SSH with secure defaults
- Uses SSH key-based authentication when available
- Implements proper user isolation and sudo access
```

## Security Best Practices

When using this provider:

### API Token Security

- **Storage**: Store Dokploy API tokens securely using environment variables or secret management
- **Environment Files**: Use `.env` files for local development (automatically gitignored)
- **Version Control**: Never commit tokens to version control
- **Rotation**: Rotate tokens regularly (recommended: every 90 days)
- **Permissions**: Use tokens with minimal required permissions
- **Monitoring**: Monitor token usage and audit access logs

### Network Security

- **HTTPS Enforcement**: All Dokploy API communications use HTTPS
- **Certificate Validation**: SSL certificates are verified for all connections
- **Private Networks**: Consider VPN or private networks for sensitive environments
- **Firewall Configuration**: Implement proper firewall rules for SSH and application ports
- **Network Monitoring**: Monitor network traffic for anomalies

### Container and SSH Security

- **Base Images**: Use trusted base images (default: `ubuntu:22.04`)
- **Security Updates**: Keep images updated with latest security patches
- **Vulnerability Scanning**: Regularly scan images for known vulnerabilities
- **Resource Limits**: Implement CPU and memory limits to prevent DoS
- **User Privileges**: Containers run with non-root users when possible
- **SSH Hardening**: SSH configured with secure defaults and key-based auth

### Access Control and Authentication

- **Principle of Least Privilege**: Use minimal required Dokploy permissions
- **User Management**: Proper user isolation in containers
- **Audit Logging**: Regular audit of access logs and user activities
- **Multi-Factor Authentication**: Enable MFA on Dokploy accounts when available
- **Session Management**: Proper session timeout and management

### Development Environment Security

- **Environment Isolation**: Use `.env` files for configuration (gitignored)
- **Secret Management**: Never expose secrets in logs or error messages
- **Development vs Production**: Separate configurations for different environments
- **Workspace Cleanup**: Regular cleanup of unused workspaces and resources
- **Code Security**: Scan development code for security vulnerabilities

## Known Security Considerations

### API Token Exposure

- **Risk**: API tokens in logs, error messages, or command history
- **Mitigation**:
  - Tokens are marked as `password: true` in provider configuration
  - Error messages sanitize sensitive information
  - Use `.env` files to avoid command-line exposure
  - Provider masks tokens in debug output

### Network Communication Security

- **Risk**: Man-in-the-middle attacks on API communications
- **Mitigation**:
  - HTTPS enforced for all Dokploy API calls
  - SSL certificate validation enabled
  - Secure SSH connections for workspace access

### Container and SSH Security

- **Risk**: Privilege escalation or container escape
- **Mitigation**:
  - Containers run with appropriate security contexts
  - SSH access limited to dedicated `devpod` user
  - Resource limits prevent resource exhaustion attacks
  - Regular security updates for base images

### Workspace Data Security

- **Risk**: Unauthorized access to workspace data
- **Mitigation**:
  - Proper user isolation in containers
  - SSH key-based authentication when available
  - Workspace cleanup on deletion
  - Volume encryption when supported by Dokploy

### Development Tool Security

- **Risk**: Malicious code execution in development environment
- **Mitigation**:
  - Isolated development containers
  - Limited network access from containers
  - Regular security scanning of development tools
  - Proper secret management in development workflows

## Security Configuration Examples

### Secure Environment Configuration

```bash
# .env file (automatically gitignored)
DOKPLOY_SERVER_URL=https://secure-dokploy.company.com
DOKPLOY_API_TOKEN=your_secure_token_here
DOKPLOY_PROJECT_NAME=secure-development
MACHINE_TYPE=small
MACHINE_IMAGE=ubuntu:22.04  # Use LTS versions for security
```

### Secure Provider Configuration

```bash
# Configure with security best practices
devpod provider set-options dokploy \
  --option DOKPLOY_SERVER_URL=https://secure-dokploy.company.com \
  --option DOKPLOY_API_TOKEN="$(cat /secure/path/to/token)" \
  --option MACHINE_IMAGE=ubuntu:22.04
```

### Secure Development Workflow

```bash
# Use Makefile for secure development
make setup-env          # Create secure environment file
# Edit .env with secure configuration
make configure-env       # Configure without exposing secrets
make test-docker         # Test with security validation
make cleanup-workspaces  # Clean up resources regularly
```

## Security Monitoring and Auditing

### Recommended Monitoring

- **API Access**: Monitor Dokploy API access logs
- **SSH Connections**: Monitor SSH access to workspaces
- **Resource Usage**: Monitor CPU/memory usage for anomalies
- **Network Traffic**: Monitor unusual network patterns
- **Workspace Activity**: Track workspace creation/deletion patterns

### Audit Checklist

- [ ] Regular API token rotation
- [ ] Review workspace access logs
- [ ] Update base images for security patches
- [ ] Scan for vulnerable dependencies
- [ ] Review user permissions and access
- [ ] Monitor resource usage patterns
- [ ] Check for unused workspaces and cleanup

## Security Updates and Patches

Security updates will be:

- **Released**: As patch versions following semantic versioning
- **Documented**: In release notes with security impact assessment
- **Announced**: Via GitHub releases and security advisories
- **Tagged**: With appropriate security labels and CVSS scores
- **Tested**: With comprehensive security validation

### Update Process

1. **Assessment**: Evaluate security impact and affected versions
2. **Development**: Develop and test security fixes
3. **Validation**: Security testing and code review
4. **Release**: Coordinated release with documentation
5. **Communication**: Notify users via multiple channels

## Incident Response

In case of a security incident:

1. **Immediate**: Isolate affected systems and workspaces
2. **Assessment**: Evaluate scope and impact of the incident
3. **Containment**: Implement containment measures
4. **Investigation**: Conduct thorough security investigation
5. **Recovery**: Restore services with security improvements
6. **Communication**: Transparent communication with users
7. **Post-Incident**: Review and improve security measures

## Security Resources

### Documentation

- [DevPod Security Guide](https://devpod.sh/docs/security)
- [Dokploy Security Documentation](https://docs.dokploy.com/docs/security)
- [Container Security Best Practices](https://kubernetes.io/docs/concepts/security/)

### Tools and Validation

- Use `make validate` for configuration validation
- Use `make lint` for security-focused code analysis
- Use `make debug-env` for secure debugging (sanitized output)
- Regular security scanning with container security tools

### Community Resources

- [OWASP Container Security](https://owasp.org/www-project-container-security/)
- [CIS Docker Benchmark](https://www.cisecurity.org/benchmark/docker)
- [DevSecOps Best Practices](https://devsecops.org/)

## Contact and Support

For security-related questions or concerns:

- **Vulnerabilities**: Use GitHub's private vulnerability reporting
- **Security Questions**: Create GitHub discussions with security tag
- **Documentation**: Check existing security advisories and documentation
- **Community**: Engage with the DevPod and Dokploy security communities

### Emergency Contact

For critical security issues requiring immediate attention:

- Use GitHub's private vulnerability reporting with "Critical" severity
- Include detailed impact assessment and reproduction steps
- Provide contact information for coordinated disclosure

Thank you for helping keep the Dokploy DevPod Provider and the broader development community secure! 🔒
