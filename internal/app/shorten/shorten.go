package shorten

import (
	"crypto/sha256"

	"github.com/btcsuite/btcutil/base58"
)

// Shorten is a function for short URL creation
func Shorten(originalURL []byte) string {
	return encodeToBase58(hashURL(originalURL))
}

func hashURL(originalURL []byte) []byte {
	h := sha256.New()
	h.Write(originalURL)
	return h.Sum(nil)
}

func encodeToBase58(hashed []byte) string {
	return base58.Encode(hashed)
}
