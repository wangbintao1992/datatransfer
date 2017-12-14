package main

import (
	"net"
	"fmt"
	"bytes"
	"encoding/binary"
	"os"
)
var path = "D://tmp"
func main() {
	conn, e := net.Listen("tcp", "localhost:8080")
	fmt.Println(e)

	for{
		accept, _ := conn.Accept()
		handle(accept)
	}
}
func handle(accept net.Conn) {
	t := getType(accept)
	fmt.Println("reuqest come in type:",t)

	if t == 1{
		writeHead(accept)
	} else {
		writeBody(accept)
	}
}
//[ [x] [x,x,x,x]]
// type length (max int32)
func writeBody(conn net.Conn) {
	fmt.Println("receiveing ...")
	//TODO fix buf @link length

	//TODO path
	path += "//data.pdf"
	file, _ := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	defer file.Close()
	for {
		buf := make([]byte, 10240)
		num, _ := conn.Read(buf)
		if num == 0{
			break
		}

		file.Write(buf[0:num])
		file.Sync()
	}
	//flush
	//TODO check
	fmt.Println("finsh")
}
func writeHead(conn net.Conn) {
	th := make([]byte, 4)
	conn.Read(th)
	i := byteToInt32(th)
	fmt.Println("length:",i)
}

func getType(accept net.Conn) byte{
	t := make([]byte,1)
	_, e := accept.Read(t)

	if e != nil{
		fmt.Println("read type err",e)
	}

	return t[0]
}

func byteToInt32(d []byte) int32{
	var r int32 = 0
	buffer := bytes.NewBuffer(d)
	binary.Read(buffer,binary.BigEndian,&r)

	return r
}