// kraken is a utility package with common functions shared across other packagges.
package kraken

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// UUID creates a new random UUID and returns it as a string or panics
func UUID() string {
	return uuid.NewString()
}

// ToJSON returns the JSON string of v or panics.
func ToJSON(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}

// ToJSONIndent returns the JSON string of v with indents or panics.
func ToJSONIndent(v any) string {
	data, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(data)
}

// Maps recursively merges multiple maps of type map[string]any.
func Maps(p map[string]any, s ...map[string]any) map[string]any {
	result := make(map[string]any)
	subs := make(map[string][]map[string]any)
	for _, entry := range append([]map[string]any{p}, s...) {
		for key, value := range entry {
			switch value := value.(type) {
			case map[string]any:
				subs[key] = append(subs[key], value)
			default:
				result[key] = value
			}
		}
	}
	for key, sub := range subs {
		if len(sub) == 0 {
			continue
		}
		result[key] = Maps(sub[0], sub[1:]...)
	}
	return result
}

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

// Must handles conversion.
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
