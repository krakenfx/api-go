package kraken

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/krakenfx/api-go/v2/pkg/callback"
)

// ErrAlreadyConnected is returned by [WebSocket.Connect] when a connection is
// already established, or when another caller is in the middle of establishing
// one.
//
// Opening a second connection is worse than a no-op. The reader of the first
// one keeps running, both readers end up on whichever connection was stored
// last, and two readers on a single gorilla connection share its buffered
// reader. Their frames interleave, the connection fails with errors such as
// "RSV2 set" or "bad opcode", and the shared buffer can be corrupted badly
// enough to panic with a slice bounds error.
var ErrAlreadyConnected = errors.New("already connected")

// WebSocket implements a common structure for the WebSocket APIs.
type WebSocket struct {
	Reconnect     func()
	ReconnectWait time.Duration
	DoReconnect   bool

	OnConnected    *callback.Manager[any]
	OnDisconnected *callback.Manager[error]
	OnSent         *callback.Manager[*WebSocketMessage]
	OnReceived     *callback.Manager[*WebSocketMessage]

	URL      string
	Insecure bool

	// mux guards the connection state below. It is never held while a callback
	// runs or while the network is touched.
	mux sync.Mutex
	// connecting is held for the duration of a dial so that racing callers are
	// turned away before they open a connection that all but one of them would
	// immediately close.
	connecting bool
	conn       *websocket.Conn
	active     bool
	writeMux   sync.Mutex
}

// NewWebSocket creates a new [WebSocket] object with default values.
func NewWebSocket() *WebSocket {
	ws := &WebSocket{
		ReconnectWait:  2 * time.Second,
		OnConnected:    callback.NewManager[any](),
		OnDisconnected: callback.NewManager[error](),
		OnSent:         callback.NewManager[*WebSocketMessage](),
		OnReceived:     callback.NewManager[*WebSocketMessage](),
	}
	ws.Reconnect = func() {
		for {
			err := ws.Connect()
			if err == nil {
				return
			}
			// Another caller has already restored the connection, which is the
			// whole purpose of this loop. Retrying would open a second one.
			if errors.Is(err, ErrAlreadyConnected) {
				return
			}
			time.Sleep(ws.ReconnectWait)
		}
	}
	ws.OnDisconnected.Recurring(func(e *callback.Event[error]) {
		if ws.Reconnect != nil && !websocket.IsCloseError(e.Data, websocket.CloseNormalClosure) && ws.DoReconnect {
			ws.Reconnect()
		}
	})
	return ws
}

// Connect establishes a connection.
//
// It is safe for concurrent use and returns [ErrAlreadyConnected] instead of
// opening a second connection over a live one, so that a caller with its own
// reconnect logic on top of the built-in handler cannot end up with two readers
// on the same connection.
func (ws *WebSocket) Connect() error {
	ws.mux.Lock()
	if ws.active || ws.connecting {
		ws.mux.Unlock()
		return ErrAlreadyConnected
	}
	ws.connecting = true
	insecure, url := ws.Insecure, ws.URL
	ws.mux.Unlock()

	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
	}
	if insecure {
		dialer.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
			VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				return nil
			},
		}
	}
	connection, _, err := dialer.Dial(url, nil)

	ws.mux.Lock()
	ws.connecting = false
	if err != nil {
		ws.mux.Unlock()
		return fmt.Errorf("dial failed: %s", err)
	}
	// A reader may have failed and reconnected while this dial was in flight.
	if ws.active {
		ws.mux.Unlock()
		_ = connection.Close()
		return ErrAlreadyConnected
	}
	ws.conn = connection
	ws.active = true
	ws.DoReconnect = true
	ws.mux.Unlock()

	go ws.read(connection)
	ws.OnConnected.Call(nil)
	return nil
}

type WebSocketMessage struct {
	data   []byte
	mapped map[string]any
	mux    sync.Mutex
}

func NewWebSocketMessage(d []byte) *WebSocketMessage {
	return &WebSocketMessage{
		data: d,
	}
}

func (m *WebSocketMessage) JSON(v any) error {
	decoder := json.NewDecoder(bytes.NewReader(m.data))
	decoder.UseNumber()
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("json unmarshal \"%s\": %w", m.data, err)
	}
	return nil
}

func (m *WebSocketMessage) Bytes() []byte {
	return m.data
}

func (m *WebSocketMessage) String() string {
	return string(m.data)
}

func (m *WebSocketMessage) Map() (map[string]any, error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	if m.mapped != nil {
		return m.mapped, nil
	}
	var dataMapped map[string]any
	if err := m.JSON(&dataMapped); err != nil {
		return nil, err
	}
	m.mapped = dataMapped
	return dataMapped, nil
}

// read pumps a single connection. The connection is a parameter instead of a
// field read on every iteration, so that a reader can never drift onto a
// connection that replaced the one it was started for.
func (ws *WebSocket) read(conn *websocket.Conn) {
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			_ = conn.Close()
			ws.mux.Lock()
			// Only retire the connection if this reader still owns it. Clearing
			// the flag here is what allows the handler below to reconnect.
			if ws.conn == conn {
				ws.active = false
			}
			ws.mux.Unlock()
			ws.OnDisconnected.Call(err)
			return
		}
		ws.OnReceived.Call(NewWebSocketMessage(data))
	}
}

// Disconnect stops the connection. It is a no-op when nothing is connected.
func (ws *WebSocket) Disconnect() error {
	ws.mux.Lock()
	ws.DoReconnect = false
	conn := ws.conn
	ws.mux.Unlock()
	if conn == nil {
		return nil
	}
	// Buffered, and sent without blocking: the handler runs on the reader
	// goroutine, and it must neither block on a channel nobody reads after the
	// wait below times out, nor send on one that has been closed.
	done := make(chan struct{}, 1)
	cb := ws.OnDisconnected.Recurring(func(e *callback.Event[error]) {
		select {
		case done <- struct{}{}:
		default:
		}
	})
	defer ws.OnDisconnected.Deregister(cb)
	defer func() {
		_ = conn.Close()
	}()
	message := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	if err := ws.WriteMessage(websocket.CloseMessage, message); err != nil {
		return fmt.Errorf("write close failed: %s", err)
	}
	select {
	case <-done:
	case <-time.After(time.Second):
	}
	return nil
}

// IsActive returns the status of the connection.
func (ws *WebSocket) IsActive() bool {
	ws.mux.Lock()
	defer ws.mux.Unlock()
	return ws.active
}

// WriteJSON submits a message to the connection.
func (ws *WebSocket) WriteJSON(message any) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("json marshal failed: %s", err)
	}
	if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("write message failed: %s", err)
	}
	return nil
}

// WriteMessage submits a raw message to the connection.
func (ws *WebSocket) WriteMessage(messageType int, data []byte) error {
	ws.mux.Lock()
	conn := ws.conn
	ws.mux.Unlock()
	if conn == nil {
		return fmt.Errorf("no connection")
	}
	ws.writeMux.Lock()
	defer ws.writeMux.Unlock()
	if err := conn.WriteMessage(messageType, data); err != nil {
		return fmt.Errorf("write message failed: %s", err)
	}
	ws.OnSent.Call(NewWebSocketMessage(data))
	return nil
}
