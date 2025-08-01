package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const LEADERBOARD_KEY = "game:leaderboard"

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient() *RedisClient {
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")
	password := getEnv("REDIS_PASSWORD", "")
	dbStr := getEnv("REDIS_DB", "0")

	
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		db = 0 
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisClient{
		client: rdb,
		ctx:    context.Background(),
	}
}



func (r *RedisClient) Ping() error {
	_, err := r.client.Ping(r.ctx).Result()
	return err
}

func (r *RedisClient) GetPlayerRank(playerName string) (int64, float64, error) {
	rank, err := r.client.ZRevRank(r.ctx, LEADERBOARD_KEY, playerName).Result()
	if err != nil {
		return 0, 0, err
	}

	score, err := r.client.ZScore(r.ctx, LEADERBOARD_KEY, playerName).Result()
	if err != nil {
		return 0, 0, err
	}

	return rank + 1, score, nil // Convert to 1-based ranking
}

func (r *RedisClient) GetTopPlayers(count int64) ([]redis.Z, int64, error) {
	results, err := r.client.ZRevRangeWithScores(r.ctx, LEADERBOARD_KEY, 0, count-1).Result()
	if err != nil {
		return nil, 0, err
	}

	total, err := r.client.ZCard(r.ctx, LEADERBOARD_KEY).Result()
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (r *RedisClient) UpdatePlayerScore(playerName string, score float64) error {
	return r.client.ZAdd(r.ctx, LEADERBOARD_KEY, redis.Z{
		Score:  score,
		Member: playerName,
	}).Err()
}