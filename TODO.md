# TODO: Dokploy DevPod Provider

## High Priority

### 🔐 SSH Authentication Alternatives to sshpass

**Current Issue**: The provider currently requires `sshpass` for automated password authentication, which:

- Is not available by default on most systems
- Requires manual installation
- May not be available in some environments (Windows, restricted systems)
- Is considered less secure than key-based authentication

**Potential Solutions to Investigate**:

#### 1. **SSH Key-Based Authentication** (Recommended)

- **Approach**: Automatically generate and inject SSH keys during container creation
- **Benefits**: More secure, no external dependencies, standard SSH practice
- **Implementation**:
  - Generate SSH key pair during `create` phase
  - Inject public key into container's `authorized_keys`
  - Use private key for SSH connections in `command` phase
- **Challenges**: Key management, storage location

#### 2. **Expect Script Alternative**

- **Approach**: Use `expect` or similar tools for password automation
- **Benefits**: More widely available than `sshpass`
- **Challenges**: Still requires external dependency

#### 3. **SSH Agent Integration**

- **Approach**: Leverage existing SSH agent for key management
- **Benefits**: Uses existing SSH infrastructure
- **Challenges**: Complexity, agent availability

#### 4. **Container-Based SSH Key Injection**

- **Approach**: Modify container setup to accept SSH keys via environment variables
- **Benefits**: Clean separation, no password needed
- **Implementation**:
  - Accept `DEVPOD_SSH_PUBLIC_KEY` environment variable
  - Inject into container during startup
  - Use corresponding private key for connections

#### 5. **DevPod Native SSH Key Support**

- **Approach**: Investigate if DevPod has built-in SSH key management
- **Benefits**: Leverage DevPod's existing capabilities
- **Research**: Check DevPod documentation for SSH key handling

### 🔧 Implementation Priority

1. **Phase 1**: Research DevPod's native SSH key capabilities
2. **Phase 2**: Implement SSH key-based authentication
3. **Phase 3**: Add fallback mechanisms for different environments
4. **Phase 4**: Remove sshpass dependency entirely

### 📋 Research Tasks

- [ ] Study official DevPod SSH provider implementation
- [ ] Check if DevPod automatically handles SSH key injection
- [ ] Research container SSH key injection best practices
- [ ] Test SSH key generation and injection workflow
- [ ] Evaluate security implications of different approaches

### 🎯 Success Criteria

- [ ] No external dependencies required (no sshpass)
- [ ] Works out-of-the-box on macOS, Linux, Windows
- [ ] Maintains security best practices
- [ ] Compatible with DevPod's agent injection mechanism
- [ ] Provides clear error messages for troubleshooting

## Medium Priority

### 🚀 Performance Optimizations

- [ ] Reduce API calls in command section
- [ ] Cache application lookup results
- [ ] Optimize port mapping discovery

### 📚 Documentation Improvements

- [ ] Add troubleshooting guide for SSH issues
- [ ] Document SSH key setup process
- [ ] Add examples for different authentication methods

## Low Priority

### 🔄 Feature Enhancements

- [ ] Support for custom SSH ports
- [ ] Multiple authentication method fallbacks
- [ ] SSH connection pooling/reuse

---

**Note**: The current sshpass implementation works but should be considered a temporary solution until a more robust, dependency-free authentication method is implemented.
