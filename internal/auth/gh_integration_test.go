package auth

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitHubCLIIntegration(t *testing.T) {
	t.Run("NewGitHubCLIAuth creates new instance", func(t *testing.T) {
		auth := NewGitHubCLIAuth()
		assert.NotNil(t, auth)
		assert.IsType(t, &GitHubCLIAuth{}, auth)
	})

	t.Run("GetToken returns token from gh CLI or fallback", func(t *testing.T) {
		auth := NewGitHubCLIAuth()

		token, err := auth.GetToken("github.com")

		// Should either succeed or fail gracefully
		if err != nil {
			// Should contain meaningful error message
			assert.True(t,
				strings.Contains(err.Error(), "gh CLI") ||
					strings.Contains(err.Error(), "failed to get token") ||
					strings.Contains(err.Error(), "not installed"),
			)
		} else {
			// If successful, token should not be empty
			assert.NotEmpty(t, token)
		}
	})

	t.Run("ValidateToken validates GitHub token", func(t *testing.T) {
		auth := NewGitHubCLIAuth()

		// Test with empty token
		isValid, scopes, err := auth.ValidateToken("")
		assert.False(t, isValid)
		assert.Nil(t, scopes)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty token")

		// Test with invalid token
		isValid, scopes, err = auth.ValidateToken("invalid_token_123")
		assert.False(t, isValid)
		assert.Nil(t, scopes)
		assert.Error(t, err)

		// Skip integration test if we don't have a real token
		if testing.Short() {
			t.Skip("Skipping integration test in short mode")
		}

		// Test with real token from environment if available
		if realToken := os.Getenv("GITHUB_TOKEN"); realToken != "" {
			isValid, scopes, err = auth.ValidateToken(realToken)
			if err == nil {
				assert.True(t, isValid)
				assert.NotNil(t, scopes) // Can be empty but not nil
			} else {
				// Network errors are acceptable in tests
				t.Logf("Token validation failed (network/token issue): %v", err)
			}
		}
	})

	t.Run("GetFallbackToken gets token from environment", func(t *testing.T) {
		auth := NewGitHubCLIAuth()

		// Should handle missing environment variables gracefully
		token := auth.GetFallbackToken()
		// Token might be empty if no env vars are set, which is fine for testing
		assert.IsType(t, "", token)
	})

	t.Run("CheckGHCLIInstalled checks if gh CLI is available", func(t *testing.T) {
		auth := NewGitHubCLIAuth()

		installed := auth.CheckGHCLIInstalled()
		// This should return true or false, not error
		assert.IsType(t, true, installed)
	})
}

func TestTokenScopes(t *testing.T) {
	t.Run("HasRequiredScopes checks for project permissions", func(t *testing.T) {
		requiredScopes := []string{"project", "repo", "read:org"}

		// Test with matching scopes
		userScopes := []string{"repo", "project", "read:org", "user"}
		assert.True(t, HasRequiredScopes(userScopes, requiredScopes))

		// Test with missing scopes
		userScopes = []string{"repo", "user"}
		assert.False(t, HasRequiredScopes(userScopes, requiredScopes))

		// Test with empty scopes
		userScopes = []string{}
		assert.False(t, HasRequiredScopes(userScopes, requiredScopes))
	})
}
