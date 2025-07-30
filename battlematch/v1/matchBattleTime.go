package battlematch

import (
	"context"
	"fmt"
	"path"
	"strconv"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// MatchBattleTimeSet 匹配战斗时间设置
func MatchBattleTimeSet(ctx context.Context, groupId string, statTime int64) error {
	_, err := etcdClient.Client.Put(ctx, path.Join("/", projectName, match_battle_group_time, groupId), strconv.FormatInt(statTime, 10), clientv3.WithLease(etcdClient.NewLease(ctx, match_success_timeout)))
	if err != nil {
		return fmt.Errorf("MatchBattleTimeSet etcdClient.Client.Put error, group: %v, Err: %v", groupId, err)
	}
	return nil
}

// MatchBattleTimeGet 匹配战斗时间获取
func MatchBattleTimeGet(ctx context.Context, groupId string) (int64, error) {
	resp, err := etcdClient.Client.Get(ctx, path.Join("/", projectName, match_battle_group_time, groupId))
	if err != nil {
		return 0, fmt.Errorf("MatchBattleTimeGet etcdClient.Client.Get error, group: %v, Err: %v", groupId, err)
	}
	if len(resp.Kvs) == 0 {
		return 0, nil
	}
	return strconv.ParseInt(string(resp.Kvs[0].Value), 10, 64)
}

// matchBattleTimeDel 匹配战斗时间删除
func matchBattleTimeDel(ctx context.Context, groupId string) error {
	_, err := etcdClient.Client.Delete(ctx, path.Join("/", projectName, match_battle_group_time, groupId))
	if err != nil {
		return fmt.Errorf("MatchBattleTimeDel etcdClient.Client.Delete error, group: %v, Err: %v", groupId, err)
	}
	return nil
}
