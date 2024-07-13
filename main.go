package main

import (
	"fmt"
	"github.com/zelophant/habari/database"
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
			return
		}

		for {
			// Read message from client
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Print(err)
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			// Get a reply from the database
			reply := database.HandleMsg(msg)

			// Write message back to client
			if err = conn.WriteMessage(1, reply); err != nil {
				fmt.Print(err)
				return
			}
		}
	})

	http.Handle("/", http.FileServer(http.Dir("./src")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
