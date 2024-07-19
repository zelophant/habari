package main

import (
	"log"
	"net"
	"net/http"

	"github.com/zelophant/habari/pkg"
)

func main() {
	// http server for local rendering
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", habari.HandleClient)
	go log.Fatal(http.ListenAndServe(":8080", nil))

	// peer network

	bytesFromPeer := make(chan habari.Content)

	// temporary just to test if can write to connection
	var content = habari.Content{
		Bytes: []byte("hello"),
	}

	pathToAddressBook := "my/path/"
	// NOTE: insert some code to extract addresses from addressbook
	for _, peer := range habari.DialAddressesForPeers(pathToAddressBook) {
		habari.ListenToPeer(peer, bytesFromPeer)

		// temporary test
		content.ID = peer
	}

	for {
		select {
		case peer := <-habari.ListenForNewPeers(":8080"):
			habari.ListenToPeer(peer, bytesFromPeer)
		case content := <-bytesFromPeer:
			// currently just send the bytes to the client
			habari.BytesToClient <- content.Bytes // content also includes an ID field of the net.Conn sending the data
		case bytes := <-habari.BytesFromClient:
			// currently just send the bytes to last peer to send message
			content.ID.Write(bytes)
		}
	}
}
