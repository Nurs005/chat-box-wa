package usecase

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type client struct {
	conn      *websocket.Conn
	writeChan chan *payload
}

type payload struct {
	data      []byte
	typeOfMsg int
}

type Hub struct {
	clients map[string]*client
	mu      sync.RWMutex
}

type IHub interface {
	HandleConnection(token string, conn *websocket.Conn)
	Send(token string, data []byte, typeOfMsg int) error
	DeleteConnection(token string)
}

func NewHub() IHub {
	return &Hub{
		clients: make(map[string]*client),
	}
}

func (h *Hub) HandleConnection(token string, conn *websocket.Conn) {
	writeChan := make(chan *payload, 16)
	c := &client{
		conn:      conn,
		writeChan: writeChan,
	}

	h.mu.Lock()
	h.clients[token] = c
	h.mu.Unlock()

	// Writer loop
	go msgSender(h.clients[token])
}

func (h *Hub) Send(token string, data []byte, typeOfMsg int) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	c, ok := h.clients[token]
	if !ok {
		return fmt.Errorf("no active websocket for token: %s", token)
	}

	select {
	case c.writeChan <- &payload{
		data:      data,
		typeOfMsg: typeOfMsg,
	}:
		return nil
	default:
		return fmt.Errorf("write channel full for token: %s", token)
	}
}

func (h *Hub) DeleteConnection(token string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if c, ok := h.clients[token]; ok {
		close(c.writeChan)
		c.conn.Close()
		delete(h.clients, token)
	}
}

func msgSender(c *client) {
	for msg := range c.writeChan {
		if err := c.conn.WriteMessage(msg.typeOfMsg, msg.data); err != nil {
			break
		}
	}
}
