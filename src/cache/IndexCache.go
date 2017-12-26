package cache

import (
	"log"
	"bufio"
	"io"
	"encoding/json"
	"os"
	"bytes"
	"util"
	"path"
	"gvar"
)
var cache *IndexCache
type IndexCache struct {
	index map[string]util.Head
	file *os.File
}

func Init(path string){
	//TODO 断点续传
	log.Println("cache 初始化")
	cache = &IndexCache{
		index: make(map[string]util.Head, 0)}

	cache.loadToMermory(path)
}
func GetCache() *IndexCache{
	return cache
}
func (i *IndexCache)loadToMermory(path string) {
	if util.PathExists(path) {
		r := bufio.NewReader(util.GetR(path))
		for {
			line, _, err := r.ReadLine()

			if err != nil && err == io.EOF {
				break
			}

			head := &util.Head{}
			json.Unmarshal(bytes.TrimSpace(line),head)

			log.Println("恢复：", string(line))
			i.Put(head.Hash, *head)
		}
	}
}
func (i *IndexCache) SaveCache(h *util.Head){
	if i.file == nil{
		i.file = util.GetRW(path.Join(gvar.RtPath, h.MD5))
	}
	writer := bufio.NewWriter(i.file)
	data, _ := json.Marshal(h)
	writer.Write(data)
	writer.WriteString("\r\n")
	writer.Flush()
	log.Println("index保存成功,hash:",h.Hash)
}

func (i *IndexCache) CheckMerge() bool{

	return false
}

func (i *IndexCache) Put(k string,v util.Head){
	i.index[k] = v
}

func (i *IndexCache) Get(k string) (util.Head,bool){
	head ,ok := i.index[k]
	return head,ok
}

func (i *IndexCache) GetCache() []util.Head{
	array := make([]util.Head,0)
	for k,v := range i.index {
		log.Println("merge index:",k," file order：",v.Order)
		array = append(array, v)
	}

	return array
}
func (i *IndexCache) Close() {
	i.file.Close()
}

