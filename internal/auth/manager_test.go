package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthManager(t *testing.T) {
	t.Run("NewAuthManager creates new manager", func(t *testing.T) {
		manager := NewAuthManager()
		assert.NotNil(t, manager)
		assert.NotNil(t, manager.ghAuth)
	})

	t.Run("GetTokenWithoutValidation returns token from available sources", func(t *testing.T) {
		manager := NewAuthManager()

		token, err := manager.GetTokenWithoutValidation()

		// Should either succeed or fail gracefully
		if err != nil {
			assert.Contains(t, err.Error(), "no GitHub token found")
		} else {
			assert.NotEmpty(t, token)
		}
	})

	t.Run("GetAuthenticationStatus returns detailed status", func(t *testing.T) {
		manager := NewAuthManager()

		status := manager.GetAuthenticationStatus()

		// Should always provide status information
		assert.NotNil(t, status)

		// GH CLI installation status should be boolean
		assert.IsType(t, true, status.GHCLIInstalled)
		assert.IsType(t, true, status.HasEnvToken)
		assert.IsType(t, true, status.TokenAvailable)

		// Should provide recommendation
		recommendation := status.GetRecommendation()
		assert.NotEmpty(t, recommendation)

		// IsReady should be deterministic based on status
		isReady := status.IsReady()
		expectedReady := status.TokenAvailable && status.TokenValid && status.HasRequiredScopes
		assert.Equal(t, expectedReady, isReady)
	})

	t.Run("CheckAuthentication works without token", func(t *testing.T) {
		manager := NewAuthManager()

		err := manager.CheckAuthentication()

		// Should return error if no valid token
		if err != nil {
			assert.Contains(t, err.Error(), "no valid GitHub token found")
		}
		// If no error, authentication is working
	})
}

func TestAuthStatus(t *testing.T) {
	t.Run("IsReady returns correct status", func(t *testing.T) {
		// Test fully ready status
		readyStatus := AuthStatus{
			TokenAvailable:    true,
			TokenValid:        true,
			HasRequiredScopes: true,
		}
		assert.True(t, readyStatus.IsReady())

		// Test not ready status
		notReadyStatus := AuthStatus{
			TokenAvailable:    true,
			TokenValid:        false,
			HasRequiredScopes: true,
		}
		assert.False(t, notReadyStatus.IsReady())
	})

	t.Run("GetRecommendation provides helpful guidance", func(t *testing.T) {
		// Test no GH CLI
		status := AuthStatus{GHCLIInstalled: false}
		rec := status.GetRecommendation()
		assert.Contains(t, rec, "Install GitHub CLI")

		// Test no token
		status = AuthStatus{
			GHCLIInstalled: true,
			TokenAvailable: false,
		}
		rec = status.GetRecommendation()
		assert.Contains(t, rec, "gh auth login")

		// Test invalid token
		status = AuthStatus{
			GHCLIInstalled: true,
			TokenAvailable: true,
			TokenValid:     false,
		}
		rec = status.GetRecommendation()
		assert.Contains(t, rec, "Re-authenticate")

		// Test missing scopes
		status = AuthStatus{
			GHCLIInstalled:    true,
			TokenAvailable:    true,
			TokenValid:        true,
			HasRequiredScopes: false,
		}
		rec = status.GetRecommendation()
		assert.Contains(t, rec, "Grant additional scopes")

		// Test all good
		status = AuthStatus{
			GHCLIInstalled:    true,
			TokenAvailable:    true,
			TokenValid:        true,
			HasRequiredScopes: true,
		}
		rec = status.GetRecommendation()
		assert.Contains(t, rec, "properly configured")
	})
}

func TestGetMissingScopes(t *testing.T) {
	t.Run("Returns missing scopes correctly", func(t *testing.T) {
		available := []string{"repo", "user"}
		required := []string{"repo", "project", "user"}

		missing := getMissingScopes(available, required)

		assert.Equal(t, []string{"project"}, missing)
	})

	t.Run("Returns empty when all scopes available", func(t *testing.T) {
		available := []string{"repo", "project", "user"}
		required := []string{"repo", "project"}

		missing := getMissingScopes(available, required)

		assert.Empty(t, missing)
	})

	t.Run("Returns all required when none available", func(t *testing.T) {
		available := []string{"user"}
		required := []string{"repo", "project"}

		missing := getMissingScopes(available, required)

		assert.Equal(t, []string{"repo", "project"}, missing)
	})
}
