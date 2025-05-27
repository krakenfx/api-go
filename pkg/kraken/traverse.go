package kraken

import (
	"fmt"
	"reflect"
)

// Traverse returns the value in a nested structure based on a sequence of keys.
// Iteration-based map indexes are supported. However, they are non deterministic as Go maps are unordered.
func Traverse[V any](m any, keys ...any) (*V, error) {
	if len(keys) == 0 {
		m, ok := m.(V)
		if ok {
			return &m, nil
		}
		return nil, fmt.Errorf("no key provided")
	}
	cursor := m
	for _, key := range keys {
		switch c := cursor.(type) {
		case map[any]any:
			var ok bool
			cursor, ok = c[key]
			if !ok {
				return nil, fmt.Errorf("%s not found", key)
			}
		case map[string]any:
			var keyString string
			keyInt, ok := key.(int)
			if ok {
				if keyInt < 0 || keyInt >= len(c) {
					return nil, fmt.Errorf("index %d out of range (max %d)", key, len(c)-1)
				}
				var i int
				for k := range c {
					if i == keyInt {
						keyString = k
						break
					}
					i++
				}
			} else {
				keyString = fmt.Sprintf("%v", key)
			}
			cursor, ok = c[keyString]
			if !ok {
				return nil, fmt.Errorf("%s not found", key)
			}
		case []any:
			key, ok := key.(int)
			if !ok {
				return nil, fmt.Errorf("%v key assertion failed", key)
			}
			if key < 0 || key >= len(c) {
				return nil, fmt.Errorf("index %d out of range (max %d)", key, len(c)-1)
			}
			cursor = c[key]
		default:
			return nil, fmt.Errorf("%s cannot be traversed as it is a %s", key, reflect.TypeOf(c))
		}
	}
	result, ok := cursor.(V)
	if !ok {
		return nil, fmt.Errorf("final value assertion failed, expected %s but got %s (%v)", reflect.TypeOf(result), reflect.TypeOf(cursor), cursor)
	}
	return &result, nil
}
