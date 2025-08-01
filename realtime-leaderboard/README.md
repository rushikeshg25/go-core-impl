# Gaming Leaderboard API

A realtime gaming leaderboard built with Go and Redis sorted sets. Provides fast ranking operations for competitive gaming applications.

## Installation

1. **Set up environment variables:**
   ```bash
   cp .env.example .env
   ```

2. **Configure your Redis connection** (edit `.env`):
   ```env
   REDIS_HOST=localhost
   REDIS_PORT=6379
   REDIS_PASSWORD=your-password
   REDIS_DB=0
   SERVER_PORT=8080
   ```

## Running Redis

### Using Docker (Recommended)
```bash
docker run -d --name redis-leaderboard -p 6379:6379 redis:alpine
```

## Running the Application

1. **Start the server:**
   ```bash
   go run .
   ```

2. **Verify it's running:**
   ```
   Leaderboard server starting on :8080
   Connected to Redis successfully
   ```

## API Endpoints

### 1. Add/Update Player Score
```bash
curl -X POST http://localhost:8080/add \
  -d "player=alice&score=1500"
```

**Response:**
```json
{
  "name": "alice",
  "score": 1500,
  "rank": 1
}
```

### 2. Get Player Rank
```bash
curl "http://localhost:8080/rank?player=alice"
```

**Response:**
```json
{
  "player": {
    "name": "alice",
    "score": 1500,
    "rank": 1
  },
  "found": true
}
```

### 3. Get Top 10 Leaderboard
```bash
curl http://localhost:8080/all
```

**Response:**
```json
{
  "players": [
    {
      "name": "alice",
      "score": 1500,
      "rank": 1
    },
    {
      "name": "bob",
      "score": 1200,
      "rank": 2
    }
  ],
  "total": 2
}
```


### Running Tests

**Prerequisites for testing:**
- Redis server running on localhost:6379
- Tests use database 15 to avoid conflicts with your data

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test file
go test -v redis_test.go redis.go types.go test_helper.go

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...
```



### Building for Production
```bash
go build -o leaderboard .
./leaderboard
```


