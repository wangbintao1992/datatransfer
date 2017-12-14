package main

import (
	"net"
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
	"os"
)

func main() {
	conn, _ := net.Dial("tcp", "localhost:8080")

	//TODO pool
	//TODO asych timeout
	sendData(12, 122, conn)
	time.Sleep(30e9)
}



func sendData(order int, lenght int32, conn net.Conn){
	head := getHead(order, lenght)
	conn.Write(head)
	writeBody2(nil, order,conn)

}

func writeBody2(data []byte,order int,conn net.Conn) []byte{
	var path = "D:\\go.pdf"

	buf := make([]byte, 30000)
	file, _ := os.Open(path)

	conn.Write([]byte{2})
	for{
		n, _ := file.Read(buf)

		if n ==0 {
			break
		}

		conn.Write(buf[0:n])
	}

	return nil
}
//[ [x] [x,x,x,x]]
func getHead(order int,l int32) []byte{
	//1 head 2 body
	var t byte = 1
	head := make([]byte,1)
	head[0] = t
	i := append(head, int32ToByte(l)...)
	fmt.Println("head",i)
	return i
	//TODO body length
}

func int32ToByte(l int32) []byte{
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer,binary.BigEndian,l)
	return buffer.Bytes()
}