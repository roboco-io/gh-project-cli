package auth

import (
	"errors"
	"os"
	"os/exec"
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

	// TODO: Implement actual token retrieval from gh CLI
	// This will involve reading from:
	// - ~/.config/gh/hosts.yml (config file)
	// - System keychain/credential store
	return "", errors.New("not implemented: gh CLI token retrieval")
}

// ValidateToken validates the given token with GitHub API and returns scopes
func (g *GitHubCLIAuth) ValidateToken(token string) (bool, []string, error) {
	if token == "" {
		return false, nil, errors.New("empty token provided")
	}

	// TODO: Implement GitHub API validation
	// Make a request to https://api.github.com/user with the token
	// Parse the X-OAuth-Scopes header to get permissions
	return false, nil, errors.New("not implemented: token validation")
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