package manager

import (
	"encoding/json"

	"github.com/gorilla/websocket"

	"github.com/farhanramadhan/app-chat/message"
)

// Client is a struct for client
type Client struct {
	ID     string
	Socket *websocket.Conn
	Send   chan []byte
}

func (c *Client) Read(clientManager *ClientManager) {
	defer func() {
		clientManager.unregister <- c
		c.Socket.Close()
	}()

	for {
		_, messageByte, err := c.Socket.ReadMessage()
		if err != nil {
			clientManager.unregister <- c
			c.Socket.Close()
			break
		}
		jsonMessage, _ := json.Marshal(&message.Message{Sender: c.ID, Content: string(messageByte)})
		clientManager.broadcast <- jsonMessage
	}
}

func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// ClientManager is a struct for client manager
type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	Register   chan *Client
	unregister chan *Client
}

// CreateClientManager is constructor for client managers
func CreateClientManager() ClientManager {
	return ClientManager{
		broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Start is to start client manager
func (cm *ClientManager) Start() {
	// While True
	for {
		// Case of for go routine
		select {
		case conn := <-cm.Register:
			cm.clients[conn] = true
			jsonMessage, _ := json.Marshal(&message.Message{Content: "/A new socket has connected."})
			cm.Send(jsonMessage, conn)
		case conn := <-cm.unregister:
			if _, ok := cm.clients[conn]; ok {
				close(conn.Send)
				delete(cm.clients, conn)
				jsonMessage, _ := json.Marshal(&message.Message{Content: "/A socket has disconnected."})
				cm.Send(jsonMessage, conn)
			}
		case message := <-cm.broadcast:
			for conn := range cm.clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(cm.clients, conn)
				}
			}
		}
	}
}

// Send is to send to all client
func (cm *ClientManager) Send(message []byte, ignore *Client) {
	for conn := range cm.clients {
		if conn != ignore {
			conn.Send <- message
		}
	}
}
