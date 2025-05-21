package spot

import (
	"fmt"

	"github.com/krakenfx/api-go/pkg/kraken"
)

// WebSocketBase is the underlying of the [WebSocket] client.
type WebSocketBase struct {
	REST            *REST
	Token           string
	OnAuthenticated *kraken.CallbackManager[string]
	*kraken.WebSocket
}

// NewWebSocketBase constructs a [WebSocketBase] struct with default values.
func NewWebSocketBase() *WebSocketBase {
	b := &WebSocketBase{
		REST:            NewREST(),
		OnAuthenticated: kraken.NewCallbackManager[string](),
		WebSocket:       kraken.NewWebSocket(),
	}
	b.URL = "wss://ws.kraken.com/v2"
	return b
}

// Authenticate retrieves a new WebSocket token from the REST API.
func (b *WebSocketBase) Authenticate() error {
	resp, err := b.REST.GetWebSocketsToken()
	if err != nil {
		return fmt.Errorf("get websockets token: %w", err)
	}
	b.Token = resp.Result.Token
	b.OnAuthenticated.Call(b.Token)
	return nil
}

// SendPublic submits a JSON-encoded map.
func (b *WebSocketBase) SendPublic(m map[string]any) error {
	return b.WriteJSON(m)
}

// SendPrivate submits a JSON-encoded map with the token included.
func (b *WebSocketBase) SendPrivate(m map[string]any) error {
	return b.WriteJSON(kraken.Maps(map[string]any{
		"params": map[string]any{
			"token": b.Token,
		},
	}, m))
}

// SubPublic submits a subscription request.
func (b *WebSocketBase) SubPublic(channel string, options ...map[string]any) error {
	return b.SendPublic(kraken.Maps(map[string]any{
		"method": "subscribe",
		"params": map[string]any{
			"channel": channel,
		},
	}, options...))
}

// SubPrivate submits a subscription request with the token included.
func (b *WebSocketBase) SubPrivate(channel string, options ...map[string]any) error {
	return b.SendPrivate(kraken.Maps(map[string]any{
		"method": "subscribe",
		"params": map[string]any{
			"channel": channel,
		},
	}, options...))
}
