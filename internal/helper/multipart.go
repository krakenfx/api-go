package helper

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"reflect"
)

// Form is a combined representation of [bytes.Buffer] and [multipart.Writer].
type Form struct {
	*bytes.Buffer
	*multipart.Writer
}

// NewForm constructs a [Form] from the given map[string]any.
// See [AddFormField] for accepted values.
func NewForm(m map[string]any) (*Form, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for k, v := range m {
		if err := AddFormField(writer, "", k, v); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %s", err)
	}
	return &Form{
		Buffer: body,
		Writer: writer,
	}, nil
}

// MultipartFile is an interface to attach a file into the multipart form.
// [os.File] implements this interface.
type MultipartFile interface {
	Name() string
	io.ReadCloser
}

// AddFormField fills a form field with a key defined by parent[child] or child if parent is empty and writes them into [multipart.Writer].
// Accepted values: string, []byte, func() (MultipartFile, error), and map[string]any
func AddFormField(writer *multipart.Writer, parent string, child string, v any) error {
	var key string
	if parent == "" {
		key = child
	} else {
		key = fmt.Sprintf("%s[%s]", parent, child)
	}
	switch assertedValue := v.(type) {
	case string:
		if err := writer.WriteField(key, assertedValue); err != nil {
			return fmt.Errorf("write field %s: %w", key, err)
		}
	case []byte:
		subwriter, err := writer.CreateFormField(key)
		if err != nil {
			return fmt.Errorf("create form field %s: %w", key, err)
		}
		if _, err := subwriter.Write(assertedValue); err != nil {
			return fmt.Errorf("write form field: %s: %w", key, err)
		}
	case func() (MultipartFile, error):
		f, err := assertedValue()
		if err != nil {
			return fmt.Errorf("open form file: %s", err)
		}
		defer func() {
			_ = f.Close()
		}()
		filename := filepath.Base(f.Name())
		subwriter, err := writer.CreateFormFile(key, filename)
		if err != nil {
			return fmt.Errorf("create form file %s: %w", key, err)
		}
		if _, err := io.Copy(subwriter, f); err != nil {
			return fmt.Errorf("copy to form file: %w", err)
		}
	case map[string]any:
		for sub, v := range assertedValue {
			if err := AddFormField(writer, key, sub, v); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unsupported data type %s for key %s", reflect.TypeOf(v), key)
	}
	return nil
}
