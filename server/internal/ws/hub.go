package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"pxx/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
var jwtSecret = []byte("pxx-jwt-secret-2026")

type Client struct {
	UserID uint
	Conn   *websocket.Conn
	Send   chan []byte
}

type Hub struct {
	clients    map[uint]*Client
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

var DefaultHub *Hub

func InitHub() {
	DefaultHub = &Hub{
		clients:    make(map[uint]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go DefaultHub.run()
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("[WS] user %d connected", client.UserID)
		case client := <-h.unregister:
			h.mu.Lock()
			if c, ok := h.clients[client.UserID]; ok {
				close(c.Send)
				delete(h.clients, client.UserID)
			}
			h.mu.Unlock()
			log.Printf("[WS] user %d disconnected", client.UserID)
		}
	}
}

func (h *Hub) SendToUser(userID uint, msg interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if client, ok := h.clients[userID]; ok {
		data, _ := json.Marshal(msg)
		select {
		case client.Send <- data:
		default:
		}
	}
}

func (h *Hub) IsOnline(userID uint) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

func HandleWS(c *gin.Context) {
	userID := c.GetUint("user_id")
	// Fallback: parse token from query
	if userID == 0 {
		tokenStr := c.Query("token")
		if tokenStr != "" {
			claims := &middleware.Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})
			if err == nil && token.Valid {
				userID = claims.UserID
			}
		}
	}
	if userID == 0 {
		c.JSON(401, gin.H{"code": 401, "msg": "unauthorized"})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil { return }
	client := &Client{UserID: userID, Conn: conn, Send: make(chan []byte, 64)}
	DefaultHub.register <- client

	go func() {
		defer func() { DefaultHub.unregister <- client; conn.Close() }()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil { break }
			// Echo for now; real processing in service layer
			log.Printf("[WS] received from %d: %s", userID, string(msg))
		}
	}()
	go func() {
		for msg := range client.Send {
			conn.WriteMessage(websocket.TextMessage, msg)
		}
	}()
}
