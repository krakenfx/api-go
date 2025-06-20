package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

// MarshalOptions contains the parameters for [MarshalFunction].
type MarshalOptions struct {
	MediaType string
	Object    any
}

// MarshalResult contains the result of [MarshalFunction].
type MarshalResult struct {
	Data        []byte
	ContentType string
}

// MarshalFunction encodes the provided object according to the media type.
type MarshalFunction func(opts MarshalOptions) (result MarshalResult, err error)

// Marshal implements [MarshalFunction] by encoding the provided object according to the media type:
//
// - "multipart/form-data": expects object to be a map[string]any and sets it as a multipart form.
//
// - "application/x-www-form-urlencoded": expects object to be a map[string]any and encodes it as form values.
//
// - "application/json": marshals object as JSON.
//
// Returns an error if the media type is not among the list or the data doesn't match the expected format.
func Marshal(opts MarshalOptions) (result MarshalResult, err error) {
	obj := opts.Object
	if GetDirectReflection(opts.Object).Type.Kind() == reflect.Struct {
		var err error
		obj, err = StructToMap(obj)
		if err != nil {
			return result, fmt.Errorf("struct to map: %w", err)
		}
	}
	result.ContentType = opts.MediaType
	switch opts.MediaType {
	case "multipart/form-data":
		m, ok := obj.(map[string]any)
		if !ok {
			return result, fmt.Errorf("MarshalOptions.Object must be map[string]any")
		}
		form, err := NewForm(m)
		if err != nil {
			return result, err
		}
		result.Data = form.Bytes()
		result.ContentType = form.FormDataContentType()
	case "application/x-www-form-urlencoded":
		values, err := ToURLValues(obj)
		if err != nil {
			return result, fmt.Errorf("to url values: %w", err)
		}
		result.Data = []byte(values.Encode())

	case "application/json":
		result.Data, err = json.Marshal(obj)
		if err != nil {
			return result, fmt.Errorf("MarshalOptions.Object \"%+v\" json marshal failed: %w", result.Data, err)
		}
	default:
		if opts.MediaType == "" {
			return result, fmt.Errorf("unspecified content type")
		} else {
			return result, fmt.Errorf("content type \"%s\" not supported", opts.MediaType)
		}
	}
	return
}

// CreateReadCloser constructs an [io.ReadCloser] from a byte slice.
func CreateReadCloser(b []byte) io.ReadCloser {
	return io.NopCloser(bytes.NewReader(b))
}
