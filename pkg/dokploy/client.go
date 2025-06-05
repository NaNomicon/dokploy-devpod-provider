package dokploy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/NaNomicon/dokploy-devpod-provider/pkg/client"
	"github.com/NaNomicon/dokploy-devpod-provider/pkg/options"
	"github.com/sirupsen/logrus"
)

// Client represents a Dokploy API client
type Client struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
	logger     *logrus.Logger
}

// NewClient creates a new Dokploy API client
func NewClient(opts *options.Options, logger *logrus.Logger) *Client {
	return &Client{
		baseURL:  strings.TrimSuffix(opts.DokployServerURL, "/"),
		apiToken: opts.DokployAPIToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// Project represents a Dokploy project
type Project struct {
	ProjectID   string        `json:"projectId"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Applications []Application `json:"applications"`
	Composes    []Compose     `json:"compose"`
}

// Application represents a Dokploy application
type Application struct {
	ApplicationID string   `json:"applicationId"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	ProjectID     string   `json:"projectId"`
	Status        string   `json:"applicationStatus"`
	Domains       []Domain `json:"domains"`
	Ports         []Port   `json:"ports"`
}

// Domain represents a domain/port mapping in Dokploy
type Domain struct {
	DomainID    string `json:"domainId"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Path        string `json:"path"`
	HTTPS       bool   `json:"https"`
	DomainType  string `json:"domainType"`
	ServiceName string `json:"serviceName"`
}

// Port represents a port mapping (for backward compatibility)
type Port struct {
	PortID        string `json:"portId"`
	PublishedPort int    `json:"publishedPort"`
	TargetPort    int    `json:"targetPort"`
	Protocol      string `json:"protocol"`
	ApplicationID string `json:"applicationId"`
}

// CreateProjectRequest represents a project creation request
type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateApplicationRequest represents an application creation request
type CreateApplicationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ProjectID   string `json:"projectId"`
}

// DockerProviderRequest represents a Docker provider configuration request
type DockerProviderRequest struct {
	ApplicationID string `json:"applicationId"`
	DockerImage   string `json:"dockerImage"`
}

// EnvironmentRequest represents an environment configuration request
type EnvironmentRequest struct {
	ApplicationID string `json:"applicationId"`
	Env           string `json:"env"`
}

// UpdateApplicationRequest represents an application update request
type UpdateApplicationRequest struct {
	ApplicationID string `json:"applicationId"`
	Command       string `json:"command"`
}

// DeployRequest represents a deployment request
type DeployRequest struct {
	ApplicationID string `json:"applicationId"`
}

// CreatePortRequest represents a port creation request
type CreatePortRequest struct {
	PublishedPort int    `json:"publishedPort"`
	TargetPort    int    `json:"targetPort"`
	Protocol      string `json:"protocol"`
	ApplicationID string `json:"applicationId"`
}

// Compose represents a Dokploy Docker Compose service
type Compose struct {
	ComposeID   string `json:"composeId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ProjectID   string `json:"projectId"`
	Status      string `json:"composeStatus"`
	ComposeType string `json:"composeType"`
}

// CreateComposeRequest represents a Docker Compose creation request
type CreateComposeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ProjectID   string `json:"projectId"`
	ComposeType string `json:"composeType"` // "docker-compose" or "stack"
}

// SaveComposeFileRequest represents a request to save docker-compose.yml content
type SaveComposeFileRequest struct {
	ComposeID     string `json:"composeId"`
	DockerCompose string `json:"dockerCompose"`
}

// UpdateComposeRequest represents a request to update Docker Compose configuration
type UpdateComposeRequest struct {
	ComposeID    string `json:"composeId"`
	ComposeFile  string `json:"composeFile"`
	SourceType   string `json:"sourceType"`
	ComposePath  string `json:"composePath"`
}

// DeployComposeRequest represents a Docker Compose deployment request
type DeployComposeRequest struct {
	ComposeID string `json:"composeId"`
}

// HealthCheck checks if the Dokploy server is accessible
func (c *Client) HealthCheck() error {
	resp, err := c.makeRequest("GET", "/api/settings.health", nil)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	c.logger.Debug("Health check successful")
	return nil
}

// GetAllProjects retrieves all projects
func (c *Client) GetAllProjects() ([]Project, error) {
	resp, err := c.makeRequest("GET", "/api/project.all", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}
	defer resp.Body.Close()

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, fmt.Errorf("failed to decode projects response: %w", err)
	}

	return projects, nil
}

// CreateProject creates a new project
func (c *Client) CreateProject(req CreateProjectRequest) (*Project, error) {
	resp, err := c.makeRequest("POST", "/api/project.create", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	defer resp.Body.Close()

	var project Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, fmt.Errorf("failed to decode project response: %w", err)
	}

	return &project, nil
}

// CreateApplication creates a new application
func (c *Client) CreateApplication(req CreateApplicationRequest) (*Application, error) {
	resp, err := c.makeRequest("POST", "/api/application.create", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}
	defer resp.Body.Close()

	var app Application
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, fmt.Errorf("failed to decode application response: %w", err)
	}

	return &app, nil
}

// SaveDockerProvider configures Docker provider for an application
func (c *Client) SaveDockerProvider(req DockerProviderRequest) error {
	resp, err := c.makeRequest("POST", "/api/application.saveDockerProvider", req)
	if err != nil {
		return fmt.Errorf("failed to save Docker provider: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to save Docker provider, status: %d", resp.StatusCode)
	}

	return nil
}

// SaveEnvironment configures environment variables for an application
func (c *Client) SaveEnvironment(req EnvironmentRequest) error {
	resp, err := c.makeRequest("POST", "/api/application.saveEnvironment", req)
	if err != nil {
		return fmt.Errorf("failed to save environment: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to save environment, status: %d", resp.StatusCode)
	}

	return nil
}

// UpdateApplication updates an application
func (c *Client) UpdateApplication(req UpdateApplicationRequest) error {
	resp, err := c.makeRequest("POST", "/api/application.update", req)
	if err != nil {
		return fmt.Errorf("failed to update application: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update application, status: %d", resp.StatusCode)
	}

	return nil
}

// DeployApplication deploys an application
func (c *Client) DeployApplication(req DeployRequest) error {
	resp, err := c.makeRequest("POST", "/api/application.deploy", req)
	if err != nil {
		return fmt.Errorf("failed to deploy application: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to deploy application, status: %d", resp.StatusCode)
	}

	return nil
}

// CreatePort creates a port mapping
func (c *Client) CreatePort(req CreatePortRequest) error {
	resp, err := c.makeRequest("POST", "/api/port.create", req)
	if err != nil {
		return fmt.Errorf("failed to create port: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create port, status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetApplication retrieves an application by ID
func (c *Client) GetApplication(applicationID string) (*Application, error) {
	endpoint := fmt.Sprintf("/api/application.one?applicationId=%s", url.QueryEscape(applicationID))
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body for debugging
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	c.logger.Debugf("Dokploy API response for application %s: %s", applicationID, string(body))

	// Check for error responses
	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Message string `json:"message"`
			Code    string `json:"code"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil {
			if errorResp.Code == "NOT_FOUND" {
				return nil, fmt.Errorf("application not found: %s", errorResp.Message)
			}
		}
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var app Application
	if err := json.Unmarshal(body, &app); err != nil {
		return nil, fmt.Errorf("failed to decode application response: %w", err)
	}

	return &app, nil
}

// DeleteApplication deletes an application
func (c *Client) DeleteApplication(applicationID string) error {
	req := map[string]string{"applicationId": applicationID}
	resp, err := c.makeRequest("DELETE", "/api/application.remove", req)
	if err != nil {
		return fmt.Errorf("failed to delete application: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete application, status: %d", resp.StatusCode)
	}

	return nil
}

// StartApplication starts an application
func (c *Client) StartApplication(applicationID string) error {
	req := map[string]string{"applicationId": applicationID}
	resp, err := c.makeRequest("POST", "/api/application.start", req)
	if err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to start application, status: %d", resp.StatusCode)
	}

	return nil
}

// StopApplication stops an application
func (c *Client) StopApplication(applicationID string) error {
	req := map[string]string{"applicationId": applicationID}
	resp, err := c.makeRequest("POST", "/api/application.stop", req)
	if err != nil {
		return fmt.Errorf("failed to stop application: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to stop application, status: %d", resp.StatusCode)
	}

	return nil
}

// GetApplicationStatus returns the DevPod-compatible status of an application
func (c *Client) GetApplicationStatus(applicationName string) (client.Status, error) {
	app, err := c.GetApplicationByName(applicationName)
	if err != nil {
		return client.StatusNotFound, err
	}

	c.logger.Debugf("Dokploy application status for %s: '%s'", applicationName, app.Status)

	switch app.Status {
	case "done", "running":
		return client.StatusRunning, nil
	case "idle", "stopped":
		return client.StatusStopped, nil
	case "error", "failed":
		return client.StatusNotFound, nil
	case "building", "deploying", "restarting":
		return client.StatusBusy, nil
	default:
		c.logger.Warnf("Unknown Dokploy status '%s' for application %s, treating as busy", app.Status, applicationName)
		return client.StatusBusy, nil
	}
}

// GetApplicationByName retrieves an application by name
func (c *Client) GetApplicationByName(applicationName string) (*Application, error) {
	// Get all projects to find the application
	projects, err := c.GetAllProjects()
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	// Find the application with matching name
	for _, project := range projects {
		for _, application := range project.Applications {
			if application.Name == applicationName {
				return &application, nil
			}
		}
	}

	return nil, fmt.Errorf("application with name '%s' not found", applicationName)
}

// DeleteApplicationByName deletes an application by name
func (c *Client) DeleteApplicationByName(applicationName string) error {
	app, err := c.GetApplicationByName(applicationName)
	if err != nil {
		return fmt.Errorf("failed to find application: %w", err)
	}
	
	return c.DeleteApplication(app.ApplicationID)
}

// StartApplicationByName starts an application by name
func (c *Client) StartApplicationByName(applicationName string) error {
	app, err := c.GetApplicationByName(applicationName)
	if err != nil {
		return fmt.Errorf("failed to find application: %w", err)
	}
	
	return c.StartApplication(app.ApplicationID)
}

// StopApplicationByName stops an application by name
func (c *Client) StopApplicationByName(applicationName string) error {
	app, err := c.GetApplicationByName(applicationName)
	if err != nil {
		return fmt.Errorf("failed to find application: %w", err)
	}
	
	return c.StopApplication(app.ApplicationID)
}

// CreateCompose creates a new Docker Compose service
func (c *Client) CreateCompose(req CreateComposeRequest) (*Compose, error) {
	resp, err := c.makeRequest("POST", "/api/compose.create", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create compose service: %w", err)
	}
	defer resp.Body.Close()

	var compose Compose
	if err := json.NewDecoder(resp.Body).Decode(&compose); err != nil {
		return nil, fmt.Errorf("failed to decode compose response: %w", err)
	}

	return &compose, nil
}

// SaveComposeFile saves the docker-compose.yml content using the update endpoint
func (c *Client) SaveComposeFile(req SaveComposeFileRequest) error {
	// Convert to UpdateComposeRequest format
	updateReq := UpdateComposeRequest{
		ComposeID:    req.ComposeID,
		ComposeFile:  req.DockerCompose,
		SourceType:   "raw", // Use raw source type for direct compose content
		ComposePath:  "./docker-compose.yml", // Default compose path
	}

	resp, err := c.makeRequest("POST", "/api/compose.update", updateReq)
	if err != nil {
		return fmt.Errorf("failed to save compose file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to save compose file, status: %d", resp.StatusCode)
	}

	return nil
}

// DeployCompose deploys a Docker Compose service
func (c *Client) DeployCompose(req DeployComposeRequest) error {
	resp, err := c.makeRequest("POST", "/api/compose.deploy", req)
	if err != nil {
		return fmt.Errorf("failed to deploy compose service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to deploy compose service, status: %d", resp.StatusCode)
	}

	return nil
}

// GetCompose retrieves a Docker Compose service by ID
func (c *Client) GetCompose(composeID string) (*Compose, error) {
	endpoint := fmt.Sprintf("/api/compose.one?composeId=%s", url.QueryEscape(composeID))
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get compose service: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body for debugging
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	c.logger.Debugf("Dokploy API response for compose %s: %s", composeID, string(body))

	// Check for error responses
	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Message string `json:"message"`
			Code    string `json:"code"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil {
			if errorResp.Code == "NOT_FOUND" {
				return nil, fmt.Errorf("compose service not found: %s", errorResp.Message)
			}
		}
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var compose Compose
	if err := json.Unmarshal(body, &compose); err != nil {
		return nil, fmt.Errorf("failed to decode compose response: %w", err)
	}

	return &compose, nil
}

// DeleteCompose deletes a Docker Compose service
func (c *Client) DeleteCompose(composeID string) error {
	req := map[string]string{"composeId": composeID}
	resp, err := c.makeRequest("DELETE", "/api/compose.remove", req)
	if err != nil {
		return fmt.Errorf("failed to delete compose service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete compose service, status: %d", resp.StatusCode)
	}

	return nil
}

// StartCompose starts a Docker Compose service
func (c *Client) StartCompose(composeID string) error {
	req := map[string]string{"composeId": composeID}
	resp, err := c.makeRequest("POST", "/api/compose.start", req)
	if err != nil {
		return fmt.Errorf("failed to start compose service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to start compose service, status: %d", resp.StatusCode)
	}

	return nil
}

// StopCompose stops a Docker Compose service
func (c *Client) StopCompose(composeID string) error {
	req := map[string]string{"composeId": composeID}
	resp, err := c.makeRequest("POST", "/api/compose.stop", req)
	if err != nil {
		return fmt.Errorf("failed to stop compose service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to stop compose service, status: %d", resp.StatusCode)
	}

	return nil
}

// GetComposeStatus returns the DevPod-compatible status of a Docker Compose service
func (c *Client) GetComposeStatus(composeName string) (client.Status, error) {
	compose, err := c.GetComposeByName(composeName)
	if err != nil {
		return client.StatusNotFound, err
	}

	c.logger.Debugf("Dokploy compose status for %s: '%s'", composeName, compose.Status)

	switch compose.Status {
	case "done", "running":
		return client.StatusRunning, nil
	case "idle", "stopped":
		return client.StatusStopped, nil
	case "error", "failed":
		return client.StatusNotFound, nil
	case "building", "deploying", "restarting":
		return client.StatusBusy, nil
	default:
		c.logger.Warnf("Unknown Dokploy status '%s' for compose %s, treating as busy", compose.Status, composeName)
		return client.StatusBusy, nil
	}
}

// GetComposeByName retrieves a Docker Compose service by name
func (c *Client) GetComposeByName(composeName string) (*Compose, error) {
	// Get all projects to find the compose service
	projects, err := c.GetAllProjects()
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	// Find the compose service with matching name
	for _, project := range projects {
		for _, compose := range project.Composes {
			if compose.Name == composeName {
				return &compose, nil
			}
		}
	}

	return nil, fmt.Errorf("compose service with name '%s' not found", composeName)
}

// DeleteComposeByName deletes a Docker Compose service by name
func (c *Client) DeleteComposeByName(composeName string) error {
	compose, err := c.GetComposeByName(composeName)
	if err != nil {
		return fmt.Errorf("failed to find compose service: %w", err)
	}
	
	return c.DeleteCompose(compose.ComposeID)
}

// StartComposeByName starts a Docker Compose service by name
func (c *Client) StartComposeByName(composeName string) error {
	compose, err := c.GetComposeByName(composeName)
	if err != nil {
		return fmt.Errorf("failed to find compose service: %w", err)
	}
	
	return c.StartCompose(compose.ComposeID)
}

// StopComposeByName stops a Docker Compose service by name
func (c *Client) StopComposeByName(composeName string) error {
	compose, err := c.GetComposeByName(composeName)
	if err != nil {
		return fmt.Errorf("failed to find compose service: %w", err)
	}
	
	return c.StopCompose(compose.ComposeID)
}

// makeRequest makes an HTTP request to the Dokploy API with comprehensive debug logging
func (c *Client) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	var requestBodyStr string
	
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
		requestBodyStr = string(jsonBody)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("x-api-key", c.apiToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Debug log the request
	c.logger.Debugf("=== API REQUEST ===")
	c.logger.Debugf("Making %s request to %s", method, url)
	if requestBodyStr != "" {
		c.logger.Debugf("Request body: %s", requestBodyStr)
	} else {
		c.logger.Debugf("Request body: (empty)")
	}
	c.logger.Debugf("Headers: x-api-key=[REDACTED], Content-Type=%s", req.Header.Get("Content-Type"))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Debugf("=== API REQUEST FAILED ===")
		c.logger.Debugf("Error: %v", err)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	// Read response body for logging (we'll need to recreate it for the caller)
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Debugf("=== API RESPONSE (READ ERROR) ===")
		c.logger.Debugf("Status: %d %s", resp.StatusCode, resp.Status)
		c.logger.Debugf("Error reading response body: %v", err)
		resp.Body.Close()
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	resp.Body.Close()

	// Debug log the response
	c.logger.Debugf("=== API RESPONSE ===")
	c.logger.Debugf("Status: %d %s", resp.StatusCode, resp.Status)
	c.logger.Debugf("Response body: %s", string(responseBody))

	// Recreate the response body for the caller
	resp.Body = io.NopCloser(bytes.NewReader(responseBody))

	return resp, nil
} 