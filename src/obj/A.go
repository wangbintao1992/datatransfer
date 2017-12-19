package main

import (
	"fmt"
	"path"
)

func main() {

	join := path.Join("d：", "tmp","go.pdf",".tmp")

	fmt.Println(join)
}
