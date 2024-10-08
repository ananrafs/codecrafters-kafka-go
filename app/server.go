package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	buff := make([]byte, 1024)
	_, err = conn.Read(buff)
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
	fmt.Println(string(buff), "buffers")
	fmt.Printf("Receiving buff : %v (%d) | key version : %d : %d",
		buff[8:12],
		int32(binary.BigEndian.Uint32(buff[8:12])),
		int32(binary.BigEndian.Uint16(buff[0:4])),
		int32(binary.BigEndian.Uint16(buff[4:8])))

	correlationID := buff[8:12]
	version := buff[4:8]
	if int(binary.BigEndian.Uint16(version)) > 4 || int(binary.BigEndian.Uint16(version)) < 0 {
		version = []byte{0, 35}
	}
	// apiKey := buff[0:4]
	res := make([]byte, len(buff))
	copy(res[4:], correlationID)
	copy(res[8:], version)
	conn.Write(res)
}
