package options

import (
	"fmt"
	"os"
)

// Options represents the configuration options for the Dokploy provider
type Options struct {
	// Required options
	DokployServerURL string `json:"dokployServerURL"`
	DokployAPIToken  string `json:"dokployAPIToken"`

	// Optional options
	DokployProjectName string `json:"dokployProjectName"`
	DokployServerID    string `json:"dokployServerID"`
	MachineType        string `json:"machineType"`

	// Machine identification
	MachineID string `json:"machineID"`
}

// LoadFromEnv loads options from environment variables
func LoadFromEnv() (*Options, error) {
	opts := &Options{
		DokployServerURL:   os.Getenv("DOKPLOY_SERVER_URL"),
		DokployAPIToken:    os.Getenv("DOKPLOY_API_TOKEN"),
		DokployProjectName: getEnvWithDefault("DOKPLOY_PROJECT_NAME", "devpod-workspaces"),
		DokployServerID:    os.Getenv("DOKPLOY_SERVER_ID"),
		MachineType:        getEnvWithDefault("MACHINE_TYPE", "small"),
		MachineID:          os.Getenv("MACHINE_ID"),
	}

	// Validate required options
	if opts.DokployServerURL == "" {
		return nil, fmt.Errorf("DOKPLOY_SERVER_URL is required")
	}
	if opts.DokployAPIToken == "" {
		return nil, fmt.Errorf("DOKPLOY_API_TOKEN is required")
	}

	return opts, nil
}

// getEnvWithDefault returns the environment variable value or a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 