package bloomfilter

import (
	"hash"

	"github.com/twmb/murmur3"
)

type BloomFilter struct {
	bitset []bool
	size   uint32
	k      uint32
	hashes []hash.Hash32
}

func New(size uint32, k uint32) *BloomFilter {
	bf := &BloomFilter{
		bitset: make([]bool, size),
		size:   size,
		k:      k,
		hashes: make([]hash.Hash32, k),
	}
	var i uint32
	for i = 0; i < k; i++ {
		bf.hashes[i] = murmur3.SeedNew32(i)
	}
	return bf
}

func (bf *BloomFilter) Add(val []byte) {
	for _, h := range bf.hashes {
		h.Reset()
		_, _ = h.Write(val)
		idx := h.Sum32() % bf.size
		bf.bitset[idx] = true
	}
}

func (bf *BloomFilter) Test(val []byte) bool {
	for _, h := range bf.hashes {
		h.Reset()
		_, _ = h.Write(val)
		idx := h.Sum32() % bf.size
		if !bf.bitset[idx] {
			return false
		}
	}
	return true
}

func (bf *BloomFilter) Clear() {
	for i := range bf.bitset {
		bf.bitset[i] = false
	}
}
