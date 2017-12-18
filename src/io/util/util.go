package util

import (
	"bytes"
	"encoding/binary"
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

func ByteToInt32(d []byte) int32{
	var r int32 = 0
	buffer := bytes.NewBuffer(d)
	binary.Read(buffer,binary.BigEndian,&r)

	return r
}