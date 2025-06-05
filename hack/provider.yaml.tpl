name: dokploy
version: {{ env.Getenv "VERSION" }}
description: |-
  DevPod on Dokploy (Development Version)

  This provider allows you to create and manage development workspaces on Dokploy.

  Features:
  - Automatic SSH setup with key-based authentication
  - Docker-in-Docker support for development containers
  - Port mapping for SSH access
  - Secure workspace isolation

  Requirements:
  - Dokploy server with API access
  - Valid API token with application management permissions

icon: https://raw.githubusercontent.com/Dokploy/dokploy/refs/heads/canary/apps/dokploy/logo.png
optionGroups:
  - options:
      - DOKPLOY_SERVER_URL
      - DOKPLOY_API_TOKEN
    name: "Dokploy Configuration"
    defaultVisible: true
  - options:
      - DOKPLOY_PROJECT_NAME
    name: "Advanced Configuration"
    defaultVisible: false
options:
  DOKPLOY_SERVER_URL:
    description: The URL of your Dokploy server (e.g., https://dokploy.example.com)
    required: true
    suggestions:
      - https://dokploy.example.com
      - https://your-dokploy-server.com
  DOKPLOY_API_TOKEN:
    description: Your Dokploy API token (generate from Settings > API)
    required: true
    password: true
  DOKPLOY_PROJECT_NAME:
    description: Dokploy project name for DevPod workspaces
    default: "devpod-workspaces"

binaries:
  DOKPLOY_PROVIDER_BINARY:
    - os: linux
      arch: amd64
      path: https://github.com/{{ env.Getenv "GITHUB_REPO" }}/releases/download/{{ env.Getenv "VERSION" }}/dokploy-provider-linux-amd64
      checksum: "{{ file.Read "dist/dokploy-provider-linux-amd64" | crypto.SHA256 }}"
    - os: linux
      arch: arm64
      path: https://github.com/{{ env.Getenv "GITHUB_REPO" }}/releases/download/{{ env.Getenv "VERSION" }}/dokploy-provider-linux-arm64
      checksum: "{{ file.Read "dist/dokploy-provider-linux-arm64" | crypto.SHA256 }}"
    - os: darwin
      arch: amd64
      path: https://github.com/{{ env.Getenv "GITHUB_REPO" }}/releases/download/{{ env.Getenv "VERSION" }}/dokploy-provider-darwin-amd64
      checksum: "{{ file.Read "dist/dokploy-provider-darwin-amd64" | crypto.SHA256 }}"
    - os: darwin
      arch: arm64
      path: https://github.com/{{ env.Getenv "GITHUB_REPO" }}/releases/download/{{ env.Getenv "VERSION" }}/dokploy-provider-darwin-arm64
      checksum: "{{ file.Read "dist/dokploy-provider-darwin-arm64" | crypto.SHA256 }}"
    - os: windows
      arch: amd64
      path: https://github.com/{{ env.Getenv "GITHUB_REPO" }}/releases/download/{{ env.Getenv "VERSION" }}/dokploy-provider-windows-amd64.exe
      checksum: "{{ file.Read "dist/dokploy-provider-windows-amd64.exe" | crypto.SHA256 }}"

exec:
  # Initialize the provider (test API connection)
  init: |-
    ${DOKPLOY_PROVIDER_BINARY} init

  # Create a new machine/container in Dokploy
  create: |-
    ${DOKPLOY_PROVIDER_BINARY} create

  # Delete the machine
  delete: |-
    ${DOKPLOY_PROVIDER_BINARY} delete

  # Start a stopped machine
  start: |-
    ${DOKPLOY_PROVIDER_BINARY} start

  # Stop the machine
  stop: |-
    ${DOKPLOY_PROVIDER_BINARY} stop

  # Get machine status
  status: |-
    ${DOKPLOY_PROVIDER_BINARY} status

  # Execute commands via SSH
  command: |-
    ${DOKPLOY_PROVIDER_BINARY} command 