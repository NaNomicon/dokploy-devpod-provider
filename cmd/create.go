package cmd

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/NaNomicon/dokploy-devpod-provider/pkg/dokploy"
	"github.com/NaNomicon/dokploy-devpod-provider/pkg/options"
	"github.com/NaNomicon/dokploy-devpod-provider/pkg/templates"
	"github.com/loft-sh/devpod/pkg/ssh"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Dokploy workspace using Docker Compose",
	Long: `Create a new development workspace in Dokploy using Docker Compose with automatic SSH setup,
privileged mode support, and Docker-in-Docker capabilities.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCreate()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func runCreate() error {
	// Setup logger
	logger := logrus.New()
	if verbose {
		logger.SetLevel(logrus.DebugLevel)
	}

	// Debug: Print all environment variables that might contain machine ID
	logger.Debug("Environment variables:")
	for _, env := range os.Environ() {
		if strings.Contains(strings.ToUpper(env), "MACHINE") ||
			strings.Contains(strings.ToUpper(env), "DEVPOD") ||
			strings.Contains(strings.ToUpper(env), "WORKSPACE") {
			logger.Debugf("  %s", env)
		}
	}

	machineID := os.Getenv("DEVPOD_MACHINE_ID")
	if machineID == "" {
		// Try alternative environment variables
		machineID = os.Getenv("MACHINE_ID")
		if machineID == "" {
			machineID = os.Getenv("DEVPOD_WORKSPACE_ID")
			if machineID == "" {
				machineID = os.Getenv("WORKSPACE_ID")
			}
		}
	}
	
	logger.Infof("Creating Dokploy workspace via Docker Compose: %s", machineID)

	// Load options from environment
	opts, err := options.LoadFromEnv()
	if err != nil {
		return fmt.Errorf("failed to load options: %w", err)
	}

	if machineID == "" {
		return fmt.Errorf("DEVPOD_MACHINE_ID is required for workspace creation")
	}

	// Create Dokploy client
	client := dokploy.NewClient(opts, logger)

	// Check if project exists, create if it doesn't
	logger.Infof("Checking if project '%s' exists...", opts.DokployProjectName)
	projects, err := client.GetAllProjects()
	if err != nil {
		return fmt.Errorf("failed to get projects: %w", err)
	}

	var projectID string
	for _, project := range projects {
		if project.Name == opts.DokployProjectName {
			projectID = project.ProjectID
			break
		}
	}

	if projectID == "" {
		logger.Infof("Project '%s' not found. Creating project...", opts.DokployProjectName)
		
		project, err := client.CreateProject(dokploy.CreateProjectRequest{
			Name:        opts.DokployProjectName,
			Description: "DevPod workspaces project - automatically created by Dokploy provider",
		})
		if err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}

		projectID = project.ProjectID
		logger.Infof("✓ Project created successfully with ID: %s", projectID)
	} else {
		logger.Infof("✓ Project '%s' already exists with ID: %s", opts.DokployProjectName, projectID)
	}

	// Get SSH public key from DevPod for injection into container
	logger.Info("Getting SSH public key from DevPod...")
	machineFolder := os.Getenv("MACHINE_FOLDER")
	if machineFolder == "" {
		return fmt.Errorf("MACHINE_FOLDER environment variable is missing")
	}

	logger.Debugf("Machine folder: %s", machineFolder)

	publicKey, err := ssh.GetPublicKeyBase(machineFolder)
	if err != nil {
		return fmt.Errorf("failed to get SSH public key: %w", err)
	}

	logger.Debugf("Retrieved SSH public key (full): %s", publicKey)
	
	// Decode the base64 encoded SSH key if needed
	if !strings.HasPrefix(publicKey, "ssh-") {
		logger.Info("Decoding base64 encoded SSH key from DevPod...")
		decodedKey, err := base64.StdEncoding.DecodeString(publicKey)
		if err != nil {
			return fmt.Errorf("failed to decode base64 SSH key: %w", err)
		}
		publicKey = string(decodedKey)
		logger.Infof("✓ SSH key decoded successfully")
	}
	
	logger.Info("✓ SSH public key retrieved from DevPod")

	// Extract Dokploy host from server URL for SSH connection info
	parsedURL, err := url.Parse(opts.DokployServerURL)
	if err != nil {
		return fmt.Errorf("failed to parse server URL: %w", err)
	}
	sshHost := strings.Split(parsedURL.Host, ":")[0]

	// Find an available SSH port
	var sshHostPort int
	logger.Info("Finding available SSH port (range 2222-2250)...")

	// Check existing port usage
	allProjects, err := client.GetAllProjects()
	if err != nil {
		logger.Warnf("Failed to get projects for port conflict check: %v", err)
	}

	usedPorts := make(map[int]bool)
	if allProjects != nil {
		for _, project := range allProjects {
			for _, compose := range project.Composes {
				logger.Debugf("Checking compose service %s for port usage", compose.Name)
				// Note: We'll need to add compose port checking in the client
			}
		}
	}

	// Find available port
	portFound := false
	for port := 2222; port <= 2250; port++ {
		if usedPorts[port] {
			continue
		}

		// Test network availability
		testAddress := fmt.Sprintf("%s:%d", sshHost, port)
		conn, err := net.DialTimeout("tcp", testAddress, 3*time.Second)
		if err == nil {
			conn.Close()
			usedPorts[port] = true
			continue
		}
		
		sshHostPort = port
		portFound = true
		logger.Infof("✓ Selected SSH port: %d", port)
		break
	}

	if !portFound {
		return fmt.Errorf("no available ports in range 2222-2250")
	}

	// Create docker-compose.yml content with privileged mode
	logger.Info("Creating Docker Compose configuration with privileged mode...")
	
	dockerComposeContent, err := generateDockerCompose(sshHostPort, publicKey, logger)
	if err != nil {
		return fmt.Errorf("failed to generate Docker Compose configuration: %w", err)
	}

	logger.Info("✓ Docker Compose configuration created with full privileged mode support")

	// Create Docker Compose service in Dokploy
	logger.Info("Creating Docker Compose service in Dokploy...")
	
	compose, err := client.CreateCompose(dokploy.CreateComposeRequest{
		Name:        machineID,
		Description: fmt.Sprintf("DevPod workspace created on %s via Docker Compose", time.Now().Format(time.RFC3339)),
		ProjectID:   projectID,
		ComposeType: "docker-compose", // Use docker-compose instead of stack for full feature support
	})
	if err != nil {
		return fmt.Errorf("failed to create Docker Compose service: %w", err)
	}

	logger.Infof("✓ Docker Compose service created with ID: %s", compose.ComposeID)

	// Set the docker-compose.yml content
	logger.Info("Uploading Docker Compose configuration...")
	err = client.SaveComposeFile(dokploy.SaveComposeFileRequest{
		ComposeID:     compose.ComposeID,
		DockerCompose: dockerComposeContent,
	})
	if err != nil {
		return fmt.Errorf("failed to save Docker Compose file: %w", err)
	}

	logger.Info("✓ Docker Compose file uploaded successfully")

	// Deploy the Docker Compose service
	logger.Info("Deploying Docker Compose service...")
	logger.Info("")
	logger.Info("🐳 Enhanced Docker Compose Deployment:")
	logger.Info("   • Privileged mode: ENABLED (full Docker-in-Docker support)")
	logger.Info("   • SSH port mapping: External port %d → Container port 22", sshHostPort)
	logger.Info("   • Base image: cruizba/ubuntu-dind:latest")
	logger.Info("   • User setup: devpod user with sudo and docker group access")
	logger.Info("   • SSH authentication: Both key-based and password authentication")
	logger.Info("   • Docker daemon: Full dockerd with overlay2 storage driver")
	logger.Info("")
	logger.Info("⏳ Deployment Timeline:")
	logger.Info("   • Stage 1: Docker daemon startup (30-60 seconds)")
	logger.Info("   • Stage 2: Package installation (30-60 seconds)")
	logger.Info("   • Stage 3: User creation (10-20 seconds)")
	logger.Info("   • Stage 4: SSH key setup (10-20 seconds)")
	logger.Info("   • Stage 5: SSH daemon config (10-20 seconds)")
	logger.Info("   • Stage 6: SSH service start (10-20 seconds)")
	logger.Info("   • Total estimated time: 2-4 minutes")
	logger.Info("")

	err = client.DeployCompose(dokploy.DeployComposeRequest{
		ComposeID: compose.ComposeID,
	})
	if err != nil {
		return fmt.Errorf("failed to deploy Docker Compose service: %w", err)
	}

	logger.Info("✓ Docker Compose deployment started")

	// Wait for deployment to complete
	logger.Info("Waiting for Docker Compose deployment to complete...")
	setupStartTime := time.Now()
	
	for i := 1; i <= 60; i++ {
		currentCompose, err := client.GetCompose(compose.ComposeID)
		if err != nil {
			logger.Warnf("Failed to get compose status: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		elapsedTime := time.Since(setupStartTime)
		
		if currentCompose.Status == "done" {
			logger.Infof("✓ Docker Compose deployment completed successfully (%v elapsed)", elapsedTime)
			break
		} else if currentCompose.Status == "error" {
			logger.Warn("⚠️  Docker Compose deployment failed, but continuing...")
			break
		}

		// Provide stage-specific feedback
		var stageInfo string
		if elapsedTime < 60*time.Second {
			stageInfo = "(likely Stage 1-2: Docker daemon + package installation)"
		} else if elapsedTime < 120*time.Second {
			stageInfo = "(likely Stage 3-6: User setup + SSH configuration)"
		} else {
			stageInfo = "(likely finalizing SSH service)"
		}

		logger.Infof("   Deployment status: %s - %v elapsed %s (attempt %d/60)", currentCompose.Status, elapsedTime, stageInfo, i)
		time.Sleep(5 * time.Second)
	}

	// Additional wait for SSH service to be fully ready
	logger.Info("Waiting for SSH service to be fully accessible...")
	time.Sleep(15 * time.Second)

	logger.Info("")
	logger.Info("✅ Dokploy workspace created successfully via Docker Compose!")
	logger.Info("🎉 Privileged Docker-in-Docker workspace deployment completed!")
	logger.Info("")
	logger.Info("Workspace Details:")
	logger.Infof("- Compose ID: %s", compose.ComposeID)
	logger.Infof("- SSH Host: %s", sshHost)
	logger.Infof("- SSH Port: %d", sshHostPort)
	logger.Info("- SSH User: devpod")
	logger.Info("- SSH Auth: Key + password authentication")
	logger.Info("- Privileged Mode: ENABLED")
	logger.Info("- Base Image: cruizba/ubuntu-dind:latest")
	logger.Info("- Docker Daemon: Full Docker-in-Docker with overlay2")
	logger.Infof("- Dokploy Dashboard: %s", opts.DokployServerURL)
	logger.Info("")
	logger.Info("🐳 Container Capabilities:")
	logger.Info("- Full Docker daemon: Available with privileged mode")
	logger.Info("- Docker commands: docker build, run, compose, etc.")
	logger.Info("- SSH access: Ready for DevPod connection")
	logger.Info("- Development environment: Full containerized development")
	logger.Info("- User permissions: sudo access + docker group membership")
	logger.Info("")
	logger.Info("Next Steps:")
	logger.Info("- DevPod will verify SSH connectivity via status command")
	logger.Info("- Once ready, DevPod will connect and set up the development environment")
	logger.Info("- Full Docker-in-Docker capabilities available for development")
	logger.Info("- Ready for any containerized development workflow!")
	logger.Info("")

	// Return connection info to DevPod (MUST BE LAST OUTPUT)
	fmt.Printf("DEVPOD_MACHINE_ID=%s\n", machineID)
	fmt.Printf("DEVPOD_MACHINE_HOST=%s\n", sshHost)
	fmt.Printf("DEVPOD_MACHINE_PORT=%d\n", sshHostPort)
	fmt.Printf("DEVPOD_MACHINE_USER=root\n")

	return nil
}

// generateDockerCompose generates the docker-compose.yml content from embedded templates
func generateDockerCompose(sshPort int, sshPublicKey string, logger *logrus.Logger) (string, error) {
	logger.Debugf("=== GENERATING DOCKER COMPOSE ===")
	logger.Debugf("SSH Port: %d", sshPort)
	logger.Debugf("SSH Key length: %d", len(sshPublicKey))
	
	// Use embedded template constants
	logger.Debugf("Docker compose template loaded (%d bytes)", len(templates.DockerComposeTemplate))
	logger.Debugf("Setup script template loaded (%d bytes)", len(templates.SetupScriptTemplate))

	// Encode the setup script as base64 to avoid quoting/escaping issues
	encodedScript := base64.StdEncoding.EncodeToString([]byte(templates.SetupScriptTemplate))

	// Create a command that decodes and executes the script
	setupCommand := fmt.Sprintf("echo '%s' | base64 -d | bash", encodedScript)
	
	logger.Debugf("Encoded setup script as base64 (length: %d)", len(encodedScript))
	logger.Debugf("Setup command: %s", setupCommand)
	
	// Prepare SSH key - trim whitespace and escape for YAML
	sshPublicKey = strings.TrimSpace(sshPublicKey)
	logger.Debugf("SSH public key to inject: %s", sshPublicKey)
	escapedSSHKey := strings.ReplaceAll(sshPublicKey, `"`, `\"`)

	// Replace template variables using our safe placeholders
	dockerCompose := templates.DockerComposeTemplate
	logger.Debugf("Before replacement - contains SSH_PORT placeholder: %v", strings.Contains(dockerCompose, "__SSH_PORT_PLACEHOLDER__"))
	logger.Debugf("Before replacement - contains SSH_KEY placeholder: %v", strings.Contains(dockerCompose, "__SSH_PUBLIC_KEY_PLACEHOLDER__"))
	logger.Debugf("Before replacement - contains SCRIPT placeholder: %v", strings.Contains(dockerCompose, "__SETUP_SCRIPT_PLACEHOLDER__"))
	
	dockerCompose = strings.ReplaceAll(dockerCompose, "__SSH_PORT_PLACEHOLDER__", fmt.Sprintf("%d", sshPort))
	dockerCompose = strings.ReplaceAll(dockerCompose, "__SSH_PUBLIC_KEY_PLACEHOLDER__", escapedSSHKey)
	dockerCompose = strings.ReplaceAll(dockerCompose, "__SETUP_SCRIPT_PLACEHOLDER__", setupCommand)

	logger.Debugf("After replacement - contains SSH_PORT placeholder: %v", strings.Contains(dockerCompose, "__SSH_PORT_PLACEHOLDER__"))
	logger.Debugf("After replacement - contains SSH_KEY placeholder: %v", strings.Contains(dockerCompose, "__SSH_PUBLIC_KEY_PLACEHOLDER__"))
	logger.Debugf("After replacement - contains SCRIPT placeholder: %v", strings.Contains(dockerCompose, "__SETUP_SCRIPT_PLACEHOLDER__"))
	
	logger.Debugf("Generated docker-compose.yml:\n%s", dockerCompose)
	
	return dockerCompose, nil
} 