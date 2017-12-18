package util

import (
	"fmt"
	"os"
	"strconv"
)

func GetBlockArr(path string,blockSize int64) []Block{
	fmt.Println("分块大小byte:",blockSize)
	fmt.Println("源文件路径:" + path)

	info, e := os.Stat(path)

	if e != nil{
		panic(e)
	}
	size := info.Size()
	fmt.Println("源文件大小:",size," unit:byte")

	blockNum, remain := divFile(size, path, blockSize)

	fmt.Println("分块:" + strconv.FormatInt(blockNum,10))
	fmt.Println("剩余:" + strconv.FormatInt(remain,10))

	blocks := make([]Block, blockNum)
	//TODO int ?
	blockNum2 := int(blockNum)
	for i := 0; i < blockNum2; i ++{
		blocks[i] = &Block{Offset:  i * blockNum2,Blength: blockNum2,Order:i}
	}

	return nil
}

func divFile(size int64,path string,packet int64)  (int64,int64){
	blockNum := size / packet
	mod := size % packet

	var	remain int64
	if mod != 0{
		remain = size - blockNum * packet

	}

	return blockNum,remain
}