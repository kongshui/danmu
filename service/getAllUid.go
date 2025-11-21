package service

import (
	"fmt"
	"path"

	"github.com/kongshui/danmu/common"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// GetAllUid 获取所有主播的ID
func GetAllUid() ([]string, error) {
	// 从数据库中查询所有主播ID
	var uids []string
	res, err := etcdClient.Client.Get(first_ctx, path.Join("/", cfg.Project, common.Uuid_Online_key), clientv3.WithPrefix())
	// 处理查询结果
	if err != nil {
		return nil, fmt.Errorf("GetAllUid 查询主播ID失败：%v", err)
	}
	for _, kv := range res.Kvs {
		uids = append(uids, string(kv.Value))
	}
	return uids, nil
}
