package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/NaNomicon/dokploy-devpod-provider/pkg/options"
	"github.com/NaNomicon/dokploy-devpod-provider/pkg/dokploy"
	"github.com/loft-sh/devpod/pkg/ssh"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// commandCmd represents the command execution command
var commandCmd = &cobra.Command{
	Use:   "command",
	Short: "Execute a command on a Dokploy workspace via SSH",
	Long:  `Execute a command on a remote development workspace in Dokploy via SSH.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCommand()
	},
}

func init() {
	rootCmd.AddCommand(commandCmd)
}

func runCommand() error {
	// Setup logger with stderr output to avoid interfering with command output
	logger := logrus.New()
	logger.SetOutput(os.Stderr)
	if verbose {
		logger.SetLevel(logrus.DebugLevel)
	}

	// Debug: Print all environment variables
	logger.Debug("All environment variables:")
	for _, env := range os.Environ() {
		logger.Debugf("  %s", env)
	}

	// Get machine ID from environment (set by DevPod)
	machineID := os.Getenv("MACHINE_ID")
	if machineID == "" {
		return fmt.Errorf("MACHINE_ID environment variable is missing")
	}

	// Get command from environment (set by DevPod)
	command := os.Getenv("COMMAND")
	if command == "" {
		return fmt.Errorf("COMMAND environment variable is missing")
	}

	logger.Debugf("Executing command on machine %s: %s", machineID, command)

	// Load options from environment
	opts, err := options.LoadFromEnv()
	if err != nil {
		return fmt.Errorf("failed to load options: %w", err)
	}

	// Create Dokploy client to get SSH connection details
	dokployClient := dokploy.NewClient(opts, logger)

	// Get application details by finding it by name (machineID is the application name)
	// First get all projects to find the application
	projects, err := dokployClient.GetAllProjects()
	if err != nil {
		return fmt.Errorf("failed to get projects: %w", err)
	}

	// Find the application with matching name
	var app *dokploy.Application
	for _, project := range projects {
		for _, application := range project.Applications {
			if application.Name == machineID {
				app = &application
				break
			}
		}
		if app != nil {
			break
		}
	}

	if app == nil {
		return fmt.Errorf("application with name '%s' not found", machineID)
	}

	logger.Debugf("Found application: %s (ID: %s)", app.Name, app.ApplicationID)

	// Get full application details using the ApplicationID to ensure we have complete port information
	// This is important because the project listing might not include all port details
	fullApp, err := dokployClient.GetApplication(app.ApplicationID)
	if err != nil {
		logger.Warnf("Failed to get full application details, using basic info: %v", err)
		fullApp = app
	} else {
		logger.Debugf("Retrieved full application details with %d ports", len(fullApp.Ports))
	}

	// Find SSH port from application ports with retry logic
	var sshPort string
	maxRetries := 10
	retryDelay := 2 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		logger.Debugf("Attempt %d/%d: Looking for SSH port in %d ports", attempt, maxRetries, len(fullApp.Ports))
		
		for i, port := range fullApp.Ports {
			logger.Debugf("  Port %d: %d -> %d (%s)", i, port.PublishedPort, port.TargetPort, port.Protocol)
			if port.TargetPort == 22 {
				sshPort = fmt.Sprintf("%d", port.PublishedPort)
				logger.Debugf("Found SSH port: %s", sshPort)
				break
			}
		}

		if sshPort != "" {
			break
		}

		if attempt < maxRetries {
			logger.Debugf("SSH port not found, waiting %v before retry %d/%d...", retryDelay, attempt+1, maxRetries)
			time.Sleep(retryDelay)
			
			// Refresh application details
			freshApp, err := dokployClient.GetApplication(app.ApplicationID)
			if err != nil {
				logger.Warnf("Failed to refresh application details: %v", err)
			} else {
				fullApp = freshApp
				logger.Debugf("Refreshed application details, now has %d ports", len(fullApp.Ports))
			}
		}
	}

	if sshPort == "" {
		return fmt.Errorf("SSH port not found for application %s", machineID)
	}

	// Get machine folder for SSH keys
	machineFolder := os.Getenv("MACHINE_FOLDER")
	if machineFolder == "" {
		return fmt.Errorf("MACHINE_FOLDER environment variable is missing")
	}

	// Get private key for SSH authentication
	privateKey, err := ssh.GetPrivateKeyRawBase(machineFolder)
	if err != nil {
		return fmt.Errorf("failed to load private key: %w", err)
	}

	// Create SSH client using DevPod's SSH utilities
	// Extract hostname from ServerURL (remove https:// prefix if present)
	hostname := opts.DokployServerURL
	if strings.HasPrefix(hostname, "https://") {
		hostname = strings.TrimPrefix(hostname, "https://")
	}
	if strings.HasPrefix(hostname, "http://") {
		hostname = strings.TrimPrefix(hostname, "http://")
	}
	
	sshAddress := hostname + ":" + sshPort
	sshClient, err := ssh.NewSSHClient("devpod", sshAddress, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create SSH client: %w", err)
	}
	defer sshClient.Close()

	// Execute the command via SSH
	return ssh.Run(context.Background(), sshClient, command, os.Stdin, os.Stdout, os.Stderr, nil)
} 