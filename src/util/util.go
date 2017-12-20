package util

import (
	"bytes"
	"encoding/binary"
	"net"
)

var(
	HeadSize  = 150
)

//TODO deprecated
func GetSpace(space int) []byte{
	t := make([]byte, space)
	for i := 0;i < space; i ++{
		t[i] = 0
	}

	return t
}
//TODO buf size
func SetTCPOption(conn net.Conn) {
	tcpConn := conn.(*net.TCPConn)
	tcpConn.SetNoDelay(false)
	tcpConn.SetWriteBuffer(300000)
	tcpConn.SetReadBuffer(300000)
}

func ByteToInt32(d []byte) int32{
	var r int32 = 0
	buffer := bytes.NewBuffer(d)
	binary.Read(buffer,binary.BigEndian,&r)

	return r
}

func int32ToByte(l int32) []byte{
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer,binary.BigEndian,l)
	return buffer.Bytes()
}
