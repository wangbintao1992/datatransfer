package main

import (
	"net"
	"fmt"
	"os"
	"io/util"
	"sort"
	"time"
	"bufio"
	"io"
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
	time.Sleep(10e9)
	fmt.Println("finash")
}
func mergeFile() {
	fmt.Println("merge ...")
	index := make([]util.Head,0)
	for k,v := range indexFile {
		fmt.Println("merge index:",k," file order：",v.Order)
		index = append(index, v)
	}

	sort.Sort(util.HeadIndex(index))

	writer := getFinalFile(index[0])

	for _, v := range index {
		reader := getReader(v.Path)
		io.Copy(writer,reader)
		writer.Flush()
	}

	fmt.Println("merge finash...")
}
func getReader(s string) *bufio.Reader{
	file, _ := os.Open(s)
	return bufio.NewReader(file)
}
func getFinalFile(h util.Head) *bufio.Writer{
	f, _ := os.OpenFile(h.GetFilePath(path), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	return bufio.NewWriter(f)
}
func clean(){
	for k,_ := range indexFile{
		fmt.Println("delete index:",k)
	}
}

func deleteFile(f os.File){

}

func writeBody(conn net.Conn,h *util.Head) {
	var length = h.Length
	var path = h.Path

	write := getWrite(path)

	var part []byte
	var readNum int
	for {
		part,readNum,length = getPartBody(length, conn)

		if readNum == 0{
			break
		}

		write.Write(part)
		write.Flush()
	}
}
func getWrite(s string) *bufio.Writer{
	file, _ := os.OpenFile(s,os.O_CREATE | os.O_RDWR|os.O_APPEND, 0666)
	return bufio.NewWriter(file)
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
		//TODO continuingly-transferring
		fmt.Println("重复读 err")
		head.Path = h.Path
	} else {
		head.Path = head.GetTmpFilePath(path)
		indexFile[head.Hash] = *head
	}
}