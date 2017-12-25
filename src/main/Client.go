package main

import (
	"flag"
	"sender"
	"time"
)

func main() {
	inputPath := flag.String("source", "D://go2.pdf", "source")
	blockSize := flag.Int("bsize", 30000, "block size")

	sender := &sender.Sender{InputPath: *inputPath, BlockSize: *blockSize}

	sender.Send()

	time.Sleep(time.Second * 5)
}
