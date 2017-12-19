package main

import (
	"os"
	"fmt"
	"net"
	"io/ioutil"
	"io"
)

func main() {
	file, _ := os.OpenFile("", 1, 0666)
	head := &Head{F: *file}
	fmt.Println(head)
	ioutil.ReadAll(&head.F)

	io.Copy()
}
func readPart(offset int64, file *os.File,conn net.Conn) {
	fmt.Println("=========")
	buf := make([]byte, 100)
	file.ReadAt(buf, offset)

	//conn.Write(buf)
	fmt.Println(string(buf),"===")
}

type Head struct {
	F os.File
}