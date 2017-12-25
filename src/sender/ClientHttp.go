package sender

import (
	"util"
	"os"
	"net/http"
	"io/ioutil"
	"strings"
	"path"
)
func CheckComlete0(s string,md5 string) bool{
	return util.PathExists(path.Join(s,md5))
}
func QueryAbsentBlockClient(s string) []string{
	file, e := os.Open(s)
	util.CheckErr(e)
	md5 := util.GetFileMD5(file)
	resp, err := http.Get("localhost:12345/qab/?md5=" + md5)

	util.CheckErr(err)
	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)
	util.CheckErr(err2)

	arr := strings.Split(string(body), ",")
	return arr
}