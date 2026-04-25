package main

import (
	"sort"

	"github.com/zeebo/blake3"
)

func sortByName(files ...string) []string {
	filesArray := []string(files)
	sort.Strings(filesArray)
	return filesArray
}

func hashData(toHash string) [32]byte {
	data := []byte(toHash)
	hash := blake3.Sum256(data)
	return hash

}
