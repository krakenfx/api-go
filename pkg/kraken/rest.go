package kraken

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"reflect"
	"strings"
)

// REST API wrapper for the HTTP client.
type REST struct {
	BaseURL            string
	DefaultContentType string
	DefaultUserAgent   string
	*http.Client
}

// NewREST constructs a new [REST] struct with default values.
func NewREST() *REST {
	return &REST{
		DefaultContentType: "application/x-www-form-urlencoded",
		DefaultUserAgent:   "krakenfx/api-go",
		Client:             http.DefaultClient,
	}
}

// NewRequestOptions contains the parameters for [REST.NewRequest].
type NewRequestOptions struct {
	Method     string
	Path       any
	PathValues map[string]any
	Query      map[string]any
	Headers    map[string]any
	Body       any
}

// NewRequest creates a new [Request].
func (r *REST) NewRequest(opts *NewRequestOptions) (*Request, error) {
	request := NewRequest()
	request.Method = opts.Method
	switch path := opts.Path.(type) {
	case []any:
		if err := request.SetURL(r.BaseURL, path...); err != nil {
			return nil, fmt.Errorf("set url: %w", err)
		}
	default:
		if err := request.SetURL(r.BaseURL, fmt.Sprint(path)); err != nil {
			return nil, fmt.Errorf("set url: %w", err)
		}
	}
	request.SetQuery(opts.Query)
	request.SetHeader("Content-Type", r.DefaultContentType)
	request.SetHeader("User-Agent", r.DefaultUserAgent)
	request.SetHeaders(opts.Headers)
	if opts.Body != nil {
		if err := request.SetBody(opts.Body); err != nil {
			return nil, err
		}
	}
	return request, nil
}

// Do submits a [Request] struct and returns a [Response].
func (r *REST) Do(req *Request) (*Response, error) {
	response, err := r.Client.Do(req.Request)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("io read all failed: %s", err)
	}
	responseWrapped := &Response{
		Request:  req,
		Body:     body,
		Response: response,
	}
	return responseWrapped, nil
}

// Request helps to call [REST.NewRequest] and [REST.Do] in one line.
func (r *REST) Request(cfg *NewRequestOptions) (*Response, error) {
	request, err := r.NewRequest(cfg)
	if err != nil {
		return nil, err
	}
	return r.Do(request)
}

// Request is a wrapper around [http.Request] to assist with internal functions.
type Request struct {
	Body          []byte `json:"body,omitempty"`
	*http.Request `json:"-"`
}

// NewRequest initializes a [Request] object with default values.
func NewRequest() *Request {
	return &Request{
		Request: &http.Request{
			Method:     "GET",
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header: http.Header{
				"Content-Type": []string{"application/x-www-form-urlencoded"},
			},
			Body:    http.NoBody,
			GetBody: func() (io.ReadCloser, error) { return http.NoBody, nil },
		},
	}
}

// GetMediaType retrieves the Content-Type header without additional parameters.
func (r *Request) GetMediaType() string {
	contentType := r.Header.Get("Content-Type")
	mediaType, _, _ := strings.Cut(contentType, ";")
	return strings.ToLower(mediaType)
}

// SetURL creates a [url.URL] from the given base and path parameters and sets it as the URL.
func (r *Request) SetURL(base string, s ...any) error {
	u, err := url.Parse(base)
	if err != nil {
		return fmt.Errorf("url parse \"%s\": %w", base, err)
	}
	if u.Path == "" {
		u.Path = "/"
	}
	for _, element := range s {
		u = u.JoinPath(fmt.Sprint(element))
	}
	r.URL = u
	r.Host = r.URL.Host
	return nil
}

// JoinPath adds a path parameter to the request URL.
func (r *Request) JoinPath(p any) {
	r.URL = r.URL.JoinPath(fmt.Sprint(p))
}

// SetQuery converts a map[string]any into a [url.Values] object and sets it as the URL query.
func (r *Request) SetQuery(q map[string]any) {
	if len(q) == 0 {
		return
	}
	query := r.URL.Query()
	for k, v := range q {
		switch v := v.(type) {
		case string:
			query[k] = []string{v}
		case []string:
			query[k] = v
		default:
			s, _ := json.Marshal(v)
			query[k] = []string{(string(s))}
		}
	}
	r.URL.RawQuery = query.Encode()
}

// SetHeader sets the value of a header field.
func (r *Request) SetHeader(key string, value any) {
	switch value := value.(type) {
	case string:
		r.Header.Set(key, value)
	case []string:
		for _, item := range value {
			r.Header.Add(key, item)
		}
	default:
		s, _ := json.Marshal(value)
		r.Header.Set(key, string(s))
	}
}

// SetHeaders ranges over the hash map and calls [Request.SetHeader].
func (r *Request) SetHeaders(h map[string]any) {
	for k, v := range h {
		r.SetHeader(k, v)
	}
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

// CreateReadCloser constructs an [io.ReadCloser] from the a slice of byte characters.
func CreateReadCloser(b []byte) io.ReadCloser {
	return io.NopCloser(bytes.NewReader(b))
}

// MultipartForm is a combined representation of [bytes.Buffer] and [multipart.Writer], both essential for setting the body.
type MultipartForm struct {
	*bytes.Buffer
	*multipart.Writer
}

// CreateMultipartForm constructs a [MultipartForm] from the given map[string]any.
// See [AddFormField] for accepted values.
func CreateMultipartForm(m map[string]any) (*MultipartForm, error) {
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
	return &MultipartForm{
		Buffer: body,
		Writer: writer,
	}, nil
}

// SetBody sets the request body based on the current Content-Type of the Request.
//
// It automatically encodes the provided value `v` according to the media type:
//
// - "multipart/form-data": expects v to be a map[string]any and sets it as a multipart form.
//
// - "application/x-www-form-urlencoded": expects v to be a map[string]any and encodes it as form values.
//
// - "application/json": marshals v as JSON.
//
// Returns an error if the Content-Type is unsupported or if the input type doesn't match the expected format.
func (r *Request) SetBody(v any) error {
	mediaType := r.GetMediaType()
	var body []byte
	switch mediaType {
	case "multipart/form-data":
		m, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("b must be map[string]any")
		}
		form, err := CreateMultipartForm(m)
		if err != nil {
			return err
		}
		body = form.Bytes()
		r.SetHeader("Content-Type", form.FormDataContentType())
	case "application/x-www-form-urlencoded":
		values := make(url.Values)
		m, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("b must be map[string]any")
		}
		for k, v := range m {
			switch v := v.(type) {
			case string:
				values[k] = []string{v}
			case []string:
				values[k] = v
			default:
				s, _ := json.Marshal(v)
				values[k] = []string{string(s)}
			}
		}
		body = []byte(values.Encode())
	case "application/json":
		body, _ = json.Marshal(v)
	default:
		if mediaType == "" {
			return fmt.Errorf("unspecified content type")
		} else {
			return fmt.Errorf("content type \"%s\" not supported", mediaType)
		}
	}
	r.Body = body
	r.ContentLength = int64(len(r.Body))
	r.Request.Body = CreateReadCloser(r.Body)
	r.GetBody = func() (io.ReadCloser, error) { return CreateReadCloser(r.Body), nil }
	return nil
}

// Response is a wrapper around [http.Response] with an already read body.
type Response struct {
	Request        *Request       `json:"request,omitempty"`
	Body           []byte         `json:"body,omitempty"`
	BodyMap        map[string]any `json:"-"`
	*http.Response `json:"-"`
}

// Map decodes the body into map[string]any.
func (r *Response) Map() (map[string]any, error) {
	if r.BodyMap != nil {
		return r.BodyMap, nil
	}
	var b map[string]any
	if err := r.JSON(&b); err != nil {
		return nil, err
	}
	r.BodyMap = b
	return b, nil
}

// JSON decodes the body and stores it into the value pointed by v.
func (r *Response) JSON(v any) error {
	decoder := json.NewDecoder(bytes.NewReader(r.Body))
	decoder.UseNumber()
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("json decode \"%s\": %w", string(r.Body), err)
	}
	return nil
}
