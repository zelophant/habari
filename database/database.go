package database

import "fmt"

func HandleMsg(msg []byte) []byte {
	fmt.Println(msg)
	return []byte("some reply")
}
