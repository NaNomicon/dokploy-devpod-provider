package cmd

import (
	"fmt"

	"github.com/NaNomicon/dokploy-devpod-provider/pkg/dokploy"
	"github.com/NaNomicon/dokploy-devpod-provider/pkg/options"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize and validate the Dokploy provider",
	Long: `Initialize the Dokploy provider by validating configuration options
and testing connectivity to the Dokploy server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit() error {
	// Setup logger
	logger := logrus.New()
	if verbose {
		logger.SetLevel(logrus.DebugLevel)
	}

	logger.Info("Initializing Dokploy provider...")

	// Load options from environment
	opts, err := options.LoadFromEnv()
	if err != nil {
		return fmt.Errorf("failed to load options: %w", err)
	}

	logger.Debug("Configuration loaded successfully")

	// Create Dokploy client
	client := dokploy.NewClient(opts, logger)

	// Test connection to Dokploy server
	logger.Info("Testing Dokploy server connection...")
	if err := client.HealthCheck(); err != nil {
		return fmt.Errorf("Dokploy server connection failed: %w", err)
	}

	logger.Info("✓ Dokploy server connection successful")

	// Test SSH connection if we have a machine ID (for existing workspaces)
	if opts.MachineID != "" {
		logger.Infof("Testing SSH connection to existing workspace: %s", opts.MachineID)
		
		// Get application details to test SSH connectivity
		app, err := client.GetApplication(opts.MachineID)
		if err != nil {
			logger.Warnf("Could not retrieve application details: %v", err)
		} else {
			// Find SSH port mapping
			var sshPort int
			for _, port := range app.Ports {
				if port.TargetPort == 22 {
					sshPort = port.PublishedPort
					break
				}
			}

			if sshPort > 0 {
				logger.Infof("Found SSH port mapping: %d -> 22", sshPort)
				logger.Info("✓ SSH port mapping configured")
			} else {
				logger.Warn("⚠️  No SSH port mapping found")
			}
		}
	}

	logger.Info("Dokploy provider initialized successfully")
	return nil
} 