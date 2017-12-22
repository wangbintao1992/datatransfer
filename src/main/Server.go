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
var indexCahce *util.IndexCache
var once sync.Once
var tmpPath string

func init()  {
	fmt.Println("初始化...")
	//indexPath := path.Join(p,"index.log")
	//file, _ := os.OpenFile(indexPath,os.O_CREATE | os.O_RDWR|os.O_APPEND, 0666)
	//
	//loadToCache(file)
	indexCahce = util.Init(path.Join(p,"index.log"))
}
func loadToCache(file *os.File) {
	
}

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

		saveCache()
	}

	fmt.Println("finash")
}
func isContinuTransferr() bool{
	return false
}
func saveCache() {
	indexCahce.SaveCache()
}
func handle(accept net.Conn) {

	fmt.Println("reuqest come in")

	readPacket(accept)
}

//不同协程 并行读，所以每读一包，包含包头包体
func readPacket(conn net.Conn) {
	fmt.Println("start receive")
	//TODO fix buf @link length
	//TODO p
	fmt.Println("receiveing ...")
	for{
		head, n := readHead(conn)

		if n == 0{
			break
		}

		once.Do(func() {
			tmpPath = path.Join(p,"\\",head.Name + "tmp")
			os.Mkdir(tmpPath,0666)
		})

		//TODO index to large ?
		addToIndex(head)

		fmt.Println("body length",head.Length)

		readAndWriteBody(conn, head)
	}
}
func readAndWriteBody(conn net.Conn,h *util.Head) {
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
	//TODO bufio write 写 0字节
	writer := bufio.NewWriter(file)
	writer.Write(part)
	writer.Flush()
}
func getPartBody(length int, conn net.Conn) ([]byte,int,int){
	buf := make([]byte, length)

	num, _ := bufio.NewReader(conn).Read(buf)

	fmt.Println("read num:", num)

	remain := 0
	if remain = length - num; remain != 0{
		fmt.Println("resize buf",remain)
	}

	return buf[:num],num,remain
}
func mergeFile() {
	if indexCahce.CheckMerge(){
		fmt.Println("数据包缺失")
		fmt.Println("缺失:")
		return
	}

	fmt.Println("merge ...")

	index := indexCahce.GetCache()

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
func checkMerge() {

}
func getFinalFile(h util.Head) *os.File{
	return util.GetRW(h.GetFilePath(p))
}
func clean(){
	fmt.Println("清理临时文件")
	r := os.RemoveAll(tmpPath)
	if r != nil{
		fmt.Println(r)
	}
	fmt.Println("清理临时文件结束")
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

	if h ,ok:= indexCahce.Get(head.Hash); ok{
		fmt.Println("重复读 err")
		head.Path = h.Path
	} else {
		head.Path = util.GetRW(head.GetTmpFilePath(p))
		indexCahce.Put(head.Hash,*head)
	}
}
