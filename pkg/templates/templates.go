package templates

import _ "embed"

// DockerComposeTemplate contains the docker-compose.yml template
//
//go:embed docker-compose.yml
var DockerComposeTemplate string

// SetupScriptTemplate contains the setup-root.sh template
//
//go:embed setup-root.sh
var SetupScriptTemplate string
