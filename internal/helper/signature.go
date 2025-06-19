package helper

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
)

// Sign encodes a message into HMAC-SHA-512 hashed with a base64-encoded key.
func Sign(key string, message []byte) (string, error) {
	keyDecoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", fmt.Errorf("decode private key failed: %s", err)
	}
	hmacHash := hmac.New(sha512.New, keyDecoded)
	hmacHash.Write(message)
	return base64.StdEncoding.EncodeToString(hmacHash.Sum(nil)), nil
}
