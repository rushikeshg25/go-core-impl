package main

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) *RedisClient {
	// Use a test database to avoid conflicts
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       15, // Use DB 15 for testing
	})

	ctx := context.Background()
	
	// Clear the test database
	err := client.FlushDB(ctx).Err()
	require.NoError(t, err)

	redisClient := &RedisClient{
		client: client,
		ctx:    ctx,
	}

	// Test connection
	err = redisClient.Ping()
	require.NoError(t, err, "Redis connection failed - make sure Redis is running")

	return redisClient
}

func teardownTestRedis(t *testing.T, client *RedisClient) {
	// Clean up test data
	err := client.client.FlushDB(client.ctx).Err()
	require.NoError(t, err)
	
	err = client.client.Close()
	require.NoError(t, err)
}

func TestRedisClient_UpdatePlayerScore(t *testing.T) {
	client := setupTestRedis(t)
	defer teardownTestRedis(t, client)

	tests := []struct {
		name       string
		playerName string
		score      float64
		wantErr    bool
	}{
		{
			name:       "valid player and score",
			playerName: "alice",
			score:      1500.0,
			wantErr:    false,
		},
		{
			name:       "update existing player",
			playerName: "alice",
			score:      2000.0,
			wantErr:    false,
		},
		{
			name:       "negative score",
			playerName: "bob",
			score:      -100.0,
			wantErr:    false,
		},
		{
			name:       "zero score",
			playerName: "charlie",
			score:      0.0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.UpdatePlayerScore(tt.playerName, tt.score)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				
				// Verify the score was set correctly
				score, err := client.client.ZScore(client.ctx, LEADERBOARD_KEY, tt.playerName).Result()
				assert.NoError(t, err)
				assert.Equal(t, tt.score, score)
			}
		})
	}
}

func TestRedisClient_GetPlayerRank(t *testing.T) {
	client := setupTestRedis(t)
	defer teardownTestRedis(t, client)

	// Setup test data
	testPlayers := map[string]float64{
		"alice":   2000.0,
		"bob":     1500.0,
		"charlie": 1800.0,
		"david":   1200.0,
	}

	for name, score := range testPlayers {
		err := client.UpdatePlayerScore(name, score)
		require.NoError(t, err)
	}

	tests := []struct {
		name         string
		playerName   string
		expectedRank int64
		expectedScore float64
		wantErr      bool
	}{
		{
			name:         "highest score player",
			playerName:   "alice",
			expectedRank: 1,
			expectedScore: 2000.0,
			wantErr:      false,
		},
		{
			name:         "middle score player",
			playerName:   "charlie",
			expectedRank: 2,
			expectedScore: 1800.0,
			wantErr:      false,
		},
		{
			name:         "lowest score player",
			playerName:   "david",
			expectedRank: 4,
			expectedScore: 1200.0,
			wantErr:      false,
		},
		{
			name:        "non-existent player",
			playerName:  "eve",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rank, score, err := client.GetPlayerRank(tt.playerName)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRank, rank)
				assert.Equal(t, tt.expectedScore, score)
			}
		})
	}
}

func TestRedisClient_GetTopPlayers(t *testing.T) {
	client := setupTestRedis(t)
	defer teardownTestRedis(t, client)

	// Setup test data
	testPlayers := map[string]float64{
		"alice":   2000.0,
		"bob":     1500.0,
		"charlie": 1800.0,
		"david":   1200.0,
		"eve":     2200.0,
		"frank":   1000.0,
	}

	for name, score := range testPlayers {
		err := client.UpdatePlayerScore(name, score)
		require.NoError(t, err)
	}

	tests := []struct {
		name          string
		count         int64
		expectedLen   int
		expectedTotal int64
		expectedFirst string
		expectedLast  string
	}{
		{
			name:          "get top 3",
			count:         3,
			expectedLen:   3,
			expectedTotal: 6,
			expectedFirst: "eve",   // 2200
			expectedLast:  "charlie", // 1800
		},
		{
			name:          "get top 10 (more than available)",
			count:         10,
			expectedLen:   6,
			expectedTotal: 6,
			expectedFirst: "eve",   // 2200
			expectedLast:  "frank", // 1000
		},
		{
			name:          "get top 1",
			count:         1,
			expectedLen:   1,
			expectedTotal: 6,
			expectedFirst: "eve", // 2200
			expectedLast:  "eve", // 2200
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, total, err := client.GetTopPlayers(tt.count)
			
			assert.NoError(t, err)
			assert.Len(t, results, tt.expectedLen)
			assert.Equal(t, tt.expectedTotal, total)
			
			if len(results) > 0 {
				assert.Equal(t, tt.expectedFirst, results[0].Member)
				assert.Equal(t, tt.expectedLast, results[len(results)-1].Member)
			}
		})
	}
}

func TestRedisClient_EmptyLeaderboard(t *testing.T) {
	client := setupTestRedis(t)
	defer teardownTestRedis(t, client)

	// Test getting top players from empty leaderboard
	results, total, err := client.GetTopPlayers(10)
	assert.NoError(t, err)
	assert.Empty(t, results)
	assert.Equal(t, int64(0), total)

	// Test getting rank for non-existent player
	_, _, err = client.GetPlayerRank("nonexistent")
	assert.Error(t, err)
}