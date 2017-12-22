package main

import (
	"fmt"
	"os"
)

func main() {
	a := &A{Name: "jack"}

	a.setName("rose")
	os.OpenFile("D://test.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
}

func (a *A) setName(s string){
	a.Name = s
	fmt.Println(a.Name)
}


func (a A) setName2(s string){
	a.Name = s
	fmt.Println(a.Name)
}

type A struct {
	Name string
}