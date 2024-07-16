package main

import (
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
		log.Print(err)
		return
	}

	for {
		// Read message from client
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Print(err)
			return
		}

		// Print the message to the console
		log.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		// Get a reply from the database
		reply := database.HandleMsg(msg)

		// Write message back to client
		if err = conn.WriteMessage(1, reply); err != nil {
			log.Print(err)
			return
		}
	}
}

func acceptIncoming(newConnections chan net.Conn) {
	port := ":8080"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Couldn't listen on port ", port, err)
	}
	for {
		if conn, err := listener.Accept(); err != nil {
			log.Println("Couldn't accept new incoming connection: ", err)
		} else {
			newConnections <- conn
		}
	}
}

func dialAddresses(pathToAddresses string, newConnections chan net.Conn) {
	localAddr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		log.Fatal("Couldn't get localhost address: ", err)
	}

	//  FIXME: read address book from file in the future
	addressBook := make([]*net.TCPAddr, 0)
	for _, addr := range addressBook {
		// Dial address
		if conn, err := net.DialTCP("tcp", localAddr, addr); err != nil {
			log.Println("couldn't dial address ", addr, ": ", err)
		} else {
			newConnections <- conn
		}
	}
}

func listen(conn net.Conn, messages chan []byte) {
	buffer := make([]byte, 1024)
	for {
		if _, err := conn.Read(buffer); err != nil {
			log.Println("couldn't read connection: ", err)
		} else {
			messages <- buffer
		}
	}
}

func main() {
	// http server for serving data for frontend javascript rendering
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", handleFrontend)
	go log.Fatal(http.ListenAndServe(":8080", nil))

	// peer network

	newConnections := make(chan net.Conn)
	in := make(chan []byte)
	out := make(chan []byte)

	go acceptIncoming(newConnections)

	pathToAddressBook := "my/path/"
	dialAddresses(pathToAddressBook, newConnections)

	//  NOTE: write go routine for handling new incoming messages
	// needs to be a value telling us if they are wanting us to serve a file or telling us to update a file???
	// or does this need to be thought out a bit more?
	go func() {
		for msg := range <-in {

		}
	}()

	conns := make([]net.Conn, 0)
	for {
		select {
		case conn := <-newConnections:
			conns = append(conns, conn)
			go listen(conn, in)
		case msg := <-out:
			for _, conn := range conns {
				conn.Write(msg)
			}
		}
	}

}
