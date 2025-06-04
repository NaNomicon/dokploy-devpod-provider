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
	
	// Always enable debug logging for command execution to understand the issue
	logger.SetLevel(logrus.DebugLevel)
	
	logger.Debug("=== COMMAND EXECUTION DEBUG START ===")
	logger.Debugf("Current working directory: %s", func() string {
		if cwd, err := os.Getwd(); err == nil {
			return cwd
		}
		return "unknown"
	}())
	
	// Debug: Print ALL environment variables
	logger.Debug("=== ALL ENVIRONMENT VARIABLES ===")
	envVars := os.Environ()
	for i, env := range envVars {
		logger.Debugf("ENV[%d]: %s", i, env)
	}
	logger.Debugf("Total environment variables: %d", len(envVars))

	// Debug: Print command line arguments
	logger.Debug("=== COMMAND LINE ARGUMENTS ===")
	for i, arg := range os.Args {
		logger.Debugf("ARG[%d]: %s", i, arg)
	}
	logger.Debugf("Total arguments: %d", len(os.Args))

	// Debug: Check stdin status
	logger.Debug("=== STDIN STATUS ===")
	stat, err := os.Stdin.Stat()
	if err != nil {
		logger.Debugf("Failed to stat stdin: %v", err)
	} else {
		logger.Debugf("Stdin mode: %v", stat.Mode())
		logger.Debugf("Stdin size: %d", stat.Size())
		logger.Debugf("Is character device: %v", (stat.Mode() & os.ModeCharDevice) != 0)
		logger.Debugf("Is pipe: %v", (stat.Mode() & os.ModeNamedPipe) != 0)
		logger.Debugf("Is regular file: %v", stat.Mode().IsRegular())
	}

	// Debug: Check specific environment variables that DevPod might use
	logger.Debug("=== DEVPOD-SPECIFIC ENVIRONMENT VARIABLES ===")
	devpodVars := []string{
		"MACHINE_ID", "DEVPOD_MACHINE_ID", "COMMAND", "DEVPOD_COMMAND",
		"MACHINE_FOLDER", "DEVPOD_MACHINE_FOLDER", "WORKSPACE_ID", "DEVPOD_WORKSPACE_ID",
		"DEVPOD_PROVIDER_DEV", "DOKPLOY_PROVIDER_DEV", "DEVPOD_DEBUG", "DEBUG",
	}
	for _, varName := range devpodVars {
		value := os.Getenv(varName)
		if value != "" {
			logger.Debugf("%s = %s", varName, value)
		} else {
			logger.Debugf("%s = <not set>", varName)
		}
	}

	// Get machine ID from environment (set by DevPod)
	machineID := os.Getenv("MACHINE_ID")
	if machineID == "" {
		machineID = os.Getenv("DEVPOD_MACHINE_ID")
		if machineID == "" {
			logger.Error("=== MACHINE ID NOT FOUND ===")
			logger.Error("Neither MACHINE_ID nor DEVPOD_MACHINE_ID environment variables are set")
			logger.Error("This indicates DevPod is not properly calling the provider")
			return fmt.Errorf("MACHINE_ID or DEVPOD_MACHINE_ID environment variable is missing")
		}
	}
	logger.Debugf("Using machine ID: %s", machineID)

	// Get command from environment (set by DevPod)
	command := os.Getenv("COMMAND")
	if command == "" {
		// Try alternative environment variable names
		command = os.Getenv("DEVPOD_COMMAND")
		if command == "" {
			logger.Error("=== COMMAND NOT FOUND ===")
			logger.Error("Neither COMMAND nor DEVPOD_COMMAND environment variables are set")
			logger.Error("This indicates DevPod is not properly passing the command to execute")
			
			// Try to read from stdin as a last resort
			logger.Debug("Attempting to read command from stdin...")
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				// Data might be available on stdin
				logger.Debug("Stdin appears to have data, attempting to read...")
				var buf [4096]byte
				n, err := os.Stdin.Read(buf[:])
				if err != nil {
					logger.Debugf("Failed to read from stdin: %v", err)
				} else if n > 0 {
					stdinContent := strings.TrimSpace(string(buf[:n]))
					logger.Debugf("Read %d bytes from stdin: %q", n, stdinContent)
					if stdinContent != "" {
						command = stdinContent
						logger.Debugf("Using command from stdin: %s", command)
					}
				} else {
					logger.Debug("No data available on stdin")
				}
			} else {
				logger.Debug("Stdin is a character device (terminal), no data to read")
			}
			
			if command == "" {
				logger.Error("=== COMMAND EXECUTION FAILED ===")
				logger.Error("No command found in environment variables or stdin")
				logger.Error("Expected: DevPod should set COMMAND environment variable")
				logger.Error("This is likely a DevPod provider integration issue")
				return fmt.Errorf("COMMAND environment variable is missing and no command available on stdin")
			}
		}
	}

	logger.Debugf("=== COMMAND TO EXECUTE ===")
	logger.Debugf("Machine ID: %s", machineID)
	logger.Debugf("Command: %s", command)
	logger.Debugf("Command length: %d characters", len(command))

	// Load options from environment
	opts, err := options.LoadFromEnv()
	if err != nil {
		logger.Errorf("Failed to load options: %v", err)
		return fmt.Errorf("failed to load options: %w", err)
	}
	logger.Debug("✓ Options loaded successfully")

	// Create Dokploy client to get SSH connection details
	logger.Debug("=== CREATING DOKPLOY CLIENT ===")
	dokployClient := dokploy.NewClient(opts, logger)
	logger.Debug("✓ Dokploy client created")

	// Get application details by finding it by name (machineID is the application name)
	// First get all projects to find the application
	logger.Debug("=== FINDING APPLICATION ===")
	logger.Debugf("Searching for application with name: %s", machineID)
	projects, err := dokployClient.GetAllProjects()
	if err != nil {
		logger.Errorf("Failed to get projects: %v", err)
		return fmt.Errorf("failed to get projects: %w", err)
	}
	logger.Debugf("Retrieved %d projects", len(projects))

	// Find the application with matching name
	var app *dokploy.Application
	for i, project := range projects {
		logger.Debugf("Checking project %d: %s (ID: %s) with %d applications", i, project.Name, project.ProjectID, len(project.Applications))
		for j, application := range project.Applications {
			logger.Debugf("  Application %d: %s (ID: %s)", j, application.Name, application.ApplicationID)
			if application.Name == machineID {
				app = &application
				logger.Debugf("✓ Found matching application: %s", application.Name)
				break
			}
		}
		if app != nil {
			break
		}
	}

	if app == nil {
		logger.Errorf("Application with name '%s' not found", machineID)
		logger.Error("Available applications:")
		for _, project := range projects {
			for _, application := range project.Applications {
				logger.Errorf("  - %s (ID: %s)", application.Name, application.ApplicationID)
			}
		}
		return fmt.Errorf("application with name '%s' not found", machineID)
	}

	logger.Debugf("✓ Found application: %s (ID: %s)", app.Name, app.ApplicationID)

	// Get full application details using the ApplicationID to ensure we have complete port information
	logger.Debug("=== GETTING APPLICATION DETAILS ===")
	fullApp, err := dokployClient.GetApplication(app.ApplicationID)
	if err != nil {
		logger.Warnf("Failed to get full application details, using basic info: %v", err)
		fullApp = app
	} else {
		logger.Debugf("✓ Retrieved full application details with %d ports", len(fullApp.Ports))
	}

	// Find SSH port from application ports with retry logic
	logger.Debug("=== FINDING SSH PORT ===")
	var sshPort string
	maxRetries := 10
	retryDelay := 2 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		logger.Debugf("Attempt %d/%d: Looking for SSH port in %d ports", attempt, maxRetries, len(fullApp.Ports))
		
		for i, port := range fullApp.Ports {
			logger.Debugf("  Port %d: %d -> %d (%s)", i, port.PublishedPort, port.TargetPort, port.Protocol)
			if port.TargetPort == 22 {
				sshPort = fmt.Sprintf("%d", port.PublishedPort)
				logger.Debugf("✓ Found SSH port: %s", sshPort)
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
		logger.Errorf("SSH port not found for application %s", machineID)
		logger.Error("Available ports:")
		for i, port := range fullApp.Ports {
			logger.Errorf("  Port %d: %d -> %d (%s)", i, port.PublishedPort, port.TargetPort, port.Protocol)
		}
		return fmt.Errorf("SSH port not found for application %s", machineID)
	}

	// Get machine folder for SSH keys
	logger.Debug("=== GETTING SSH KEYS ===")
	machineFolder := os.Getenv("MACHINE_FOLDER")
	if machineFolder == "" {
		logger.Error("MACHINE_FOLDER environment variable is missing")
		return fmt.Errorf("MACHINE_FOLDER environment variable is missing")
	}
	logger.Debugf("Machine folder: %s", machineFolder)

	// Get private key for SSH authentication
	privateKey, err := ssh.GetPrivateKeyRawBase(machineFolder)
	if err != nil {
		logger.Errorf("Failed to load private key: %v", err)
		return fmt.Errorf("failed to load private key: %w", err)
	}
	logger.Debugf("✓ Private key loaded (length: %d bytes)", len(privateKey))

	// Create SSH client using DevPod's SSH utilities
	logger.Debug("=== CREATING SSH CONNECTION ===")
	// Extract hostname from ServerURL (remove https:// prefix if present)
	hostname := opts.DokployServerURL
	if strings.HasPrefix(hostname, "https://") {
		hostname = strings.TrimPrefix(hostname, "https://")
	}
	if strings.HasPrefix(hostname, "http://") {
		hostname = strings.TrimPrefix(hostname, "http://")
	}
	
	sshAddress := hostname + ":" + sshPort
	logger.Debugf("SSH address: %s", sshAddress)
	logger.Debugf("SSH user: devpod")
	
	sshClient, err := ssh.NewSSHClient("devpod", sshAddress, privateKey)
	if err != nil {
		logger.Errorf("Failed to create SSH client: %v", err)
		return fmt.Errorf("failed to create SSH client: %w", err)
	}
	defer sshClient.Close()
	logger.Debug("✓ SSH client created successfully")

	// Execute the command via SSH using the correct signature for DevPod v0.6.16-alpha.2
	logger.Debug("=== EXECUTING COMMAND VIA SSH ===")
	logger.Debugf("About to execute command: %s", command)
	logger.Debug("Using empty environment map to avoid SSH setenv errors")
	
	err = ssh.Run(context.Background(), sshClient, command, os.Stdin, os.Stdout, os.Stderr, map[string]string{})
	if err != nil {
		logger.Errorf("SSH command execution failed: %v", err)
		return fmt.Errorf("SSH command execution failed: %w", err)
	}
	
	logger.Debug("✓ SSH command executed successfully")
	logger.Debug("=== COMMAND EXECUTION DEBUG END ===")
	return nil
} 