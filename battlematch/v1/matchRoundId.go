package battlematch

import (
	"context"
	"errors"
	"path"
	"strconv"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// MatchV1GroupRoundIdSet 设置roundId
func MatchV1GroupRoundIdSet(ctx context.Context, groupId string, roundId int64) error {
	// 获取所有用户的键值对
	_, err := etcdClient.Client.Put(ctx, path.Join("/", projectName, match_battle_roundid_v1, groupId), strconv.FormatInt(roundId, 10), clientv3.WithLease(etcdClient.NewLease(ctx, match_success_timeout)))
	if err != nil {
		return errors.New("MatchV1SetGroupRoundId 设置round失败: " + err.Error() + ", groupId: " + groupId + ", roundId: " + strconv.FormatInt(roundId, 10))
	}
	return nil
}

// MatchV1GroupRoundIdGet 获取roundId
func MatchV1GroupRoundIdGet(ctx context.Context, groupId string) (int64, error) {
	// 获取所有用户的键值对
	resp, err := etcdClient.Client.Get(ctx, path.Join("/", projectName, match_battle_roundid_v1, groupId))
	if err != nil {
		return 0, errors.New("MatchV1GroupRoundIdGet 获取round失败: " + err.Error() + ", groupId: " + groupId)
	}
	if resp.Count == 0 { // 没有用户掉线注册
		return 0, errors.New("MatchV1GroupRoundIdGet 获取round失败, 长度为0 " + ", groupId: " + groupId)
	}
	return strconv.ParseInt(string(resp.Kvs[0].Value), 10, 64)
}

// matchV1GroupRoundIdDel 删除roundId
func matchV1GroupRoundIdDel(ctx context.Context, groupId string) error {
	// 获取所有用户的键值对
	_, err := etcdClient.Client.Delete(ctx, path.Join("/", projectName, match_battle_roundid_v1, groupId))
	if err != nil {
		return errors.New("MatchV1GroupRoundIdDel 删除round失败: " + err.Error() + ", groupId: " + groupId)
	}
	return nil
}
