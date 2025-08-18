package battlematch

import (
	"context"
	"errors"
	"fmt"
	"path"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 设置是否匿名，匿名则为1，否则则为0
func MatchBattleAnonymousSet(ctx context.Context, openId string) error {
	_, err := etcdClient.Client.Put(ctx, path.Join("/", projectName, match_battle_anonymous_status, openId), openId, clientv3.WithLease(etcdClient.NewLease(ctx, 30)))
	if err != nil {
		return errors.New("MatchBattleAnonymousSet err: " + err.Error())
	}
	return nil
}

// 查询是否匿名
func MatchBattleAnonymousGet(ctx context.Context, openId string) (bool, error) {
	res, err := etcdClient.Client.Get(ctx, path.Join("/", projectName, match_battle_anonymous_status, openId), clientv3.WithCountOnly())
	if err != nil {
		return false, errors.New("MatchBattleAnonymousSet err: " + err.Error())
	}
	return res.Count >= 1, nil
}

// 查询是否匿名
func MatchBattleAnonymousGetByGroupId(ctx context.Context, groupId string) (bool, error) {
	openIdList, _ := QueryVs1GroupInfo(ctx, groupId)
	for _, openId := range openIdList {
		res, err := etcdClient.Client.Get(ctx, path.Join("/", projectName, match_battle_anonymous_status, openId), clientv3.WithCountOnly())
		if err != nil {
			continue
		}
		if res.Count >= 1 {
			return true, nil
		}
	}
	return false, nil
}

// 匿名删除通过openId
func MatchBattleAnonymousDelByOpenId(ctx context.Context, openId string) error {
	_, err := etcdClient.Client.Delete(ctx, path.Join("/", projectName, match_battle_anonymous_status, openId), clientv3.WithCountOnly())
	if err != nil {
		return errors.New("MatchBattleAnonymousDelByOpenId err: " + err.Error())
	}
	return nil
}

// 删除匿名通过组
func matchBattleAnonymousDelByGroupId(ctx context.Context, groupId string) error {
	var err error
	openIdList, _ := QueryVs1GroupInfo(ctx, groupId)
	for _, openId := range openIdList {
		ok, _ := MatchBattleAnonymousGet(ctx, openId)
		if ok {
			if mErr := MatchBattleAnonymousDelByOpenId(ctx, openId); mErr != nil {
				if err != nil {
					err = fmt.Errorf("%v,  %v", err, mErr)
				} else {
					err = fmt.Errorf("matchBattleAnonymousDelByGroupId  err: %v", mErr)
				}
			}
		}
	}
	return err
}
