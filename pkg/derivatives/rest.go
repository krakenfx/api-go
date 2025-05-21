package derivatives

import (
	"fmt"
	"time"

	"github.com/krakenfx/api-go/pkg/kraken"
)

// REST wraps [RESTBase] with functions to call common endpoints.
type REST struct {
	*RESTBase
}

// REST constructs a new [REST] object with default values.
//
// For authentication, store the derivatives API key on the PublicKey and PrivateKey fields.
func NewREST() *REST {
	rest := &REST{
		RESTBase: NewRESTBase(),
	}
	return rest
}

type InstrumentResponse struct {
	Instruments []Instrument `json:"instruments,omitempty"`
	DerivativesResponse
}

// Instruments retrieves the specifications of all available contract pairs.
//
// https://docs.kraken.com/api/docs/futures-api/trading/get-instruments
func (r *REST) Instruments() (*InstrumentResponse, error) {
	wrappedResponse := &InstrumentResponse{}
	response, err := r.Issue(&RequestOptions{
		Method: "GET",
		Path:   "/derivatives/api/v3/instruments",
	})
	wrappedResponse.Response = response
	if err != nil {
		return nil, err
	}
	return wrappedResponse, response.JSON(&wrappedResponse)
}

// InstrumentSymbol calls [REST.Instruments] and returns the first [Instrument] with matching symbol.
func (r *REST) InstrumentSymbol(s string) (*Instrument, error) {
	resp, err := r.Instruments()
	if err != nil {
		return nil, err
	}
	for _, instrument := range resp.Instruments {
		if instrument.Symbol == s {
			return &instrument, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

type TickersResponse struct {
	Tickers []TickerData `json:"tickers,omitempty"`
	DerivativesResponse
}

// Tickers retrieves the ticker information of all available contract pairs and indices.
//
// https://docs.kraken.com/api/docs/futures-api/trading/get-tickers
func (r *REST) Tickers() (*TickersResponse, error) {
	wrappedResponse := &TickersResponse{}
	resp, err := r.Issue(&RequestOptions{
		Method: "GET",
		Path:   "/derivatives/api/v3/tickers",
	})
	wrappedResponse.Response = resp
	if err != nil {
		return wrappedResponse, err
	}
	return wrappedResponse, resp.JSON(&wrappedResponse)
}

type TickersSingleResponse struct {
	Ticker TickerData `json:"ticker,omitempty"`
	Errors []any      `json:"errors,omitempty"`
	Error  any        `json:"error,omitempty"`
	DerivativesResponse
}

// TickerSymbol retrieves the ticker information of a specific contract pair or indice.
//
// https://docs.kraken.com/api/docs/futures-api/trading/get-ticker
func (r *REST) TickerSymbol(symbol string) (*TickersSingleResponse, error) {
	wrappedResponse := &TickersSingleResponse{}
	resp, err := r.Issue(&RequestOptions{
		Method: "GET",
		Path:   []any{"/derivatives/api/v3/tickers/", symbol},
		Auth:   false,
	})
	wrappedResponse.Response = resp
	if err != nil {
		return wrappedResponse, err
	}
	return wrappedResponse, resp.JSON(&wrappedResponse)
}

type OrderBookRequest struct {
	Symbol string `json:"symbol,omitempty"`
}

type OrderBookResponse struct {
	OrderBook OrderBook `json:"orderBook,omitempty"`
	DerivativesResponse
}

// OrderBook retrieves the top bid and ask records of a contract pair.
//
// https://docs.kraken.com/api/docs/futures-api/trading/get-orderbook
func (r *REST) OrderBook(opts *OrderBookRequest) (*OrderBookResponse, error) {
	wrappedResponse := &OrderBookResponse{}
	var query map[string]any
	if opts != nil {
		var err error
		query, err = kraken.StructToMap(opts)
		if err != nil {
			return wrappedResponse, fmt.Errorf("query: %w", err)
		}
	}
	resp, err := r.Issue(&RequestOptions{
		Method: "GET",
		Path:   []any{"/derivatives/api/v3/orderbook"},
		Query:  query,
	})
	wrappedResponse.Response = resp
	if err != nil {
		return wrappedResponse, err
	}
	return wrappedResponse, resp.JSON(&wrappedResponse)
}

type TradeHistoryRequest struct {
	Symbol   string `json:"symbol,omitempty"`
	LastTime string `json:"lastTime,omitempty"`
}

type TradeHistoryResponse struct {
	History []Trade `json:"history,omitempty"`
	DerivativesResponse
}

// TradeHistory retrieves the most recent trade events in a futures market.
//
// https://docs.kraken.com/api/docs/futures-api/trading/get-history
func (r *REST) TradeHistory(opts *TradeHistoryRequest) (*TradeHistoryResponse, error) {
	wrappedResponse := &TradeHistoryResponse{}
	var query map[string]any
	if opts != nil {
		var err error
		query, err = kraken.StructToMap(opts)
		if err != nil {
			return wrappedResponse, fmt.Errorf("query: %w", err)
		}
	}
	resp, err := r.Issue(&RequestOptions{
		Method: "GET",
		Path:   "/derivatives/api/v3/history",
		Query:  query,
	})
	wrappedResponse.Response = resp
	if err != nil {
		return wrappedResponse, err
	}
	return wrappedResponse, resp.JSON(&wrappedResponse)
}

type AccountsResponse struct {
	Accounts map[string]any `json:"accounts,omitempty"`
	DerivativesResponse
}

// Accounts retrieves balances, margin requirements, margin trigger estimates, and other auxilary information of all futures cash and margin accounts.
//
// https://docs.kraken.com/api/docs/futures-api/trading/get-accounts
func (r *REST) Accounts() (*AccountsResponse, error) {
	wrappedResponse := &AccountsResponse{}
	resp, err := r.Issue(&RequestOptions{
		Method: "GET",
		Path:   "/derivatives/api/v3/accounts",
		Auth:   true,
	})
	wrappedResponse.Response = resp
	if err != nil {
		return wrappedResponse, err
	}
	return wrappedResponse, resp.JSON(&wrappedResponse)
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

type SendOrderResponse struct {
	SendStatus OrderStatus `json:"sendStatus,omitempty"`
	DerivativesResponse
}

// SendOrder places a new order.
//
// https://docs.kraken.com/api/docs/futures-api/trading/send-order
func (r *REST) SendOrder(opts *OrderRequest) (*SendOrderResponse, error) {
	wrappedResponse := &SendOrderResponse{}
	var body map[string]any
	if opts != nil {
		var err error
		body, err = kraken.StructToMap(opts)
		if err != nil {
			return wrappedResponse, fmt.Errorf("body: %w", err)
		}
	}
	resp, err := r.Issue(&RequestOptions{
		Method: "POST",
		Path:   []any{"/derivatives/api/v3/sendorder"},
		Body:   body,
		Auth:   true,
	})
	wrappedResponse.Response = resp
	if err != nil {
		return wrappedResponse, err
	}
	return wrappedResponse, resp.JSON(&wrappedResponse)
}

type BatchOrderRequest struct {
	ProcessBefore time.Time       `json:"processBefore,omitempty"`
	JSON          *BatchOrderJson `json:"json,omitempty" map:"stringify"`
}

type BatchOrderResponse struct {
	BatchStatus []BatchStatusInfo `json:"batchStatus,omitempty"`
	DerivativesResponse
}

// BatchOrder allows placing an order, cancelling an open order, or editing an existing order in a single request.
//
// https://docs.kraken.com/api/docs/futures-api/trading/send-batch-order
func (r *REST) BatchOrder(opts *BatchOrderRequest) (*BatchOrderResponse, error) {
	wrappedResponse := &BatchOrderResponse{}
	var body map[string]any
	if opts != nil {
		var err error
		if body, err = kraken.StructToMap(opts); err != nil {
			return wrappedResponse, fmt.Errorf("body: %w", err)
		}
	}
	resp, err := r.Issue(&RequestOptions{
		Method: "POST",
		Path:   "/derivatives/api/v3/batchorder",
		Body:   body,
		Auth:   true,
	})
	wrappedResponse.Response = resp
	if err != nil {
		return wrappedResponse, err
	}
	return wrappedResponse, resp.JSON(&wrappedResponse)
}

type EditOrderResponse struct {
	EditStatus OrderStatus
	DerivativesResponse
}

// EditOrder edits an existing order.
//
// https://docs.kraken.com/api/docs/futures-api/trading/edit-order-spring
func (r *REST) EditOrder(opts *OrderRequest) (*EditOrderResponse, error) {
	wrappedResponse := &EditOrderResponse{}
	var body map[string]any
	if opts != nil {
		var err error
		if body, err = kraken.StructToMap(opts); err != nil {
			return wrappedResponse, fmt.Errorf("body: %w", err)
		}
	}
	resp, err := r.Issue(&RequestOptions{
		Method: "POST",
		Path:   "/derivatives/api/v3/editorder",
		Body:   body,
		Auth:   true,
	})
	wrappedResponse.Response = resp
	if err != nil {
		return wrappedResponse, err
	}
	return wrappedResponse, resp.JSON(&wrappedResponse)
}

type CancelOrderRequest struct {
	ProcessBefore time.Time `json:"processBefore,omitempty"`
	OrderID       string    `json:"order_id,omitempty"`
	ClientOrderID string    `json:"cliOrdId,omitempty"`
}

type CancelOrderResponse struct {
	CancelStatus OrderStatus `json:"cancelStatus,omitempty"`
	DerivativesResponse
}

// CancelOrder cancels an open order.
//
// https://docs.kraken.com/api/docs/futures-api/trading/cancel-order
func (r *REST) CancelOrder(opts *CancelOrderRequest) (*CancelOrderResponse, error) {
	wrappedResponse := &CancelOrderResponse{}
	var body map[string]any
	if opts != nil {
		var err error
		if body, err = kraken.StructToMap(opts); err != nil {
			return wrappedResponse, fmt.Errorf("body: %w", err)
		}
	}
	resp, err := r.Issue(&RequestOptions{
		Method: "POST",
		Path:   []any{"/derivatives/api/v3/cancelorder"},
		Body:   body,
		Auth:   true,
	})
	wrappedResponse.Response = resp
	if err != nil {
		return wrappedResponse, err
	}
	return wrappedResponse, resp.JSON(&wrappedResponse)
}

type CancelAllRequest struct {
	Symbol string `json:"symbol,omitempty"`
}

type CancelAllResponse struct {
	CancelStatus CancelStatus `json:"cancelStatus,omitempty"`
	DerivativesResponse
}

// CancelAll cancels all open orders on a specific contract pair.
// If symbol is unspecified, all open orders regardless of the contract pair will be cancelled.
//
// https://docs.kraken.com/api/docs/futures-api/trading/cancel-all-orders
func (r *REST) CancelAll(opts *CancelAllRequest) (*CancelAllResponse, error) {
	wrappedResponse := &CancelAllResponse{}
	var body map[string]any
	if opts != nil {
		var err error
		if body, err = kraken.StructToMap(opts); err != nil {
			return wrappedResponse, fmt.Errorf("body: %w", err)
		}
	}
	resp, err := r.Issue(&RequestOptions{
		Method: "POST",
		Path:   "/derivatives/api/v3/cancelallorders",
		Body:   body,
		Auth:   true,
	})
	wrappedResponse.Response = resp
	if err != nil {
		return wrappedResponse, err
	}
	return wrappedResponse, resp.JSON(&wrappedResponse)
}

type OpenOrdersResponse struct {
	OpenOrders []OpenOrder `json:"openOrders,omitempty"`
	DerivativesResponse
}

// OpenOrders retrieves information regarding all open orders on the futures account.
//
// https://docs.kraken.com/api/docs/futures-api/trading/get-open-orders
func (r *REST) OpenOrders() (*OpenOrdersResponse, error) {
	wrappedResponse := &OpenOrdersResponse{}
	resp, err := r.Issue(&RequestOptions{
		Method: "GET",
		Path:   "/derivatives/api/v3/openorders",
		Auth:   true,
	})
	wrappedResponse.Response = resp
	if err != nil {
		return wrappedResponse, err
	}
	return wrappedResponse, resp.JSON(&wrappedResponse)
}
