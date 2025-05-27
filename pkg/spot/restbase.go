package spot

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"maps"
	"sync"

	"github.com/krakenfx/api-go/pkg/kraken"
)

// RESTBase wraps [kraken.REST] with methods for creating authenticated requests in the Spot API.
type RESTBase struct {
	PublicKey  string
	PrivateKey string
	Nonce      func() string
	OTP        func() string
	Ordered    bool
	requestMux sync.Mutex
	*kraken.REST
}

// NewRESTBase construct a new [RESTBase] struct with default values.
func NewRESTBase() *RESTBase {
	client := &RESTBase{
		Nonce:   kraken.NewEpochCounter().Get,
		Ordered: true,
		REST:    kraken.NewREST(),
	}
	client.BaseURL = "https://api.kraken.com"
	client.DefaultContentType = "application/json"
	return client
}

// RequestOptions contain the parameters for [RESTBase.NewRequest].
type RequestOptions struct {
	Auth    bool
	Method  string
	Path    any
	Query   map[string]any
	Headers map[string]any
	Body    map[string]any
	Version int
}

// NewRequest creates a [kraken.Request] struct for submission to the Spot API.
//
// The placement of Nonce and OTP is determined by the Version option:
//
// - [0] sets the nonce and otp in the body.
//
// - [1] sets the nonce and otp in the header.
//
// Authentication algorithm: https://docs.kraken.com/api/docs/guides/spot-rest-auth
func (b *RESTBase) NewRequest(opts *RequestOptions) (*kraken.Request, error) {
	body := make(map[string]any)
	if len(opts.Body) > 0 {
		maps.Copy(body, opts.Body)
	}
	var nonce any
	if opts.Auth && opts.Version == 0 {
		var ok bool
		if nonce, ok = body["nonce"]; !ok {
			nonce = b.Nonce()
			body["nonce"] = nonce
		}
		if _, ok := body["otp"]; !ok && b.OTP != nil {
			body["otp"] = b.OTP()
		}
	}
	request, err := b.REST.NewRequest(&kraken.NewRequestOptions{
		Method:  opts.Method,
		Path:    opts.Path,
		Query:   opts.Query,
		Headers: opts.Headers,
		Body:    body,
	})
	if err != nil {
		return request, err
	}
	if opts.Auth {
		if opts.Version == 1 {
			nonceString := request.Header.Get("API-Nonce")
			otp := request.Header.Get("API-OTP")
			if len(nonceString) == 0 {
				nonceString = b.Nonce()
				request.SetHeader("API-Nonce", nonceString)
			}
			if len(otp) == 0 && b.OTP != nil {
				otp := b.OTP()
				request.SetHeader("API-OTP", otp)
			}
			nonce = nonceString
		}
		var bodyReader io.ReadCloser
		bodyReader, err := request.GetBody()
		if err != nil {
			return request, fmt.Errorf("get body: %s", err)
		}
		defer func() {
			_ = bodyReader.Close()
		}()
		signature, err := Sign(b.PrivateKey, request.URL.RequestURI(), fmt.Sprint(nonce), bodyReader)
		if err != nil {
			return request, fmt.Errorf("sign: %s", err)
		}
		request.SetHeader("API-Key", b.PublicKey)
		request.SetHeader("API-Sign", signature)
	}
	return request, nil
}

// Issue is a helper method to call [RESTBase.NewRequest] and [RESTBase.Do] in one line and check for API errors.
//
// If Ordered is true, then the mutex is locked at the beginning and unlocked at end of function.
// This ensures that the nonce is ordered even in a concurrent setting.
func (b *RESTBase) Issue(opts *RequestOptions) (*kraken.Response, error) {
	if b.Ordered {
		b.requestMux.Lock()
		defer b.requestMux.Unlock()
	}
	request, err := b.NewRequest(opts)
	if err != nil {
		return nil, err
	}
	response, err := b.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Sign hashes path, nonce, and body using the given private key.
// Returns the base64-encoded result.
func Sign(privateKey string, path string, nonce string, body io.ReadCloser) (string, error) {
	sha256Hash := sha256.New()
	sha256Hash.Write([]byte((nonce)))
	if body != nil {
		if _, err := io.Copy(sha256Hash, body); err != nil {
			return "", fmt.Errorf("copy body to hash: %w", err)
		}
	}
	message := path + string(sha256Hash.Sum(nil))
	return kraken.Sign(privateKey, []byte(message))
}

// Response wraps [kraken.Response] with fields expected of that of the Spot API
type Response[T any] struct {
	Error            []any `json:"error,omitempty"`
	Result           T     `json:"result,omitempty"`
	*kraken.Response `json:"-"`
}

// GetError returns the API error message if it exists on the body.
func (r *Response[T]) GetError() error {
	if len(r.Error) == 0 {
		return nil
	}
	var err error
	for _, errorEntry := range r.Error {
		err = errors.Join(err, fmt.Errorf("%v", errorEntry))
	}
	if traceID := r.Header.Get("x-trace-id"); len(traceID) > 0 {
		err = errors.Join(err, fmt.Errorf("trace id \"%s\"", traceID))
	}
	return err
}

// Whether the request parameters should be in the body or URL query.
type ParamMode uint8

const (
	BodyMode ParamMode = iota
	QueryMode
)

// Issuer refers to a struct with a Issue function.
// it accepts [RequestOptions] as a parameter and returning a [kraken.Response].
type Issuer interface {
	Issue(opts *RequestOptions) (*kraken.Response, error)
}

// APIFunctionOptions contains the parameters for [NewAPIFunctionWithNoParams] and [NewAPIFunction].
type APIFunctionOptions struct {
	REST       Issuer
	Auth       bool
	Method     string
	Path       any
	Headers    map[string]any
	Version    int
	BodyField  bool
	QueryField bool
	ParamMode  ParamMode
}

// NewAPIFunctionWithNoParams returns a request function with no parameters.
func NewAPIFunctionWithNoParams[R any](setup *APIFunctionOptions) func() (*Response[R], error) {
	f := NewAPIFunction[any, R](setup)
	return func() (*Response[R], error) {
		return f(nil)
	}
}

// NewAPIFunction returns a request function for an endpoint.
func NewAPIFunction[P any, R any](setup *APIFunctionOptions) func(opts *P) (*Response[R], error) {
	return func(opts *P) (*Response[R], error) {
		wrappedResponse := &Response[R]{}
		var query map[string]any
		var body map[string]any
		if opts != nil {
			optsMap, err := kraken.StructToMap(opts)
			if err != nil {
				return wrappedResponse, fmt.Errorf("options: %w", err)
			}
			if setup.BodyField {
				body, _ = optsMap["body"].(map[string]any)
			}
			if setup.QueryField {
				query, _ = optsMap["query"].(map[string]any)
			}
			if setup.ParamMode == BodyMode && !setup.BodyField {
				if body == nil {
					body = make(map[string]any)
				}
				maps.Copy(body, optsMap)
			}
			if setup.ParamMode == QueryMode && !setup.QueryField {
				if query == nil {
					query = make(map[string]any)
				}
				maps.Copy(query, optsMap)
			}
		}
		response, err := setup.REST.Issue(&RequestOptions{
			Auth:    setup.Auth,
			Method:  setup.Method,
			Path:    setup.Path,
			Headers: setup.Headers,
			Query:   query,
			Body:    body,
			Version: setup.Version,
		})
		wrappedResponse.Response = response
		if err != nil {
			return wrappedResponse, err
		}
		if err := response.JSON(&wrappedResponse); err != nil {
			return wrappedResponse, err
		}
		return wrappedResponse, wrappedResponse.GetError()
	}
}
