package util

import (
	"log"
	"os"
	"path"
)

func GetFileName(p string) string{
	return path.Base(p)
}

func GetBlockArr(path string,blockSize int64) []Block{
	log.Println("分块大小byte:",blockSize)
	log.Println("源文件路径:" + path)

	fileName := GetFileName(path)
	log.Println("文件名:" + fileName)
	info, e := os.Stat(path)

	if e != nil{
		panic(e)
	}
	size := info.Size()
	log.Println("源文件大小:",size," unit:byte")

	blockNum, remain := divFile(size, path, blockSize)

	log.Println("分块:",blockNum)
	log.Println("剩余:",remain)

	blocks := make([]Block, blockNum + 1)
	//TODO int ?
	blockNum2 := int(blockNum)
	bs := int(blockSize)
	order := 0
	for i := 0; i < blockNum2; i ++{
		blocks[i] = *&Block{Offset:  i * bs,Blength: bs,Order:order,Name:fileName,Max:int(blockNum)}
		order ++
	}

	log.Println(cap(blocks))
	blocks[blockNum] = *&Block{Offset: blockNum2 * int(blockSize),Blength:int(remain),Order:order,Name:fileName,Max:int(blockNum)}

	return blocks
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