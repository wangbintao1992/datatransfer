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

	defer conn.Close()
	//TODO when stop?
	time.Sleep(30e9)
}


func sendData(path string, conn net.Conn){
	file, _ := os.Open(path)

	defer file.Close()

	//TODO blockSize?
	blocks := getBlocks(path, 300000)

	startCalcTime()

	fmt.Println("write start")
	for i := 0; i < len(blocks); i ++ {
		msg := getPacket(file, blocks[i])
		n, _ := conn.Write(msg)

		fmt.Println("write:", n, "write...")
	}
	endCalcTime()
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
