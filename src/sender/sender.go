package sender

import (
	"net"
	"log"
	"time"
	"os"
	"util"
	"bufio"
	"path"
)

var startTime time.Time
type Sender struct {
	InputPath string
	BlockSize int
	MD5 string
	Root string
}
func NewSender(s string,blockSize int) *Sender{
	file, e := os.Open(s)
	util.CheckErr(e)
	sender := &Sender{InputPath: s, BlockSize: blockSize,MD5:util.GetFileMD5(file),Root:path.Dir(s)}

	return sender
}
func (this *Sender) Send() {

	if(!this.checkComplete()) {
		//TOOD 查询报文
		QueryAbsentBlockClient(this.MD5)
	}

	conn, _ := net.Dial("tcp", "localhost:8080")

	util.SetTCPOption(conn,this.BlockSize)

	//TODO pool
	//TODO asych timeout


	this.sendData(conn)
}
func (this *Sender) checkComplete() bool{
	return CheckComlete0(this.InputPath,this.MD5)
}

func (this *Sender)sendData(conn net.Conn){
	file, _ := os.Open(this.InputPath)
	this.setStatus(false)

	this.sendData_0(file,conn)

	this.setStatus(true)
}
func (this *Sender) setStatus(b bool) {
	sf := path.Join(this.Root, this.MD5)
	log.Println(sf)
	if b{
		os.Remove(sf)
	} else {
		rw := util.GetRW(sf)
		rw.WriteString(this.MD5)
		rw.Sync()
		defer rw.Close()
	}
}
func (this *Sender)sendData_0(file *os.File,conn net.Conn){

	blocks := this.getBlocks(this.InputPath, int64(this.BlockSize))

	this.startCalcTime()

	log.Println("write start")

	//TODO 非常重要，待解决，socket 缓冲区可能第一次发包，不满

	for i := 0; i < len(blocks); i ++ {
		go this.writePacket(file, blocks, i, conn)
	}
	this.endCalcTime()
}
func (this *Sender)writePacket(file *os.File, blocks []util.Block, i int, conn net.Conn) {
	msg := this.getPacket(file, blocks[i])
	n, e := bufio.NewWriter(conn).Write(msg)
	if e != nil {
		log.Println(e)
	}
	log.Println("write:", n, "write...")
}
//TODO common
func (this *Sender)endCalcTime() {
	since := time.Since(startTime)
	log.Println("用时 s", since)
}
//TODO common
func (this *Sender)startCalcTime() {
	startTime = time.Now()
}

func (this *Sender)getBlocks(path string,blockSize int64) []util.Block{
	return util.GetBlockArr(path,blockSize)
}

func (this *Sender)getPacket(file *os.File,b util.Block) []byte{
	buf := make([]byte, b.Blength)

	n, e := this.readFile(file, buf,b)

	if e != nil{
		log.Println(e)
	}
	head := this.getHead(b)

	return this.mergePacket(head, buf[:n])
}

func (this *Sender)readFile(file *os.File,buf []byte,block util.Block)(int,error) {
	return file.ReadAt(buf,int64(block.Offset))
}
func (this *Sender)getHead(b util.Block) []byte {
	return util.EncodeHead(b)
}
func (this *Sender)mergePacket(head []byte, body []byte) []byte{
	return append(head, body...)
}
