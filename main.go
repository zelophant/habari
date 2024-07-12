package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Print(err)
		}

		for {
			// Read message from client
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			// Write message back to client
			if err = conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})

	http.Handle("/", http.FileServer(http.Dir("./src")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
