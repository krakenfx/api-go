package derivatives

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/krakenfx/api-go/pkg/kraken"
)

// RESTBase wraps [kraken.REST] with methods for creating authenticated requests in the Derivatives API.
type RESTBase struct {
	PublicKey  string
	PrivateKey string
	Nonce      func() string
	*kraken.REST
}

// NewRESTBase construct a new [RESTBase] object with default values.
func NewRESTBase() *RESTBase {
	restBase := &RESTBase{
		REST: kraken.NewREST(),
	}
	restBase.BaseURL = "https://futures.kraken.com"
	restBase.DefaultContentType = "application/x-www-form-urlencoded"
	return restBase
}

// RequestOptions contain the parameters for [RESTBase.NewRequest].
type RequestOptions struct {
	Auth    bool
	Method  string
	Path    any
	Query   map[string]any
	Headers map[string]any
	Body    map[string]any
}

// NewRequest creates a [kraken.Request] struct for submission to the Derivatives API.
//
// Authentication algorithm: https://docs.kraken.com/api/docs/guides/futures-rest
func (b *RESTBase) NewRequest(opts *RequestOptions) (*kraken.Request, error) {
	request, err := b.REST.NewRequest(&kraken.NewRequestOptions{
		Method:  opts.Method,
		Path:    opts.Path,
		Query:   opts.Query,
		Headers: opts.Headers,
		Body:    opts.Body,
	})
	if err != nil {
		return nil, err
	}
	if opts.Auth {
		var data io.Reader
		if request.Method == "POST" {
			bodyReader, err := request.GetBody()
			if err != nil {
				return nil, fmt.Errorf("get body: %s", err)
			}
			data = bodyReader
		} else {
			data = strings.NewReader(request.URL.Query().Encode())
		}
		nonce := request.Header.Get("Nonce")
		if b.Nonce != nil {
			if nonce == "" {
				nonce = b.Nonce()
			}
			request.SetHeader("Nonce", nonce)
		}
		authent, err := Sign(b.PrivateKey, data, nonce, request.URL.Path)
		if err != nil {
			return nil, fmt.Errorf("sign failed: %s", err)
		}
		request.SetHeader("APIKey", b.PublicKey)
		request.SetHeader("Authent", authent)
	}
	return request, nil
}

// Issue is a helper method to call [RESTBase.NewRequest] and [RESTBase.Do] in one line and check for API errors.
func (b *RESTBase) Issue(opts *RequestOptions) (*kraken.Response, error) {
	request, err := b.NewRequest(opts)
	if err != nil {
		return nil, err
	}
	response, err := b.Do(request)
	if err != nil {
		return response, err
	}
	responseMap, err := response.Map()
	if err != nil {
		return response, err
	}
	if errs, traverseErr := kraken.Traverse[[]any](responseMap, "errors"); traverseErr == nil {
		var err error
		for _, errorEntry := range *errs {
			errors.Join(err, fmt.Errorf("%v", errorEntry))
		}
		return response, err
	}
	if err, traverseErr := kraken.Traverse[string](responseMap, "error"); traverseErr == nil {
		return response, errors.New(*err)
	}
	return response, nil
}

// Sign hashes path, nonce, and body using the given private key.
// Returns the base64-encoded result.
func Sign(privateKey string, data io.Reader, nonce string, endpointPath string) (string, error) {
	sha256Hash := sha256.New()
	io.Copy(sha256Hash, data)
	sha256Hash.Write([]byte(nonce + strings.TrimPrefix(endpointPath, "/derivatives")))
	return kraken.Sign(privateKey, sha256Hash.Sum(nil))
}
