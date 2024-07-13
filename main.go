package main

import (
	"log"
	"net"
	"net/http"
	"sync"

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

func handleNetworkConnection(conn net.Conn) {
	// Should set a deadline with conn.SetReadDeadline()
	buffer := make([]byte, 100)
	conn.Read(buffer)
}

func main() {
	// http server for serving data for frontend javascript rendering
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", handleFrontend)
	go log.Fatal(http.ListenAndServe(":8080", nil))

	// peer network
	addr := &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 8000,
		Zone: "",
	}

	// Incoming calls are made by peers who have my address in their "address book"

	// Currently storing all such connections in an array
	// This is bad because ideally each connection is handled in its own goroutine
	// each routine would write the connection's messages to a channel
	// then all messages would be processed in a separate routine which reads from that channel.
	incomingCalls := struct {
		mu    sync.Mutex
		conns []*net.UDPConn
	}{
		mu:    sync.Mutex{},
		conns: make([]*net.UDPConn, 0),
	}

	go func() {
		for {
			conn, err := net.ListenUDP("udp", addr)
			if err != nil {
				return
			}
			incomingCalls.mu.Lock()
			incomingCalls.conns = append(incomingCalls.conns, conn)
			incomingCalls.mu.Unlock()
		}
	}()

	// Outgoing calls made by looking up addresses in a locally stored "address book"
	outgoingCalls := make([]net.Conn, 0)

	//  FIXME: read address book from file in the future
	addressBook := make([]*net.UDPAddr, 0)

	for _, val := range addressBook {
		// Dial to the address with UDP
		conn, err := net.DialUDP("udp", nil, val)
		if err != nil {
			log.Println(err)
		}
		outgoingCalls = append(outgoingCalls, conn)
	}

	// Stupid bad idea ranging over all connections sequentially
	for {
		for i, v := range outgoingCalls {
			// do something
		}
		incomingCalls.mu.Lock()
		for i, v := range incomingCalls.conns {
			// do something
		}
		incomingCalls.mu.Unlock()
	}
}
