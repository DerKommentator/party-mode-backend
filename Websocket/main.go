package main

import (
	"Bergflix-Partymode/pkg/websocket"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)


func wsHandler(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")

	ws, err := websocket.NewWebSocket(w, r)
	if err != nil {
		log.Printf("Error creating websocket connection: %v\n", err)
		return
	}

	client := &websocket.Client{
		ID:       uuid.New(),
		Conn:     ws.Conn,
		Pool:     pool,
		Ws: 	  ws,
	}

	ws.On("register", func(event *websocket.Event) {

		client = &websocket.Client{
			ID:       uuid.New(),
			Username: event.Data.(string),
			Conn:     ws.Conn,
			Pool:     pool,
			Ws: 	  ws,
		}

		pool.Register <- client
		//go client.Read()

		ws.Out <- (&websocket.Event{
			Event: "registered",
			Data:  client.ID,
		}).Raw()
	})

	ws.On("message", func(event *websocket.Event) {
		fmt.Printf("Message received: %s\n", event.Data.(string))

		/*ws.Out <- (&websocket.Event{
			Event: "response",
			Data: event.Data,
		}).Raw()*/

		if client.Username != "" {
			message := websocket.NewMessage{
				ClientID:       client.ID,
				ClientUsername: client.Username,
				Message:        event.Data.(string),
			}

			pool.Broadcast <- message
		}
	})
}


var router = mux.NewRouter()

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	router.Path("/ws").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		wsHandler(pool, writer, request)
	})
}

func main() {
	fmt.Println("Websocket Test")
	setupRoutes()
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		return 
	}
}