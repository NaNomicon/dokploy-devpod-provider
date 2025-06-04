package cmd

import (
	"fmt"

	"github.com/NaNomicon/dokploy-devpod-provider/pkg/dokploy"
	"github.com/NaNomicon/dokploy-devpod-provider/pkg/options"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Dokploy workspace",
	Long: `Delete an existing development workspace from Dokploy.
This will remove the application and all associated resources.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDelete()
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func runDelete() error {
	// Setup logger
	logger := logrus.New()
	if verbose {
		logger.SetLevel(logrus.DebugLevel)
	}

	machineID, err := getMachineIDFromContext()
	if err != nil {
		return fmt.Errorf("failed to get machine ID: %w", err)
	}

	logger.Infof("Deleting Dokploy workspace: %s", machineID)

	// Load options from environment
	opts, err := options.LoadFromEnv()
	if err != nil {
		return fmt.Errorf("failed to load options: %w", err)
	}

	// Create Dokploy client
	client := dokploy.NewClient(opts, logger)

	// Delete the application
	err = client.DeleteApplicationByName(machineID)
	if err != nil {
		return fmt.Errorf("failed to delete application: %w", err)
	}

	logger.Info("✓ Dokploy workspace deleted")
	return nil
} 