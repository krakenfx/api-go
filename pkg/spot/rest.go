package spot

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"reflect"

	"github.com/krakenfx/api-go/pkg/kraken"
)

// REST wraps [RESTBase] with functions to call common endpoints.
type REST struct {
	Nonce      func() string
	OTP        func() string
	PublicKey  string
	PrivateKey string
	BaseURL    string
	UserAgent  string
	Executor   kraken.ExecutorFunction
}

// CreateUserRequest contains the parameters for [REST.CreateUser].
type CreateUserRequest struct {
	*UserInfo
}

// CreateUserResult contains the result of a [REST.CreateUser] response.
type CreateUserResult struct {
	User string `json:"user,omitempty"`
}

// UpdateUserQuery contains the query parameters for [REST.UpdateUser]
type UpdateUserQuery struct {
	User string `json:"user,omitempty"`
}

// UpdateUserBody contains the body parameters for [REST.UpdateUser]
type UpdateUserBody struct {
	*UserInfo
}

// UpdateUserRequest contains the parameters for [REST.UpdateUser]
type UpdateUserRequest struct {
	UpdateUserQuery *UpdateUserQuery `json:"query,omitempty"`
	UpdateUserBody  *UpdateUserBody  `json:"body,omitempty"`
}

// GetUserRequest contains the parameters for [REST.GetUser]
type GetUserRequest struct {
	User string `json:"user,omitempty"`
}

// GetUserResult contains the result of a [REST.GetUser] response.
type GetUserResult struct {
	User       string      `json:"user,omitempty"`
	ExternalID string      `json:"external_id,omitempty"`
	UserType   string      `json:"user_type,omitempty"`
	Status     *UserStatus `json:"status,omitempty"`
	CreatedAt  string      `json:"created_at,omitempty"`
}

// VerificationRequestQuery contains the query parameters for [REST.VerifyUser].
type VerificationRequestQuery struct {
	User string `json:"user,omitempty"`
}

// VerificationRequestBody contains the body parameters for [REST.VerifyUser].
type VerificationRequestBody struct {
	Type                       string                               `json:"type,omitempty"`
	Metadata                   *VerificationMetadata                `json:"metadata,omitempty" map:"stringify"`
	SanctionsVendorResponse    string                               `json:"sanctions_vendor_response,omitempty"`
	NegativeNewsVendorResponse string                               `json:"negative_news_vendor_response,omitempty"`
	PepVendorResponse          string                               `json:"pep_vendor_response,omitempty"`
	Selfie                     func() (kraken.MultipartFile, error) `json:"selfie,omitempty"`
	VendorResponse             func() (kraken.MultipartFile, error) `json:"vendor_response,omitempty"`
	Document                   func() (kraken.MultipartFile, error) `json:"document,omitempty"`
	Front                      func() (kraken.MultipartFile, error) `json:"front,omitempty"`
	Back                       func() (kraken.MultipartFile, error) `json:"back,omitempty"`
}

// SubmitVerificationRequest contains the parameters for [REST.VerifyUser].
type SubmitVerificationRequest struct {
	VerificationRequestQuery *VerificationRequestQuery `json:"query,omitempty"`
	VerificationRequestBody  *VerificationRequestBody  `json:"body,omitempty"`
}

// SubmitVerificationResult contains the result of a [REST.VerifyUser] response.
type SubmitVerificationResult struct {
	VerificationID string `json:"verification_id,omitempty"`
}

type ServerTimeResult struct {
	UnixTime int    `json:"unixtime,omitempty"`
	RFC1123  string `json:"rfc1123,omitempty"`
}

type TradesHistoryRequest struct {
	Type             string `json:"type,omitempty"`
	Trades           bool   `json:"trades,omitempty"`
	Start            int    `json:"start,omitempty"`
	End              int    `json:"end,omitempty"`
	Ofs              int    `json:"ofs,omitempty"`
	ConsolidateTaker bool   `json:"consolidate_taker,omitempty"`
	Ledgers          bool   `json:"ledgers,omitempty"`
}

type TradesHistoryResult struct {
	Count  json.Number      `json:"count,omitempty"`
	Trades map[string]Trade `json:"trades,omitempty"`
}

type OpenOrdersRequest struct {
	Trades  bool   `json:"trades,omitempty"`
	Userref int    `json:"userref,omitempty"`
	ClOrdID string `json:"cl_ord_id,omitempty"`
}

type OpenOrdersResult struct {
	Open map[string]Order `json:"open,omitempty"`
}

type ClosedOrdersRequest struct {
	Trades           bool   `json:"trades,omitempty"`
	Userref          int    `json:"userref,omitempty"`
	ClOrdID          string `json:"cl_ord_id,omitempty"`
	Start            int    `json:"start,omitempty"`
	End              int    `json:"end,omitempty"`
	Ofs              int    `json:"ofs,omitempty"`
	CloseTime        string `json:"closeTime,omitempty"`
	ConsolidateTaker bool   `json:"consolidate_taker,omitempty"`
	WithoutCount     bool   `json:"without_count,omitempty"`
}

type ClosedOrdersResult struct {
	Closed map[string]ClosedOrder `json:"closed,omitempty"`
}

type QueryOrdersRequest struct {
	Trades           bool   `json:"trades,omitempty"`
	Userref          int    `json:"userref,omitempty"`
	TxID             string `json:"txid,omitempty"`
	ConsolidateTaker bool   `json:"consolidate_taker,omitempty"`
}

type AddOrderRequest struct {
	UserRef             int    `json:"userref,omitempty"`
	ClOrdId             string `json:"cl_ord_id,omitempty"`
	OrderType           string `json:"ordertype,omitempty"`
	Type                string `json:"type,omitempty"`
	Volume              string `json:"volume,omitempty"`
	DisplayVol          string `json:"displayvol,omitempty"`
	Pair                string `json:"pair,omitempty"`
	Price               string `json:"price,omitempty"`
	SecondaryPrice      string `json:"price2,omitempty"`
	Trigger             string `json:"trigger,omitempty"`
	Leverage            string `json:"leverage,omitempty"`
	ReduceOnly          bool   `json:"reduce_only,omitempty"`
	StpType             string `json:"stptype,omitempty"`
	OrderFlags          string `json:"oflags,omitempty"`
	TimeInForce         string `json:"timeinforce,omitempty"`
	StartTm             string `json:"starttm,omitempty"`
	ExpireTm            string `json:"expiretm,omitempty"`
	CloseOrderType      string `json:"close[ordertype],omitempty"`
	ClosePrice          string `json:"close[price],omitempty"`
	CloseSecondaryPrice string `json:"close[price2],omitempty"`
	Deadline            string `json:"deadline,omitempty"`
	Validate            bool   `json:"validate,omitempty"`
}

type AddOrderResult struct {
	OrderPlacementSingle
}

type AddBatchRequest struct {
	Orders   []*OrderRequest `json:"orders,omitempty"`
	Pair     string          `json:"pair,omitempty"`
	Deadline string          `json:"deadline,omitempty"`
	Validate bool            `json:"bool,omitempty"`
}

type AddBatchResult struct {
	Orders []OrderPlacementBatch `json:"orders,omitempty"`
}

type AmendOrderRequest struct {
	TxID            string `json:"txid,omitempty"`
	ClOrdID         string `json:"cl_ord_id,omitempty"`
	OrderQuantity   string `json:"order_qty,omitempty"`
	DisplayQuantity string `json:"display_qty,omitempty"`
	LimitPrice      string `json:"limit_price,omitempty"`
	TriggerPrice    string `json:"trigger_price,omitempty"`
	PostOnly        bool   `json:"post_only,omitempty"`
	Deadline        string `json:"deadline,omitempty"`
}

type AmendOrderResult struct {
	AmendID string `json:"amend_id,omitempty"`
}

type CancelResult struct {
	Count   int  `json:"count,omitempty"`
	Pending bool `json:"pending,omitempty"`
}

type CancelOrderRequest struct {
	TxID    any    `json:"txid,omitempty"`
	ClOrdID string `json:"cl_ord_id,omitempty"`
}

type AssetsRequest struct {
	Asset        string `json:"asset,omitempty"`
	AssetClass   string `json:"aclass,omitempty"`
	AssetVersion int    `json:"assetVersion,omitempty"`
}

type AssetPairsRequest struct {
	Pair         string `json:"pair,omitempty"`
	Info         string `json:"info,omitempty"`
	CountryCode  string `json:"country_code,omitempty"`
	AssetVersion int    `json:"assetVersion,omitempty"`
}

type TickerRequest struct {
	Pair string `json:"pair,omitempty"`
}

type OrderBookRequest struct {
	Pair  string `json:"pair,omitempty"`
	Count int    `json:"count,omitempty"`
}

type RecentTradesRequest struct {
	Pair  string `json:"pair,omitempty"`
	Since string `json:"since,omitempty"`
	Count int    `json:"count,omitempty"`
}

type OHLCRequest struct {
	Pair     string `json:"pair,omitempty"`
	Interval int    `json:"interval,omitempty"`
	Since    int    `json:"since,omitempty"`
}

type GetWebSocketsTokenResult struct {
	Token   string `json:"token,omitempty"`
	Expires int    `json:"expires,omitempty"`
}

// REST constructs a new [REST] struct with default values.
//
// For authentication, store the Spot API key on the PublicKey and PrivateKey fields.
func NewREST() *REST {
	return &REST{
		Nonce:   kraken.NewEpochCounter().Get,
		BaseURL: "https://api.kraken.com",
	}
}

// NewRequest creates a [kraken.Request] with the parameters specified in [REST].
func (r *REST) NewRequest(opts RequestOptions) (*kraken.Request, error) {
	return NewRequest(RequestOptions{
		Auth:        opts.Auth,
		Version:     opts.Version,
		PublicKey:   r.PublicKey,
		PrivateKey:  r.PrivateKey,
		Nonce:       r.Nonce,
		OTP:         r.OTP,
		Method:      opts.Method,
		BaseURL:     r.BaseURL,
		Headers:     opts.Headers,
		Path:        opts.Path,
		Query:       opts.Query,
		Body:        opts.Body,
		ContentType: opts.ContentType,
		UserAgent:   opts.UserAgent,
		Executor:    r.Executor,
	})
}

type Requestor interface {
	NewRequest(RequestOptions) (*kraken.Request, error)
}

// Call creates a request, checks for errors, and returns a generic response.
func Call[T any](r Requestor, opts RequestOptions) (resp *Response[T], err error) {
	req, err := r.NewRequest(opts)
	if err != nil {
		return resp, err
	}
	krakenResponse, err := req.Do()
	if err != nil {
		return resp, err
	}
	resp = &Response[T]{Http: krakenResponse}
	if err = resp.Http.JSON(&resp); err != nil {
		return resp, err
	} else if err = resp.GetError(); err != nil {
		return resp, err
	} else {
		return resp, nil
	}
}

// Call creates a request, checks for errors, and returns a generic response.
func (r *REST) Call(opts RequestOptions) (*Response[any], error) {
	return Call[any](r, opts)
}

// CreateUser creates a new user account in the Kraken system.
//
// https://docs.kraken.com/api/docs/embed-api/create-embed-user
func (r *REST) CreateUser(opts *CreateUserRequest) (*Response[CreateUserResult], error) {
	return Call[CreateUserResult](r, RequestOptions{
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/CreateUser",
		Body:   opts,
	})
}

// UpdateUser updates an existing user's profile.
//
// https://docs.kraken.com/api/docs/embed-api/update-embed-user
func (r *REST) UpdateUser(opts *UpdateUserRequest) (*Response[string], error) {
	return Call[string](r, RequestOptions{
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/UpdateUser",
		Body:   opts,
	})
}

// GetUser retrieves a previously created user.
//
// https://docs.kraken.com/api/docs/embed-api/get-embed-user
func (r *REST) GetUser(opts *GetUserRequest) (*Response[GetUserResult], error) {
	return Call[GetUserResult](r, RequestOptions{
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/GetUser",
		Query:  opts,
	})
}

// VerifyUser submits a verification for a user with documents and details.
//
// https://docs.kraken.com/api/docs/embed-api/submit-embed-verification
func (r *REST) VerifyUser(opts *SubmitVerificationRequest) (*Response[SubmitVerificationResult], error) {
	return Call[SubmitVerificationResult](r, RequestOptions{
		Auth:        true,
		Version:     1,
		Method:      "POST",
		Path:        "/0/private/VerifyUser",
		Query:       opts.VerificationRequestQuery,
		Body:        opts.VerificationRequestBody,
		ContentType: "multipart/form-data",
	})
}

// Balances retrieves the balances on the spot wallet.
//
// https://docs.kraken.com/api/docs/rest-api/get-account-balance
func (r *REST) Balances() (*Response[map[string]*kraken.Money], error) {
	return Call[map[string]*kraken.Money](r, RequestOptions{
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/Balance",
	})
}

// ServerTime retrieves the current server time.
//
// https://docs.kraken.com/api/docs/rest-api/get-server-time
func (r *REST) ServerTime() (*Response[ServerTimeResult], error) {
	return Call[ServerTimeResult](r, RequestOptions{
		Method: "GET",
		Path:   "/0/public/Time",
	})
}

// TradesHistory retrieves the trade events of the user.
//
// https://docs.kraken.com/api/docs/rest-api/get-trade-history
func (r *REST) TradesHistory(opts *TradesHistoryRequest) (*Response[TradesHistoryResult], error) {
	return Call[TradesHistoryResult](r, RequestOptions{
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/TradesHistory",
		Body:   opts,
	})
}

// OpenOrders retrieves information about currently open orders.
//
// https://docs.kraken.com/api/docs/rest-api/get-open-orders
func (r *REST) OpenOrders(opts *OpenOrdersRequest) (*Response[OpenOrdersResult], error) {
	return Call[OpenOrdersResult](r, RequestOptions{
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/OpenOrders",
		Body:   opts,
	})
}

// ClosedOrders retrieves information about orders that have been closed.
//
// https://docs.kraken.com/api/docs/rest-api/get-closed-orders
func (r *REST) ClosedOrders(opts *ClosedOrdersRequest) (*Response[ClosedOrdersResult], error) {
	return Call[ClosedOrdersResult](r, RequestOptions{
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/ClosedOrders",
		Body:   opts,
	})
}

// QueryOrders retrieves information about specific orders.
//
// https://docs.kraken.com/api/docs/rest-api/get-orders-info
func (r *REST) QueryOrders(opts *QueryOrdersRequest) (*Response[map[string]ClosedOrder], error) {
	return Call[map[string]ClosedOrder](r, RequestOptions{
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/QueryOrders",
		Body:   opts,
	})
}

// AddOrder places a new order.
//
// https://docs.kraken.com/api/docs/rest-api/add-order
func (r *REST) AddOrder(opts *AddOrderRequest) (*Response[AddOrderResult], error) {
	return Call[AddOrderResult](r, RequestOptions{
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/AddOrder",
		Body:   opts,
	})
}

// AddBatch places a collection of orders.
//
// https://docs.kraken.com/api/docs/rest-api/add-order-batch
func (r *REST) AddBatch(opts *AddBatchRequest) (*Response[AddBatchResult], error) {
	return Call[AddBatchResult](r, RequestOptions{
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/AddOrderBatch",
		Body:   opts,
	})
}

// AmendOrder changes the properties of an open order.
//
// https://docs.kraken.com/api/docs/rest-api/amend-order
func (r *REST) AmendOrder(opts *AmendOrderRequest) (*Response[AmendOrderResult], error) {
	return Call[AmendOrderResult](r, RequestOptions{
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/AmendOrder",
		Body:   opts,
	})
}

// CancelOrder cancels an open order.
//
// https://docs.kraken.com/api/docs/rest-api/cancel-order
func (r *REST) CancelOrder(opts *CancelOrderRequest) (*Response[CancelResult], error) {
	return Call[CancelResult](r, RequestOptions{
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/CancelOrder",
		Body:   opts,
	})
}

// CancelAll cancels all open orders.
//
// https://docs.kraken.com/api/docs/rest-api/cancel-all-orders
func (r *REST) CancelAll() (*Response[CancelResult], error) {
	return Call[CancelResult](r, RequestOptions{
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/CancelAll",
	})
}

// Assets retrieves information about available assets on spot.
//
// https://docs.kraken.com/api/docs/rest-api/get-asset-info
func (r *REST) Assets(opts *AssetsRequest) (*Response[map[string]AssetInfo], error) {
	return Call[map[string]AssetInfo](r, RequestOptions{
		Method: "GET",
		Path:   "/0/public/Assets",
		Query:  opts,
	})
}

// AssetPairs retrieves information about tradeable asset pairs on spot.
//
// https://docs.kraken.com/api/docs/rest-api/get-tradable-asset-pairs
func (r *REST) AssetPairs(opts *AssetPairsRequest) (*Response[map[string]AssetPair], error) {
	return Call[map[string]AssetPair](r, RequestOptions{
		Method: "GET",
		Path:   "/0/public/AssetPairs",
		Query:  opts,
	})
}

// Ticker retrieves information of all spot markets, or a specific market if `pair` is specified.
//
// https://docs.kraken.com/api/docs/rest-api/get-ticker-information
func (r *REST) Ticker(opts *TickerRequest) (*Response[map[string]AssetTickerInfo], error) {
	return Call[map[string]AssetTickerInfo](r, RequestOptions{
		Method: "GET",
		Path:   "/0/public/Ticker",
		Query:  opts,
	})
}

// OrderBook retrieves the bid and ask records of a specific spot market.
//
// https://docs.kraken.com/api/docs/rest-api/get-order-book
func (r *REST) OrderBook(opts *OrderBookRequest) (*Response[map[string]OrderBook], error) {
	return Call[map[string]OrderBook](r, RequestOptions{
		Method: "GET",
		Path:   "/0/public/Depth",
		Query:  opts,
	})
}

// RecentTrades retrieves the recent trade records of a specified spot market.
//
// https://docs.kraken.com/api/docs/rest-api/get-recent-trades
func (r *REST) RecentTrades(opts *RecentTradesRequest) (*Response[map[string]any], error) {
	return Call[map[string]any](r, RequestOptions{
		Method: "GET",
		Path:   "/0/public/Trades",
		Query:  opts,
	})
}

// OHLC retrieves recent open, high, low, close, and volume records of a specified spot market.
//
// https://docs.kraken.com/api/docs/rest-api/get-ohlc-data
func (r *REST) OHLC(opts *OHLCRequest) (*Response[map[string]any], error) {
	return Call[map[string]any](r, RequestOptions{
		Method: "GET",
		Path:   "/0/public/OHLC",
		Query:  opts,
	})
}

// GetWebSocketsToken generates an authentication token for WebSocket, which must be used within 15 minutes of creation to prevent expiration.
//
// https://docs.kraken.com/api/docs/rest-api/get-websockets-token
func (r *REST) GetWebSocketsToken() (*Response[GetWebSocketsTokenResult], error) {
	return Call[GetWebSocketsTokenResult](r, RequestOptions{
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/GetWebSocketsToken",
	})
}

// TickerSingle calls [REST.Ticker] and returns the first [AssetTickerInfo] item from the result.
func (r *REST) TickerSingle(pair string) (*AssetTickerInfo, error) {
	ticker, err := r.Ticker(&TickerRequest{
		Pair: pair,
	})
	if err != nil {
		return nil, err
	}
	for _, v := range ticker.Result {
		return &v, nil
	}
	return nil, fmt.Errorf("not found")
}

// RequestOptions contain the parameters for [NewRequest].
type RequestOptions struct {
	Auth        bool
	Version     int
	PublicKey   string
	PrivateKey  string
	Nonce       func() string
	OTP         func() string
	Method      string
	BaseURL     string
	Headers     map[string]any
	Path        any
	Query       any
	Body        any
	ContentType string
	UserAgent   string
	Executor    kraken.ExecutorFunction
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
func NewRequest(opts RequestOptions) (req *kraken.Request, err error) {
	body := make(map[string]any)
	if opts.Body != nil {
		bodyReflection := kraken.GetDirectReflection(opts.Body)
		if bodyReflection.Type.Kind() == reflect.Struct {
			body, err = kraken.StructToMap(opts.Body)
			if err != nil {
				return req, fmt.Errorf("body struct to map: %w", err)
			}
		} else if bodyReflection.Type.Kind() == reflect.Map && !bodyReflection.Value.IsZero() {
			src, ok := opts.Body.(map[string]any)
			if !ok {
				return req, fmt.Errorf("unsupported body type of %s", reflect.TypeOf(body).Name())
			}
			maps.Copy(body, src)
		}
	}
	var nonce string
	if opts.Auth && opts.Version == 0 {
		if nonceAny, ok := body["nonce"]; ok {
			nonce = fmt.Sprint(nonceAny)
		} else if !ok && opts.Nonce != nil {
			nonce = opts.Nonce()
			body["nonce"] = nonce
		}
		if _, ok := body["otp"]; !ok && opts.OTP != nil {
			body["otp"] = opts.OTP()
		}
	}
	contentType := opts.ContentType
	if contentType == "" {
		contentType = "application/json"
	}
	req, err = kraken.NewRequestWithOptions(kraken.RequestOptions{
		Method:      opts.Method,
		URL:         opts.BaseURL,
		Headers:     opts.Headers,
		Path:        opts.Path,
		Query:       opts.Query,
		Body:        body,
		ContentType: contentType,
		UserAgent:   opts.UserAgent,
		Executor:    opts.Executor,
	})
	if err != nil {
		return
	}
	if opts.Auth {
		if opts.Version == 1 {
			nonce = req.Header.Get("API-Nonce")
			if nonce == "" {
				nonce = opts.Nonce()
				req.SetHeader("API-Nonce", nonce)
			}
			otp := req.Header.Get("API-OTP")
			if otp == "" && opts.OTP != nil {
				otp := opts.OTP()
				req.SetHeader("API-OTP", otp)
			}
		}
		var bodyReader io.ReadCloser
		bodyReader, err := req.GetBody()
		if err != nil {
			return req, fmt.Errorf("get body: %s", err)
		}
		defer func() {
			_ = bodyReader.Close()
		}()
		signature, err := Sign(opts.PrivateKey, req.URL.RequestURI(), fmt.Sprint(nonce), bodyReader)
		if err != nil {
			return req, fmt.Errorf("sign: %s", err)
		}
		req.SetHeader("API-Key", opts.PublicKey)
		req.SetHeader("API-Sign", signature)
	}
	return
}

// Response wraps [kraken.Response] with fields expected of that of the Spot API
type Response[T any] struct {
	Error  []any            `json:"error,omitempty"`
	Result T                `json:"result,omitempty"`
	Http   *kraken.Response `json:"-"`
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
	return err
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
