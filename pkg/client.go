package habari

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

/*
NOTE: There is no mechanism in place to guarantee that a websocket has been established before
the application starts sending and recieving on BytesFromClient and BytesToClient.
*/

// BytesFromClient is written to by habari.HandleClient when bytes are received on a websocket connection.
var BytesFromClient chan []byte = make(chan []byte)

// BytesToClient is recieved by habari.HandleClient which writes bytes to a websocket connection.
var BytesToClient chan []byte = make(chan []byte)

/*
HandleClient is passed to http.handleFunc
it establishes a websocket connection with the client
and handles that connection in a loop
*/
func HandleClient(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}

	bytes := []byte{}
	for {
		select {
		case bytes = <-BytesToClient:
			// Write message back to client
			if err = conn.WriteMessage(1, bytes); err != nil {
				log.Println("Couldn't write to websocket: ", err)
			}
		default:
			// Read message from client
			_, bytes, err := conn.ReadMessage()
			if err != nil {
				log.Print("Couldn't read from websocket: ", err)
			}
			// Print the message to the console
			log.Printf("%s sent: %s\n", conn.RemoteAddr(), string(bytes))
			BytesFromClient <- bytes
		}
	}
}
