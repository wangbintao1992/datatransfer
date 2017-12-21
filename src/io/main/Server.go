package main

import (
	"net"
	"fmt"
	"os"
	"io/util"
	"sort"
	"bufio"
	"io"
	"sync"
	"path"
)
var p = "D://tmp"
var indexFile = make(map[string]util.Head,10)
var once sync.Once
var tmpPath string
func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "localhost:8080")

	fmt.Println(err)

	server, e := net.ListenTCP("tcp", tcpAddr)
	fmt.Println(e)
	for {
		conn, _ := server.AcceptTCP()
		util.SetTCPOption(conn)
		handle(conn)

		mergeFile()

		clean()
	}

	//TODO md5 check
	//TODO wait

	fmt.Println("finash")
}
func handle(accept net.Conn) {

	fmt.Println("reuqest come in")

	readPacket(accept)
}
//不同协程 并行读，所以每读一包，包含包头包体
func readPacket(conn net.Conn) {

	for {
		fmt.Println("receiveing ...")
		head, n := readHead(conn)

		once.Do(func() {
			tmpPath = path.Join(p,"\\",head.Name + "tmp")
			os.Mkdir(tmpPath,0666)
		})

		if n == 0{
			break
		}

		//TODO index to large ?
		addToIndex(head)

		fmt.Println("body length",head.Length)

		writeBody(conn, head)
	}
}
func writeBody(conn net.Conn,h *util.Head) {
	var length = h.Length
	var file = h.Path

	tmpBuf := make([]byte, 0)
	for{
		num := 0
		data,num,remain := getPartBody(length, conn)

		if num == 0{
			break
		}

		length = remain

		tmpBuf = append(tmpBuf, data...)
	}
	go flushToDisk(file, tmpBuf)
}
func flushToDisk(file *os.File, part []byte) {
	writer := bufio.NewWriter(file)
	writer.Write(part)
	writer.Flush()
}
func getPartBody(length int, conn net.Conn) ([]byte,int,int){
	buf := make([]byte, length)

	reader := bufio.NewReader(conn)

	num, _ := reader.Read(buf)

	fmt.Println("read num:", num)

	remain := 0
	//TODO recover
	if remain = length - num; remain != 0{
		fmt.Println("resize buf",remain)
	}

	return buf[:num],num,remain
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

	defer writer.Close()

	bufWrite := bufio.NewWriter(writer)

	var totalSize int64 = 0
	for _, v := range index {
		//!important reset read pointer
		v.Path.Seek(0,io.SeekStart)
		reader := bufio.NewReader(v.Path)
		written, _ := io.Copy(bufWrite, reader)
		fmt.Println("merge size:",written)
		totalSize += written
		bufWrite.Flush()
		v.Path.Close()
	}

	fmt.Println("totalSize:",totalSize)
	fmt.Println("merge finash...")
}
func getFinalFile(h util.Head) *os.File{
	f, _ := os.OpenFile(h.GetFilePath(p), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	return f
}
func clean(){
	fmt.Println("清理临时文件")
	r := os.RemoveAll(tmpPath)
	if r != nil{
		fmt.Println(r)
	}
	fmt.Println("清理临时文件结束")
	indexFile = nil
}


func readHead(conn net.Conn) (h *util.Head, num int){
	t := make([]byte,util.HeadSize)
	n, e := conn.Read(t)

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
		file, _ := os.OpenFile(head.GetTmpFilePath(p),os.O_CREATE | os.O_RDWR|os.O_APPEND, 0666)
		head.Path = file
		indexFile[head.Hash] = *head
	}
}