package cmd

import (
	"fmt"
	"os"

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
This will restart the application container.`,
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

	machineID := os.Getenv("DEVPOD_MACHINE_ID")
	if machineID == "" {
		return fmt.Errorf("DEVPOD_MACHINE_ID is required")
	}

	logger.Infof("Starting Dokploy workspace: %s", machineID)

	// Load options from environment
	opts, err := options.LoadFromEnv()
	if err != nil {
		return fmt.Errorf("failed to load options: %w", err)
	}

	// Create Dokploy client
	client := dokploy.NewClient(opts, logger)

	// Start the application
	err = client.StartApplication(machineID)
	if err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}

	logger.Info("✓ Dokploy workspace started")
	return nil
} 