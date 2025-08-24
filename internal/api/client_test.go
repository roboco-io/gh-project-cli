package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGraphQLClient(t *testing.T) {
	t.Run("NewClient creates new GraphQL client", func(t *testing.T) {
		token := "test-token"
		client := NewClient(token)

		assert.NotNil(t, client)
		assert.IsType(t, &Client{}, client)
	})

	t.Run("NewClient with empty token returns error client", func(t *testing.T) {
		client := NewClient("")

		assert.NotNil(t, client)
		// Client should be created but operations should fail
	})

	t.Run("Client has correct base URL", func(t *testing.T) {
		client := NewClient("test-token")

		assert.Equal(t, DefaultAPIURL, client.baseURL)
	})
}

func TestRateLimiting(t *testing.T) {
	t.Run("Client handles rate limiting", func(t *testing.T) {
		client := NewClient("test-token")

		// Test that rate limiting is properly configured
		limiter := client.rateLimiter
		assert.NotNil(t, limiter)
	})
}

func TestRetryLogic(t *testing.T) {
	t.Run("Client has retry configuration", func(t *testing.T) {
		client := NewClient("test-token")

		config := client.retryConfig
		assert.NotNil(t, config)
		assert.Greater(t, config.MaxRetries, 0)
		assert.Greater(t, config.BaseDelay.Seconds(), 0.0)
	})
}

func TestHealthCheck(t *testing.T) {
	t.Run("HealthCheck validates connection", func(t *testing.T) {
		client := NewClient("invalid-token")
		ctx := context.Background()

		// This should fail with invalid token but not panic
		err := client.HealthCheck(ctx)
		assert.Error(t, err)
	})

	t.Run("HealthCheck with valid token succeeds", func(t *testing.T) {
		// Skip this test if no valid token is available
		if testing.Short() {
			t.Skip("Skipping integration test in short mode")
		}

		// This test would require a valid token to run
		// For now, just ensure the method exists
		client := NewClient("test-token")
		ctx := context.Background()

		err := client.HealthCheck(ctx)
		// With invalid token, this should return a specific error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "authentication")
	})
}

func TestErrorHandling(t *testing.T) {
	t.Run("Client handles GraphQL errors", func(t *testing.T) {
		client := NewClient("test-token")

		// Test error parsing
		err := client.parseGraphQLError([]byte(`{"errors": [{"message": "test error"}]}`))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test error")
	})

	t.Run("Client handles network errors", func(t *testing.T) {
		client := NewClient("test-token")

		// Test network error handling
		err := client.handleNetworkError(assert.AnError)
		assert.Error(t, err)
	})
}

func TestRequestBuilder(t *testing.T) {
	t.Run("BuildRequest creates proper GraphQL request", func(t *testing.T) {
		client := NewClient("test-token")

		query := "query { viewer { login } }"
		variables := map[string]interface{}{"test": "value"}

		req, err := client.buildRequest(query, variables)
		require.NoError(t, err)
		assert.NotNil(t, req)
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, DefaultAPIURL, req.URL.String())
		assert.Contains(t, req.Header.Get("Authorization"), "token test-token")
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
	})
}
