package main

import (
	"bufio"
	"log"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

func main() {
	start := false

	in := "D://log.txt"
	file, err := os.Open(in)

	if err != nil {
		log.Print(err)
	}

	read := bufio.NewReader(file)

	e := xlsx.NewFile()
	sheet, _ := e.AddSheet("Sheet1")

	var o = new(Obj)
	for {
		d, _, err := read.ReadLine()

		if err != nil {
			if err == io.EOF {
				break
			}
		}

		tmp := string(d[:])

		if start {
			if len(tmp) == 0 {
				//new line
				sheet.AddRow()
				row := sheet.AddRow()
				row.SetHeightCM(1) //设置每行的高度
				cell := row.AddCell()
				cell.Value = o.fileName
				cell = row.AddCell()
				cell.Value = o.fullName
				cell = row.AddCell()
				cell.Value = o.problem
				cell = row.AddCell()
				cell.Value = o.result
				cell = row.AddCell()
				cell.Value = o.action
				o = new(Obj)

			} else {
				if strings.Contains(tmp, "File Name") {
					log.Println(tmp)
					index := strings.Index(tmp, ":")
					o.fileName = tmp[index+1:]
				}
				if strings.Contains(tmp, "Full Name") {
					index := strings.Index(tmp, ":")
					o.fullName = tmp[index+1:]
				}
				if strings.Contains(tmp, "Action") {
					index := strings.Index(tmp, ":")
					o.action = tmp[index+1:]
				}
				if strings.Contains(tmp, "Result") {
					index := strings.Index(tmp, ":")
					o.result = tmp[index+1:]
				}
				if strings.Contains(tmp, "Problem") {
					index := strings.Index(tmp, ":")
					o.problem = tmp[index+1:]
				}
			}

			continue
		}
		if strings.Contains(tmp, "Deploy Results") {
			start = true
		}
	}
	e.Save("D://file.xlsx")
}

func handle() int {

	return 0
}

func getExcel() {
	file2 := xlsx.NewFile()
	sheet, _ := file2.AddSheet("Sheet1")
	row := sheet.AddRow()
	row.SetHeightCM(1) //设置每行的高度
	cell := row.AddCell()
	cell.Value = "haha"
	cell = row.AddCell()
	cell.Value = "xixi"

	err2 := file2.Save("D://file.xlsx")
	if err2 != nil {
		panic(err2)
	}
}

func byteToString(b []byte) string {
	s := make([]string, len(b))

	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}

	return strings.Join(s, ",")
}

type Obj struct {
	fileName string
	fullName string
	action   string
	result   string
	problem  string
}
