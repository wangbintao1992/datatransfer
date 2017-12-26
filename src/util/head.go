package util

import (
	"bytes"
	"strconv"
	"encoding/json"
	"log"
	"path"
	"os"
	"gvar"
)

type Head struct {
	Order int
	Length int
	Hash string
	Name string
	Path *os.File
	MD5 string
	Max int
}

func EncodeHead(b Block,md5 string) []byte{
	head := &Head{
		Order:  b.Order,
		Length: b.Blength,
		Name:   b.Name,
		Hash:   b.GetHash(),
		MD5:md5}

	marshal, e := json.Marshal(head)

	if e != nil{
		log.Println(e)
	}

	log.Println("head json",string(marshal),"offset",b.Offset)

	if space := gvar.HeadSize - len(marshal); space > 0{
		marshal = append(marshal,GetSpace(space)...)
	}

	return marshal
}

func DecodeHead(data []byte) *Head{
	i := new(Head)

	right := bytes.TrimRight(data, "\x00")

	log.Println("head json",string(right))
	e2 := json.Unmarshal(right, &i)

	if e2 != nil{
		log.Println(e2)
	}
	return i
}

func (this *Head) GetFilePath(s string) string{
	return path.Join(s,this.Name)
}

//example path = data.pdf_1.tmp
func (this *Head) GetTmpFilePath(s string) string{
	buffer := bytes.Buffer{}
	buffer.WriteString(this.Name)
	buffer.WriteString("_")
	buffer.WriteString(strconv.Itoa(this.Order))
	buffer.WriteString(".tmp")

	return path.Join(s,this.Name + "tmp",buffer.String())
}



type HeadIndex []Head

func(h HeadIndex)Len() int{
	return len(h)
}
func(h HeadIndex) Less(i, j int) bool{
	return h[i].Order < h[j].Order
}
func(h HeadIndex)Swap(i, j int){
	h[i],h[j] = h[j],h[i]
}