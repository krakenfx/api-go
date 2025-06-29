package helper

import "encoding/json"

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
