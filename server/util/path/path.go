package path

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var osType = runtime.GOOS
var rootPath string

func init() {
	rootPath, _ = os.Getwd()
}

func GetPath(path string) string {
	pathNames := strings.Split(path, "/")
	absPath := filepath.Join(rootPath, filepath.Join(pathNames...))
	if runtime.GOOS == "windows" {
		absPath = filepath.FromSlash(absPath)
	}
	return absPath
}

func JoinPath(basePath, filePath string) string {
	if basePath == "" {
		return GetPath(filePath)
	}
	pathNames := strings.Split(filePath, "/")
	absPath := filepath.Join(basePath, filepath.Join(pathNames...))
	if runtime.GOOS == "windows" {
		absPath = filepath.FromSlash(absPath)
	}
	return absPath
}
