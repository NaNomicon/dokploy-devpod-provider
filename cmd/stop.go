package cmd

import (
	"fmt"

	"github.com/NaNomicon/dokploy-devpod-provider/pkg/dokploy"
	"github.com/NaNomicon/dokploy-devpod-provider/pkg/options"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a running Dokploy workspace",
	Long: `Stop a currently running development workspace in Dokploy.
This will stop the application container but preserve the workspace.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runStop()
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func runStop() error {
	// Setup logger
	logger := logrus.New()
	if verbose {
		logger.SetLevel(logrus.DebugLevel)
	}

	machineID, err := getMachineIDFromContext()
	if err != nil {
		return fmt.Errorf("failed to get machine ID: %w", err)
	}

	logger.Infof("Stopping Dokploy workspace: %s", machineID)

	// Load options from environment
	opts, err := options.LoadFromEnv()
	if err != nil {
		return fmt.Errorf("failed to load options: %w", err)
	}

	// Create Dokploy client
	client := dokploy.NewClient(opts, logger)

	// Stop the application
	err = client.StopApplicationByName(machineID)
	if err != nil {
		return fmt.Errorf("failed to stop application: %w", err)
	}

	logger.Info("✓ Dokploy workspace stopped")
	return nil
} 