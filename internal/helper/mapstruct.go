package helper

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
)

// StructToMapFunction converts a struct into map[string]any.
type StructToMapFunction func(s any) (map[string]any, error)

// StructToMap converts a struct into map[string]any.
func StructToMap(s any) (map[string]any, error) {
	m := make(map[string]any)
	r := GetDirectReflection(s)
	if r.Type.Kind() != reflect.Struct {
		return nil, fmt.Errorf("s is not a struct")
	}
	for i := range r.Type.NumField() {
		field := r.Type.Field(i)
		if !field.IsExported() {
			continue
		}
		value := GetDirectReflection(r.Value.Interface()).Value.FieldByName(field.Name)
		mapTag := field.Tag.Get("map")
		jsonTag := field.Tag.Get("json")
		// Ignore this field if `json` or `map` is "-".
		if jsonTag == "-" || mapTag == "-" {
			continue
		}
		// Split `json` tag into a slice.
		jsonTagParts := strings.Split(jsonTag, ",")
		// Ignore this field if `json` has "omitempty" and value is zero.
		if slices.Contains(jsonTagParts, "omitempty") && value.IsZero() {
			continue
		}
		// Define the field name.
		name := field.Name
		if len(jsonTagParts) > 0 && len(jsonTagParts[0]) > 0 {
			name = jsonTagParts[0]
		}
		// Split `map` tag into a slice.
		mapTagParts := strings.Split(mapTag, ",")
		// If `map` has "stringify", convert to a JSON string.
		if slices.Contains(mapTagParts, "stringify") {
			value, err := json.Marshal(value.Interface())
			if err != nil {
				return nil, fmt.Errorf("stringify: %w", err)
			}
			m[name] = string(value)
			continue
		}
		// Convert field of map type into struct.
		if value.Kind() == reflect.Struct {
			var err error
			fieldMap, err := StructToMap(value.Interface())
			if err != nil {
				return nil, fmt.Errorf("struct to map: %s", err)
			}
			value = reflect.ValueOf(fieldMap)
		}
		// Convert field of []struct into []map[string]any.
		if (value.Kind() == reflect.Array || value.Kind() == reflect.Slice) && value.Type().Elem().Kind() == reflect.Struct {
			var arr []any
			for j := range value.Len() {
				elem := value.Index(j).Interface()
				mapped, err := StructToMap(elem)
				if err != nil {
					return nil, fmt.Errorf("struct to map: %s", err)
				}
				arr = append(arr, mapped)
			}
			value = reflect.ValueOf(arr)
		}
		// Copy embedded fields into the main map.
		if fieldValueMap, ok := value.Interface().(map[string]any); ok && field.Anonymous {
			maps.Copy(m, fieldValueMap)
			continue
		}
		m[name] = value.Interface()
	}
	return m, nil
}

// DirectReflection contains the result of [GetDirectReflection].
type DirectReflection struct {
	Value reflect.Value
	Type  reflect.Type
}

// GetDirectReflectionFunction returns the direct reflect.Value of an object without the pointer.
type GetDirectReflectionFunction func(s any) DirectReflection

// GetDirectReflection returns the direct reflect.Value of an object without the pointer.
func GetDirectReflection(s any) DirectReflection {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
		t = t.Elem()
	}
	return DirectReflection{
		Value: v,
		Type:  t,
	}
}
