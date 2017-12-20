package main

import (
	"os"
	"fmt"
	"net"
	"bufio"
)

func main() {
	bytes := make([]byte, 10)
	file, _ := os.OpenFile("D://tmp.log",os.O_CREATE | os.O_RDWR|os.O_APPEND, 0666)

	writer := bufio.NewWriter(file)

	writer.
	file.Read(bytes)

	fmt.Println(1,bytes)
	fmt.Println(2,string(bytes))

	file.Write([]byte("why why why"))
	file.Sync()

	file.Write([]byte("222222"))
	file.Sync()

	bytes2 := make([]byte, 10)
	_, err := file.Read(bytes2)
	fmt.Println(err)
	fmt.Println(3,bytes2)
	fmt.Println(4,string(bytes2))
}
func readPart(offset int64, file *os.File,conn net.Conn) {
	fmt.Println("=========")
	buf := make([]byte, 100)
	file.ReadAt(buf, offset)

	//conn.Write(buf)
	fmt.Println(string(buf),"===")
}

type Head struct {
	*os.File
}