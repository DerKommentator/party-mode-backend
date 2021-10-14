package websocket

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   		uuid.UUID
	Username 	string
	Conn 		*websocket.Conn
	Pool 		*Pool
	Ws 			*WebSocket
}

type Message struct {
	Type 			int			`json:"type"`
	ClientID 		uuid.UUID 	`json:"client_id"`
	ClientUsername 	string 		`json:"client_username"`
	Body 			string 		`json:"body"`
}

type NewMessage struct {
	ClientID 		uuid.UUID 	`json:"client_id"`
	ClientUsername 	string 		`json:"client_username"`
	Message			string 		`json:"message"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()

	/*for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			return
		}

		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		message := Message{
			Type: messageType,
			ClientID: c.ID,
			ClientUsername: c.Username,
			Body: string(p),
		}

		c.Pool.Broadcast <- message
		fmt.Printf("Message: %+v\n", message)
	}*/
}