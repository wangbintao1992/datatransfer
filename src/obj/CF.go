package main

import (
	"os"
	"fmt"
	"net"
)

func main() {


}
func readPart(offset int64, file *os.File,conn net.Conn) {
	fmt.Println("=========")
	buf := make([]byte, 100)
	file.ReadAt(buf, offset)

	//conn.Write(buf)
	fmt.Println(string(buf),"===")
}
