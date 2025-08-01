package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "environment variable exists",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
		{
			name:         "environment variable does not exist",
			key:          "NON_EXISTENT_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "empty environment variable",
			key:          "EMPTY_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up any existing env var
			os.Unsetenv(tt.key)
			
			// Set env var if specified
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnvWithRealEnvironmentVariables(t *testing.T) {
	// Test with actual environment variables that might be set
	tests := []struct {
		name         string
		key          string
		defaultValue string
	}{
		{
			name:         "REDIS_HOST with default",
			key:          "REDIS_HOST",
			defaultValue: "localhost",
		},
		{
			name:         "REDIS_PORT with default",
			key:          "REDIS_PORT",
			defaultValue: "6379",
		},
		{
			name:         "SERVER_PORT with default",
			key:          "SERVER_PORT",
			defaultValue: "8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getEnv(tt.key, tt.defaultValue)
			
			// Result should either be the env value or the default
			envValue := os.Getenv(tt.key)
			if envValue != "" {
				assert.Equal(t, envValue, result)
			} else {
				assert.Equal(t, tt.defaultValue, result)
			}
		})
	}
}