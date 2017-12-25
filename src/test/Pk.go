package main

import (
	"fmt"
	"path"
)

func main()  {
	dir := path.Dir("D://go2.pdf")
	fmt.Println(dir)
}