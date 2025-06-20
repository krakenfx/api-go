package derivatives

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/krakenfx/api-go/internal/helper"
	"github.com/krakenfx/api-go/pkg/callback"
	"github.com/krakenfx/api-go/pkg/kraken"
)

// WebSocketBase is the underlying of the [WebSocket] client.
type WebSocketBase struct {
	AuthenticateTimeout time.Duration
	PublicKey           string
	PrivateKey          string
	Challenge           string
	Signature           string
	OnAuthenticated     *callback.Manager[string]
	*kraken.WebSocket
}

// NewWebSocketBase constructs a [WebSocketBase] struct with default values.
func NewWebSocketBase() *WebSocketBase {
	b := &WebSocketBase{
		AuthenticateTimeout: 15 * time.Second,
		WebSocket:           kraken.NewWebSocket(),
		OnAuthenticated:     callback.NewManager[string](),
	}
	b.URL = "wss://futures.kraken.com/ws/v1"
	return b
}

// Authenticate submits a challenge request and retrieves the authentication fields.
//
// If contained within a [WebSocketBase] callback, this must be wrapped with a goroutine to prevent blocking.
func (b *WebSocketBase) Authenticate() error {
	if err := b.WriteJSON(map[string]any{
		"event":   "challenge",
		"api_key": b.PublicKey,
	}); err != nil {
		return fmt.Errorf("request challenge failed: %s", err)
	}
	var mainErr error
	threadStarted := time.Now()
	b.OnReceived.SleepUntilDisabled(func(e *callback.Event[*kraken.WebSocketMessage]) {
		dataMap, err := e.Data.Map()
		if err != nil {
			e.Callback.Enabled = false
			mainErr = err
			return
		}
		if event, err := helper.Traverse[string](dataMap, "event"); err != nil || *event != "challenge" {
			if time.Now().After(threadStarted.Add(b.AuthenticateTimeout)) {
				e.Callback.Enabled = false
				mainErr = fmt.Errorf("authentication timed out")
				return
			}
			return
		}
		e.Callback.Enabled = false
		message, err := helper.Traverse[string](dataMap, "message")
		if err != nil {
			mainErr = err
			return
		}
		b.Challenge = *message
	})
	if mainErr != nil {
		return fmt.Errorf("retrieve challenge: %w", mainErr)
	}
	sha256Hash := sha256.New()
	sha256Hash.Write([]byte(b.Challenge))
	signature, err := helper.Sign(b.PrivateKey, sha256Hash.Sum(nil))
	if err != nil {
		return fmt.Errorf("sign challenge failed: %s", err)
	}
	b.Signature = signature
	b.OnAuthenticated.Call(b.Challenge)
	return nil
}

// SendPrivate sends a JSON-encoded map with the authentication fields included.
func (b *WebSocketBase) SendPrivate(m map[string]any) error {
	return b.WriteJSON(helper.Maps(map[string]any{
		"api_key":            b.PublicKey,
		"original_challenge": b.Challenge,
		"signed_challenge":   b.Signature,
	}, m))
}

// SubPublic submits a subscription request.
func (b *WebSocket) SubPublic(feed string, options ...map[string]any) error {
	return b.WriteJSON(helper.Maps(map[string]any{
		"event": "subscribe",
		"feed":  feed,
	}, options...))
}

// SubPrivate submits a subscription request with the authentication fields included.
func (b *WebSocket) SubPrivate(feed string, options ...map[string]any) error {
	return b.SendPrivate(helper.Maps(map[string]any{
		"event": "subscribe",
		"feed":  feed,
	}, options...))
}
