package battlematch

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"path"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 匹配信息
func MatchV1Battle(ctx context.Context, openId, matchNum string) (userId []string, groupId string, err error) {
	t := time.NewTicker(time.Second * 3)
	defer t.Stop()
	count := 0
	for {
		<-t.C
		if isInVs1Cancel(ctx, openId) { // 查看是否在取消匹配组中
			if count == 0 {
				removeFromCancelMatchV1Battle(ctx, openId)
			} else {
				return []string{}, "", errors.New("用户已取消匹配")
			}
		}
		if IsInResiter(ctx, matchNum, openId) {
			return []string{}, "", registerBattle(ctx, matchNum, openId)
		}

		// 创建租约
		id := etcdClient.NewLease(ctx, 10)
		//检查锁                       // 创建租约
		if etcdClient.PutIfNotExist(ctx, path.Join("/", projectName, match_battle_v1_lock), openId, id) {
			defer etcdClient.Client.Delete(ctx, path.Join("/", projectName, match_battle_v1_lock)) // 删除锁
			// 查询注册用户多少
			rLen, err := getResiterUserLen(ctx, matchNum)
			if err != nil {
				return []string{}, "", registerBattle(ctx, matchNum, openId) // 注册用户
			}
			if rLen == 0 {
				return []string{}, "", registerBattle(ctx, matchNum, openId) // 注册用户
			}
			// 匹配用户
			userList, err := getUserListByNum(ctx, matchNum, 1) // 获取注册用户列表
			if err != nil {
				return []string{}, "", registerBattle(ctx, matchNum, openId) // 注册用户
			}
			UnregisterBattleByOpenId(ctx, matchNum, userList[0])
			// 如果在取消匹配组中，继续匹配
			if isInVs1Cancel(ctx, userList[0]) {
				continue
			}
			h := sha256.New()
			_, err = h.Write([]byte(openId + userList[0] + time.Now().Format("2006-01-02 15:04:05")))
			if err != nil {
				return []string{}, "", registerBattle(ctx, matchNum, openId) // 注册用户
			}
			groupId := fmt.Sprintf("%x", h.Sum(nil))
			choose := rand.Intn(2) // 生成匹配组ID
			var userChooseList []string = []string{openId, userList[0]}
			if choose == 1 {
				userChooseList = []string{userList[0], openId}
			}
			if err := registerMatchBattle(ctx, userChooseList, groupId); err != nil { // 注册匹配组
				return []string{}, "", UnregisterBattleV1ByGroupId(ctx, groupId) // 注册用户
			}
			return userChooseList, groupId, nil
		}
		count++
		if count > 40 { // 匹配超时
			return []string{}, "", errors.New("匹配超时")
		}
	}

}

// 注册到匹配组
func registerMatchBattle(ctx context.Context, openIdList []string, groupId string) error {
	dataByte, _ := json.Marshal(openIdList)
	// 创建租约
	id := etcdClient.NewLease(ctx, match_success_timeout) // 创建租约
	_, err := etcdClient.Client.Put(ctx, path.Join("/", projectName, matcd_battle_store_v1, groupId), string(dataByte), clientv3.WithLease(id))
	if err != nil {
		return errors.New("RegisterBattle 注册组失败: " + err.Error() + ", openIdList: " + string(dataByte) + ", groupId: " + groupId)
	}
	for _, openId := range openIdList {
		_, err := etcdClient.Client.Put(ctx, path.Join("/", projectName, matcd_battle_store_v1, openId), groupId, clientv3.WithLease(id))
		if err != nil {
			if err := UnregisterBattleV1ByGroupId(ctx, groupId); err != nil {
				return errors.New("RegisterBattle 注册用户失败,unregisterBattleV1ByGroupId: " + err.Error() + ", openIdList: " + string(dataByte) + ", openId: " + openId + ", groupId: " + groupId)
			}
			return errors.New("RegisterBattle 注册用户失败: " + err.Error() + ", openIdList: " + string(dataByte) + ", openId: " + openId + ", groupId: " + groupId)
		}
	}
	return nil
}

// UnregisterBattle 通过groupId注销匹配组
func UnregisterBattleV1ByGroupId(ctx context.Context, groupId string) error {
	// 注销前的善后
	matchBattleAnonymousDelByGroupId(ctx, groupId)
	// MatchGroupStatusDel(ctx, groupId)
	deleteMatchDisconnectRegister(ctx, groupId)
	matchV1GroupRoundIdDel(ctx, groupId)
	matchBattleTimeDel(ctx, groupId)
	// 查询组信息
	openIdList, err := QueryVs1GroupInfo(ctx, groupId)
	if err != nil {
		return errors.New("unregisterBattleV1ByGroupId 查询组信息失败: " + err.Error() + ", groupId: " + groupId)
	}
	// 注销用户
	for _, openId := range openIdList {

		ok, checkGroupId := IsInVs1Group(ctx, openId)
		if ok && checkGroupId == groupId {
			etcdClient.Client.Delete(ctx, path.Join("/", projectName, matcd_battle_store_v1, openId))
		}
	}
	etcdClient.Client.Delete(ctx, path.Join("/", projectName, matcd_battle_store_v1, groupId))

	return nil
}

// 通过openId注销匹配组
func UnregisterBattleV1ByOpenId(ctx context.Context, openId string) error {
	groupId, err := queryVs1GroupByOpenId(ctx, openId)
	if err != nil {
		return errors.New("UnregisterBattleV1ByOpenId 查询组信息失败: " + err.Error() + ", openId: " + openId)
	}
	if groupId == "" {
		return nil
	}
	if err := UnregisterBattleV1ByGroupId(ctx, groupId); err != nil {
		return errors.New("UnregisterBattleV1ByOpenId 注销组信息失败: " + err.Error() + ", openId: " + openId)
	}
	return nil
}

// 取消匹配
func CancelMatchV1Battle(ctx context.Context, openId string) error {
	if err := UnregisterBattleV1ByOpenId(ctx, openId); err != nil { // 注销用户
		return errors.New("CancelMatchV1Battle 注销用户失败: " + err.Error() + ", openId: " + openId)
	}
	return nil
}

// 添加到配置取消组
func AddToCancelMatchV1Battle(ctx context.Context, matchNum, openId string) error {
	// 创建租约
	id := etcdClient.NewLease(ctx, 180)
	_, err := etcdClient.Client.Put(ctx, path.Join("/", projectName, match_battle_cancel_v1, openId), openId, clientv3.WithLease(id)) // 注册到匹配组
	if err != nil {
		return errors.New("AddToCancelMatchV1Battle 注册用户失败: " + err.Error() + ", openId: " + openId)
	}
	if err := UnregisterBattleByOpenId(ctx, matchNum, openId); err != nil {
		return errors.New("CancelMatchV1Battle 注销用户失败: " + err.Error() + ", openId: " + openId)
	}
	ok, _ := IsInVs1Group(ctx, openId)
	if ok {
		return errors.New("取消失败，已经匹配成功")
	}
	return nil
}

// removeFromCancelMatchV1Battle 从取消匹配组中移除
func removeFromCancelMatchV1Battle(ctx context.Context, openId string) error {
	_, err := etcdClient.Client.Delete(ctx, path.Join("/", projectName, match_battle_cancel_v1, openId)) // 注销用户
	if err != nil {
		return errors.New("RemoveFromCancelMatchV1Battle 注销用户失败: " + err.Error() + ", openId: " + openId)
	}
	return nil
}
