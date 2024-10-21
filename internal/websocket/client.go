package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
	MaxMessageSizeInBytes = 512
)

type Clients map[*Client]bool
// this map will be used to make connection faster with sending messages, instead o iterating over all connection to check the userId 
// is equal to the wanted user to send a message, it will be mapped directly with map.
type RegistedClients map[uint]*Client

type Client struct {
	conn    *websocket.Conn
	manager *Manager
	egress  chan Event
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		conn:    conn,
		manager: manager,
		egress:  make(chan Event),
	}
}

func (c *Client) readMessage() {
	defer func() {
		c.manager.deleteClient(c)
	}()

	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Println(err)
		return
	}

	c.conn.SetReadLimit(MaxMessageSizeInBytes)

	c.conn.SetPongHandler(c.pongHandler)

	for {
		var event Event
		err := c.conn.ReadJSON(&event)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error reading message: ", err)
			}
			break
		}

		err = c.manager.routeHandler(event, c)
		if err != nil {
			log.Println(err)
			break
		}
	}

}

func (c *Client) writeMessage() {
	defer func() {
		c.manager.deleteClient(c)
	}()

	ticker := time.NewTicker(pingInterval)
	for {
		// this was used to disallow concurrency (for safety) and forbid the user from sending 100 concurrent goroutine at once
		select {
			
		case message, ok := <-c.egress:
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}

			messageBytes, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				break
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
				log.Println("an error has occurred during sending message: ", err)
				return
			}
			log.Println("message sent")

		case <-ticker.C:
			log.Println("ping...")
			err := c.conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Println("connection has failed to respond, err: ", err)
				return
			}
		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	log.Println("pong")
	return c.conn.SetReadDeadline(time.Now().Add(pongWait))
}
