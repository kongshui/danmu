package service

import (
	"os"
	"path/filepath"
	"time"
)

// DeleteStoreLog 删除存储的日志文件, deleteDay 为删除多少天之前的文件，dirPath 为日志文件所在目录
func DeleteStoreLog(dirPath string, deleteDay int) {
	// TODO 删除存储的日志文件
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.ModTime().AddDate(0, 0, deleteDay).Before(time.Now()) {
			return os.Remove(path)
		}
		return nil
	})
}

// 自动删除日志文件
func autoDeleteLogFile() {
	// 每天凌晨删除日志文件
	if cfg.App.DeleteDay <= 0 {
		return
	}
	t := time.NewTicker(24 * time.Hour)
	defer t.Stop()
	for range t.C {
		DeleteStoreLog(cfg.App.LogStoreDir, cfg.App.DeleteDay)
	}
}
