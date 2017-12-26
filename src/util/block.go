package util

import (
	"bytes"
	"strconv"
)

type Block struct {
	Offset  int
	Blength int
	Order   int
	Name string
	Max int
}

func (this *Block) GetHash() string{
	buffer := bytes.Buffer{}
	buffer.WriteString(strconv.Itoa(this.Offset))
	buffer.WriteString("_")
	buffer.WriteString(strconv.Itoa(this.Blength))
	buffer.WriteString("_")
	buffer.WriteString(strconv.Itoa(this.Order))

	return GetHash(buffer.Bytes())
}