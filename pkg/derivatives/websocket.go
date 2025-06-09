package derivatives

// WebSocket wraps a [WebSocketBase] struct with order management and subscription request functions.
type WebSocket struct {
	REST *REST
	*WebSocketBase
}

// NewWebSocket constructs a new [WebSocket] struct with default values.
//
// For authentication, store the derivatives API key on REST.PublicKey and REST.PrivateKey.
func NewWebSocket() *WebSocket {
	ws := &WebSocket{
		REST:          NewREST(),
		WebSocketBase: NewWebSocketBase(),
	}
	return ws
}

// SubBalances sends a subscription request to retrieve information for holding wallets, single collateral wallets and multi-collateral wallets.
//
// https://docs.kraken.com/api/docs/futures-api/websocket/balances
func (s *WebSocket) SubBalances() error {
	return s.SubPrivate("balances")
}

// SubOpenOrders sends a subscription request to retrieve information about user open orders.
//
// https://docs.kraken.com/api/docs/futures-api/websocket/open_orders_verbose
func (s *WebSocket) SubOpenOrders() error {
	return s.SubPrivate("open_orders_verbose")
}

// SubExecutions sends a subscription request for the user's fill events.
//
// https://docs.kraken.com/api/docs/futures-api/websocket/fills
func (s *WebSocket) SubExecutions() error {
	return s.SubPrivate("fills")
}

// SubTicker sends a subscription request for pricing information regarding available futures markets.
//
// https://docs.kraken.com/api/docs/futures-api/websocket/ticker
func (s *WebSocket) SubTicker(productID ...string) error {
	var productSubscription map[string]any
	if len(productID) > 0 {
		productSubscription = map[string]any{
			"product_ids": productID,
		}
	}
	return s.SubPublic("ticker", productSubscription)
}

// SubBook sends a subscription request to retrieve information about the order book.
//
// https://docs.kraken.com/api/docs/futures-api/websocket/book
func (s *WebSocket) SubBook(productID ...string) error {
	var productSubscription map[string]any
	if len(productID) > 0 {
		productSubscription = map[string]any{
			"product_ids": productID,
		}
	}
	return s.SubPublic("book", productSubscription)
}

// SubTrade sends a subscription request to retrieve information about executed trades.
//
// https://docs.kraken.com/api/docs/futures-api/websocket/trade
func (s *WebSocket) SubTrade(productID ...string) error {
	var productSubscription map[string]any
	if len(productID) > 0 {
		productSubscription = map[string]any{
			"product_ids": productID,
		}
	}
	return s.SubPublic("trade", productSubscription)
}
