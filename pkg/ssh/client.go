package ssh

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"

	"github.com/NaNomicon/dokploy-devpod-provider/pkg/dokploy"
	"github.com/NaNomicon/dokploy-devpod-provider/pkg/options"
	"github.com/sirupsen/logrus"
)

// Client represents an SSH client for connecting to Dokploy workspaces
type Client struct {
	dokployClient *dokploy.Client
	opts          *options.Options
	logger        *logrus.Logger
}

// NewClient creates a new SSH client
func NewClient(opts *options.Options, logger *logrus.Logger) *Client {
	return &Client{
		dokployClient: dokploy.NewClient(opts, logger),
		opts:          opts,
		logger:        logger,
	}
}

// ExecuteCommand executes a command on the remote workspace via SSH
func (c *Client) ExecuteCommand(machineID, command string) error {
	// Get all projects and applications to find the application by name
	projects, err := c.dokployClient.GetAllProjects()
	if err != nil {
		return fmt.Errorf("failed to retrieve projects: %w", err)
	}

	// Find the application with matching name (machineID is the application name)
	var applicationID string
	for _, project := range projects {
		for _, app := range project.Applications {
			if app.Name == machineID {
				applicationID = app.ApplicationID
				break
			}
		}
		if applicationID != "" {
			break
		}
	}

	if applicationID == "" {
		return fmt.Errorf("no application found with name '%s'", machineID)
	}

	// Get application details including port mappings
	app, err := c.dokployClient.GetApplication(applicationID)
	if err != nil {
		return fmt.Errorf("failed to retrieve application details: %w", err)
	}

	// Extract SSH port mapping from the ports array
	var sshPort int
	for _, port := range app.Ports {
		if port.TargetPort == 22 {
			sshPort = port.PublishedPort
			break
		}
	}

	if sshPort == 0 {
		return fmt.Errorf("no SSH port mapping found for application")
	}

	// Extract Dokploy host from server URL
	parsedURL, err := url.Parse(c.opts.DokployServerURL)
	if err != nil {
		return fmt.Errorf("failed to parse server URL: %w", err)
	}
	sshHost := strings.Split(parsedURL.Host, ":")[0]

	// Check if sshpass is available
	if _, err := exec.LookPath("sshpass"); err != nil {
		return fmt.Errorf("sshpass is required for password authentication. Please install sshpass:\n" +
			"  macOS: brew install hudochenkov/sshpass/sshpass\n" +
			"  Ubuntu/Debian: sudo apt-get install sshpass\n" +
			"  CentOS/RHEL: sudo yum install sshpass")
	}

	// Execute command via SSH using sshpass for password authentication
	sshArgs := []string{
		"-p", "devpod", // password
		"ssh",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "ConnectTimeout=30",
		"-o", "ServerAliveInterval=5",
		"-o", "ServerAliveCountMax=3",
		"-o", "PreferredAuthentications=password",
		"-o", "PasswordAuthentication=yes",
		"-o", "PubkeyAuthentication=no",
		"-p", fmt.Sprintf("%d", sshPort),
		fmt.Sprintf("devpod@%s", sshHost),
		command,
	}

	c.logger.Debugf("Executing SSH command: sshpass %s", strings.Join(sshArgs, " "))

	cmd := exec.Command("sshpass", sshArgs...)
	cmd.Stdout = nil // Let the command output go to stdout directly
	cmd.Stderr = nil // Let the command errors go to stderr directly
	cmd.Stdin = nil  // No stdin needed

	return cmd.Run()
} 