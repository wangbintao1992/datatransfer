package main

import (
	"net"
	"fmt"
	"bytes"
	"os"
	"encoding/json"
	"io/util"
)
var path = "D://tmp"
var bufSize int32 = 20240
var indexFile = make(map[string]util.Head,10)

func main() {
	conn, e := net.Listen("tcp", "localhost:8080")
	fmt.Println(e)

	for{
		accept, _ := conn.Accept()
		handle(accept)
	}
}
func handle(accept net.Conn) {

	fmt.Println("reuqest come in")

	readPacket(accept)
}
//head[x,x,x,x][data]
// type length (max int32)
func readPacket(conn net.Conn) {
	fmt.Println("start receive")
	//TODO fix buf @link length
	//TODO path
	path += "//data.pdf"

	for {
		fmt.Println("receiveing ...")
		head, hnum := readHead(conn)

		if hnum == 0{
			break
		}

		getFile(head)

		fmt.Println("body length",head.Length)

		writeBody(conn, head)
	}
	//TODO check

	fmt.Println("finsh")
}
func getFile(head *util.Head){

	if h ,ok:= indexFile[head.Hash]; ok{
		head.Path = h.Path
	} else {
		f2, _ := os.OpenFile(head.GetTmpFilePath(path), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		head.Path = *f2
		indexFile[head.Hash] = *head
	}
}
func writeBody(conn net.Conn,h *util.Head) {
	var length = h.Length
	var file = h.Path
	for {
		buf := make([]byte, bufSize)
		num, _ := conn.Read(buf)

		fmt.Println("read num:", num)
		if num == 0 {
			break
		}

		file.Write(buf[:num])
		file.Sync()

		if remain := length - int32(num); remain != 0 {
			bufSize = remain
			length = remain
			fmt.Println("resize buf:", remain)
		} else {
			break
		}
	}
}

func readHead(accept net.Conn) (h *util.Head, num int32){
	t := make([]byte,util.HeadSize)
	n, e := accept.Read(t)

	if e != nil{
		fmt.Println("read type err",e)
	}

	i := new(util.Head)

	right := bytes.TrimRight(t, "\x00")

	fmt.Println("head json",string(right))
	e2 := json.Unmarshal(right, &i)

	if e2 != nil{
		fmt.Println(e2)
	}

	//TODO tmp index

	return i,int32(n)
}