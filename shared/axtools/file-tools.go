package axtools

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func GetFirstDir(path string) string {
	dirs := strings.Split(path, "/")
	if len(dirs) < 2 {
		return ""
	}

	dir := dirs[len(dirs)-2]
	if dir == "." {
		dir = ""
	}
	return dir
}

func GetFirstDirOfDir(path string) string {
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	dirs := strings.Split(path, "/")
	if len(dirs) < 2 {
		return ""
	}

	dir := dirs[len(dirs)-1]
	if dir == "." {
		dir = ""
	}
	return dir
}

func FileNameWithoutExtension(fileName string) string {
	dirs := strings.Split(fileName, "/")
	return strings.TrimSuffix(dirs[len(dirs)-1], filepath.Ext(dirs[len(dirs)-1]))
}

func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}
