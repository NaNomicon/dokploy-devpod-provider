package cmd

import (
	"fmt"

	"github.com/NaNomicon/dokploy-devpod-provider/pkg/dokploy"
	"github.com/NaNomicon/dokploy-devpod-provider/pkg/options"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a stopped Dokploy workspace",
	Long: `Start a previously stopped development workspace in Dokploy.
This will restart the Docker Compose service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runStart()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func runStart() error {
	// Setup logger
	logger := logrus.New()
	if verbose {
		logger.SetLevel(logrus.DebugLevel)
	}

	machineID, err := getMachineIDFromContext()
	if err != nil {
		return fmt.Errorf("failed to get machine ID: %w", err)
	}

	logger.Infof("Starting Dokploy workspace via Docker Compose: %s", machineID)

	// Load options from environment
	opts, err := options.LoadFromEnv()
	if err != nil {
		return fmt.Errorf("failed to load options: %w", err)
	}

	// Create Dokploy client
	client := dokploy.NewClient(opts, logger)

	// Start the Docker Compose service
	err = client.StartComposeByName(machineID)
	if err != nil {
		return fmt.Errorf("failed to start Docker Compose service: %w", err)
	}

	logger.Info("✓ Dokploy workspace started (Docker Compose service)")
	return nil
} 