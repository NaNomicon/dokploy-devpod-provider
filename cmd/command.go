package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/NaNomicon/dokploy-devpod-provider/pkg/options"
	"github.com/NaNomicon/dokploy-devpod-provider/pkg/ssh"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// commandCmd represents the command execution command
var commandCmd = &cobra.Command{
	Use:   "command [command...]",
	Short: "Execute a command on a Dokploy workspace via SSH",
	Long: `Execute a command on a remote development workspace in Dokploy via SSH.
This command connects to the workspace and runs the specified command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCommand(args)
	},
}

func init() {
	rootCmd.AddCommand(commandCmd)
}

func runCommand(args []string) error {
	// Setup logger with stderr output to avoid interfering with command output
	logger := logrus.New()
	logger.SetOutput(os.Stderr)
	if verbose {
		logger.SetLevel(logrus.DebugLevel)
	}

	// Get machine ID from DevPod environment
	machineID := os.Getenv("DEVPOD_MACHINE_ID")
	if machineID == "" {
		return fmt.Errorf("DEVPOD_MACHINE_ID is required")
	}

	// Get command from arguments or environment
	var command string
	if len(args) > 0 {
		command = strings.Join(args, " ")
	} else {
		command = os.Getenv("DEVPOD_COMMAND")
		if command == "" {
			return fmt.Errorf("command is required (either as arguments or DEVPOD_COMMAND environment variable)")
		}
	}

	logger.Debugf("Executing command on workspace %s: %s", machineID, command)

	// Load options from environment
	opts, err := options.LoadFromEnv()
	if err != nil {
		return fmt.Errorf("failed to load options: %w", err)
	}

	// Create SSH client
	sshClient := ssh.NewClient(opts, logger)

	// Execute the command via SSH
	err = sshClient.ExecuteCommand(machineID, command)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	return nil
} 