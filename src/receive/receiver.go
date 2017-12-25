package receive

import (
	"net"
	"log"
	"os"
	"util"
	"sort"
	"bufio"
	"io"
	"sync"
	"path"
	"cache"
)
type Receiver struct {
	RtPath string
	TmpPath string
	indexCache *cache.IndexCache
	once sync.Once
}
func (this *Receiver)InitParam(rtPath string,tmpPath string,cache *cache.IndexCache)  {
	log.Println("参数初始化...")
	this.RtPath = rtPath
	this.TmpPath = tmpPath
	this.indexCache = cache
}

func (this *Receiver)Start() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "localhost:8080")

	util.CheckErr(err)

	server, e := net.ListenTCP("tcp", tcpAddr)

	util.CheckErr(e)

	for {
		conn, _ := server.AcceptTCP()
		//TODO
		util.SetTCPOption(conn,30000)
		this.handle(conn)

		//TODO 同步
		this.mergeFile()

		this.clean()
	}

	log.Println("finash")
}
func (this *Receiver)handle(accept net.Conn) {

	log.Println("reuqest come in")

	this.readPacket(accept)
}

//不同协程 并行读，所以每读一包，包含包头包体
func (this *Receiver)readPacket(conn net.Conn) {
	log.Println("start receive")
	log.Println("receiveing ...")
	for{
		head, n := this.readHead(conn)

		if n == 0{
			break
		}

		this.once.Do(func() {
			this.TmpPath = path.Join(this.RtPath,"\\",head.Name + this.TmpPath)
			os.Mkdir(this.TmpPath,0666)
		})

		//TODO asyn?
		this.addToIndex(head)

		log.Println("body length",head.Length)

		this.readAndWriteBody(conn, head)
	}
}
func (this *Receiver)readAndWriteBody(conn net.Conn,h *util.Head) {
	var length = h.Length
	var file = h.Path

	tmpBuf := make([]byte, 0)
	for{
		num := 0
		data,num,remain := this.getPartBody(length, conn)

		if num == 0{
			break
		}

		length = remain

		tmpBuf = append(tmpBuf, data...)
	}
	go this.flushToDisk(file, tmpBuf,h)
}
func (this *Receiver)flushToDisk(file *os.File, part []byte,h *util.Head) {
	//flush to disk
	writer := bufio.NewWriter(file)
	writer.Write(part)
	writer.Flush()

	//flush to index.log
	this.indexCache.SaveCache(h)
}
func (this *Receiver)getPartBody(length int, conn net.Conn) ([]byte,int,int){
	buf := make([]byte, length)

	num, _ := bufio.NewReader(conn).Read(buf)

	log.Println("read num:", num)

	remain := 0
	if remain = length - num; remain != 0{
		log.Println("resize buf",remain)
	}

	return buf[:num],num,remain
}
func (this *Receiver)mergeFile() {
	if this.indexCache.CheckMerge(){
		//TODO 完整性校验 max ?
		log.Println("数据包缺失")
		log.Println("缺失:")
		return
	}

	log.Println("merge ...")

	index := this.indexCache.GetCache()

	sort.Sort(util.HeadIndex(index))

	writer := this.getFinalFile(index[0])

	defer writer.Close()

	bufWrite := bufio.NewWriter(writer)

	var totalSize int64 = 0
	for _, v := range index {
		//!important reset read pointer
		v.Path.Seek(0,io.SeekStart)
		reader := bufio.NewReader(v.Path)
		written, _ := io.Copy(bufWrite, reader)
		log.Println("merge size:",written)
		totalSize += written
		bufWrite.Flush()
		v.Path.Close()
	}

	log.Println("totalSize:",totalSize)
	log.Println("merge finash...")
}
func (this *Receiver)checkMerge() {

}
func (this *Receiver)getFinalFile(h util.Head) *os.File{
	return util.GetRW(h.GetFilePath(this.RtPath))
}
func (this *Receiver)clean(){
	log.Println("清理临时文件")
	r := os.RemoveAll(this.TmpPath)
	if r != nil{
		log.Println(r)
	}
	this.indexCache.Close()
	log.Println("清理临时文件结束")
}

func (this *Receiver)readHead(conn net.Conn) (h *util.Head, num int){
	t := make([]byte,util.HeadSize)

	n, e := conn.Read(t)

	if e != nil{
		log.Println("read type err",e)
		return nil,0
	}

	return util.DecodeHead(t),n
}
//TODO 先考虑断线，不考虑半途中断
func (this *Receiver)addToIndex(head *util.Head){

	if h ,ok:= this.indexCache.Get(head.Hash); ok{
		log.Println("重复读 err")
		head.Path = h.Path
	} else {
		head.Path = util.GetRW(head.GetTmpFilePath(this.RtPath))
		this.indexCache.Put(head.Hash,*head)
	}
}
