package main

import (
	"net"
	"fmt"
	"os"
	"util"
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
	var lock sync.WaitGroup

	tcpAddr, err := net.ResolveTCPAddr("tcp4", "localhost:8080")

	fmt.Println(err)

	server, e := net.ListenTCP("tcp", tcpAddr)
	fmt.Println(e)
	for {
		conn, _ := server.AcceptTCP()
		util.SetTCPOption(conn)
		go handle(conn,lock)
	}

	//TODO md5 check
	lock.Wait()
	mergeFile()
	clean()
	fmt.Println("finash")
}
func handle(accept net.Conn,lock sync.WaitGroup) {

	fmt.Println("reuqest come in")

	readPacket(accept,lock)
}

func readPacket(conn net.Conn,lock sync.WaitGroup) {
	fmt.Println("start receive")
	//TODO fix buf @link length
	//TODO p
	lock.Add(1)
	fmt.Println("receiveing ...")
	head, _ := readHead(conn)

	once.Do(func() {
		tmpPath = path.Join(p,"\\",head.Name + "tmp")
		os.Mkdir(tmpPath,0666)
	})

	//TODO index to large ?
	addToIndex(head)

	fmt.Println("body length",head.Length)

	writeBody(conn, head)

	lock.Done()
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

	for _, v := range index {
		//!important reset read pointer
		v.Path.Seek(0,io.SeekStart)
		reader := bufio.NewReader(v.Path)
		io.Copy(bufWrite,reader)
		bufWrite.Flush()

		v.Path.Close()
	}

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
}

func writeBody(conn net.Conn,h *util.Head) {
	var length = h.Length
	var file = h.Path

	part := getPartBody(length, conn)

	file.Write(util.GzipDecode(part))
	file.Sync()
}

func getPartBody(length int, conn net.Conn) ([]byte){
	buf := make([]byte, length)

	num, _ := conn.Read(buf)

	fmt.Println("read num:", num)

	//TODO recover
	if length != num{
		panic("数据包异常")
	}

	return buf[:num]
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
		file, _ := os.OpenFile(head.GetTmpFilePath(p),os.O_CREATE | os.O_RDWR|os.O_APPEND, 0666)
		head.Path = file
		indexFile[head.Hash] = *head
	}
}