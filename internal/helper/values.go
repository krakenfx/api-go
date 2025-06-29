package helper

import (
	"fmt"
	"net/url"
	"reflect"
)

// ToURLValues converts v into [url.Values].
func ToURLValues(val any) (values url.Values, err error) {
	values = make(url.Values)
	if GetDirectReflection(val).Type.Kind() == reflect.Struct {
		val, err = StructToMap(val)
		if err != nil {
			return values, fmt.Errorf("struct to map: %w", err)
		}
	}
	switch m := val.(type) {
	case map[string]any:
		for k, v := range m {
			values[k], err = StringSlice(v)
			if err != nil {
				return values, fmt.Errorf("string slice: %w", err)
			}
		}
	case url.Values:
		values = m
		return values, fmt.Errorf("invalid type v of %s", reflect.TypeOf(val))
	}
	return
}
