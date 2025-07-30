package battlematch

import (
	"context"
	"encoding/json"
	"errors"
	"path"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 查询匹配组中信息
func QueryVs1GroupInfo(ctx context.Context, groupId string) ([]string, error) {
	result, err := etcdClient.Client.Get(ctx, path.Join("/", projectName, matcd_battle_store_v1, groupId), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	if result.Count == 0 {
		return nil, errors.New("QueryVs1GroupInfo is nil")
	}
	var OpenIdList []string
	for _, kv := range result.Kvs {
		json.Unmarshal(kv.Value, &OpenIdList)
	}
	return OpenIdList, nil
}

// 通过OpeenId查询匹配组信息
func QueryVs1GroupInfoByOpenId(ctx context.Context, openId string) ([]string, error) {
	groupid, err := queryVs1GroupByOpenId(ctx, openId)
	if err != nil {
		return nil, err
	}
	if groupid == "" {
		return nil, errors.New("not in group")
	}
	return QueryVs1GroupInfo(ctx, groupid)
}

// 查询匹配组
func queryVs1GroupByOpenId(ctx context.Context, openId string) (string, error) {
	result, err := etcdClient.Client.Get(ctx, path.Join("/", projectName, matcd_battle_store_v1, openId))
	if err != nil {
		return "", err
	}
	if result.Count == 0 {
		return "", nil
	}
	return string(result.Kvs[0].Value), nil
}

// 查看是否已经在分组中
func IsInVs1Group(ctx context.Context, openId string) (bool, string) {
	result, _ := etcdClient.Client.Get(ctx, path.Join("/", projectName, matcd_battle_store_v1, openId))
	if result.Count == 0 {
		return false, ""
	}
	return true, string(result.Kvs[0].Value)
}

// 查询是否在取消匹配组中
func isInVs1Cancel(ctx context.Context, openId string) bool {
	result, _ := etcdClient.Client.Get(ctx, path.Join("/", projectName, match_battle_cancel_v1, openId), clientv3.WithCountOnly())
	return result.Count != 0
}

// 查找用户是否在匹配列表
func IsInResiter(ctx context.Context, matchNum, openId string) bool {
	key := path.Join(matcd_battle_register_v1, openId)
	if matchNum != "" {
		key = path.Join("/", projectName, matcd_battle_num_register_v1, matchNum, openId)
	}
	res, _ := etcdClient.Client.Get(ctx, key, clientv3.WithCountOnly())
	return res.Count > 0
}
