package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func IsFileExist(filename string) bool {
	if _,err:=os.Stat(filename);os.IsNotExist(err){
		return  false
	}
	return  true
}

func RemoveFile(filename string) bool {

	if err := os.Remove(filename); err != nil{
		fmt.Println(err)
		return false
	}
	return true
}

func ListDir(dir, prefix string, suffix string) (files []string, err error){

	files = []string{}

	_dir, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	prefix = strings.ToLower(prefix)  //匹配前缀
	suffix = strings.ToLower(suffix)  //匹配后缀

	for _, _file := range _dir {

		if _file.IsDir() {
			continue //忽略目录
		}

		if len(prefix) == 0 || strings.HasPrefix(strings.ToLower(_file.Name()), prefix) {
			if len(suffix) == 0 || strings.HasSuffix(strings.ToLower(_file.Name()), suffix) {
				files = append(files, path.Join(dir, _file.Name()))
			}
		}
	}

	return files, nil
}