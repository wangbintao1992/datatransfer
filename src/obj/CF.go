package main

import (
	"os"
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {

	var cdl sync.WaitGroup


	go t(cdl)
	go t(cdl)
	go t(cdl)
	go t(cdl)

	cdl.Wait()
	fmt.Println("asd")
}
func t(group sync.WaitGroup) {
	group.Add(1)
	time.Sleep(2000)
	fmt.Println("zzz")
	group.Done()
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