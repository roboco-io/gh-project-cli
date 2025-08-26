package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	// HTTP status codes
	httpStatusUnauthorized = 401
	httpStatusOK           = 200
)

// GitHubCLIAuth handles authentication by integrating with GitHub CLI
type GitHubCLIAuth struct {
	// Future: can add configuration options here
}

// NewGitHubCLIAuth creates a new GitHub CLI authentication handler
func NewGitHubCLIAuth() *GitHubCLIAuth {
	return &GitHubCLIAuth{}
}

// GetToken retrieves the authentication token from GitHub CLI for the given hostname
func (g *GitHubCLIAuth) GetToken(hostname string) (string, error) {
	if !g.CheckGHCLIInstalled() {
		return "", errors.New("gh CLI is not installed or not available in PATH")
	}

	// Use gh CLI to get the token
	cmd := exec.Command("gh", "auth", "token", "--hostname", hostname)
	output, err := cmd.Output()
	if err != nil {
		// If gh CLI fails, try fallback token
		if fallbackToken := g.GetFallbackToken(); fallbackToken != "" {
			return fallbackToken, nil
		}
		return "", fmt.Errorf("failed to get token from gh CLI: %w", err)
	}

	token := strings.TrimSpace(string(output))
	if token == "" {
		// If gh CLI returns empty token, try fallback
		if fallbackToken := g.GetFallbackToken(); fallbackToken != "" {
			return fallbackToken, nil
		}
		return "", errors.New("gh CLI returned empty token")
	}

	return token, nil
}

// UserResponse represents the GitHub user API response
type UserResponse struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	ID    int    `json:"id"`
}

// ValidateToken validates the given token with GitHub API and returns scopes
func (g *GitHubCLIAuth) ValidateToken(token string) (isValid bool, scopes []string, err error) {
	if token == "" {
		return false, nil, errors.New("empty token provided")
	}

	// Create HTTP client with timeout
	const requestTimeout = 10 * time.Second
	client := &http.Client{
		Timeout: requestTimeout,
	}

	// Make request to GitHub API user endpoint
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", http.NoBody)
	if err != nil {
		return false, nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "ghp-cli")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return false, nil, fmt.Errorf("failed to validate token: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode == httpStatusUnauthorized {
		return false, nil, errors.New("invalid or expired token")
	}
	if resp.StatusCode != httpStatusOK {
		return false, nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response to ensure token works
	var user UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return false, nil, fmt.Errorf("failed to parse user response: %w", err)
	}

	// Parse scopes from X-OAuth-Scopes header
	if scopeHeader := resp.Header.Get("X-OAuth-Scopes"); scopeHeader != "" {
		scopesList := strings.Split(scopeHeader, ", ")
		for _, scope := range scopesList {
			scope = strings.TrimSpace(scope)
			if scope != "" {
				scopes = append(scopes, scope)
			}
		}
	}

	return true, scopes, nil
}

// GetFallbackToken attempts to get token from environment variables
func (g *GitHubCLIAuth) GetFallbackToken() string {
	// Try GitHub CLI standard environment variables first
	if token := os.Getenv("GH_TOKEN"); token != "" {
		return token
	}

	// Fall back to common GitHub token environment variable
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token
	}

	return ""
}

// CheckGHCLIInstalled checks if gh CLI is installed and available
func (g *GitHubCLIAuth) CheckGHCLIInstalled() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

// HasRequiredScopes checks if user has all required scopes for GitHub Projects
func HasRequiredScopes(userScopes, requiredScopes []string) bool {
	scopeMap := make(map[string]bool)
	for _, scope := range userScopes {
		scopeMap[scope] = true
	}

	for _, required := range requiredScopes {
		if !scopeMap[required] {
			return false
		}
	}

	return true
}
