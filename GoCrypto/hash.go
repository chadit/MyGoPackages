package gocrypto

import "hash/fnv"

// HashFNV1a32 accepts a string and returns a 32 bit hash based off FNV-1a, non-cryptographic hash function.
func HashFNV1a32(s string) uint32 {
	h := fnv.New32a()
	if _, err := h.Write([]byte(s)); err != nil {
		return 0 // This will never occur under current Write implementation,  https://github.com/golang/go/blob/master/src/hash/fnv/fnv.go#L84
	}
	return h.Sum32()
}

// HashFNV1a64 accepts a string and returns a 64 bit hash based off FNV-1a, non-cryptographic hash function.
func HashFNV1a64(s string) uint64 {
	h := fnv.New64a()
	if _, err := h.Write([]byte(s)); err != nil {
		return 0 // This will never occur under current Write implementation,  https://github.com/golang/go/blob/master/src/hash/fnv/fnv.go#L104
	}
	return h.Sum64()
}
