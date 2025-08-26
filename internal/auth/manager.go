package auth

import (
	"fmt"
)

// Manager handles authentication flow and provides unified access to tokens
type Manager struct {
	ghAuth *GitHubCLIAuth
}

// NewAuthManager creates a new authentication manager
func NewAuthManager() *Manager {
	return &Manager{
		ghAuth: NewGitHubCLIAuth(),
	}
}

// GetValidatedToken retrieves and validates a GitHub token from various sources
func (am *Manager) GetValidatedToken() (string, error) {
	var token string
	var err error

	// Try to get token from GitHub CLI first
	if am.ghAuth.CheckGHCLIInstalled() {
		token, err = am.ghAuth.GetToken("github.com")
		if err == nil && token != "" {
			// Validate the token
			valid, scopes, validErr := am.ghAuth.ValidateToken(token)
			if validErr != nil {
				return "", fmt.Errorf("token validation failed: %w", validErr)
			}
			if !valid {
				return "", fmt.Errorf("token is invalid")
			}

			// Check for required scopes
			requiredScopes := []string{"repo", "project"}
			if !HasRequiredScopes(scopes, requiredScopes) {
				return "", fmt.Errorf("token missing required scopes. Required: %v, Available: %v", requiredScopes, scopes)
			}

			return token, nil
		}
	}

	// If gh CLI fails, try environment variables
	if fallbackToken := am.ghAuth.GetFallbackToken(); fallbackToken != "" {
		// Validate the fallback token
		valid, scopes, validErr := am.ghAuth.ValidateToken(fallbackToken)
		if validErr != nil {
			return "", fmt.Errorf("fallback token validation failed: %w", validErr)
		}
		if !valid {
			return "", fmt.Errorf("fallback token is invalid")
		}

		// Check for required scopes
		requiredScopes := []string{"repo", "project"}
		if !HasRequiredScopes(scopes, requiredScopes) {
			return "", fmt.Errorf("fallback token missing required scopes. Required: %v, Available: %v", requiredScopes, scopes)
		}

		return fallbackToken, nil
	}

	return "", fmt.Errorf("no valid GitHub token found. Please authenticate with 'gh auth login' or set GITHUB_TOKEN environment variable")
}

// GetTokenWithoutValidation gets a token without validation (for testing)
func (am *Manager) GetTokenWithoutValidation() (string, error) {
	// Try GitHub CLI first
	if am.ghAuth.CheckGHCLIInstalled() {
		if token, err := am.ghAuth.GetToken("github.com"); err == nil && token != "" {
			return token, nil
		}
	}

	// Try environment variables
	if fallbackToken := am.ghAuth.GetFallbackToken(); fallbackToken != "" {
		return fallbackToken, nil
	}

	return "", fmt.Errorf("no GitHub token found")
}

// CheckAuthentication checks if user is properly authenticated
func (am *Manager) CheckAuthentication() error {
	_, err := am.GetValidatedToken()
	return err
}

// GetAuthenticationStatus returns detailed authentication status
func (am *Manager) GetAuthenticationStatus() Status {
	status := Status{
		GHCLIInstalled: am.ghAuth.CheckGHCLIInstalled(),
		HasEnvToken:    am.ghAuth.GetFallbackToken() != "",
	}

	// Try to get and validate token
	token, err := am.GetTokenWithoutValidation()
	if err != nil {
		status.TokenAvailable = false
		status.Error = err.Error()
		return status
	}

	status.TokenAvailable = true

	// Validate token
	valid, scopes, err := am.ghAuth.ValidateToken(token)
	if err != nil {
		status.TokenValid = false
		status.Error = err.Error()
		return status
	}

	status.TokenValid = valid
	status.Scopes = scopes

	// Check required scopes
	requiredScopes := []string{"repo", "project"}
	status.HasRequiredScopes = HasRequiredScopes(scopes, requiredScopes)
	status.RequiredScopes = requiredScopes

	if !status.HasRequiredScopes {
		status.Error = fmt.Sprintf("Missing required scopes: %v", getMissingScopes(scopes, requiredScopes))
	}

	return status
}

// Status represents the current authentication status
type Status struct {
	Error             string   `json:"error,omitempty"`
	Scopes            []string `json:"scopes"`
	RequiredScopes    []string `json:"required_scopes"`
	GHCLIInstalled    bool     `json:"gh_cli_installed"`
	HasEnvToken       bool     `json:"has_env_token"`
	TokenAvailable    bool     `json:"token_available"`
	TokenValid        bool     `json:"token_valid"`
	HasRequiredScopes bool     `json:"has_required_scopes"`
}

// IsReady returns true if authentication is fully configured
func (as *Status) IsReady() bool {
	return as.TokenAvailable && as.TokenValid && as.HasRequiredScopes
}

// GetRecommendation returns a recommendation for fixing authentication issues
func (as *Status) GetRecommendation() string {
	if !as.GHCLIInstalled {
		return "Install GitHub CLI: https://cli.github.com/manual/installation"
	}

	if !as.TokenAvailable {
		return "Authenticate with GitHub CLI: gh auth login"
	}

	if !as.TokenValid {
		return "Re-authenticate with GitHub CLI: gh auth login --force"
	}

	if !as.HasRequiredScopes {
		return "Grant additional scopes: gh auth refresh -s repo -s project"
	}

	return "Authentication is properly configured"
}

// getMissingScopes returns scopes that are required but not available
func getMissingScopes(available, required []string) []string {
	scopeMap := make(map[string]bool)
	for _, scope := range available {
		scopeMap[scope] = true
	}

	var missing []string
	for _, required := range required {
		if !scopeMap[required] {
			missing = append(missing, required)
		}
	}

	return missing
}
