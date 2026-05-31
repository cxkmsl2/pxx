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
			// 如果该用户已有连接，先关掉旧的
			if old, ok := h.clients[client.UserID]; ok {
				old.Conn.Close()
				// 这里不调 close(old.Send)，让 writePump 自然退出
			}
			h.clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("[WS] user %d connected", client.UserID)
		case client := <-h.unregister:
			h.mu.Lock()
			if c, ok := h.clients[client.UserID]; ok && c == client {
				delete(h.clients, client.UserID)
				close(c.Send)
			}
			h.mu.Unlock()
			log.Printf("[WS] user %d disconnected", client.UserID)
		}
	}
}

func (h *Hub) SendToUser(userID uint, msg interface{}) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()

	if ok {
		data, _ := json.Marshal(msg)
		// 非阻塞发送，防止慢客户端阻塞 Hub
		select {
		case client.Send <- data:
		default:
			log.Printf("[WS] user %d send buffer full, dropping message", userID)
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
				return middleware.GetJWTSecret(), nil
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

	client := &Client{UserID: userID, Conn: conn, Send: make(chan []byte, 256)}
	DefaultHub.register <- client

	// 启动读写循环
	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		DefaultHub.unregister <- c
		c.Conn.Close()
	}()
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *Client) writePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}
