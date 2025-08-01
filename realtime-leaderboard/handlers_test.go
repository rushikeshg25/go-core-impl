package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T) (*LeaderboardServer, *RedisClient) {
	redisClient := setupTestRedis(t)
	server := NewLeaderboardServer(redisClient)
	return server, redisClient
}

func TestLeaderboardServer_AddHandler(t *testing.T) {
	server, redisClient := setupTestServer(t)
	defer teardownTestRedis(t, redisClient)

	tests := []struct {
		name           string
		method         string
		formData       url.Values
		expectedStatus int
		expectedName   string
		expectedScore  float64
		expectedRank   int64
	}{
		{
			name:   "valid add request",
			method: "POST",
			formData: url.Values{
				"player": []string{"alice"},
				"score":  []string{"1500"},
			},
			expectedStatus: http.StatusOK,
			expectedName:   "alice",
			expectedScore:  1500.0,
			expectedRank:   1,
		},
		{
			name:   "update existing player",
			method: "POST",
			formData: url.Values{
				"player": []string{"alice"},
				"score":  []string{"2000"},
			},
			expectedStatus: http.StatusOK,
			expectedName:   "alice",
			expectedScore:  2000.0,
			expectedRank:   1,
		},
		{
			name:           "missing player parameter",
			method:         "POST",
			formData:       url.Values{"score": []string{"1500"}},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing score parameter",
			method:         "POST",
			formData:       url.Values{"player": []string{"alice"}},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "invalid score format",
			method: "POST",
			formData: url.Values{
				"player": []string{"alice"},
				"score":  []string{"invalid"},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "wrong HTTP method",
			method:         "GET",
			formData:       url.Values{},
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/add", strings.NewReader(tt.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			
			rr := httptest.NewRecorder()
			server.AddHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				var player Player
				err := json.Unmarshal(rr.Body.Bytes(), &player)
				require.NoError(t, err)
				
				assert.Equal(t, tt.expectedName, player.Name)
				assert.Equal(t, tt.expectedScore, player.Score)
				assert.Equal(t, tt.expectedRank, player.Rank)
			}
		})
	}
}

func TestLeaderboardServer_GetRankHandler(t *testing.T) {
	server, redisClient := setupTestServer(t)
	defer teardownTestRedis(t, redisClient)

	// Setup test data
	err := redisClient.UpdatePlayerScore("alice", 2000.0)
	require.NoError(t, err)
	err = redisClient.UpdatePlayerScore("bob", 1500.0)
	require.NoError(t, err)

	tests := []struct {
		name           string
		method         string
		queryParam     string
		expectedStatus int
		expectedFound  bool
		expectedName   string
		expectedScore  float64
		expectedRank   int64
	}{
		{
			name:           "existing player",
			method:         "GET",
			queryParam:     "alice",
			expectedStatus: http.StatusOK,
			expectedFound:  true,
			expectedName:   "alice",
			expectedScore:  2000.0,
			expectedRank:   1,
		},
		{
			name:           "non-existent player",
			method:         "GET",
			queryParam:     "charlie",
			expectedStatus: http.StatusOK,
			expectedFound:  false,
		},
		{
			name:           "missing player parameter",
			method:         "GET",
			queryParam:     "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "wrong HTTP method",
			method:         "POST",
			queryParam:     "alice",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/rank"
			if tt.queryParam != "" {
				url += "?player=" + tt.queryParam
			}
			
			req := httptest.NewRequest(tt.method, url, nil)
			rr := httptest.NewRecorder()
			
			server.GetRankHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				var response RankResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				
				assert.Equal(t, tt.expectedFound, response.Found)
				
				if tt.expectedFound {
					require.NotNil(t, response.Player)
					assert.Equal(t, tt.expectedName, response.Player.Name)
					assert.Equal(t, tt.expectedScore, response.Player.Score)
					assert.Equal(t, tt.expectedRank, response.Player.Rank)
				} else {
					assert.Nil(t, response.Player)
				}
			}
		})
	}
}

func TestLeaderboardServer_GetAllHandler(t *testing.T) {
	server, redisClient := setupTestServer(t)
	defer teardownTestRedis(t, redisClient)

	tests := []struct {
		name           string
		method         string
		setupData      map[string]float64
		expectedStatus int
		expectedCount  int
		expectedTotal  int64
		expectedFirst  string
		expectedLast   string
	}{
		{
			name:   "multiple players",
			method: "GET",
			setupData: map[string]float64{
				"alice":   2000.0,
				"bob":     1500.0,
				"charlie": 1800.0,
				"david":   1200.0,
			},
			expectedStatus: http.StatusOK,
			expectedCount:  4,
			expectedTotal:  4,
			expectedFirst:  "alice",  // 2000
			expectedLast:   "david",  // 1200
		},
		{
			name:           "empty leaderboard",
			method:         "GET",
			setupData:      map[string]float64{},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
			expectedTotal:  0,
		},
		{
			name:   "more than 10 players",
			method: "GET",
			setupData: map[string]float64{
				"p1":  1100.0, "p2": 1200.0, "p3": 1300.0, "p4": 1400.0, "p5": 1500.0,
				"p6":  1600.0, "p7": 1700.0, "p8": 1800.0, "p9": 1900.0, "p10": 2000.0,
				"p11": 2100.0, "p12": 2200.0,
			},
			expectedStatus: http.StatusOK,
			expectedCount:  10, // Should return only top 10
			expectedTotal:  12,
			expectedFirst:  "p12", // 2200
			expectedLast:   "p3",  // 1300 (10th place)
		},
		{
			name:           "wrong HTTP method",
			method:         "POST",
			setupData:      map[string]float64{},
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear and setup test data
			err := redisClient.client.FlushDB(redisClient.ctx).Err()
			require.NoError(t, err)
			
			for name, score := range tt.setupData {
				err := redisClient.UpdatePlayerScore(name, score)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(tt.method, "/all", nil)
			rr := httptest.NewRecorder()
			
			server.GetAllHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				var response LeaderboardResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				
				assert.Len(t, response.Players, tt.expectedCount)
				assert.Equal(t, tt.expectedTotal, response.Total)
				
				if tt.expectedCount > 0 {
					assert.Equal(t, tt.expectedFirst, response.Players[0].Name)
					assert.Equal(t, tt.expectedLast, response.Players[len(response.Players)-1].Name)
					
					// Verify ranking is correct
					for i, player := range response.Players {
						assert.Equal(t, int64(i+1), player.Rank)
					}
				}
			}
		})
	}
}

func TestLeaderboardServer_Integration(t *testing.T) {
	server, redisClient := setupTestServer(t)
	defer teardownTestRedis(t, redisClient)

	// Add some players
	players := []struct {
		name  string
		score string
	}{
		{"alice", "2000"},
		{"bob", "1500"},
		{"charlie", "1800"},
	}

	for _, p := range players {
		formData := url.Values{
			"player": []string{p.name},
			"score":  []string{p.score},
		}
		req := httptest.NewRequest("POST", "/add", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		
		rr := httptest.NewRecorder()
		server.AddHandler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	}

	// Test leaderboard
	req := httptest.NewRequest("GET", "/all", nil)
	rr := httptest.NewRecorder()
	server.GetAllHandler(rr, req)
	
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var response LeaderboardResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Len(t, response.Players, 3)
	assert.Equal(t, int64(3), response.Total)
	assert.Equal(t, "alice", response.Players[0].Name) // Highest score
	assert.Equal(t, int64(1), response.Players[0].Rank)

	// Test individual rank
	req = httptest.NewRequest("GET", "/rank?player=charlie", nil)
	rr = httptest.NewRecorder()
	server.GetRankHandler(rr, req)
	
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var rankResponse RankResponse
	err = json.Unmarshal(rr.Body.Bytes(), &rankResponse)
	require.NoError(t, err)
	
	assert.True(t, rankResponse.Found)
	assert.Equal(t, "charlie", rankResponse.Player.Name)
	assert.Equal(t, int64(2), rankResponse.Player.Rank) // Second place
}