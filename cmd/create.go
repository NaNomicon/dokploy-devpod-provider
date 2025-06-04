package cmd

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/NaNomicon/dokploy-devpod-provider/pkg/dokploy"
	"github.com/NaNomicon/dokploy-devpod-provider/pkg/options"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Dokploy workspace",
	Long: `Create a new development workspace in Dokploy with automatic SSH setup
and Docker-in-Docker support.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCreate()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
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
	
	logger.Infof("Creating Dokploy workspace: %s", machineID)

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

	// Create application in Dokploy
	logger.Info("Creating application...")
	app, err := client.CreateApplication(dokploy.CreateApplicationRequest{
		Name:        machineID,
		Description: fmt.Sprintf("DevPod workspace created on %s", time.Now().Format(time.RFC3339)),
		ProjectID:   projectID,
	})
	if err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}

	logger.Infof("✓ Application created with ID: %s", app.ApplicationID)

	// Configure Docker provider
	logger.Info("Configuring Docker provider...")
	err = client.SaveDockerProvider(dokploy.DockerProviderRequest{
		ApplicationID: app.ApplicationID,
		DockerImage:   "cruizba/ubuntu-dind:latest",
	})
	if err != nil {
		return fmt.Errorf("failed to configure Docker provider: %w", err)
	}

	logger.Info("✓ Docker provider configured")

	// Configure environment variables
	logger.Info("Configuring environment variables...")
	err = client.SaveEnvironment(dokploy.EnvironmentRequest{
		ApplicationID: app.ApplicationID,
		Env:           "DEVPOD_WORKSPACE=true",
	})
	if err != nil {
		return fmt.Errorf("failed to configure environment variables: %w", err)
	}

	logger.Info("✓ Environment variables configured")

	// Update application with SSH setup command
	logger.Info("Configuring application for SSH access...")
	
	sshSetupCommand := `sh -c 'echo "=== DevPod SSH Setup Starting ===" && echo "Stage 1/4: Updating package lists (this may take 1-2 minutes)..." && apt-get update -qq && echo "✓ Package lists updated" && echo "Stage 2/4: Installing SSH server and sudo..." && apt-get install -y -qq openssh-server sudo && echo "✓ SSH server installed" && echo "Stage 3/4: Creating devpod user and configuring permissions..." && useradd -m -s /bin/bash devpod && echo "devpod ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers && (getent group docker >/dev/null || groupadd docker) && usermod -aG docker devpod && mkdir -p /home/devpod/.ssh && chmod 700 /home/devpod/.ssh && chown devpod:devpod /home/devpod/.ssh && touch /home/devpod/.ssh/authorized_keys && chmod 600 /home/devpod/.ssh/authorized_keys && chown devpod:devpod /home/devpod/.ssh/authorized_keys && echo "✓ User devpod configured" && echo "Stage 4/4: Configuring SSH daemon..." && mkdir -p /run/sshd && ssh-keygen -A && echo "PubkeyAuthentication yes" >> /etc/ssh/sshd_config && echo "AuthorizedKeysFile .ssh/authorized_keys" >> /etc/ssh/sshd_config && echo "PasswordAuthentication no" >> /etc/ssh/sshd_config && echo "PermitRootLogin no" >> /etc/ssh/sshd_config && echo "Port 22" >> /etc/ssh/sshd_config && echo "ChallengeResponseAuthentication no" >> /etc/ssh/sshd_config && echo "UsePAM no" >> /etc/ssh/sshd_config && echo "X11Forwarding yes" >> /etc/ssh/sshd_config && echo "PrintMotd no" >> /etc/ssh/sshd_config && echo "AcceptEnv LANG LC_*" >> /etc/ssh/sshd_config && echo "Subsystem sftp /usr/lib/openssh/sftp-server" >> /etc/ssh/sshd_config && echo "✓ SSH daemon configured for key-only authentication" && echo "=== Starting SSH daemon ===" && /usr/sbin/sshd -D -e &
echo "SSH daemon started in background" && echo "🎉 SSH SETUP COMPLETE - DevPod status command will verify readiness" && sleep 5 && echo "Container ready for DevPod connection" && tail -f /dev/null'`

	err = client.UpdateApplication(dokploy.UpdateApplicationRequest{
		ApplicationID: app.ApplicationID,
		Command:       sshSetupCommand,
	})
	if err != nil {
		return fmt.Errorf("failed to update application with SSH setup: %w", err)
	}

	logger.Info("✓ Application configured with SSH setup command")

	// Deploy the application
	logger.Info("Deploying application with SSH configuration...")
	logger.Info("NOTE: Using Ubuntu Docker-in-Docker container with pre-installed Docker")
	logger.Info("      DevPod will connect via SSH and inject its agent automatically")
	logger.Info("")
	logger.Info("⏳ Container Setup Timeline:")
	logger.Info("   • Stage 1: Package update (1-2 minutes) - apt-get update")
	logger.Info("   • Stage 2: SSH installation (30-60 seconds) - install openssh-server")
	logger.Info("   • Stage 3: User setup (10-20 seconds) - create devpod user")
	logger.Info("   • Stage 4: SSH configuration (10-20 seconds) - configure SSH daemon")
	logger.Info("   • Total estimated time: 2-4 minutes")
	logger.Info("")

	err = client.DeployApplication(dokploy.DeployRequest{
		ApplicationID: app.ApplicationID,
	})
	if err != nil {
		return fmt.Errorf("failed to deploy application: %w", err)
	}

	logger.Info("✓ Deployment started with SSH configuration")

	// Wait for deployment to complete
	logger.Info("Waiting for container setup to complete...")
	logger.Info("ℹ️  The container is now running the 4-stage SSH setup process.")
	logger.Info("   If this takes longer than expected, the apt-get update stage may be slow.")
	logger.Info("")

	setupStartTime := time.Now()
	for i := 1; i <= 60; i++ {
		currentApp, err := client.GetApplication(app.ApplicationID)
		if err != nil {
			logger.Warnf("Failed to get application status: %v", err)
			continue
		}

		elapsedTime := time.Since(setupStartTime)
		
		if currentApp.Status == "done" {
			logger.Infof("✓ Container setup completed successfully (%v elapsed)", elapsedTime)
			break
		} else if currentApp.Status == "error" {
			logger.Warn("⚠️  Container setup failed, but continuing...")
			break
		}

		// Provide stage-specific feedback based on elapsed time
		var stageInfo string
		if elapsedTime < time.Minute {
			stageInfo = "(likely Stage 1: apt-get update)"
		} else if elapsedTime < 2*time.Minute {
			stageInfo = "(likely Stage 1-2: package update/SSH install)"
		} else if elapsedTime < 3*time.Minute {
			stageInfo = "(likely Stage 2-3: SSH install/user setup)"
		} else {
			stageInfo = "(likely Stage 3-4: user/SSH configuration)"
		}

		logger.Infof("   Setup status: %s - %v elapsed %s (attempt %d/60)", currentApp.Status, elapsedTime, stageInfo, i)
		time.Sleep(5 * time.Second)
	}

	// Configure SSH port mapping
	logger.Info("Configuring SSH port mapping via Dokploy API...")

	// Extract Dokploy host from server URL
	parsedURL, err := url.Parse(opts.DokployServerURL)
	if err != nil {
		return fmt.Errorf("failed to parse server URL: %w", err)
	}
	sshHost := strings.Split(parsedURL.Host, ":")[0]

	// Try to find an available port for SSH mapping
	var sshHostPort int
	portCreated := false
	logger.Info("Searching for available SSH port (trying ports 2222-2250)...")

	// First, get all existing applications to check for port conflicts
	allProjects, err := client.GetAllProjects()
	if err != nil {
		logger.Warnf("Failed to get projects for port conflict check: %v", err)
	}

	// Build a map of used ports from API data
	usedPorts := make(map[int]bool)
	if allProjects != nil {
		for _, project := range allProjects {
			for _, application := range project.Applications {
				for _, port := range application.Ports {
					if port.Protocol == "tcp" {
						usedPorts[port.PublishedPort] = true
						logger.Debugf("Port %d is already used by application %s", port.PublishedPort, application.Name)
					}
				}
			}
		}
	}

	// Test actual network connectivity to find truly available ports
	logger.Info("Testing network connectivity to find available ports...")
	for port := 2222; port <= 2250; port++ {
		if portCreated {
			break
		}

		// Check if port is already used by another application according to API
		if usedPorts[port] {
			logger.Infof("Port %d is already in use by another application (API check), skipping...", port)
			continue
		}

		// Test actual network connectivity to see if port is really available
		testAddress := fmt.Sprintf("%s:%d", sshHost, port)
		logger.Debugf("Testing connectivity to %s...", testAddress)
		
		conn, err := net.DialTimeout("tcp", testAddress, 3*time.Second)
		if err == nil {
			// Port is already in use by something
			conn.Close()
			logger.Infof("Port %d is already in use (network test), skipping...", port)
			usedPorts[port] = true
			continue
		}
		
		// Port appears to be available, try to create the mapping
		logger.Infof("Port %d appears available, attempting to create mapping...", port)

		err = client.CreatePort(dokploy.CreatePortRequest{
			PublishedPort: port,
			TargetPort:    22,
			Protocol:      "tcp",
			ApplicationID: app.ApplicationID,
		})

		if err != nil {
			logger.Infof("   Port %d API creation failed: %v", port, err)
			usedPorts[port] = true
			continue
		}

		logger.Infof("✅ Port %d API mapping created successfully!", port)
		logger.Infof("   Mapping: %d (host) → 22 (container)", port)
		
		// Store the port for later verification after deployment
		sshHostPort = port
		portCreated = true
		break
	}

	if !portCreated {
		return fmt.Errorf("could not configure SSH port mapping - all ports in range 2222-2250 are in use")
	}

	logger.Info("")
	logger.Info("🎉 SSH port mapping configured successfully!")
	logger.Infof("   Using port: %d", sshHostPort)

	// Redeploy to apply the port mapping
	logger.Info("Redeploying application to apply port mapping...")
	err = client.DeployApplication(dokploy.DeployRequest{
		ApplicationID: app.ApplicationID,
	})
	if err != nil {
		return fmt.Errorf("failed to redeploy application: %w", err)
	}

	// Wait for redeployment
	logger.Info("Waiting for redeployment to complete...")
	for i := 1; i <= 30; i++ {
		currentApp, err := client.GetApplication(app.ApplicationID)
		if err != nil {
			logger.Warnf("Failed to get application status: %v", err)
			continue
		}

		if currentApp.Status == "done" {
			logger.Info("✓ Redeployment completed successfully")
			break
		} else if currentApp.Status == "error" {
			logger.Warn("⚠️  Redeployment failed, but continuing...")
			break
		}

		logger.Infof("   Redeployment status: %s (attempt %d/30)", currentApp.Status, i)
		time.Sleep(5 * time.Second)
	}

	// Additional wait for container to fully initialize SSH service
	logger.Info("Waiting for container to fully initialize SSH service...")
	time.Sleep(15 * time.Second)

	logger.Info("")
	logger.Info("✅ Dokploy workspace created successfully!")
	logger.Info("🎉 Machine creation completed!")
	logger.Info("")
	logger.Info("Workspace Details:")
	logger.Infof("- Application ID: %s", app.ApplicationID)
	logger.Infof("- SSH Host: %s", sshHost)
	logger.Infof("- SSH Port: %d", sshHostPort)
	logger.Info("- SSH User: devpod")
	logger.Info("- SSH Auth: Key-only authentication (password disabled)")
	logger.Info("- Base Image: cruizba/ubuntu-dind (Ubuntu with Docker-in-Docker)")
	logger.Infof("- Dokploy Dashboard: %s", opts.DokployServerURL)
	logger.Info("")
	logger.Info("Next Steps:")
	logger.Info("- DevPod will check SSH readiness via status command")
	logger.Info("- Once ready, DevPod will connect and inject its agent")
	logger.Info("- Docker is pre-installed and ready for development containers")

	// Return connection info to DevPod (MUST BE LAST OUTPUT)
	fmt.Printf("DEVPOD_MACHINE_ID=%s\n", machineID)
	fmt.Printf("DEVPOD_MACHINE_HOST=%s\n", sshHost)
	fmt.Printf("DEVPOD_MACHINE_PORT=%d\n", sshHostPort)
	fmt.Printf("DEVPOD_MACHINE_USER=devpod\n")

	return nil
} 