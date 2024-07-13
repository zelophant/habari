package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/zelophant/habari/database"

	"github.com/gorilla/websocket"
)

func handleFrontend(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

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
}

func handleNetworkConnection(conn net.Conn) {
	conn.SetReadDeadline()
	buffer := make([]byte, 100)
	conn.Read(buffer)
}

func main() {
	// http server for serving data for frontend javascript rendering
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", handleFrontend)
	go log.Fatal(http.ListenAndServe(":8080", nil))

	// peer network
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go handleNetworkConnection(conn)
	}
}
