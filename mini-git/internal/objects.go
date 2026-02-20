package internal

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
)

func StoreObject(data []byte) (string, error) {
	hash := sha1.Sum(data)
	hashStr := fmt.Sprintf("%x", hash)

	dir := filepath.Join(".mini-git", "objects", hashStr[:2])
	path := filepath.Join(dir, hashStr[2:])

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", err
	}

	return hashStr, nil
}

func ReadObject(hashStr string) ([]byte, error) {
	if len(hashStr) < 2 {
		return nil, fmt.Errorf("invalid hash: %s", hashStr)
	}
	path := filepath.Join(".mini-git", "objects", hashStr[:2], hashStr[2:])
	return os.ReadFile(path)
}
