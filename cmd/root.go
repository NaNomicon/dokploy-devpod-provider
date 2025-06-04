package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dokploy-provider",
	Short: "DevPod provider for Dokploy",
	Long: `A DevPod provider that creates and manages development machines via Dokploy.
This provider allows you to create development environments on Dokploy infrastructure
with automatic SSH setup and Docker-in-Docker support.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dokploy-provider.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".dokploy-provider" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".dokploy-provider")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}

// getMachineIDFromContext gets the machine ID from the environment
// For machine providers, DevPod sets MACHINE_ID environment variable
func getMachineIDFromContext() (string, error) {
	// DevPod sets MACHINE_ID for machine providers
	machineID := os.Getenv("MACHINE_ID")
	if machineID != "" {
		return machineID, nil
	}
	
	// Fallback: try DEVPOD_MACHINE_ID (for some operations)
	machineID = os.Getenv("DEVPOD_MACHINE_ID")
	if machineID != "" {
		return machineID, nil
	}
	
	// Collect all relevant env vars for debugging
	var envs []string
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "MACHINE") || strings.HasPrefix(env, "DEVPOD") {
			envs = append(envs, env)
		}
	}
	
	return "", fmt.Errorf("could not determine machine ID. Available env vars: %v", envs)
} 