package cookie

import (
	"crypto/rand"
	"encoding/base64"
)

func divUp(numer, denom int) int {
	return (numer + denom - 1) / denom
}

// Generate session / api key / invite cookies securely
// with at least length bits of entropy
func Generate(length int) (string, error) {
	b := make([]byte, divUp(length, 24)*3)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
