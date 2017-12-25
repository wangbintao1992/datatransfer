package main

import (
	"log"
	"net"
	"io"
	"os"
	"strconv"
	bytes2 "bytes"
)

func main()  {
	log.Println("server start")
	listener, e := net.Listen("tcp", "localhost:8080")
	log.Println(e)
	path := "D:\\tmp\\"


	index := 0
	bytes := make([]byte, 1024)
	for{
		index ++
		accept, i := listener.Accept()
		log.Println(i)

		file, e3 := os.OpenFile(path+"t"+strconv.Itoa(index)+".data", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		log.Println(e3)
		for{
			n, err := accept.Read(bytes)
			log.Println(n)
			log.Println(err)
			log.Println(string(bytes))
			bytes2.TrimSpace(bytes)
			file.Write(bytes)

			if err != nil{
				if err == io.EOF{
					break
				}
				break
			}
		}
		file.Sync()
	}
}
