package websocket

import "fmt"

type Pool struct {
	Register	chan *Client
	Unregister 	chan *Client
	Clients		map[*Client]bool
	Broadcast	chan NewMessage
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan NewMessage),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))

			for client, _ := range pool.Clients {
				fmt.Println(client.ID)
				/*err := client.Conn.WriteJSON(Message{
					//Type: 1,
					//Body: "New User Joined...",
				})
				if err != nil {
					return 
				}*/
			}
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			
			/*for client, _ := range pool.Clients {
				err := client.Conn.WriteJSON(Message{
					//Type: 1,
					//Body: "User Disconnected...",
				})
				if err != nil {
					return 
				}
			}*/
			break
		case message := <-pool.Broadcast:
			fmt.Println("Sending Message to all clients in Pool")
			fmt.Println(message)
			for client, _ := range pool.Clients {
				client.Ws.Out <- (&Event{
					Event: "response",
					Data: message,
				}).Raw()
				/*if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}*/
			}			
		}
	}
}