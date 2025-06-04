package cmd

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/NaNomicon/dokploy-devpod-provider/pkg/client"
	"github.com/NaNomicon/dokploy-devpod-provider/pkg/dokploy"
	"github.com/NaNomicon/dokploy-devpod-provider/pkg/options"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of a Dokploy workspace",
	Long: `Get the current status of a development workspace in Dokploy.
Returns one of: Running, Stopped, Busy, or NotFound.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runStatus()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus() error {
	// Setup logger - in dev environments, log to file to avoid stdout interference
	logger := logrus.New()
	
	// Check if we're in development environment
	isDev := os.Getenv("DEVPOD_PROVIDER_DEV") == "true" || os.Getenv("DOKPLOY_PROVIDER_DEV") == "true"
	
	if isDev {
		// In development, log to file
		logFile, err := os.OpenFile("/tmp/dokploy-provider-status.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logger.SetOutput(logFile)
			defer logFile.Close()
		} else {
			// Fallback to stderr if file creation fails
			logger.SetOutput(os.Stderr)
		}
	} else {
		// In production, use stderr to avoid stdout interference
		logger.SetOutput(os.Stderr)
	}
	
	if verbose || isDev {
		logger.SetLevel(logrus.DebugLevel)
	}

	logger.Debugf("=== Status Check Started ===")
	logger.Debugf("Development mode: %v", isDev)

	machineID, err := getMachineIDFromContext()
	if err != nil {
		logger.Debugf("Failed to get machine ID: %v", err)
		fmt.Println(client.StatusNotFound)
		return nil
	}

	logger.Debugf("Starting status check for machine: %s", machineID)

	// Load options from environment
	opts, err := options.LoadFromEnv()
	if err != nil {
		logger.Debugf("Failed to load options: %v", err)
		fmt.Println(client.StatusNotFound)
		return nil
	}

	// Create Dokploy client
	dokployClient := dokploy.NewClient(opts, logger)

	// Get application status from Dokploy
	dokployStatus, err := dokployClient.GetApplicationStatus(machineID)
	if err != nil {
		logger.Debugf("Failed to get application status: %v", err)
		fmt.Println(client.StatusNotFound)
		return nil
	}

	logger.Debugf("Dokploy status: %s", dokployStatus)

	// If Dokploy says the application is not running, return that status
	if dokployStatus != client.StatusRunning {
		logger.Debugf("Returning Dokploy status: %s", dokployStatus)
		fmt.Println(dokployStatus)
		return nil
	}

	// If Dokploy says it's running, check SSH readiness
	logger.Debugf("Application is running in Dokploy, checking SSH readiness...")
	
	// Get application details to find SSH connection info
	app, err := dokployClient.GetApplicationByName(machineID)
	if err != nil {
		logger.Debugf("Failed to get application details: %v", err)
		fmt.Println(client.StatusBusy) // Return Busy if we can't check SSH
		return nil
	}

	// Find SSH port mapping
	var sshPort int
	for _, port := range app.Ports {
		if port.TargetPort == 22 && port.Protocol == "tcp" {
			sshPort = port.PublishedPort
			logger.Debugf("Found SSH port mapping: %d -> 22", sshPort)
			break
		}
	}

	if sshPort == 0 {
		logger.Debugf("No SSH port mapping found, returning Busy")
		fmt.Println(client.StatusBusy) // Still setting up
		return nil
	}

	// Extract host from server URL
	parsedURL, err := url.Parse(opts.DokployServerURL)
	if err != nil {
		logger.Debugf("Failed to parse server URL: %v", err)
		fmt.Println(client.StatusBusy)
		return nil
	}
	sshHost := strings.Split(parsedURL.Host, ":")[0]
	logger.Debugf("SSH connection target: %s:%d", sshHost, sshPort)

	// Check SSH readiness
	isSSHReady := checkSSHReadiness(sshHost, sshPort, logger)
	
	if isSSHReady {
		logger.Debugf("SSH is ready on %s:%d - returning Running", sshHost, sshPort)
		fmt.Println(client.StatusRunning)
	} else {
		logger.Debugf("SSH is not ready yet on %s:%d - returning Busy", sshHost, sshPort)
		fmt.Println(client.StatusBusy)
	}

	logger.Debugf("=== Status Check Completed ===")
	return nil
}

func checkSSHReadiness(host string, port int, logger *logrus.Logger) bool {
	testAddress := fmt.Sprintf("%s:%d", host, port)
	
	// First check if the port is accessible
	conn, err := net.DialTimeout("tcp", testAddress, 3*time.Second)
	if err != nil {
		logger.Debugf("SSH port %d not accessible: %v", port, err)
		return false
	}
	conn.Close()
	
	// Test SSH service availability (without authentication, just to see if SSH daemon responds)
	config := &ssh.ClientConfig{
		User: "devpod",
		Auth: []ssh.AuthMethod{
			// Use a dummy key method that will fail but allow us to test if SSH daemon is responding
			ssh.PublicKeys(),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         3 * time.Second,
	}
	
	// Try to connect - we expect this to fail with auth error, but it should connect to SSH daemon
	sshClient, err := ssh.Dial("tcp", testAddress, config)
	if err != nil {
		// Check if it's an authentication error (which is expected and means SSH is ready)
		if strings.Contains(err.Error(), "unable to authenticate") || strings.Contains(err.Error(), "no supported methods remain") {
			logger.Debugf("SSH daemon is responding on port %d (authentication error expected)", port)
			return true
		} else {
			logger.Debugf("SSH service not ready on port %d: %v", port, err)
			return false
		}
	} else {
		// Unexpected success - close the connection
		sshClient.Close()
		logger.Debugf("SSH connection successful on port %d", port)
		return true
	}
} 