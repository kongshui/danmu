package battlematch

import (
	"context"
	"errors"
	"path"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// registerChat 注册用户
func registerBattle(ctx context.Context, matchNum string, openId string) error {
	key := path.Join(matcd_battle_register_v1, openId)
	if matchNum != "" {
		key = path.Join("/", projectName, matcd_battle_num_register_v1, matchNum, openId)
	}
	id := etcdClient.NewLease(ctx, match_register_timeout) // 创建租约
	_, err := etcdClient.Client.Put(ctx, key, openId, clientv3.WithLease(id))
	if err != nil {
		return errors.New("registerChat 注册用户失败: " + err.Error() + ", openId: " + openId)
	}
	return nil
}

// UnregisterChat 注销用户
func UnregisterBattleByOpenId(ctx context.Context, matchNum string, openId string) error {
	// 查询组信息
	key := path.Join("/", projectName, matcd_battle_register_v1, openId)
	if matchNum != "" {
		key = path.Join("/", projectName, matcd_battle_num_register_v1, matchNum, openId)
	}
	_, err := etcdClient.Client.Delete(ctx, key)
	if err != nil {
		return errors.New("UnregisterChatByOpenId 注销用户失败: " + err.Error() + ", openId: " + openId)
	}
	return nil
}

// 按数量获取用户列表
func getUserListByNum(ctx context.Context, matchNum string, num int64) ([]string, error) {
	key := matcd_battle_register_v1
	if matchNum != "" {
		key = path.Join("/", projectName, matcd_battle_num_register_v1, matchNum)
	}
	// 获取所有用户的键值对
	resp, err := etcdClient.Client.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithLimit(num))
	if err != nil {
		return nil, errors.New("GetUserListByNum 获取用户列表失败: " + err.Error())
	}
	if resp.Count == 0 { // 没有用户掉线注册
		return nil, errors.New("GetUserListByNum 获取用户列表失败, 长度为0 ")
	}
	// 提取用户列表
	var userList []string
	for _, kv := range resp.Kvs {
		userList = append(userList, string(kv.Value))
	}
	if len(userList) < int(num) {
		return nil, errors.New("GetUserListByNum 获取用户列表失败: 用户数量不足")
	}
	// 随机选择指定数量的用户
	return userList, nil
}

// 获取注册用户总长度
func getResiterUserLen(ctx context.Context, matchNum string) (int64, error) {
	// 获取所有用户的键值对
	key := matcd_battle_register_v1
	if matchNum != "" {
		key = path.Join("/", projectName, matcd_battle_num_register_v1, matchNum)
	}
	resp, err := etcdClient.Client.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithCountOnly())
	if err != nil {
		return 0, errors.New("GetUserListLen 获取用户列表失败: " + err.Error())
	}
	return resp.Count, nil
}
