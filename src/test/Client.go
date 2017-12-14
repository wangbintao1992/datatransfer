package main

import (
	"fmt"
	"flag"
	"os"
	"strconv"
	"net"
	"encoding/json"
)

func main() {
	fmt.Println("client start")
	flag.Parse()
	path := os.Args[1]
	packet := os.Args[2]

	fmt.Println("分块大小byte:" + packet)
	fmt.Println("源文件路径:" + path)

	info, e := os.Stat(path)

	if e != nil{
		panic(e)
	}
	size := info.Size()
	fmt.Println("源文件大小:" + strconv.FormatInt(size,10) + " unit:byte")

	packet2, _ := strconv.ParseInt(packet, 10, 64)
	blockNum, remain := divFile(size, path, packet2)

	fmt.Println("分块:" + strconv.FormatInt(blockNum,10))
	fmt.Println("剩余:" + strconv.FormatInt(remain,10))

	var blockNum2 = int(blockNum)
	var remain2 = int(remain)
	packet3, _ := strconv.Atoi(packet)
	for i := 0 ; i < blockNum2; i ++{
		tmp := &Packet{Offset: i * packet3, length: packet3}
		sendData(tmp,path,packet3)
	}

	if remain2 != 0{
		t := &Packet{Offset: blockNum2 * packet3, length: remain2}
		sendData(t,path,packet3)
	}
}
func sendData(packet *Packet,path string,packetSize int) {
	fmt.Println("分块读开始")
	file, e := os.Open(path)
	fmt.Println(e)

	buf := make([]byte, packetSize)
	file.Read(buf)

	packet.data = buf
	conn, e := net.Dial("tcp", "localhost:8080")
	fmt.Println(e)

	marshal, i := json.Marshal(packet)
	fmt.Println(marshal)
	fmt.Println(i)
	conn.Write(buf)
	fmt.Println("分块写结束")
}

func divFile(size int64,path string,packet int64)  (int64,int64){
	blockNum := size / packet
	mod := size % packet

	var	remain int64
	if mod != 0{
		remain = size - blockNum * packet

	}

	return blockNum,remain
}

type Packet struct {
	Offset int
	length int
	data []byte
}