package main

import (
	"flag"
	"sender"
	"time"
)
/*1. 是否成功
2. 查询缺失
3. 续传*/
func main() {
	inputPath := flag.String("source", "D://go2.pdf", "source")
	blockSize := flag.Int("bsize", 30000, "block size")

	s := sender.NewSender(*inputPath, *blockSize)

	s.Send()

	time.Sleep(time.Second * 5)
}
