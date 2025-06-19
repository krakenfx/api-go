package helper

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// StringSlice takes a value of type any and converts them into a string slice.
func StringSlice(v any) (slice []string, err error) {
	switch v := v.(type) {
	case string:
		slice = []string{v}
	case []string:
		slice = v
	case []any:
		for _, item := range v {
			strings, err := StringSlice(item)
			if err != nil {
				return slice, err
			}
			slice = append(slice, strings...)
		}
	default:
		s, err := json.Marshal(v)
		if err != nil {
			return slice, fmt.Errorf("json marshal v of type %s: %w", reflect.TypeOf(v).Name(), err)
		}
		slice = []string{string(s)}
	}
	return
}

// MustStringSlice takes a value of type any and converts them into a string slice. Panics on error.
func MustStringSlice(v any) []string {
	return Must(StringSlice(v))
}
