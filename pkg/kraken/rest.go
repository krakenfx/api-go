package kraken

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

// Request is a wrapper around [http.Request] to assist with internal functions.
type Request struct {
	Executor      ExecutorFunction `json:"-"`
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
				"User-Agent":   []string{"krakenfx/api-go"},
			},
			Body:    http.NoBody,
			GetBody: func() (io.ReadCloser, error) { return http.NoBody, nil },
		},
		Executor: http.DefaultClient.Do,
	}
}

// RequestOptions contains the parameters for [NewRequestWithOptions].
type RequestOptions struct {
	Method      string
	URL         string
	Headers     map[string]any
	Path        any
	Query       any
	Body        any
	ContentType string
	UserAgent   string
	Executor    ExecutorFunction
}

// Whether the request parameters should be in the body or URL query.
type ParamMode uint8

const (
	BodyMode ParamMode = iota
	QueryMode
)

// NewRequestWithOptions constructs a new [Request] with [RequestOptions].
func NewRequestWithOptions(opts RequestOptions) (request *Request, err error) {
	request = NewRequest()
	if err := request.SetURL(opts.URL); err != nil {
		return request, fmt.Errorf("set URL: %w", err)
	}
	if opts.ContentType != "" {
		request.Header.Set("Content-Type", opts.ContentType)
	}
	if opts.UserAgent != "" {
		request.Header.Set("User-Agent", opts.UserAgent)
	}
	if err := request.SetHeaders(opts.Headers); err != nil {
		return request, fmt.Errorf("set headers: %w", err)
	}
	if opts.Method != "" {
		request.Method = opts.Method
	}
	if opts.Path != nil {
		if err := request.SetPath(opts.Path); err != nil {
			return request, fmt.Errorf("set path: %w", err)
		}
	}
	if opts.Query != nil {
		if err := request.SetQuery(opts.Query); err != nil {
			return request, fmt.Errorf("set query: %w", err)
		}
	}
	if opts.Body != nil {
		if err := request.SetBody(opts.Body); err != nil {
			return request, fmt.Errorf("set body: %w", err)
		}
	}
	if opts.Executor != nil {
		request.Executor = opts.Executor
	}
	return request, nil
}

func MustNewRequestWithOptions(opts RequestOptions) *Request {
	return Must(NewRequestWithOptions(opts))
}

// ExecutorFunction takes a [http.Request] and returns a [http.Response].
type ExecutorFunction func(request *http.Request) (*http.Response, error)

// Do submits the request and returns a [Response].
func (r *Request) Do() (resp *Response, err error) {
	resp = &Response{Request: r}
	response, err := r.Executor(r.Request)
	resp.Response = response
	if err != nil {
		return resp, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	body, err := io.ReadAll(response.Body)
	resp.Body = body
	if err != nil {
		return resp, fmt.Errorf("io read all: %w", err)
	}
	result := &Response{
		Request:  r,
		Body:     body,
		Response: response,
	}
	return result, nil
}

// Do submits the request and returns a [Response]. Panics on error.
func (r *Request) MustDo() *Response {
	return Must(r.Do())
}

// GetMediaType retrieves the Content-Type header without additional parameters.
func (r *Request) GetMediaType() string {
	contentType := r.Header.Get("Content-Type")
	mediaType, _, _ := strings.Cut(contentType, ";")
	return strings.ToLower(mediaType)
}

// SetURL creates a [url.URL] from the given base and path parameters and sets it as the URL.
func (r *Request) SetURL(base string) error {
	if base == "" {
		return fmt.Errorf("base is empty")
	}
	u, err := url.Parse(base)
	if err != nil {
		return fmt.Errorf("url parse \"%s\": %w", base, err)
	}
	r.URL = u
	r.Host = r.URL.Host
	return nil
}

// SetPath sets the URL path to p.
func (r *Request) SetPath(p any) error {
	path, err := StringSlice(p)
	if err != nil {
		return err
	}
	for _, item := range path {
		if err := r.JoinPath(item); err != nil {
			return err
		}
	}
	return nil
}

// JoinPath adds a path parameter to the request URL.
func (r *Request) JoinPath(p string) error {
	if r.URL == nil {
		return fmt.Errorf("request URL not initialized")
	}
	if r.URL.Path == "" {
		r.URL.Path = "/"
	}
	r.URL = r.URL.JoinPath(p)
	return nil
}

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

// SetQuery converts a map[string]any into a [url.Values] object and sets it as the URL query.
func (r *Request) SetQuery(q any) error {
	query, err := ToURLValues(q)
	if err != nil {
		return err
	}
	r.URL.RawQuery = query.Encode()
	return nil
}

// SetHeader sets the value of a header field.
func (r *Request) SetHeader(key string, value any) (err error) {
	r.Header[key], err = StringSlice(value)
	return
}

// SetHeaders ranges over the hash map and calls [Request.SetHeader].
func (r *Request) SetHeaders(h map[string]any) error {
	for k, v := range h {
		if err := r.SetHeader(k, v); err != nil {
			return fmt.Errorf("set header %s: %w", k, err)
		}
	}
	return nil
}

// SetBody sets the request body based on the media type.
func (r *Request) SetBody(v any) error {
	result, err := Marshal(MarshalOptions{
		MediaType: r.GetMediaType(),
		Object:    v,
	})
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", result.ContentType)
	r.ContentLength = int64(len(result.Data))
	r.Body = CreateReadCloser(result.Data)
	r.GetBody = func() (io.ReadCloser, error) { return CreateReadCloser(result.Data), nil }
	return nil
}

func (r *Request) MustGetBody() io.ReadCloser {
	return Must(r.GetBody())
}

func (r *Request) GetBodyBytes() ([]byte, error) {
	body, err := r.GetBody()
	if err != nil {
		return nil, err
	}
	return io.ReadAll(body)
}

func (r *Request) MustGetBodyBytes() []byte {
	return Must(r.GetBodyBytes())
}

// Response is a wrapper around [http.Response] with a read body.
type Response struct {
	Request        *Request `json:"-,omitempty"`
	Body           []byte   `json:"body,omitempty"`
	*http.Response `json:"-"`
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
