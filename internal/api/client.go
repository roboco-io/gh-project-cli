package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

const (
	// DefaultAPIURL is the GitHub GraphQL API endpoint
	DefaultAPIURL = "https://api.github.com/graphql"

	// DefaultRateLimit is the default rate limit for requests per second
	DefaultRateLimit = 10

	// DefaultTimeout is the default timeout for HTTP requests
	DefaultTimeout = 30 * time.Second
)

// Client is a GraphQL client for GitHub API
type Client struct {
	httpClient    *http.Client
	graphqlClient *graphql.Client
	rateLimiter   *RateLimiter
	retryConfig   *RetryConfig
	token         string
	baseURL       string
}

// RetryConfig holds configuration for retry logic
type RetryConfig struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

// RateLimiter handles rate limiting for API requests
type RateLimiter struct {
	lastRequest       time.Time
	requestsPerSecond int
}

// GraphQLRequest represents a GraphQL request
type GraphQLRequest struct {
	Variables map[string]interface{} `json:"variables,omitempty"`
	Query     string                 `json:"query"`
}

// GraphQLResponse represents a GraphQL response
type GraphQLResponse struct {
	Data   interface{}    `json:"data"`
	Errors []GraphQLError `json:"errors,omitempty"`
}

// GraphQLError represents a GraphQL error
type GraphQLError struct {
	Message   string                 `json:"message"`
	Locations []GraphQLErrorLocation `json:"locations,omitempty"`
	Path      []interface{}          `json:"path,omitempty"`
}

// GraphQLErrorLocation represents the location of a GraphQL error
type GraphQLErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// NewClient creates a new GraphQL client for GitHub API
func NewClient(token string) *Client {
	// Create GraphQL client with authentication
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	graphqlClient := graphql.NewClient(DefaultAPIURL, httpClient)

	return &Client{
		httpClient:    httpClient,
		graphqlClient: graphqlClient,
		token:         token,
		baseURL:       DefaultAPIURL,
		rateLimiter: &RateLimiter{
			requestsPerSecond: DefaultRateLimit,
		},
		retryConfig: &RetryConfig{
			MaxRetries: 3,
			BaseDelay:  time.Second,
			MaxDelay:   30 * time.Second,
		},
	}
}

// HealthCheck validates the connection to GitHub API
func (c *Client) HealthCheck(ctx context.Context) error {
	if c.token == "" {
		return fmt.Errorf("authentication token is required")
	}

	// Simple query to test connection
	var query struct {
		Viewer struct {
			Login string
		}
	}

	err := c.graphqlClient.Query(ctx, &query, nil)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	return nil
}

// Query executes a GraphQL query
func (c *Client) Query(ctx context.Context, query interface{}, variables map[string]interface{}) error {
	// Apply rate limiting
	c.rateLimiter.Wait()

	// Execute query with retry logic
	return c.retryOperation(func() error {
		return c.graphqlClient.Query(ctx, query, variables)
	})
}

// Mutate executes a GraphQL mutation
func (c *Client) Mutate(ctx context.Context, mutation interface{}, variables map[string]interface{}) error {
	// Apply rate limiting
	c.rateLimiter.Wait()

	// Execute mutation with retry logic
	return c.retryOperation(func() error {
		return c.graphqlClient.Mutate(ctx, mutation, variables)
	})
}

// Wait implements rate limiting
func (rl *RateLimiter) Wait() {
	now := time.Now()
	elapsed := now.Sub(rl.lastRequest)
	minInterval := time.Second / time.Duration(rl.requestsPerSecond)

	if elapsed < minInterval {
		time.Sleep(minInterval - elapsed)
	}

	rl.lastRequest = time.Now()
}

// retryOperation executes an operation with exponential backoff
func (c *Client) retryOperation(operation func() error) error {
	var lastErr error

	for attempt := 0; attempt <= c.retryConfig.MaxRetries; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't retry on last attempt
		if attempt == c.retryConfig.MaxRetries {
			break
		}

		// Check if error is retryable
		if !c.isRetryableError(err) {
			return err
		}

		// Calculate delay with exponential backoff
		delay := time.Duration(attempt+1) * c.retryConfig.BaseDelay
		if delay > c.retryConfig.MaxDelay {
			delay = c.retryConfig.MaxDelay
		}

		time.Sleep(delay)
	}

	return lastErr
}

// isRetryableError determines if an error should trigger a retry
func (c *Client) isRetryableError(_ error) bool {
	// TODO: Implement logic to identify retryable errors
	// e.g., rate limiting, temporary network errors, etc.
	return false
}

// buildRequest creates an HTTP request for GraphQL
func (c *Client) buildRequest(query string, variables map[string]interface{}) (*http.Request, error) {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", c.baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// parseGraphQLError parses GraphQL error response
func (c *Client) parseGraphQLError(responseBody []byte) error {
	var response GraphQLResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return fmt.Errorf("failed to parse error response: %w", err)
	}

	if len(response.Errors) == 0 {
		return fmt.Errorf("unknown GraphQL error")
	}

	return fmt.Errorf("GraphQL error: %s", response.Errors[0].Message)
}

// handleNetworkError handles network-related errors
func (c *Client) handleNetworkError(err error) error {
	return fmt.Errorf("network error: %w", err)
}
