package sign

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

// HMACSigner is a Gengo Signer using HMAC and SHA1
type HMACSigner struct {
	Key []byte
}

// NewHMACSigner returns a new HMACSigner with the given key
func NewHMACSigner(key []byte) *HMACSigner {
	return &HMACSigner{
		Key: key,
	}
}

// Sign signs some data using the HMAC/SHA1 algorithm
// Create a new hasher everytime so we can concurrently run requests
func (h *HMACSigner) Sign(data string) string {
	hasher := hmac.New(sha1.New, h.Key)
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}
