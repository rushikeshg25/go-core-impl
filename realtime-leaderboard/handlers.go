package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type LeaderboardServer struct {
	redis *RedisClient
}

func NewLeaderboardServer(redisClient *RedisClient) *LeaderboardServer {
	return &LeaderboardServer{
		redis: redisClient,
	}
}

func (s *LeaderboardServer) GetRankHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	playerName := r.URL.Query().Get("player")
	if playerName == "" {
		http.Error(w, "player parameter is required", http.StatusBadRequest)
		return
	}

	rank, score, err := s.redis.GetPlayerRank(playerName)
	if err == redis.Nil {
		response := RankResponse{Found: false}
		writeJSONResponse(w, response)
		return
	}
	if err != nil {
		http.Error(w, "Error getting rank", http.StatusInternalServerError)
		return
	}

	player := &Player{
		Name:  playerName,
		Score: score,
		Rank:  rank,
	}

	response := RankResponse{
		Player: player,
		Found:  true,
	}

	writeJSONResponse(w, response)
}

func (s *LeaderboardServer) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	results, total, err := s.redis.GetTopPlayers(10)
	if err != nil {
		http.Error(w, "Error getting leaderboard", http.StatusInternalServerError)
		return
	}

	players := make([]Player, len(results))
	for i, result := range results {
		players[i] = Player{
			Name:  result.Member,
			Score: result.Score,
			Rank:  int64(i + 1),
		}
	}

	response := LeaderboardResponse{
		Players: players,
		Total:   total,
	}

	writeJSONResponse(w, response)
}

func (s *LeaderboardServer) AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	playerName := r.FormValue("player")
	scoreStr := r.FormValue("score")

	if playerName == "" || scoreStr == "" {
		http.Error(w, "player and score parameters are required", http.StatusBadRequest)
		return
	}

	score, err := strconv.ParseFloat(scoreStr, 64)
	if err != nil {
		http.Error(w, "Invalid score format", http.StatusBadRequest)
		return
	}

	err = s.redis.UpdatePlayerScore(playerName, score)
	if err != nil {
		http.Error(w, "Error updating leaderboard", http.StatusInternalServerError)
		return
	}

	// Get updated rank
	rank, updatedScore, err := s.redis.GetPlayerRank(playerName)
	if err != nil {
		http.Error(w, "Error getting updated rank", http.StatusInternalServerError)
		return
	}

	player := Player{
		Name:  playerName,
		Score: updatedScore,
		Rank:  rank,
	}

	writeJSONResponse(w, player)
}

func writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}