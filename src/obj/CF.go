package main

import "net"

func main() {
	conn, _ := net.Dial("tcp", "localhost:8080")

	bytes := make([]byte, 2)

	bytes[0] = 22
	bytes[1] = 99
	conn.Write(bytes)
}
