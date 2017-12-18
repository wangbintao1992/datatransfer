package util

import (
	"bytes"
	"strconv"
	"os"
)

func (this *Head) encodeHead(head Head) []byte{

	return nil
}

func (this *Head) decodeHead(data []byte) *Head{

	return nil
}

//example path = data.pdf_1.tmp
func (this *Head) GetTmpFilePath(s string) string{
	buffer := bytes.Buffer{}
	buffer.WriteString(s)
	buffer.WriteString("_")
	buffer.WriteString(strconv.Itoa(this.Order))
	buffer.WriteString(".tmp")
	return buffer.String()
}

type Head struct {
	Order int
	Length int32
	Hash string
	Name string
	Path os.File
}