package helper

import (
	"encoding/json"
	"net/url"
	"testing"
)

func TestMarshal(t *testing.T) {
	object := struct {
		Value string `json:"child"`
	}{
		Value: "Hello, world!",
	}
	result, err := Marshal(MarshalOptions{
		MediaType: "application/json",
		Object:    object,
	})
	if err != nil {
		t.Errorf("Marshal(application/json): %s", err)
	}
	var m map[string]any
	if err := json.Unmarshal(result.Data, &m); err != nil {
		t.Errorf("json.Unmarshal(result.Data): %s", err)
	}
	result, err = Marshal(MarshalOptions{
		MediaType: "application/x-www-form-urlencoded",
		Object:    object,
	})
	if err != nil {
		t.Errorf("Marshal(application/x-www-form-urlencoded): %s", err)
	}
	if _, err := url.ParseQuery(string(result.Data)); err != nil {
		t.Errorf("url.ParseQuery(result.Data): %s", err)
	}

}
