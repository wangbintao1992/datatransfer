package main

import (
	"net/http"
	"log"
	"flag"
	"receive"
	"gvar"
	"util"
	"path"
	"cache"
)
/*1. server恢复重连任务
2. client查询重连数据
3. server接受续传
4. merge失败，重新传输
*/
//TODO 并发传输
func main() {
	gvar.RtPath  = *flag.String("rtpath", "D://tmp", "root path")
	gvar.TmpPath = *flag.String("tmppath", "tmp", "tmp path")

	r := &receive.Receiver{RtPath:gvar.RtPath,TmpPath:gvar.TmpPath}
	log.Println("service start")
	go r.Start()

	http.HandleFunc("/qab", queryAbsentBlock)

	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func queryAbsentBlock(r http.ResponseWriter, req *http.Request){
	req.ParseForm()
	md5 := req.Form.Get("md5")
	log.Println("收到md5:",md5)

	result := new(util.Result)
	if util.PathExists(path.Join(gvar.RtPath,md5)){
		result.Flag = true
		cache.GetCache()
	} else {
		result.Flag = false
	}
}