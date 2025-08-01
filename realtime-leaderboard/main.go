package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	redisClient := NewRedisClient()
	if err := redisClient.Ping(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	fmt.Println("Connected to Redis successfully")

	server := NewLeaderboardServer(redisClient)

	http.HandleFunc("/rank", server.GetRankHandler)
	http.HandleFunc("/all", server.GetAllHandler)
	http.HandleFunc("/add", server.AddHandler)

	port := getEnv("SERVER_PORT", "8080")

	fmt.Printf("Leaderboard server starting on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}