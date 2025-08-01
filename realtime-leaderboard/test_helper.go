package main

import (
	"context"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
)

// TestMain runs before all tests and can be used for global setup/teardown
func TestMain(m *testing.M) {
	// Check if Redis is available for testing
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       15, // Use DB 15 for testing
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		// Redis not available, skip tests that require it
		os.Exit(0)
	}

	// Clean up test database before running tests
	client.FlushDB(ctx)
	client.Close()

	// Run tests
	code := m.Run()

	// Clean up after tests
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       15,
	})
	client.FlushDB(ctx)
	client.Close()

	os.Exit(code)
}