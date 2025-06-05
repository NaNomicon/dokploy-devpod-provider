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

	logger.Debugf("=== Status Check Started (Docker Compose Mode) ===")
	logger.Debugf("Development mode: %v", isDev)

	machineID, err := getMachineIDFromContext()
	if err != nil {
		logger.Debugf("Failed to get machine ID: %v", err)
		fmt.Println(client.StatusNotFound)
		return nil
	}

	logger.Debugf("Starting status check for compose service: %s", machineID)

	// Load options from environment
	opts, err := options.LoadFromEnv()
	if err != nil {
		logger.Debugf("Failed to load options: %v", err)
		fmt.Println(client.StatusNotFound)
		return nil
	}

	// Create Dokploy client
	dokployClient := dokploy.NewClient(opts, logger)

	// Get Docker Compose service status from Dokploy
	dokployStatus, err := dokployClient.GetComposeStatus(machineID)
	if err != nil {
		logger.Debugf("Failed to get compose service status: %v", err)
		fmt.Println(client.StatusNotFound)
		return nil
	}

	logger.Debugf("Dokploy compose status: %s", dokployStatus)

	// If Dokploy says the compose service is not running, return that status
	if dokployStatus != client.StatusRunning {
		logger.Debugf("Returning Dokploy compose status: %s", dokployStatus)
		fmt.Println(dokployStatus)
		return nil
	}

	// If Dokploy says it's running, check SSH readiness
	logger.Debugf("Compose service is running in Dokploy, checking SSH readiness...")
	
	// First get the compose service to find its ID
	composeBasic, err := dokployClient.GetComposeByName(machineID)
	if err != nil {
		logger.Debugf("Failed to get compose service details: %v", err)
		fmt.Println(client.StatusBusy) // Return Busy if we can't check SSH
		return nil
	}

	// Now get the full compose service details
	compose, err := dokployClient.GetCompose(composeBasic.ComposeID)
	if err != nil {
		logger.Debugf("Failed to get full compose service details: %v", err)
		fmt.Println(client.StatusBusy) // Return Busy if we can't check SSH
		return nil
	}

	// For Docker Compose, the port mapping is embedded in the compose file
	// We need to extract it from the compose configuration
	// Since we know we set the port in create.go, we can try to derive it
	sshPort, err := extractSSHPortFromCompose(compose, opts, logger)
	if err != nil {
		logger.Debugf("Failed to extract SSH port from compose service: %v", err)
		fmt.Println(client.StatusBusy) // Still setting up
		return nil
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

func extractSSHPortFromCompose(compose *dokploy.Compose, opts *options.Options, logger *logrus.Logger) (int, error) {
	// For Docker Compose services, we need to find the SSH port from the existing services
	// Since we know the port was allocated during creation, we can check all projects for used ports
	// and find the one that matches our naming pattern
	
	// Get all projects to check existing port usage
	dokployClient := dokploy.NewClient(opts, logger)
	allProjects, err := dokployClient.GetAllProjects()
	if err != nil {
		return 0, fmt.Errorf("failed to get projects for port discovery: %w", err)
	}

	// Look for compose services and try to derive the SSH port
	for _, project := range allProjects {
		for _, projectCompose := range project.Composes {
			if projectCompose.ComposeID == compose.ComposeID {
				// This is our compose service
				// The SSH port was allocated in the range 2222-2250
				// We need to check each port to see which one is in use
				
				parsedURL, err := url.Parse(opts.DokployServerURL)
				if err != nil {
					continue
				}
				sshHost := strings.Split(parsedURL.Host, ":")[0]
				
				// Check ports in the range we use
				for port := 2222; port <= 2250; port++ {
					testAddress := fmt.Sprintf("%s:%d", sshHost, port)
					conn, err := net.DialTimeout("tcp", testAddress, 1*time.Second)
					if err == nil {
						conn.Close()
						// This port is accessible, check if it's SSH
						if isSSHPort(sshHost, port, logger) {
							logger.Debugf("Found SSH port for compose service: %d", port)
							return port, nil
						}
					}
				}
			}
		}
	}

	return 0, fmt.Errorf("SSH port not found for compose service")
}

func isSSHPort(host string, port int, logger *logrus.Logger) bool {
	testAddress := fmt.Sprintf("%s:%d", host, port)
	
	config := &ssh.ClientConfig{
		User: "devpod",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(), // Will fail, but allows testing SSH service
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         2 * time.Second,
	}
	
	sshClient, err := ssh.Dial("tcp", testAddress, config)
	if err != nil {
		// Check if it's an authentication error (which means SSH is responding)
		if strings.Contains(err.Error(), "unable to authenticate") || strings.Contains(err.Error(), "no supported methods remain") {
			return true
		}
		return false
	} else {
		sshClient.Close()
		return true
	}
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