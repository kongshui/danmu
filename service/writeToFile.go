package service

import (
	"fmt"
	"os"
)

// TransLog 转换日志
func writeToFile(filePath string, data []byte) error {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("writeToFile open file err: %v", err)
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("writeToFile write file err: %v", err)
	}
	return nil
}
