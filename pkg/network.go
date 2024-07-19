package habari

import (
	"log"
	"net"
)

func ListenForNewPeers(port string) chan net.Conn {
	newConns := make(chan net.Conn)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Couldn't listen on port ", port, err)
	}
	go func() {
		for {
			if conn, err := listener.Accept(); err != nil {
				log.Println("Couldn't accept new incoming connection: ", err)
			} else {
				newConns <- conn
			}
		}
	}()
	return newConns
}

func DialAddressesForPeers(pathToAddresses string) []net.Conn {
	newConns := make([]net.Conn, 0)
	localAddr, err := net.ResolveTCPAddr("tcp", ":3000")
	if err != nil {
		log.Fatal("Couldn't get local address: ", err)
	}

	//  FIXME: read address book from file in the future
	addressBook := make([]*net.TCPAddr, 0)
	for _, addr := range addressBook {
		// Dial address
		if conn, err := net.DialTCP("tcp", localAddr, addr); err != nil {
			log.Println("couldn't dial address ", addr, ": ", err)
		} else {
			newConns = append(newConns, conn)
		}
	}
	return newConns
}

// ListenToPeer listens on a peer connection and sends bytes on the returned channel
func ListenToPeer(conn net.Conn, out chan<- Content) {
	bytes := make([]byte, 1024)
	go func() {
		for {
			if _, err := conn.Read(bytes); err != nil {
				log.Println("couldn't read connection: ", err)
			} else {
				out <- Content{conn, bytes}
			}
		}
	}()
}

type Content struct {
	ID    net.Conn
	Bytes []byte
}
