package kraken

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/krakenfx/api-go/v2/pkg/callback"
)

// WebSocket implements a common structure for the WebSocket APIs.
type WebSocket struct {
	Reconnect     func()
	ReconnectWait time.Duration
	DoReconnect   bool

	OnConnected    *callback.Manager[any]
	OnDisconnected *callback.Manager[error]
	OnSent         *callback.Manager[*WebSocketMessage]
	OnReceived     *callback.Manager[*WebSocketMessage]

	conn     *websocket.Conn
	URL      string
	writeMux sync.Mutex
	Insecure bool
	active   bool
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
			if err := ws.Connect(); err == nil {
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
func (ws *WebSocket) Connect() error {
	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
	}
	if ws.Insecure {
		dialer.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
			VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				return nil
			},
		}
	}
	connection, _, err := dialer.Dial(ws.URL, nil)
	if err != nil {
		return fmt.Errorf("dial failed: %s", err)
	}
	ws.conn = connection
	ws.DoReconnect = true
	go ws.read()
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

func (ws *WebSocket) read() {
	ws.active = true
	defer func() {
		ws.active = false
	}()
	for {
		_, data, err := ws.conn.ReadMessage()
		if err != nil {
			_ = ws.conn.Close()
			ws.OnDisconnected.Call(err)
			return
		}
		ws.OnReceived.Call(NewWebSocketMessage(data))
	}
}

// Disconnect stops the connection.
func (ws *WebSocket) Disconnect() error {
	ws.DoReconnect = false
	done := make(chan bool)
	defer close(done)
	defer func() {
		_ = ws.conn.Close()
	}()
	callback := ws.OnDisconnected.Recurring(func(e *callback.Event[error]) {
		done <- true
	})
	message := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	if err := ws.WriteMessage(websocket.CloseMessage, message); err != nil {
		return fmt.Errorf("write close failed: %s", err)
	}
	select {
	case <-done:
	case <-time.After(time.Second):
	}
	ws.OnDisconnected.Deregister(callback)
	return nil
}

// IsActive returns the status of the connection.
func (ws *WebSocket) IsActive() bool {
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
	if ws.conn == nil {
		return fmt.Errorf("no connection")
	}
	ws.writeMux.Lock()
	defer ws.writeMux.Unlock()
	if err := ws.conn.WriteMessage(messageType, data); err != nil {
		return fmt.Errorf("write message failed: %s", err)
	}
	ws.OnSent.Call(NewWebSocketMessage(data))
	return nil
}
