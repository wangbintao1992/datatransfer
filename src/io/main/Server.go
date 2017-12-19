package main

import (
	"net"
	"fmt"
	"os"
	"io/util"
	"sort"
	"io"
	"io/ioutil"
	"bufio"
	"bytes"
)
var path = "D://tmp"
var bufSize = 20240
var indexFile = make(map[string]util.Head,10)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "localhost:8080")

	fmt.Println(err)

	server, e := net.ListenTCP("tcp", tcpAddr)

	fmt.Println(e)
	for {
		conn, _ := server.AcceptTCP()
		util.SetTCPOption(conn)
		handle(conn)
	}
}
func handle(accept net.Conn) {

	fmt.Println("reuqest come in")

	readPacket(accept)
}

func readPacket(conn net.Conn) {
	fmt.Println("start receive")
	//TODO fix buf @link length
	//TODO path

	for {
		fmt.Println("receiveing ...")
		head, n := readHead(conn)

		if n == 0{
			break
		}

		//TODO index to large ?
		addToIndex(head)

		fmt.Println("body length",head.Length)

		writeBody(conn, head)
	}
	//TODO check
	mergeFile()
	fmt.Println("finsh")
}
func mergeFile() {

	index := make([]util.Head, len(indexFile))
	for _,v := range indexFile {
		index = append(index, v)
	}

	sort.Sort(util.HeadIndex(index))

	f, _ := os.OpenFile(head.GetFilePath(path), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	write := bufio.NewWriter(f)
	for _,v := range index{
		io.Copy(write,&v.Path)
	}
}
func clean(){

	for k,v := range indexFile{
		fmt.Println("delete index:",k)
		go deleteFile(v.Path)
	}
}

func deleteFile(f os.File){

}

func writeBody(conn net.Conn,h *util.Head) {
	var length = h.Length
	var file = h.Path

	var part []byte
	var readNum int
	for {
		part,readNum,length = getPartBody(length, conn)

		if readNum == 0{
			break
		}

		file.Write(part)
		file.Sync()
	}
}
func getPartBody(length int, conn net.Conn) ([]byte,int,int){
	buf := make([]byte, length)
	num, _ := conn.Read(buf)
	fmt.Println("read num:", num)

	remain := length - num

	if remain != 0 {
		length = remain
		fmt.Println("resize buf:", remain)
	}

	return buf[:num],num,remain
}
func readHead(accept net.Conn) (h *util.Head, num int){
	t := make([]byte,util.HeadSize)
	n, e := accept.Read(t)

	if e != nil{
		fmt.Println("read type err",e)
		return nil,0
	}

	//TODO tmp index
	return util.DecodeHead(t),n
}
func addToIndex(head *util.Head){

	if h ,ok:= indexFile[head.Hash]; ok{
		head.Path = h.Path
	} else {
		f2, _ := os.OpenFile(head.GetTmpFilePath(path), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		head.Path = *f2
		indexFile[head.Hash] = *head
	}
}