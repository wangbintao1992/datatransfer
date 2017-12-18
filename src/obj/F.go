package main

import (
	"fmt"
	"bytes"
	"encoding/binary"
	"net"
)

func main() {
	conn, _ := net.Listen("tcp", "localhost:8080")

	for{
		accept, _ := conn.Accept()
		fmt.Println("come in")
		go read2(accept)
	}
}
func read2(accept net.Conn) {
	i := make([]byte, 1000)
	accept.Read(i)
	fmt.Println(string(i))

	accept.Read(i)
	fmt.Println(string(i))

	accept.Read(i)
	fmt.Println(string(i))
}

func intToByte(d int32) []byte{
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, d)

	return buffer.Bytes()
}