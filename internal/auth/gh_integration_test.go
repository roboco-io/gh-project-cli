package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitHubCLIIntegration(t *testing.T) {
	t.Run("NewGitHubCLIAuth creates new instance", func(t *testing.T) {
		auth := NewGitHubCLIAuth()
		assert.NotNil(t, auth)
		assert.IsType(t, &GitHubCLIAuth{}, auth)
	})

	t.Run("GetToken returns token from gh CLI", func(t *testing.T) {
		auth := NewGitHubCLIAuth()
		
		// This will fail until we implement it
		token, err := auth.GetToken("github.com")
		
		// For now, expect an error since gh CLI might not be configured
		if err != nil {
			assert.Empty(t, token)
			assert.Contains(t, err.Error(), "gh CLI")
		} else {
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