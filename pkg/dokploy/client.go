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
}

// Application represents a Dokploy application
type Application struct {
	ApplicationID string `json:"applicationId"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	ProjectID     string `json:"projectId"`
	Status        string `json:"applicationStatus"`
	Ports         []Port `json:"ports"`
}

// Port represents a port mapping
type Port struct {
	PublishedPort int    `json:"publishedPort"`
	TargetPort    int    `json:"targetPort"`
	Protocol      string `json:"protocol"`
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

	var app Application
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
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
func (c *Client) GetApplicationStatus(applicationID string) (client.Status, error) {
	app, err := c.GetApplication(applicationID)
	if err != nil {
		return client.StatusNotFound, err
	}

	switch app.Status {
	case "done", "running":
		return client.StatusRunning, nil
	case "idle":
		return client.StatusStopped, nil
	case "error":
		return client.StatusNotFound, nil
	default:
		return client.StatusBusy, nil
	}
}

// makeRequest makes an HTTP request to the Dokploy API
func (c *Client) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
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

	c.logger.Debugf("Making %s request to %s", method, url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	return resp, nil
} 