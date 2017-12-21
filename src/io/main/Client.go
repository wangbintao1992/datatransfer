package main

import (
	"net"
	"fmt"
	"time"
	"os"
	"flag"
	"io/util"
)

var startTime time.Time
func main() {
	flag.Parse()

	path := flag.Arg(0)
	fmt.Println("inputPath:",path)
	conn, _ := net.Dial("tcp", "localhost:8080")

	util.SetTCPOption(conn)

	//TODO pool
	//TODO asych timeout
	sendData(path, conn)

	time.Sleep(time.Second * 3)
	conn.Close()
}


func sendData(path string, conn net.Conn){
	file, _ := os.Open(path)

	//TODO blockSize?
	blocks := getBlocks(path, 300000)

	startCalcTime()

	fmt.Println("write start")

	//TODO 非常重要，待解决，socket 缓冲区可能第一次发包，不满

	for i := 0; i < len(blocks); i ++ {
		go writePacket(file, blocks, i, conn)
	}
	endCalcTime()
}
func writePacket(file *os.File, blocks []util.Block, i int, conn net.Conn) {
	msg := getPacket(file, blocks[i])
	n, e := conn.Write(msg)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println("write:", n, "write...")
}
//TODO common
func endCalcTime() {
	since := time.Since(startTime)
	fmt.Println("用时 s", since)
}
//TODO common
func startCalcTime() {
	startTime = time.Now()
}

func getBlocks(path string,blockSize int64) []util.Block{
	return util.GetBlockArr(path,blockSize)
}

func getPacket(file *os.File,b util.Block) []byte{
	buf := make([]byte, b.Blength)

	n, e := readFile(file, buf,b)

	if e != nil{
		fmt.Println(e)
	}
	head := getHead(b)

	if head != nil{}
	return mergePacket(head, buf[:n])
}

func readFile(file *os.File,buf []byte,block util.Block)(int,error) {
	return file.ReadAt(buf,int64(block.Offset))
}
func getHead(b util.Block) []byte {
	return util.EncodeHead(b)
}
func mergePacket(head []byte, body []byte) []byte{
	return append(head, body...)
}
