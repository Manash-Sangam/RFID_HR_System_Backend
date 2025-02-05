package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rfid_backend/db"
	"rfid_backend/models"

	"github.com/gorilla/websocket"
)

var HubInstance *Hub

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	id   string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
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

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade connection:", err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), id: "ESP8266-01"}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		var data models.RFIDData
		err = json.Unmarshal(message, &data)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err)
			continue
		}

		// Process the RFID data
		if db.VerifyPerson(data.TagID) {
			db.LogRFIDData(data.TagID, data.DeviceID)
			response := []byte("Access Granted")
			c.send <- response
		} else {
			response := []byte("Access Denied")
			c.send <- response
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Println("Error writing message:", err)
				return
			}
		}
	}
}

func (h *Hub) SendToClient(clientID string, message []byte) {
	for client := range h.clients {
		if client.id == clientID {
			client.send <- message
			return
		}
	}
}
