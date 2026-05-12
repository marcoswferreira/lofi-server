package ws

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestNewHub(t *testing.T) {
	h := NewHub()
	if h.stations == nil {
		t.Error("Expected stations map to be initialized")
	}
}

func TestHubRun(t *testing.T) {
	h := NewHub()
	go h.Run()

	c1 := &Client{hub: h, send: make(chan []byte, 1), stationID: "s1"}
	c2 := &Client{hub: h, send: make(chan []byte, 1), stationID: "s1"}
	c3 := &Client{hub: h, send: make(chan []byte, 1), stationID: "s2"}

	// Register
	h.register <- c1
	h.register <- c2
	h.register <- c3

	// Allow some time for processing
	time.Sleep(50 * time.Millisecond)

	h.mu.RLock()
	if len(h.stations["s1"]) != 2 {
		t.Errorf("Expected 2 clients in s1, got %d", len(h.stations["s1"]))
	}
	if len(h.stations["s2"]) != 1 {
		t.Errorf("Expected 1 client in s2, got %d", len(h.stations["s2"]))
	}
	h.mu.RUnlock()

	// Broadcast to s1
	h.broadcast <- Message{Type: "chat", StationID: "s1", Payload: "hello s1"}

	select {
	case msg := <-c1.send:
		if string(msg) == "" {
			t.Error("Client 1 should have received a message")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Timeout waiting for message on c1")
	}

	// Unregister
	h.unregister <- c1
	time.Sleep(50 * time.Millisecond)

	h.mu.RLock()
	if len(h.stations["s1"]) != 1 {
		t.Errorf("Expected 1 client in s1 after unregister, got %d", len(h.stations["s1"]))
	}
	h.mu.RUnlock()

	// Test full buffer broadcast
	c4 := &Client{hub: h, send: make(chan []byte, 1), stationID: "s1"}
	h.register <- c4
	time.Sleep(50 * time.Millisecond)

	// Send 2 messages, second one should trigger default (disconnect)
	h.broadcast <- Message{Type: "chat", StationID: "s1", Payload: "msg 1"}
	h.broadcast <- Message{Type: "chat", StationID: "s1", Payload: "msg 2"}

	time.Sleep(50 * time.Millisecond)
	h.mu.RLock()
	if _, ok := h.stations["s1"][c4]; ok {
		t.Error("Client 4 should have been disconnected due to full buffer")
	}
	h.mu.RUnlock()
}

func TestServeWs_MissingStation(t *testing.T) {
	h := NewHub()
	req := httptest.NewRequest("GET", "/ws", nil)
	rr := httptest.NewRecorder()

	ServeWs(h, rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", rr.Code)
	}
}

func TestServeWs_Success(t *testing.T) {
	h := NewHub()
	go h.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWs(h, w, r)
	}))
	defer server.Close()

	// Convert http URL to ws URL
	wsURL := "ws" + server.URL[4:] + "?station=s1"

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	// Wait for registration
	time.Sleep(50 * time.Millisecond)

	// Test writing to client (writePump)
	msg := Message{Type: "chat", StationID: "s1", Payload: "hello"}
	h.broadcast <- msg

	_, received, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}
	if !strings.Contains(string(received), "hello") {
		t.Errorf("Expected message to contain 'hello', got %s", string(received))
	}

	// Test reading from client (readPump)
	clientMsg := Message{Type: "chat", Payload: "from client"}
	data, _ := json.Marshal(clientMsg)
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		t.Fatalf("Failed to write message: %v", err)
	}

	// The message should be broadcast back to us (since we are in s1)
	_, received, err = conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read back message: %v", err)
	}
	if !strings.Contains(string(received), "from client") {
		t.Errorf("Expected message to contain 'from client', got %s", string(received))
	}
}
