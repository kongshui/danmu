package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// FileMD5 计算文件的MD5值，返回十六进制字符串
func FileMD5(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer func() {
		// 延迟关闭文件，并检查关闭错误（非必需，但建议处理）
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("关闭文件失败: %v\n", closeErr)
		}
	}()

	// 创建MD5哈希器
	hash := md5.New()

	// 将文件内容流式复制到哈希器（适合大文件，避免占用过多内存）
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("读取文件内容失败: %w", err)
	}

	// 计算最终哈希值（字节数组），并转换为十六进制字符串
	md5Bytes := hash.Sum(nil)
	return hex.EncodeToString(md5Bytes), nil
}
