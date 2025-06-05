#!/bin/bash
set -e

# Get SSH public key from environment variable
if [ -z "$SSH_PUBLIC_KEY" ]; then
  echo "ERROR: SSH_PUBLIC_KEY environment variable is not set"
  exit 1
fi

echo "🐳 DOKPLOY DEVPOD PROVIDER - Docker Compose with Privileged Mode (ROOT MODE)"
echo "============================================================================"

echo "Stage 1/4: Starting Docker daemon using DinD built-in script..."
# Start docker using the built-in DinD script
start-docker.sh &

# Wait for Docker daemon to be ready
echo "Waiting for Docker daemon to start..."
for i in $(seq 1 30); do
  if docker info >/dev/null 2>&1; then
    echo "✓ Docker daemon started successfully"
    break
  fi
  if [ $i -eq 30 ]; then
    echo "ERROR: Docker daemon failed to start"
    exit 1
  fi
  sleep 1
done

echo "Stage 2/4: Installing SSH server and tools..."
apt-get update -qq
apt-get install -y -qq openssh-server sudo curl wget ca-certificates gnupg
echo "✓ SSH server and tools installed"

echo "Stage 3/4: Setting up SSH keys for root user..."
mkdir -p /root/.ssh
echo "$SSH_PUBLIC_KEY" > /root/.ssh/authorized_keys
chmod 700 /root/.ssh
chmod 600 /root/.ssh/authorized_keys
echo "✓ SSH keys configured for root"

echo "Stage 4/4: Configuring SSH daemon..."
echo "Port 22" > /etc/ssh/sshd_config
echo "PubkeyAuthentication yes" >> /etc/ssh/sshd_config
echo "AuthorizedKeysFile .ssh/authorized_keys" >> /etc/ssh/sshd_config
echo "PasswordAuthentication yes" >> /etc/ssh/sshd_config
echo "PermitRootLogin yes" >> /etc/ssh/sshd_config
echo "ChallengeResponseAuthentication no" >> /etc/ssh/sshd_config
echo "UsePAM no" >> /etc/ssh/sshd_config
echo "X11Forwarding yes" >> /etc/ssh/sshd_config
echo "PrintMotd no" >> /etc/ssh/sshd_config
echo "AcceptEnv LANG LC_*" >> /etc/ssh/sshd_config
echo "Subsystem sftp /usr/lib/openssh/sftp-server" >> /etc/ssh/sshd_config
echo "✓ SSH daemon configured (root access enabled)"

service ssh start
echo "✓ SSH daemon started"

echo ""
echo "🎉 WORKSPACE READY (ROOT MODE)!"
echo "✓ Docker daemon: Running (privileged mode)"
echo "✓ SSH daemon: Running on port 22"
echo "✓ User: root with full access"
echo "✓ Docker access: Full Docker-in-Docker capability"
echo "✓ Development environment: Ready for DevPod"
echo ""

# Keep container running
tail -f /dev/null 