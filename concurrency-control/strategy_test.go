package main

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	dsn := "root:@tcp(localhost:3306)/"
	var err error
	testDB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	_, err = testDB.Exec("DROP DATABASE IF EXISTS occ_test")
	if err != nil {
		log.Fatalf("Failed to drop test database: %v", err)
	}
	_, err = testDB.Exec("CREATE DATABASE occ_test")
	if err != nil {
		log.Fatalf("Failed to create test database: %v", err)
	}
	testDB.Close()

	testDB, err = sql.Open("mysql", "root:@tcp(localhost:3306)/occ_test?parseTime=true")
	if err != nil {
		log.Fatalf("Failed to connect to occ_test: %v", err)
	}

	_, err = testDB.Exec(`
	CREATE TABLE IF NOT EXISTS test (
		id         INT AUTO_INCREMENT PRIMARY KEY,
		name       VARCHAR(255),
		a          TEXT,
		b          TEXT,
		version    INT DEFAULT 1,
		checksum   VARCHAR(64),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		log.Fatalf("Failed to create test table: %v", err)
	}

	code := m.Run()

	testDB.Exec("DROP DATABASE IF EXISTS occ_test")
	testDB.Close()
	os.Exit(code)
}

func truncateTable(t *testing.T) {
	t.Helper()
	_, err := testDB.Exec("TRUNCATE TABLE test")
	if err != nil {
		t.Fatalf("Failed to truncate table: %v", err)
	}
}

// --- Version Strategy Tests ---

func TestVersionInsertAndGet(t *testing.T) {
	truncateTable(t)
	s := NewVersionStrategy(testDB)

	id, token, err := s.Insert("alice", "hello", "world")
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}
	if token != "1" {
		t.Fatalf("Expected version token '1', got '%s'", token)
	}

	record, err := s.GetByID(id)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if record.Name != "alice" || record.A != "hello" || record.B != "world" {
		t.Fatalf("Unexpected record data: %+v", record)
	}
	if record.ConflictToken != "1" {
		t.Fatalf("Expected conflict token '1', got '%s'", record.ConflictToken)
	}
}

func TestVersionUpdateSuccess(t *testing.T) {
	truncateTable(t)
	s := NewVersionStrategy(testDB)

	id, token, err := s.Insert("alice", "hello", "world")
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	newToken, err := s.Update(id, "alice-updated", "foo", "bar", token)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if newToken != "2" {
		t.Fatalf("Expected new version '2', got '%s'", newToken)
	}

	record, err := s.GetByID(id)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if record.Name != "alice-updated" || record.A != "foo" || record.B != "bar" {
		t.Fatalf("Unexpected updated data: %+v", record)
	}
	if record.ConflictToken != "2" {
		t.Fatalf("Expected conflict token '2', got '%s'", record.ConflictToken)
	}
}

func TestVersionUpdateConflict(t *testing.T) {
	truncateTable(t)
	s := NewVersionStrategy(testDB)

	id, _, err := s.Insert("alice", "hello", "world")
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	// Update with wrong version
	_, err = s.Update(id, "bob", "x", "y", "99")
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("Expected ErrConflict, got: %v", err)
	}
}

func TestVersionConcurrentUpdates(t *testing.T) {
	truncateTable(t)
	s := NewVersionStrategy(testDB)

	id, token, err := s.Insert("alice", "hello", "world")
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	var wg sync.WaitGroup
	results := make(chan error, 2)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := s.Update(id, fmt.Sprintf("writer-%d", i), "a", "b", token)
			results <- err
		}(i)
	}

	wg.Wait()
	close(results)

	var successes, conflicts int
	for err := range results {
		if err == nil {
			successes++
		} else if errors.Is(err, ErrConflict) {
			conflicts++
		} else {
			t.Fatalf("Unexpected error: %v", err)
		}
	}

	if successes != 1 || conflicts != 1 {
		t.Fatalf("Expected 1 success and 1 conflict, got %d successes and %d conflicts", successes, conflicts)
	}
}

// --- Checksum Strategy Tests ---

func TestChecksumInsertAndGet(t *testing.T) {
	truncateTable(t)
	s := NewChecksumStrategy(testDB)

	id, token, err := s.Insert("alice", "hello", "world")
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	// Verify it's a valid hex-encoded SHA-256 (64 hex chars)
	if len(token) != 64 {
		t.Fatalf("Expected 64-char hex checksum, got '%s' (len=%d)", token, len(token))
	}
	if _, err := hex.DecodeString(token); err != nil {
		t.Fatalf("Checksum is not valid hex: %v", err)
	}

	record, err := s.GetByID(id)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if record.Name != "alice" || record.A != "hello" || record.B != "world" {
		t.Fatalf("Unexpected record data: %+v", record)
	}
	if record.ConflictToken != token {
		t.Fatalf("Expected conflict token '%s', got '%s'", token, record.ConflictToken)
	}
}

func TestChecksumUpdateSuccess(t *testing.T) {
	truncateTable(t)
	s := NewChecksumStrategy(testDB)

	id, token, err := s.Insert("alice", "hello", "world")
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	newToken, err := s.Update(id, "alice-updated", "foo", "bar", token)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if newToken == token {
		t.Fatal("Expected checksum to change after update")
	}

	record, err := s.GetByID(id)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if record.Name != "alice-updated" || record.A != "foo" || record.B != "bar" {
		t.Fatalf("Unexpected updated data: %+v", record)
	}
	if record.ConflictToken != newToken {
		t.Fatalf("Expected conflict token '%s', got '%s'", newToken, record.ConflictToken)
	}
}

func TestChecksumUpdateConflict(t *testing.T) {
	truncateTable(t)
	s := NewChecksumStrategy(testDB)

	id, _, err := s.Insert("alice", "hello", "world")
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	// Update with wrong checksum
	_, err = s.Update(id, "bob", "x", "y", "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("Expected ErrConflict, got: %v", err)
	}
}

func TestChecksumConcurrentUpdates(t *testing.T) {
	truncateTable(t)
	s := NewChecksumStrategy(testDB)

	id, token, err := s.Insert("alice", "hello", "world")
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	var wg sync.WaitGroup
	results := make(chan error, 2)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := s.Update(id, fmt.Sprintf("writer-%d", i), "a", "b", token)
			results <- err
		}(i)
	}

	wg.Wait()
	close(results)

	var successes, conflicts int
	for err := range results {
		if err == nil {
			successes++
		} else if errors.Is(err, ErrConflict) {
			conflicts++
		} else {
			t.Fatalf("Unexpected error: %v", err)
		}
	}

	if successes != 1 || conflicts != 1 {
		t.Fatalf("Expected 1 success and 1 conflict, got %d successes and %d conflicts", successes, conflicts)
	}
}
