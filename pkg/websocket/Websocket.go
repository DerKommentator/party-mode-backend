package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader {
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocket struct {
	Conn *websocket.Conn
	Out chan []byte
	In chan []byte
	Events map[string]EventHandler
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}

func Reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func Writer(conn *websocket.Conn) {
	for {
		fmt.Println("Sending")
		messageType, r, err := conn.NextReader()
		if err != nil {
			log.Println(err)
			return
		}

		w, err := conn.NextWriter(messageType)
		if err != nil {
			log.Println(err)
			return
		}

		if _, err := io.Copy(w, r); err != nil {
			log.Println(err)
			return
		}

		if err := w.Close(); err != nil {
			log.Println(err)
			return
		}
	}
}