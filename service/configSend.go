package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kongshui/danmu/config"
	"github.com/kongshui/danmu/model/pmsg"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/proto"
)

// ConfigFileSend 发送配置文件
func ConfigFileSend(uidList []string, openId string) error {
	fileNames := make([]string, 0)
	filepath.Walk(cfg.App.ConfigDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		data := &pmsg.ConfigFileSendMessage{
			OpenId:   openId,
			FileName: info.Name(),
		}
		fileNames = append(fileNames, data.GetFileName())
		// 如果大于512k则分片读取发送
		if info.Size() > 512*1024 {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("ConfigFileSend 打开文件失败：%v", err)
			}
			defer file.Close()
			buf := make([]byte, 512*1024)
			labelCount := 0
			for {
				n, err := file.Read(buf)
				if err != nil && err.Error() != io.EOF.Error() {
					return fmt.Errorf("ConfigFileSend 读取文件失败：%v", err)
				}
				if n == 0 {
					break
				}
				// 发送文件内容
				labelCount++
				data.Content = buf[:n]
				data.SendId = int64(labelCount)
				dataByte, err := proto.Marshal(data)
				if err != nil {
					return fmt.Errorf("ConfigFileSend 序列化文件内容失败：%v", err)
				}
				if err := sendMessage(pmsg.MessageId_ConfigFileSend, uidList, dataByte); err != nil {
					return fmt.Errorf("ConfigFileSend 发送文件内容失败：%v", err)
				}
				if info.Size() > int64(labelCount)*512*1024 {
					sCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()
					var isOk bool
					res := etcdClient.Client.Watch(sCtx, fmt.Sprintf("/config/%s/%s/%s/%s", cfg.Project, openId, data.GetFileName(),
						strconv.FormatInt(data.GetSendId(), 10)))
					for event := range res {
						for _, ev := range event.Events {
							if ev.Type == mvccpb.PUT {
								isOk = true
							}
						}
						if isOk {
							break
						}
					}
					cancel()
					if !isOk {
						return fmt.Errorf("ConfigFileSend 等待确认超时")
					}
					continue
				}
			}
		} else {
			// 读取文件内容
			content, err := os.ReadFile(path)
			if err != nil && err.Error() != io.EOF.Error() {
				return fmt.Errorf("ConfigFileSend 读取文件失败：%v", err)
			}
			data.Content = content
			dataByte, err := proto.Marshal(data)
			if err != nil {
				return fmt.Errorf("ConfigFileSend 序列化文件内容失败：%v", err)
			}
			if err := sendMessage(pmsg.MessageId_ConfigFileSend, uidList, dataByte); err != nil {
				return fmt.Errorf("ConfigFileSend 发送文件内容失败：%v", err)
			}
		}
		return nil
	})
	endData := &pmsg.ConfigFileSendEndMessage{
		FileNames: fileNames,
		OpenId:    openId,
	}
	endDataByte, err := proto.Marshal(endData)
	if err != nil {
		return fmt.Errorf("ConfigFileSendEndMessage 序列化文件内容失败：%v", err)
	}
	if err := sendMessage(pmsg.MessageId_ConfigFileSendEnd, uidList, endDataByte); err != nil {
		return fmt.Errorf("ConfigFileSendEndMessage 发送文件内容失败：%v", err)
	}
	return nil
}

// ConfigFileSendAck 配置文件发送确认
func ConfigFileSendAck(data *pmsg.ConfigFileSendAckMessage) error {
	leaseId := etcdClient.NewLease(context.Background(), 5)
	// 确认文件发送
	if _, err := etcdClient.Client.Put(context.Background(), fmt.Sprintf("/config/%s/%s/%s/%s", cfg.Project, data.GetOpenId(), data.GetFileName(),
		strconv.FormatInt(data.GetSendId(), 10)), "1", clientv3.WithLease(leaseId)); err != nil {
		return fmt.Errorf("ConfigFileSendAck 确认文件发送失败：%v", err)
	}
	return nil
}

// Config Map Read 读取配置映射
func configMapRead() error {
	if cfg.App.ConfigDir == "" {
		return nil
	}
	var (
		hasErr     = false
		configList = make([]string, 0)
	)
	filepath.Walk(cfg.App.ConfigDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			hasErr = true
			ziLog.Error(fmt.Sprintf("ConfigMapRead 遍历配置文件失败：%v", err), debug)
			return nil
		}
		if info.IsDir() {
			return nil
		}
		// 计算文件的MD5值
		md5, err := FileMD5(path)
		if err != nil {
			hasErr = true
			configList = append(configList, info.Name())
			ziLog.Error(fmt.Sprintf("ConfigMapRead 计算文件MD5失败：%v", err), debug)
			return nil
		}
		okMd5, ok := cfgConfig.FileMd5[strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))]
		if ok && okMd5 == md5 {
			return nil
		}
		// 读取配置文件内容
		fileConfig, err := config.ReadCfgConfig(path)
		if err != nil {
			hasErr = true
			configList = append(configList, info.Name())
			ziLog.Error(fmt.Sprintf("ConfigMapRead 读取配置文件失败：%v", err), debug)
			return nil
		}
		// 编写日志
		ziLog.Info(fmt.Sprintf("ConfigMapRead 读取配置文件 %s 成功", info.Name()), debug)
		cfgConfig.FileMd5[strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))] = md5
		cfgConfig.Config[strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))] = *fileConfig
		return nil
	})
	if hasErr {
		return fmt.Errorf("ConfigMapRead 读取配置文件失败,失败列表：%v", configList)
	}
	return nil
}

// 自动检测配置文件变更
func autoDetectConfigChange() {
	// 定时检测配置文件变更
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		<-ticker.C
		if err := configMapRead(); err != nil {
			ziLog.Error(fmt.Sprintf("autoDetectConfigChange 自动检测配置文件变更失败：%v", err), debug)
		}
	}
}
