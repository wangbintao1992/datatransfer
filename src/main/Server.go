package main

import (
	"net/http"
	"log"
	"flag"
	"receive"
	"cache"
	"path"
)
/*1. server恢复重连任务
2. client查询重连数据
3. server接受续传
*/
//TODO 并发传输
func main() {
	rtPath := flag.String("rtpath", "D://tmp", "root path")
	tmpPath := flag.String("tmppath", "tmp", "tmp path")

	cache := cache.Init(path.Join(*rtPath, "index.log"))
	r := new(receive.Receiver)
	r.InitParam(*rtPath,*tmpPath,cache)

	go r.Start()

	http.HandleFunc("/qab", queryAbsentBlock)

	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}


	log.Println("服务启动")
}

func queryAbsentBlock(w http.ResponseWriter, req *http.Request){
}