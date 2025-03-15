package file

import (
	"fmt"
	"os"
)

func ReadDataFromPath(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Read file [%s] fail, error is %+v", path, err))
	}
	return data
}

// 判断所给路径文件/文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}
