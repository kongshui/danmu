package battlematch

import (
	"context"
	"errors"
	"path"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 设置匹配组状态
func MatchGroupStatusSet(ctx context.Context, groupId string, status int) error {
	t := time.NewTicker(time.Second * 1)
	defer t.Stop()
	lockPath := path.Join("/", projectName, match_battle_group_v1_lock, groupId)
	for {
		<-t.C
		if etcdClient.PutIfNotExist(ctx, lockPath, strconv.Itoa(status), clientv3.LeaseID(etcdClient.NewLease(ctx, 3))) {
			defer etcdClient.Client.Delete(ctx, lockPath) // 删除锁
			oldStatus, err := MatchGroupStatusGet(ctx, groupId)
			if err != nil {
				return errors.New("MatchGroupStatusSet 设置匹配组状态失败: " + err.Error() + ", groupId: " + groupId)
			}
			if oldStatus == status { // 状态相同
				return errors.New("equal")
			}
			id := etcdClient.NewLease(ctx, match_success_timeout) // 创建租约
			// 设置匹配组状态
			if _, err := etcdClient.Client.Put(ctx, path.Join("/", projectName, match_battle_group_status, groupId), strconv.Itoa(status), clientv3.WithLease(id)); err != nil {
				return errors.New("MatchGroupStatusSet 设置匹配组状态失败: " + err.Error() + ", groupId: " + groupId)
			}
			return nil
		}
		continue
	}
}

// 删除匹配组状态
func MatchGroupStatusDel(ctx context.Context, groupId string) error {
	// 删除匹配组状态
	_, err := etcdClient.Client.Delete(ctx, path.Join("/", projectName, match_battle_group_status, groupId))
	if err != nil {
		return errors.New("MatchGroupStatusDel 删除匹配组状态失败: " + err.Error() + ", groupId: " + groupId)
	}
	return nil
}

// 获取匹配组状态
func MatchGroupStatusGet(ctx context.Context, groupId string) (int, error) {
	// 获取匹配组状态
	resp, err := etcdClient.Client.Get(ctx, path.Join("/", projectName, match_battle_group_status, groupId))
	if err != nil {
		return 0, errors.New("MatchGroupStatusGet 获取匹配组状态失败: " + err.Error() + ", groupId: " + groupId)
	}
	if len(resp.Kvs) == 0 {
		return 0, nil
	}
	status, err := strconv.Atoi(string(resp.Kvs[0].Value))
	if err != nil {
		return 0, errors.New("MatchGroupStatusGet 获取匹配组状态失败: " + err.Error() + ", groupId: " + groupId)
	}
	return status, nil
}
