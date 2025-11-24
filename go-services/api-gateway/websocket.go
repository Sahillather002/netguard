package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// WebSocket handler for real-time updates
func websocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// Create client
	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()

	// Send initial connection message
	client.send <- []byte(`{"type":"connected","message":"WebSocket connection established"}`)

	// Simulate real-time updates
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Send periodic updates
			update := `{
				"type": "stats_update",
				"data": {
					"packets_captured": 15234,
					"alerts_count": 42,
					"timestamp": "` + time.Now().Format(time.RFC3339) + `"
				}
			}`
			client.send <- []byte(update)
		}
	}
}

// Client represents a WebSocket client
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		log.Printf("Received message: %s", message)

		// Echo back for now
		c.send <- message
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// Stats WebSocket handler - sends real-time statistics
func statsWebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Stats WebSocket client connected")

	// Send stats updates every 2 seconds
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Channel to signal when to stop
	done := make(chan struct{})

	// Read pump to detect client disconnect
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				close(done)
				return
			}
		}
	}()

	for {
		select {
		case <-done:
			log.Println("Stats WebSocket client disconnected")
			return
		case <-ticker.C:
			stats := map[string]interface{}{
				"packets_per_second": 1234.5 + float64(time.Now().Second()%100),
				"bytes_per_second":   5242880 + (time.Now().Unix() % 1000000),
				"active_connections": 42 + (time.Now().Second() % 20),
				"threats_detected":   7 + (time.Now().Second() % 10),
				"timestamp":          time.Now().Unix(),
			}

			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteJSON(stats); err != nil {
				log.Printf("Error sending stats: %v", err)
				return
			}
		}
	}
}

// Threats WebSocket handler - sends real-time threat alerts
func threatsWebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Threats WebSocket client connected")

	// Send threat updates every 10 seconds
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Channel to signal when to stop
	done := make(chan struct{})

	// Read pump to detect client disconnect
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				close(done)
				return
			}
		}
	}()

	threatTypes := []string{"Port Scan", "Brute Force", "DDoS Attack", "SQL Injection", "XSS Attack"}
	severities := []string{"low", "medium", "high", "critical"}
	sourceIPs := []string{"192.168.1.100", "10.0.0.50", "172.16.0.25", "203.0.113.42"}

	for {
		select {
		case <-done:
			log.Println("Threats WebSocket client disconnected")
			return
		case <-ticker.C:
			threat := map[string]interface{}{
				"threat_type":    threatTypes[time.Now().Second()%len(threatTypes)],
				"severity":       severities[time.Now().Second()%len(severities)],
				"source_ip":      sourceIPs[time.Now().Second()%len(sourceIPs)],
				"destination_ip": "192.168.1.1",
				"port":           22 + (time.Now().Second() % 100),
				"timestamp":      time.Now().Format(time.RFC3339),
			}

			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteJSON(threat); err != nil {
				log.Printf("Error sending threat: %v", err)
				return
			}
		}
	}
}
