package battlematch

import (
	"context"
	"errors"
	"path"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 匹配组掉线注册
func DisconnectMatchRegister(ctx context.Context, openId string) error {
	if openId == "" {
		return nil
	}
	groupId, err := queryVs1GroupByOpenId(ctx, openId)
	if err != nil {
		return errors.New("DisconnectRegister 查询组信息失败: " + err.Error() + ", openId: " + openId)
	}
	if groupId == "" {
		return nil
	}
	id := etcdClient.NewLease(ctx, 7200)                                                                                                       // 创建租约
	_, err = etcdClient.Client.Put(ctx, path.Join("/", projectName, matcd_battle_disconnect, groupId, openId), openId, clientv3.WithLease(id)) // 注册到匹配组
	if err != nil {
		return errors.New("DissconnectRegister 注册组失败: " + err.Error() + ", groupId: " + openId)
	}
	openIdList, err := QueryVs1GroupInfo(ctx, groupId)
	if err != nil {
		return errors.New("DisconnectRegister 查询组信息失败: " + err.Error() + ", groupId: " + groupId)
	}
	for _, openIdStr := range openIdList {
		if !QueryOpenIdInMatchDisconnect(ctx, groupId, openIdStr) {
			return nil
		}
	}
	// 设置断线
	MatchGroupStatusSet(first_ctx, groupId, 10)
	// 注销用户
	return UnregisterBattleV1ByGroupId(ctx, groupId)
}

// 获取匹配组掉线注册用户
func QueryOpenIdInMatchDisconnect(ctx context.Context, groupId, openId string) bool {
	res, err := etcdClient.Client.Get(ctx, path.Join("/", projectName, matcd_battle_disconnect, groupId, openId), clientv3.WithCountOnly()) // 注册到匹配组
	if err != nil {
		return true
	}
	return res.Count != 0
}

// 删除匹配组掉线注册
func deleteMatchDisconnectRegister(ctx context.Context, groupId string) {
	r, _ := etcdClient.Client.Get(ctx, path.Join("/", projectName, matcd_battle_disconnect, groupId), clientv3.WithCountOnly())
	if r.Count == 0 { // 没有用户掉线注册
		return
	}
	etcdClient.Client.Delete(ctx, path.Join("/", projectName, matcd_battle_disconnect, groupId), clientv3.WithPrefix()) // 删除匹配组掉线注册
}

// 删除匹配组掉线注册用户
func DeleteMatchDisconnectRegisterUser(ctx context.Context, openId string) error {
	groupId, err := queryVs1GroupByOpenId(ctx, openId)
	if err != nil {
		return errors.New("DisconnectRegister 查询组信息失败: " + err.Error() + ", openId: " + openId)
	}
	if groupId == "" {
		return nil
	}
	_, err = etcdClient.Client.Delete(ctx, path.Join("/", projectName, matcd_battle_disconnect, groupId, openId)) // 删除匹配组掉线注册用户
	return err
}
