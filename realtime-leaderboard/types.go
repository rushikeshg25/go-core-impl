package main

type Player struct {
	Name  string  `json:"name"`
	Score float64 `json:"score"`
	Rank  int64   `json:"rank"`
}

type RankResponse struct {
	Player *Player `json:"player"`
	Found  bool    `json:"found"`
}

type LeaderboardResponse struct {
	Players []Player `json:"players"`
	Total   int64    `json:"total"`
}