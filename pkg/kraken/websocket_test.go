package kraken

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/krakenfx/api-go/v2/pkg/callback"
)

// wsServer is a WebSocket server that counts the connections it accepts and can
// drop the live one on demand.
type wsServer struct {
	server      *httptest.Server
	connections atomic.Int64
	mux         sync.Mutex
	live        *websocket.Conn
}

func newWSServer(t *testing.T) *wsServer {
	t.Helper()
	s := &wsServer{}
	upgrader := websocket.Upgrader{}
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		s.connections.Add(1)
		s.mux.Lock()
		s.live = conn
		s.mux.Unlock()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				_ = conn.Close()
				return
			}
		}
	}))
	t.Cleanup(s.server.Close)
	return s
}

func (s *wsServer) url() string {
	return "ws" + strings.TrimPrefix(s.server.URL, "http")
}

// drop closes the connection the client is currently using.
func (s *wsServer) drop(t *testing.T) {
	t.Helper()
	deadline := time.After(2 * time.Second)
	for {
		s.mux.Lock()
		conn := s.live
		s.mux.Unlock()
		if conn != nil {
			_ = conn.Close()
			return
		}
		select {
		case <-deadline:
			t.Fatal("no live connection to drop")
		case <-time.After(5 * time.Millisecond):
		}
	}
}

func (s *wsServer) waitForConnections(t *testing.T, want int64) {
	t.Helper()
	deadline := time.After(5 * time.Second)
	for {
		if s.connections.Load() >= want {
			return
		}
		select {
		case <-deadline:
			t.Fatalf("server accepted %d connections, want %d", s.connections.Load(), want)
		case <-time.After(10 * time.Millisecond):
		}
	}
}

func newTestWebSocket(t *testing.T, s *wsServer) *WebSocket {
	t.Helper()
	ws := NewWebSocket()
	ws.URL = s.url()
	ws.ReconnectWait = 10 * time.Millisecond
	ws.DoReconnect = true
	t.Cleanup(func() {
		ws.mux.Lock()
		ws.DoReconnect = false
		ws.mux.Unlock()
		_ = ws.Disconnect()
	})
	return ws
}

func TestConnectRefusesASecondConnection(t *testing.T) {
	server := newWSServer(t)
	ws := newTestWebSocket(t, server)

	if err := ws.Connect(); err != nil {
		t.Fatalf("first connect failed: %v", err)
	}
	server.waitForConnections(t, 1)

	if err := ws.Connect(); !errors.Is(err, ErrAlreadyConnected) {
		t.Fatalf("second connect returned %v, want ErrAlreadyConnected", err)
	}
	time.Sleep(100 * time.Millisecond)
	if got := server.connections.Load(); got != 1 {
		t.Fatalf("server accepted %d connections, want 1", got)
	}
	if !ws.IsActive() {
		t.Fatal("IsActive reports false while connected")
	}
}

// TestDuplicateReconnectHandlerIsHarmless covers the mistake this guard is for.
// NewWebSocket already registers an OnDisconnected handler that reconnects, so
// a caller that registers its own reconnecting handler used to open a second
// connection on every drop, leaving two readers on one connection.
func TestDuplicateReconnectHandlerIsHarmless(t *testing.T) {
	server := newWSServer(t)
	ws := newTestWebSocket(t, server)
	ws.OnDisconnected.Recurring(func(e *callback.Event[error]) {
		if ws.Reconnect != nil && !websocket.IsCloseError(e.Data, websocket.CloseNormalClosure) {
			ws.Reconnect()
		}
	})

	if err := ws.Connect(); err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	server.waitForConnections(t, 1)

	for i := 0; i < 3; i++ {
		want := int64(i + 2)
		server.drop(t)
		server.waitForConnections(t, want)
		// Allow a duplicate connection to appear if one is going to.
		time.Sleep(200 * time.Millisecond)
		if got := server.connections.Load(); got != want {
			t.Fatalf("after %d drops the server accepted %d connections, want %d", i+1, got, want)
		}
	}
}

func TestReconnectsAfterDrop(t *testing.T) {
	server := newWSServer(t)
	ws := newTestWebSocket(t, server)
	var disconnects atomic.Int64
	ws.OnDisconnected.Recurring(func(e *callback.Event[error]) {
		disconnects.Add(1)
	})

	if err := ws.Connect(); err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	server.waitForConnections(t, 1)
	server.drop(t)
	server.waitForConnections(t, 2)

	if got := disconnects.Load(); got != 1 {
		t.Fatalf("%d disconnect events, want 1", got)
	}
}

// TestConcurrentConnectOpensOne covers the race between two dials in flight.
func TestConcurrentConnectOpensOne(t *testing.T) {
	server := newWSServer(t)
	ws := newTestWebSocket(t, server)

	var wg sync.WaitGroup
	errs := make([]error, 8)
	for i := range errs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs[i] = ws.Connect()
		}()
	}
	wg.Wait()

	var connected int
	for _, err := range errs {
		switch {
		case err == nil:
			connected++
		case errors.Is(err, ErrAlreadyConnected):
		default:
			t.Fatalf("unexpected connect error: %v", err)
		}
	}
	if connected != 1 {
		t.Fatalf("%d callers connected, want exactly 1", connected)
	}
	time.Sleep(100 * time.Millisecond)
	// The losing callers must be turned away before they dial, not after.
	if got := server.connections.Load(); got != 1 {
		t.Fatalf("server accepted %d connections, want 1", got)
	}
}

func TestDisconnectDoesNotReconnect(t *testing.T) {
	server := newWSServer(t)
	ws := newTestWebSocket(t, server)

	if err := ws.Connect(); err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	server.waitForConnections(t, 1)
	if err := ws.Disconnect(); err != nil {
		t.Fatalf("disconnect failed: %v", err)
	}
	time.Sleep(300 * time.Millisecond)
	if got := server.connections.Load(); got != 1 {
		t.Fatalf("server accepted %d connections after disconnect, want 1", got)
	}
	if ws.IsActive() {
		t.Fatal("IsActive reports true after disconnect")
	}
}

func TestDisconnectWithoutAConnection(t *testing.T) {
	ws := NewWebSocket()
	if err := ws.Disconnect(); err != nil {
		t.Fatalf("Disconnect on an unconnected WebSocket returned %v, want nil", err)
	}
	if ws.IsActive() {
		t.Fatal("IsActive reports true on an unconnected WebSocket")
	}
}

func TestWriteMessageWithoutAConnection(t *testing.T) {
	ws := NewWebSocket()
	if err := ws.WriteMessage(websocket.TextMessage, []byte("ping")); err == nil {
		t.Fatal("WriteMessage without a connection succeeded, want an error")
	}
}
