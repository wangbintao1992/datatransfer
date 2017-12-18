package util

import (
	"crypto/sha1"
	"bytes"
	"strconv"
)

type Block struct {
	Offset  int
	Blength int
	Order   int
}

func (this *Block) getHash() string{
	buffer := bytes.Buffer{}
	buffer.WriteString(strconv.Itoa(this.Offset))
	buffer.WriteString("_")
	buffer.WriteString(strconv.Itoa(this.Blength))
	buffer.WriteString("_")
	buffer.WriteString(strconv.Itoa(this.Order))

	t := sha1.New()
	t.Write(buffer.Bytes())
	return string(t.Sum(nil))
}