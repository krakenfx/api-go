package spot

import (
	"github.com/krakenfx/api-go/pkg/kraken"
)

// WebSocket wraps a [WebSocketBase] struct with order management and subscription request functions.
type WebSocket struct {
	*WebSocketBase
}

// NewWebSocket constructs a new [WebSocket] struct with default values.
//
// For authentication, store the spot API key on REST.PublicKey and REST.PrivateKey.
func NewWebSocket() *WebSocket {
	s := &WebSocket{
		WebSocketBase: NewWebSocketBase(),
	}
	return s
}

// SubExecutions sends a subscription requestffor the order and trade events of the authenticated user.
//
// https://docs.kraken.com/api/docs/websocket-v2/executions
func (s *WebSocket) SubExecutions(options ...map[string]any) error {
	return s.SubPrivate("executions", options...)
}

// SubBalances sends a subscription request for the asset balances and ledger entries of the authenticated user.
//
// https://docs.kraken.com/api/docs/websocket-v2/balances
func (s *WebSocket) SubBalances(options ...map[string]any) error {
	return s.SubPrivate("balances", options...)
}

// SubTicker sends a subscription request for level 1 market data.
// These are the top of the book prices and most recent trade data.
//
// https://docs.kraken.com/api/docs/websocket-v2/ticker
func (s *WebSocket) SubTicker(symbols []string, options ...map[string]any) error {
	return s.SubPublic("ticker", append([]map[string]any{{
		"params": map[string]any{
			"symbol": symbols,
		}}}, options...)...)
}

// SubBook sends a subscription request to stream level 2 market data.
// These are the individual price levels and aggregated order quantities at each level adjacent to the best bid or best ask price.
//
// https://docs.kraken.com/api/docs/websocket-v2/book
func (s *WebSocket) SubBook(symbols []string, depth int, options ...map[string]any) error {
	return s.SubPublic("book", append([]map[string]any{{
		"params": map[string]any{
			"symbol": symbols,
			"depth":  depth,
		},
	}}, options...)...)
}

// SubL3 sends a subscription request to stream level 3 market data, which are information about individual orders in the book.
//
// https://docs.kraken.com/api/docs/websocket-v2/level3
func (s *WebSocket) SubL3(symbols []string, depth int, options ...map[string]any) error {
	return s.SubPrivate("level3", append([]map[string]any{{
		"params": map[string]any{
			"symbol": symbols,
		},
	}}, options...)...)
}

// SubCandles sends a subscription request to stream the open, high, low, and close prices and volume of the specified spot markets.
//
// https://docs.kraken.com/api/docs/websocket-v2/ohlc
func (s *WebSocket) SubCandles(symbols []string, options ...map[string]any) error {
	return s.SubPublic("ohlc", append([]map[string]any{{"params": map[string]any{"symbol": symbols}}}, options...)...)
}

// SubTrades sends a subscription request to stream the aggregated trade events of the specified spot markets.
//
// https://docs.kraken.com/api/docs/websocket-v2/trade
func (s *WebSocket) SubTrades(symbols []string, options ...map[string]any) error {
	return s.SubPublic("trade", append([]map[string]any{{"params": map[string]any{"symbol": symbols}}}, options...)...)
}

// SubInstruments sends a subscription request to stream asset and asset pair information.
//
// https://docs.kraken.com/api/docs/websocket-v2/instrument
func (s *WebSocket) SubInstruments(options ...map[string]any) error {
	return s.SubPublic("instrument", options...)
}

// AddOrder places a new order.
//
// https://docs.kraken.com/api/docs/websocket-v2/add_order
func (s *WebSocket) AddOrder(orderType string, side string, quantity float64, symbol string, options ...map[string]any) error {
	return s.SendPrivate(kraken.Maps(map[string]any{
		"method": "add_order",
		"params": map[string]any{
			"order_type": orderType,
			"side":       side,
			"order_qty":  quantity,
			"symbol":     symbol,
		},
	}, options...))
}

// AmendOrder changes the properties of an open order without cancelling the existing one and creating a new one.
// This enables order ID and queue priority to be maintained where possible.
//
// https://docs.kraken.com/api/docs/websocket-v2/amend_order
func (s *WebSocket) AmendOrder(options ...map[string]any) error {
	return s.SendPrivate(kraken.Maps(map[string]any{
		"method": "amend_order",
	}, options...))
}

// CancelAllOrders cancels all open orders.
//
// https://docs.kraken.com/api/docs/websocket-v2/cancel_all
func (s *WebSocket) CancelAllOrders(options ...map[string]any) error {
	return s.SendPrivate(kraken.Maps(map[string]any{
		"method": "cancel_all",
	}, options...))
}

// CancelOrder cancels an individual or set of open orders provided by `order_id`, `cl_ord_id`, or `order_userref`
//
// https://docs.kraken.com/api/docs/websocket-v2/cancel_order
func (s *WebSocket) CancelOrder(options ...map[string]any) error {
	return s.SendPrivate(kraken.Maps(map[string]any{
		"method": "cancel_order",
	}, options...))
}
