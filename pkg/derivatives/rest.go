package derivatives

import (
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/krakenfx/api-go/v2/internal/helper"
	"github.com/krakenfx/api-go/v2/pkg/kraken"
)

// REST wraps [RESTBase] with functions to call common endpoints.
type REST struct {
	PublicKey  string
	PrivateKey string
	Nonce      func() string
	BaseURL    string
	Executor   kraken.ExecutorFunction
}

// REST constructs a new [REST] object with default values.
//
// For authentication, store the derivatives API key on the PublicKey and PrivateKey fields.
func NewREST() *REST {
	return &REST{
		BaseURL: "https://futures.kraken.com",
	}
}

func (r *REST) NewRequest(opts RequestOptions) (*kraken.Request, error) {
	return NewRequest(RequestOptions{
		Auth:       opts.Auth,
		PublicKey:  r.PublicKey,
		PrivateKey: r.PrivateKey,
		Nonce:      opts.Nonce,
		Method:     opts.Method,
		URL:        r.BaseURL,
		Path:       opts.Path,
		Query:      opts.Query,
		Headers:    opts.Headers,
		Body:       opts.Body,
		UserAgent:  opts.UserAgent,
		Executor:   r.Executor,
	})
}

type InstrumentResult struct {
	Instruments []Instrument `json:"instruments,omitempty"`
	DerivativesResponse
}

// Instruments retrieves the specifications of all available contract pairs.
//
// https://docs.kraken.com/api/docs/futures-api/trading/get-instruments
func (r *REST) Instruments() (*Response[InstrumentResult], error) {
	return Call[InstrumentResult](r, RequestOptions{
		Method: "GET",
		Path:   "/derivatives/api/v3/instruments",
	})
}

// InstrumentSymbol calls [REST.Instruments] and returns the first [Instrument] with matching symbol.
func (r *REST) InstrumentSymbol(s string) (*Instrument, error) {
	resp, err := r.Instruments()
	if err != nil {
		return nil, err
	}
	for _, instrument := range resp.Result.Instruments {
		if instrument.Symbol == s {
			return &instrument, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

type TickersResult struct {
	Tickers []TickerData `json:"tickers,omitempty"`
	DerivativesResponse
}

// Tickers retrieves the ticker information of all available contract pairs and indices.
//
// https://docs.kraken.com/api/docs/futures-api/trading/get-tickers
func (r *REST) Tickers() (*Response[TickersResult], error) {
	return Call[TickersResult](r, RequestOptions{
		Method: "GET",
		Path:   "/derivatives/api/v3/tickers",
	})
}

type TickersSingleResult struct {
	Data   TickerData `json:"ticker,omitempty"`
	Errors []any      `json:"errors,omitempty"`
	Error  any        `json:"error,omitempty"`
	DerivativesResponse
}

// TickerSymbol retrieves the ticker information of a specific contract pair or indice.
//
// https://docs.kraken.com/api/docs/futures-api/trading/get-ticker
func (r *REST) TickerSymbol(symbol string) (*Response[TickersSingleResult], error) {
	return Call[TickersSingleResult](r, RequestOptions{
		Method: "GET",
		Path:   []any{"/derivatives/api/v3/tickers/", symbol},
		Auth:   false,
	})
}

type OrderBookRequest struct {
	Symbol string `json:"symbol,omitempty"`
}

type OrderBookResult struct {
	OrderBook OrderBook `json:"orderBook,omitempty"`
	DerivativesResponse
}

// OrderBook retrieves the top bid and ask records of a contract pair.
//
// https://docs.kraken.com/api/docs/futures-api/trading/get-orderbook
func (r *REST) OrderBook(opts *OrderBookRequest) (*Response[OrderBookResult], error) {
	return Call[OrderBookResult](r, RequestOptions{
		Method: "GET",
		Path:   []any{"/derivatives/api/v3/orderbook"},
		Query:  opts,
	})
}

type TradeHistoryRequest struct {
	Symbol   string `json:"symbol,omitempty"`
	LastTime string `json:"lastTime,omitempty"`
}

type TradeHistoryResult struct {
	History []Trade `json:"history,omitempty"`
	DerivativesResponse
}

// TradeHistory retrieves the most recent trade events in a futures market.
//
// https://docs.kraken.com/api/docs/futures-api/trading/get-history
func (r *REST) TradeHistory(opts *TradeHistoryRequest) (*Response[TradeHistoryResult], error) {
	return Call[TradeHistoryResult](r, RequestOptions{
		Method: "GET",
		Path:   "/derivatives/api/v3/history",
		Query:  opts,
	})
}

type AccountsResult struct {
	Accounts map[string]any `json:"accounts,omitempty"`
	DerivativesResponse
}

// Accounts retrieves balances, margin requirements, margin trigger estimates, and other auxilary information of all futures cash and margin accounts.
//
// https://docs.kraken.com/api/docs/futures-api/trading/get-accounts
func (r *REST) Accounts() (*Response[AccountsResult], error) {
	return Call[AccountsResult](r, RequestOptions{
		Method: "GET",
		Path:   "/derivatives/api/v3/accounts",
		Auth:   true,
	})
}

type OrderRequest struct {
	ProcessBefore             time.Time `json:"processBefore,omitempty"`
	OrderType                 string    `json:"orderType,omitempty"`
	Symbol                    string    `json:"symbol,omitempty"`
	Side                      string    `json:"side,omitempty"`
	Size                      string    `json:"size,omitempty"`
	LimitPrice                string    `json:"limitPrice,omitempty"`
	StopPrice                 string    `json:"stopPrice,omitempty"`
	ClientOrderID             string    `json:"cliOrdId,omitempty"`
	TriggerSignal             string    `json:"triggerSignal,omitempty"`
	ReduceOnly                bool      `json:"reduceOnly,omitempty"`
	TrailingStopMaxDeviation  string    `json:"trailingStopMaxDeviation,omitempty"`
	TrailingStopDeviationUnit string    `json:"trailingStopDeviationUnit,omitempty"`
	LimitPriceOffsetValue     string    `json:"limitPriceOffsetValue,omitempty"`
	LimitPriceOffsetUnit      string    `json:"limitPriceOffsetUnit,omitempty"`
}

type SendOrderResult struct {
	SendStatus OrderStatus `json:"sendStatus,omitempty"`
	DerivativesResponse
}

// SendOrder places a new order.
//
// https://docs.kraken.com/api/docs/futures-api/trading/send-order
func (r *REST) SendOrder(opts *OrderRequest) (*Response[SendOrderResult], error) {
	return Call[SendOrderResult](r, RequestOptions{
		Method: "POST",
		Path:   "/derivatives/api/v3/sendorder",
		Body:   opts,
		Auth:   true,
	})
}

type BatchOrderRequest struct {
	ProcessBefore time.Time       `json:"processBefore,omitempty"`
	JSON          *BatchOrderJson `json:"json,omitempty" map:"stringify"`
}

type BatchOrderResult struct {
	BatchStatus []BatchStatusInfo `json:"batchStatus,omitempty"`
	DerivativesResponse
}

// BatchOrder allows placing an order, cancelling an open order, or editing an existing order in a single request.
//
// https://docs.kraken.com/api/docs/futures-api/trading/send-batch-order
func (r *REST) BatchOrder(opts *BatchOrderRequest) (*Response[BatchOrderResult], error) {
	return Call[BatchOrderResult](r, RequestOptions{
		Method: "POST",
		Path:   "/derivatives/api/v3/batchorder",
		Body:   opts,
		Auth:   true,
	})
}

type EditOrderResult struct {
	EditStatus OrderStatus
	DerivativesResponse
}

// EditOrder edits an existing order.
//
// https://docs.kraken.com/api/docs/futures-api/trading/edit-order-spring
func (r *REST) EditOrder(opts *OrderRequest) (*Response[EditOrderResult], error) {
	return Call[EditOrderResult](r, RequestOptions{
		Method: "POST",
		Path:   "/derivatives/api/v3/editorder",
		Body:   opts,
		Auth:   true,
	})
}

type CancelOrderRequest struct {
	ProcessBefore time.Time `json:"processBefore,omitempty"`
	OrderID       string    `json:"order_id,omitempty"`
	ClientOrderID string    `json:"cliOrdId,omitempty"`
}

type CancelOrderResult struct {
	CancelStatus OrderStatus `json:"cancelStatus,omitempty"`
	DerivativesResponse
}

// CancelOrder cancels an open order.
//
// https://docs.kraken.com/api/docs/futures-api/trading/cancel-order
func (r *REST) CancelOrder(opts *CancelOrderRequest) (*Response[CancelOrderResult], error) {
	return Call[CancelOrderResult](r, RequestOptions{
		Method: "POST",
		Path:   []any{"/derivatives/api/v3/cancelorder"},
		Body:   opts,
		Auth:   true,
	})
}

type CancelAllRequest struct {
	Symbol string `json:"symbol,omitempty"`
}

type CancelAllResult struct {
	CancelStatus CancelStatus `json:"cancelStatus,omitempty"`
	DerivativesResponse
}

// CancelAll cancels all open orders on a specific contract pair.
//
// https://docs.kraken.com/api/docs/futures-api/trading/cancel-all-orders
func (r *REST) CancelAll(opts *CancelAllRequest) (*Response[CancelAllResult], error) {
	return Call[CancelAllResult](r, RequestOptions{
		Method: "POST",
		Path:   "/derivatives/api/v3/cancelallorders",
		Body:   opts,
		Auth:   true,
	})
}

type OpenOrdersResult struct {
	OpenOrders []OpenOrder `json:"openOrders,omitempty"`
	DerivativesResponse
}

// OpenOrders retrieves information regarding all open orders on the futures account.
//
// https://docs.kraken.com/api/docs/futures-api/trading/get-open-orders
func (r *REST) OpenOrders() (*Response[OpenOrdersResult], error) {
	return Call[OpenOrdersResult](r, RequestOptions{
		Method: "GET",
		Path:   "/derivatives/api/v3/openorders",
		Auth:   true,
	})
}

type Requestor interface {
	NewRequest(RequestOptions) (*kraken.Request, error)
}

// Call creates a request and returns a generic response.
func Call[T any](r Requestor, opts RequestOptions) (resp *Response[T], err error) {
	resp = &Response[T]{}
	req, err := r.NewRequest(opts)
	if err != nil {
		return resp, err
	}
	resp.Http, err = req.Do()
	if err != nil {
		return resp, err
	}
	return resp, resp.Http.JSON(&resp.Result)
}

// RequestOptions contains the parameters for [NewRequest].
type RequestOptions struct {
	Auth       bool
	PublicKey  string
	PrivateKey string
	Nonce      func() string
	Method     string
	URL        string
	Path       any
	Query      any
	Headers    map[string]any
	Body       any
	UserAgent  string
	Executor   kraken.ExecutorFunction
}

// NewRequest creates a [kraken.Request] struct for submission to the Derivatives API.
//
// Authentication algorithm: https://docs.kraken.com/api/docs/guides/futures-rest
func NewRequest(opts RequestOptions) (*kraken.Request, error) {
	request, err := kraken.NewRequestWithOptions(kraken.RequestOptions{
		Method:    opts.Method,
		URL:       opts.URL,
		Headers:   opts.Headers,
		Path:      opts.Path,
		Query:     opts.Query,
		Body:      opts.Body,
		UserAgent: opts.UserAgent,
		Executor:  opts.Executor,
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
		if opts.Nonce != nil {
			if nonce == "" {
				nonce = opts.Nonce()
			}
			request.Header.Set("Nonce", nonce)
		}
		authent, err := Sign(opts.PrivateKey, data, nonce, request.URL.Path)
		if err != nil {
			return nil, fmt.Errorf("sign failed: %s", err)
		}
		request.Header.Set("APIKey", opts.PublicKey)
		request.Header.Set("Authent", authent)
	}
	return request, nil
}

// Sign hashes path, nonce, and body using the given private key.
// Returns the base64-encoded result.
func Sign(privateKey string, data io.Reader, nonce string, endpointPath string) (string, error) {
	sha256Hash := sha256.New()
	if _, err := io.Copy(sha256Hash, data); err != nil {
		return "", fmt.Errorf("copy data to hash: %w", err)
	}
	sha256Hash.Write([]byte(nonce + strings.TrimPrefix(endpointPath, "/derivatives")))
	return helper.Sign(privateKey, sha256Hash.Sum(nil))
}

type Response[T any] struct {
	Result T                `json:"result,omitempty"`
	Http   *kraken.Response `json:"-"`
}
