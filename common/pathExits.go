package common

import (
	"io"
	"os"
)

// PathExists 判断文件或文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 判断文件夹是否为空
func DirIsEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return true
	}
	defer f.Close()
	_, err = f.Readdirnames(1) // 获取一级子目录
	if err == nil {
		return false // 存在子目录
	}
	if err == io.EOF {
		return true // 文件为空
	}
	return false
}

// 判断是否是文件
func IsFile(path string) bool {
	if !PathExists(path) {
		return false
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}

// 判断是否是文件夹
func IsDir(path string) bool {
	if !PathExists(path) {
		return false
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}
