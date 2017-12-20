package main

import (
	"os"
	"fmt"
	"net"
	"compress/gzip"
	"bytes"
	"io/ioutil"
)

func main() {

	i := make([]byte, 5)

	i[0] = 1
	i[1] = 2
	i[2] = 3
	i[3] = 4
	i[4] = 5
	fmt.Println(i)
		fmt.Println("raw size:", len(i))

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()
	w.Write(i)
	w.Flush()
	fmt.Println("gzip size:", len(b.Bytes()))

	r, _ := gzip.NewReader(&b)
	defer r.Close()
	undatas, _ := ioutil.ReadAll(r)
	fmt.Println("ungzip size:", len(undatas),b.Bytes())

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