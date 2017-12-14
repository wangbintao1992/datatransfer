package main

import (
	"fmt"
	"bytes"
	"encoding/binary"
	"net"
)

func main() {
	conn, _ := net.Listen("tcp", "localhost:8080")

	accept, _ := conn.Accept()

	i := make([]byte, 1)

	var t = 0
	for{
		if t == 2{
			break
		}
		accept.Read(i)
		fmt.Println(i[0])
		t ++

		}
}


func (f *F) Say(){
	fmt.Println(f.Name)
}

type F struct {
	Name string
}

type C struct {
	*F
}

func intToByte(d int32) []byte{
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, d)

	return buffer.Bytes()
}