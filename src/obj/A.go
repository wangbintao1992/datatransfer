package main

import (
	"io/util"
	"os"
)

func main() {
	m := make(map[string]util.Head,10)
	f, _ := os.OpenFile("D:\\tmp\\t.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	defer f.Close()

	head := &util.Head{Path: f}
	m["a"] = *head

	a := m["a"]
	a.Path.WriteString("test ssss")
}
