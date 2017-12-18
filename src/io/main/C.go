package main

import (
	"net"
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
	"os"
	"flag"
	"encoding/json"
	"io/util"
	"strconv"
)

func main() {
	flag.Parse()

	path := flag.Arg(0)
	fmt.Println("inputPath:",path)
	conn, _ := net.Dial("tcp", "localhost:8080")

	//TODO pool
	//TODO asych timeout
	sendData(path,conn)
	time.Sleep(30e9)
}



func sendData(path string, conn net.Conn){
	file, _ := os.Open(path)

	defer file.Close()

	writePacket(file,conn)

}
func writePacket(file *os.File,conn net.Conn) {
	buf := make([]byte, 30000)
	fmt.Println("write start")

	order := 0
	for{
		n, _ := file.Read(buf)
		if n == 0{
			break
		}

		head := getHead(order, int32(n),key,"data.pdf")

		packet := getPacket(head, buf[0:n])

		w , _ := conn.Write(packet)
		fmt.Println("write:",w)
		fmt.Println("writing...")
	}

	fmt.Println("write end")

}
func getHead(o int, l int32, h string, n string) []byte {
	head := &util.Head{
		Order:  o,
		Length: l,
		Name:   n,
		Hash:   h}
	marshal, e := json.Marshal(head)

	if e != nil{
		fmt.Println(e)
	}

	fmt.Println("head json",string(marshal))

	if space := util.HeadSize - len(marshal); space > 0{
		marshal = append(marshal,util.GetSpace(space)...)
	}

	return marshal
}
func getPacket(head []byte, body []byte) []byte{
	return append(head, body...)
}


//[ [x] [x,x,x,x]]


func int32ToByte(l int32) []byte{
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer,binary.BigEndian,l)
	return buffer.Bytes()
}