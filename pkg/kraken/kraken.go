// kraken is a utility package with common functions shared across other packagges.
package kraken

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"

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

// StructToMap converts a struct into a map[string]any object.
func StructToMap(s any) (map[string]any, error) {
	m := make(map[string]any)
	sType := reflect.TypeOf(s)
	sValue := reflect.ValueOf(s)
	if sType.Kind() == reflect.Ptr {
		sType = sType.Elem()
		sValue = sValue.Elem()
	}
	if sType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("s is not a struct")
	}
	for i := range sType.NumField() {
		field := sType.Field(i)
		if !field.IsExported() {
			continue
		}
		key := field.Name
		fieldValue := sValue.FieldByName(key)
		mapTag := field.Tag.Get("map")
		mapTagParts := strings.Split(mapTag, ",")
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}
		jsonTagParts := strings.Split(jsonTag, ",")
		if slices.Contains(jsonTagParts, "omitempty") && fieldValue.IsZero() {
			continue
		}
		if len(jsonTagParts) > 0 && len(jsonTagParts[0]) > 0 {
			key = jsonTagParts[0]
		}
		// Convert fields with `map:"stringify"` to a JSON string.
		if slices.Contains(mapTagParts, "stringify") {
			value, err := json.Marshal(fieldValue.Interface())
			if err != nil {
				return nil, fmt.Errorf("stringify: %w", err)
			}
			m[key] = string(value)
			continue
		}
		// Dereference fields with pointer.
		for fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
			fieldValue = fieldValue.Elem()
		}
		value := fieldValue.Interface()
		// Process fields with structs.
		if fieldValue.Kind() == reflect.Struct {
			var err error
			value, err = StructToMap(value)
			if err != nil {
				return nil, fmt.Errorf("struct to map: %s", err)
			}
		}
		// Convert slice of structs to a slice of maps.
		if (fieldValue.Kind() == reflect.Array || fieldValue.Kind() == reflect.Slice) && fieldValue.Type().Elem().Kind() == reflect.Struct {
			var arr []any
			for j := range fieldValue.Len() {
				elem := fieldValue.Index(j).Interface()
				mapped, err := StructToMap(elem)
				if err != nil {
					return nil, fmt.Errorf("struct to map: %s", err)
				}
				arr = append(arr, mapped)
			}
			value = arr
		}
		// Copy embedded fields into the main map.
		if fieldValueMap, ok := value.(map[string]any); ok && field.Anonymous {
			maps.Copy(m, fieldValueMap)
			continue
		}
		m[key] = value
	}
	return m, nil
}
