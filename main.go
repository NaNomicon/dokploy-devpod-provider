package main

import (
	"os"

	"github.com/NaNomicon/dokploy-devpod-provider/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
} 