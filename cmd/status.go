package cmd

import (
	"fmt"
	"os"

	"github.com/NaNomicon/dokploy-devpod-provider/pkg/client"
	"github.com/NaNomicon/dokploy-devpod-provider/pkg/dokploy"
	"github.com/NaNomicon/dokploy-devpod-provider/pkg/options"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
	// Setup logger with stderr output to avoid interfering with status output
	logger := logrus.New()
	logger.SetOutput(os.Stderr)
	if verbose {
		logger.SetLevel(logrus.DebugLevel)
	}

	machineID := os.Getenv("DEVPOD_MACHINE_ID")
	if machineID == "" {
		logger.Debug("No machine ID provided, returning NotFound")
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

	// Get application status
	status, err := dokployClient.GetApplicationStatus(machineID)
	if err != nil {
		logger.Debugf("Failed to get application status: %v", err)
		fmt.Println(client.StatusNotFound)
		return nil
	}

	logger.Debugf("Status check completed: %s", status)
	fmt.Println(status)
	return nil
} 