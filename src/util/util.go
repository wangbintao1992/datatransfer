package util

import (
	"bytes"
	"encoding/binary"
	"net"
	"compress/gzip"
	"io/ioutil"
	"os"
)

var(
	HeadSize  = 150
)

func GzipEncode(data []byte) []byte{
	buf := bytes.Buffer{}
	writer := gzip.NewWriter(&buf)
	writer.Write(data)
	return buf.Bytes()
}

func GzipDecode(data []byte) []byte{
	buf := bytes.Buffer{}
	r, _ := gzip.NewReader(&buf)
	ioutil.ReadAll(r)
	return buf.Bytes()
}
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
	tcpConn.SetWriteBuffer(30000)
	tcpConn.SetReadBuffer(30000)
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

func GetRW(p string) *os.File{
	file, _ := os.OpenFile(p, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	return file
}
func PathExists(path string)bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}