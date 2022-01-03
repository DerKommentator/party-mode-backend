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

	ws.On("register", func(event *websocket.Event) {

		ws.Client = &websocket.Client{
			ID:       uuid.New(),
			Username: event.Data.(string),
			Conn:     ws.Conn,
			Pool:     pool,
			Ws:       ws,
		}

		pool.Register <- ws.Client
		//go client.Read()

		ws.Out <- (&websocket.Event{
			Event: "registered",
			Data:  ws.Client.ID,
		}).Raw()
	})

	ws.On("message", func(event *websocket.Event) {
		fmt.Printf("Message received: %s\n", event.Data.(string))

		/*ws.Out <- (&websocket.Event{
			Event: "response",
			Data: event.Data,
		}).Raw()*/

		if ws.Client.Username != "" {
			message := websocket.NewMessage{
				ClientID:       ws.Client.ID,
				ClientUsername: ws.Client.Username,
				Message:        event.Data.(string),
			}

			pool.Broadcast <- message
		}
	})

	ws.On("videoSyncProcess", func(event *websocket.Event) {
		fmt.Printf("videoSyncProcess received: %s\n", event.Data)

		/*ws.Out <- (&websocket.Event{
			Event: "response",
			Data: event.Data,
		}).Raw()*/

		/*if ws.Client.Username != "" {
			message := websocket.NewMessage{
				ClientID:       ws.Client.ID,
				ClientUsername: ws.Client.Username,
				Message:        event.Data.(string),
			}

			pool.Broadcast <- message
		}*/
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
