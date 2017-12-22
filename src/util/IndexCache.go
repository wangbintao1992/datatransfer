package util

import (
	"fmt"
	"bufio"
	"io"
	"encoding/json"
	"os"
)

type IndexCache struct {
	index map[string]Head
	file *os.File
}

func Init(path string) *IndexCache {

	cache := &IndexCache{
		index: make(map[string]Head, 0)}

	cache.loadToMermory(path)

	return cache
}
func (i *IndexCache)loadToMermory(path string) {
	i.file = GetRW(path)
	if PathExists(path) {
		r := bufio.NewReader(i.file)
		for {
			line, _, err := r.ReadLine()

			if err != nil && err == io.EOF {
				break
			}

			head := &Head{}
			json.Unmarshal(line,head)

			fmt.Println("恢复：", string(line))
			i.Put(head.Hash, *head)
		}
	}
}

func (i *IndexCache) SaveCache(){
	writer := bufio.NewWriter(i.file)
	for _,v := range i.index{
		bytes, _ := json.Marshal(v)
		writer.Write(bytes)
		writer.WriteString("\r\n")
		writer.Flush()
	}
	fmt.Println("保存成功")
}

func (i *IndexCache) CheckMerge() bool{

	return false
}

func (i *IndexCache) Put(k string,v Head){
	i.index[k] = v
}

func (i *IndexCache) Get(k string) (Head,bool){
	head ,ok := i.index[k]
	return head,ok
}

func (i *IndexCache) GetCache() []Head{
	array := make([]Head,0)
	for k,v := range i.index {
		fmt.Println("merge index:",k," file order：",v.Order)
		array = append(array, v)
	}

	return array
}

