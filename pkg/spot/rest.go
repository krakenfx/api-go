package spot

import (
	"encoding/json"
	"fmt"

	"github.com/krakenfx/api-go/pkg/kraken"
)

// REST wraps [RESTBase] with functions to call common endpoints.
type REST struct {
	// CreateUser creates a new user account in the Kraken system.
	//
	// https://docs.kraken.com/api/docs/embed-api/create-embed-user
	CreateUser func(opts *CreateUserRequest) (*Response[CreateUserResult], error)

	// UpdateUser updates an existing user's profile.
	//
	// https://docs.kraken.com/api/docs/embed-api/update-embed-user
	UpdateUser func(opts *UpdateUserRequest) (*Response[string], error)

	// VerifyUser submits a verification for a user with documents and details.
	//
	// https://docs.kraken.com/api/docs/embed-api/submit-embed-verification
	VerifyUser func(opts *SubmitVerificationRequest) (*Response[SubmitVerificationResult], error)

	// GetUser retrieves a previously created user.
	//
	// https://docs.kraken.com/api/docs/embed-api/get-embed-user
	GetUser func(opts *GetUserRequest) (*Response[GetUserResult], error)

	// ServerTime retrieves the current server time.
	//
	// https://docs.kraken.com/api/docs/rest-api/get-server-time
	ServerTime func() (*Response[ServerTimeResult], error)

	// Balances retrieves the balances on the spot wallet.
	//
	// https://docs.kraken.com/api/docs/rest-api/get-account-balance
	Balances func() (*Response[map[string]*kraken.Money], error)

	// TradesHistory retrieves the trade events of the user.
	//
	// https://docs.kraken.com/api/docs/rest-api/get-trade-history
	TradesHistory func(opts *TradesHistoryRequest) (*Response[TradesHistoryResult], error)

	// OpenOrders retrieves information about currently open orders.
	//
	// https://docs.kraken.com/api/docs/rest-api/get-open-orders
	OpenOrders func(opts *OpenOrdersRequest) (*Response[OpenOrdersResult], error)

	// ClosedOrders retrieves information about orders that have been closed, either filled or cancelled.
	//
	// https://docs.kraken.com/api/docs/rest-api/get-closed-orders
	ClosedOrders func(opts *ClosedOrdersRequest) (*Response[ClosedOrdersResult], error)

	// QueryOrders retrieves information about specific orders.
	//
	// https://docs.kraken.com/api/docs/rest-api/get-orders-info
	QueryOrders func(opts *QueryOrdersRequest) (*Response[map[string]ClosedOrder], error)

	// AddOrder places a new order.
	//
	// https://docs.kraken.com/api/docs/rest-api/add-order
	AddOrder func(opts *AddOrderRequest) (*Response[AddOrderResult], error)

	// AddBatch places a collection of orders, with the minimum of 2 and maximum of 15.
	//
	// https://docs.kraken.com/api/docs/rest-api/add-order-batch
	AddBatch func(opts *AddBatchRequest) (*Response[AddBatchResult], error)

	// AmendOrder changes the properties of an open order without cancelling the existing one and creating a new one, enabling order ID and queue priority to be maintained where possible.
	//
	// https://docs.kraken.com/api/docs/rest-api/amend-order
	AmendOrder func(opts *AmendOrderRequest) (*Response[AmendOrderResult], error)

	// CancelAll cancels all open orders.
	//
	// https://docs.kraken.com/api/docs/rest-api/cancel-all-orders
	CancelAll func() (*Response[CancelResult], error)

	// CancelOrder cancels an individual or set of open orders provided by `txid`, `userref“, or `cl_ord_id“.
	//
	// https://docs.kraken.com/api/docs/rest-api/cancel-order
	CancelOrder func(opts *CancelOrderRequest) (*Response[CancelResult], error)

	// Assets retrieves information about assets that are available for deposit, withdrawal, trading, and earn.
	//
	// https://docs.kraken.com/api/docs/rest-api/get-asset-info
	Assets func(opts *AssetsRequest) (*Response[map[string]AssetInfo], error)

	// AssetPairs retrieves information about tradeable asset pairs on spot.
	//
	// https://docs.kraken.com/api/docs/rest-api/get-tradable-asset-pairs
	AssetPairs func(opts *AssetPairsRequest) (*Response[map[string]AssetPair], error)

	// Ticker retrieves information of all spot markets, or a specific market if `pair` is specified.
	//
	// https://docs.kraken.com/api/docs/rest-api/get-ticker-information
	Ticker func(opts *TickerRequest) (*Response[map[string]AssetTickerInfo], error)

	// OrderBook retrieves the bid and ask records of a specific spot market.
	//
	// https://docs.kraken.com/api/docs/rest-api/get-order-book
	OrderBook func(opts *OrderBookRequest) (*Response[map[string]OrderBook], error)

	// RecentTrades retrieves the recent trade records of a specified spot market.
	//
	// https://docs.kraken.com/api/docs/rest-api/get-recent-trades
	RecentTrades func(opts *RecentTradesRequest) (*Response[map[string]any], error)

	// OHLC retrieves recent open, high, low, close, and volume records of a specified spot market.
	//
	// https://docs.kraken.com/api/docs/rest-api/get-ohlc-data
	OHLC func(opts *OHLCRequest) (*Response[map[string]any], error)

	// GetWebSocketTokens generates an authentication token for WebSocket, which must be used within 15 minutes of creation to prevent expiration.
	//
	// https://docs.kraken.com/api/docs/rest-api/get-websockets-token
	GetWebSocketsToken func() (*Response[GetWebSocketsTokenResult], error)

	*RESTBase
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
	rest := &REST{
		RESTBase: NewRESTBase(),
	}
	rest.CreateUser = NewAPIFunction[CreateUserRequest, CreateUserResult](&APIFunctionOptions{
		REST:   rest,
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/CreateUser",
	})
	rest.UpdateUser = NewAPIFunction[UpdateUserRequest, string](&APIFunctionOptions{
		REST:       rest,
		Auth:       true,
		Method:     "POST",
		Path:       "/0/private/UpdateUser",
		BodyField:  true,
		QueryField: true,
	})
	rest.GetUser = NewAPIFunction[GetUserRequest, GetUserResult](&APIFunctionOptions{
		REST:      rest,
		Auth:      true,
		Method:    "POST",
		Path:      "/0/private/GetUser",
		ParamMode: QueryMode,
	})
	rest.VerifyUser = NewAPIFunction[SubmitVerificationRequest, SubmitVerificationResult](&APIFunctionOptions{
		REST:       rest,
		Auth:       true,
		Method:     "POST",
		Path:       "/0/private/VerifyUser",
		Headers:    map[string]any{"Content-Type": "multipart/form-data"},
		BodyField:  true,
		QueryField: true,
		Version:    1,
	})
	rest.ServerTime = NewAPIFunctionWithNoParams[ServerTimeResult](&APIFunctionOptions{
		REST:   rest,
		Method: "GET",
		Path:   "/0/public/Time",
	})
	rest.Balances = NewAPIFunctionWithNoParams[map[string]*kraken.Money](&APIFunctionOptions{
		REST:   rest,
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/Balance",
	})
	rest.TradesHistory = NewAPIFunction[TradesHistoryRequest, TradesHistoryResult](&APIFunctionOptions{
		REST:   rest,
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/TradesHistory",
	})
	rest.OpenOrders = NewAPIFunction[OpenOrdersRequest, OpenOrdersResult](&APIFunctionOptions{
		REST:   rest,
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/OpenOrders",
	})
	rest.ClosedOrders = NewAPIFunction[ClosedOrdersRequest, ClosedOrdersResult](&APIFunctionOptions{
		REST:   rest,
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/ClosedOrders",
	})
	rest.QueryOrders = NewAPIFunction[QueryOrdersRequest, map[string]ClosedOrder](&APIFunctionOptions{
		REST:   rest,
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/QueryOrders",
	})
	rest.AddOrder = NewAPIFunction[AddOrderRequest, AddOrderResult](&APIFunctionOptions{
		REST:   rest,
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/AddOrder",
	})
	rest.AddBatch = NewAPIFunction[AddBatchRequest, AddBatchResult](&APIFunctionOptions{
		REST:   rest,
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/AddOrderBatch",
	})
	rest.AmendOrder = NewAPIFunction[AmendOrderRequest, AmendOrderResult](&APIFunctionOptions{
		REST:   rest,
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/AmendOrder",
	})
	rest.CancelAll = NewAPIFunctionWithNoParams[CancelResult](&APIFunctionOptions{
		REST:   rest,
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/CancelAll",
	})
	rest.CancelOrder = NewAPIFunction[CancelOrderRequest, CancelResult](&APIFunctionOptions{
		REST:   rest,
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/CancelOrder",
	})
	rest.Assets = NewAPIFunction[AssetsRequest, map[string]AssetInfo](&APIFunctionOptions{
		REST:      rest,
		Method:    "GET",
		Path:      "/0/public/Assets",
		ParamMode: QueryMode,
	})
	rest.AssetPairs = NewAPIFunction[AssetPairsRequest, map[string]AssetPair](&APIFunctionOptions{
		REST:      rest,
		Method:    "GET",
		Path:      "/0/public/AssetPairs",
		ParamMode: QueryMode,
	})
	rest.Ticker = NewAPIFunction[TickerRequest, map[string]AssetTickerInfo](&APIFunctionOptions{
		REST:      rest,
		Method:    "GET",
		Path:      "/0/public/Ticker",
		ParamMode: QueryMode,
	})
	rest.OrderBook = NewAPIFunction[OrderBookRequest, map[string]OrderBook](&APIFunctionOptions{
		REST:      rest,
		Method:    "GET",
		Path:      "/0/public/Depth",
		ParamMode: QueryMode,
	})
	rest.RecentTrades = NewAPIFunction[RecentTradesRequest, map[string]any](&APIFunctionOptions{
		REST:      rest,
		Method:    "GET",
		Path:      "/0/public/Trades",
		ParamMode: QueryMode,
	})
	rest.OHLC = NewAPIFunction[OHLCRequest, map[string]any](&APIFunctionOptions{
		REST:      rest,
		Method:    "GET",
		Path:      "/0/public/OHLC",
		ParamMode: QueryMode,
	})
	rest.GetWebSocketsToken = NewAPIFunctionWithNoParams[GetWebSocketsTokenResult](&APIFunctionOptions{
		REST:   rest,
		Auth:   true,
		Method: "POST",
		Path:   "/0/private/GetWebSocketsToken",
	})
	return rest
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
