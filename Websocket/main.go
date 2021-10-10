package main

import (
	"Bergflix-Partymode/pkg/websocket"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)


func wsHandler(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}

	client := &websocket.Client {
		ID: uuid.New(),
		Username: "",
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
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