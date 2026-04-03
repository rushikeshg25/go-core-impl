package main

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var ErrConflict = errors.New("conflict: record was modified by another transaction")

type Record struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	A             string    `json:"a"`
	B             string    `json:"b"`
	ConflictToken string    `json:"conflict_token"`
	CreatedAt     time.Time `json:"created_at"`
}

type ConcurrencyStrategy interface {
	Insert(name, a, b string) (int64, string, error)
	GetByID(id int64) (*Record, error)
	Update(id int64, name, a, b string, conflictToken string) (string, error)
}

// --- Version Strategy (MVCC) ---

type VersionStrategy struct {
	DB *sql.DB
}

func NewVersionStrategy(database *sql.DB) *VersionStrategy {
	return &VersionStrategy{DB: database}
}

func (s *VersionStrategy) Insert(name, a, b string) (int64, string, error) {
	result, err := s.DB.Exec(
		"INSERT INTO test (name, a, b, version) VALUES (?, ?, ?, 1)",
		name, a, b,
	)
	if err != nil {
		return 0, "", fmt.Errorf("insert failed: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, "", fmt.Errorf("failed to get insert id: %w", err)
	}

	return id, "1", nil
}

func (s *VersionStrategy) GetByID(id int64) (*Record, error) {
	row := s.DB.QueryRow(
		"SELECT id, name, a, b, version, created_at FROM test WHERE id = ?", id,
	)

	var r Record
	var version int
	err := row.Scan(&r.ID, &r.Name, &r.A, &r.B, &version, &r.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("record not found")
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	r.ConflictToken = strconv.Itoa(version)
	return &r, nil
}

func (s *VersionStrategy) Update(id int64, name, a, b string, conflictToken string) (string, error) {
	version, err := strconv.Atoi(conflictToken)
	if err != nil {
		return "", fmt.Errorf("invalid version token: %w", err)
	}

	result, err := s.DB.Exec(
		"UPDATE test SET name = ?, a = ?, b = ?, version = version + 1 WHERE id = ? AND version = ?",
		name, a, b, id, version,
	)
	if err != nil {
		return "", fmt.Errorf("update failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return "", ErrConflict
	}

	return strconv.Itoa(version + 1), nil
}

// --- Checksum Strategy ---

type ChecksumStrategy struct {
	DB *sql.DB
}

func NewChecksumStrategy(database *sql.DB) *ChecksumStrategy {
	return &ChecksumStrategy{DB: database}
}

func computeChecksum(name, a, b string) string {
	data := name + a + b
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func (s *ChecksumStrategy) Insert(name, a, b string) (int64, string, error) {
	checksum := computeChecksum(name, a, b)

	result, err := s.DB.Exec(
		"INSERT INTO test (name, a, b, checksum) VALUES (?, ?, ?, ?)",
		name, a, b, checksum,
	)
	if err != nil {
		return 0, "", fmt.Errorf("insert failed: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, "", fmt.Errorf("failed to get insert id: %w", err)
	}

	return id, checksum, nil
}

func (s *ChecksumStrategy) GetByID(id int64) (*Record, error) {
	row := s.DB.QueryRow(
		"SELECT id, name, a, b, checksum, created_at FROM test WHERE id = ?", id,
	)

	var r Record
	err := row.Scan(&r.ID, &r.Name, &r.A, &r.B, &r.ConflictToken, &r.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("record not found")
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &r, nil
}

func (s *ChecksumStrategy) Update(id int64, name, a, b string, conflictToken string) (string, error) {
	newChecksum := computeChecksum(name, a, b)

	result, err := s.DB.Exec(
		"UPDATE test SET name = ?, a = ?, b = ?, checksum = ? WHERE id = ? AND checksum = ?",
		name, a, b, newChecksum, id, conflictToken,
	)
	if err != nil {
		return "", fmt.Errorf("update failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return "", ErrConflict
	}

	return newChecksum, nil
}
