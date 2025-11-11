package service

import (
	"fmt"
	"path"
	"strconv"

	"github.com/kongshui/danmu/common"
)

// 设置roomid和主播id
func setRoomIdToAnchorOpenId(roomId string, openId string) bool {
	_, err := rdb.HSetNX(roomdid_to_anchoropenid_db, roomId, openId)
	return err == nil
}

// 删除roomid和主播id
// 设置roomid和主播id
func delRoomIdToAnchorOpenId(roomId string) bool {
	return rdb.HDel(roomdid_to_anchoropenid_db, roomId) == nil
}

func queryAnchorOpenIdByRoomId(roomId string) (string, error) {
	openId, err := rdb.HGet(roomdid_to_anchoropenid_db, roomId)
	if err != nil {
		return "", err
	}
	return openId, nil
}

// 通过roomid获取主播id
func QueryRoomIdInterconvertAnchorOpenId(roomId string) string {
	res, err := etcdClient.Client.Get(first_ctx, path.Join("/", cfg.Project, common.RoomId_OpenId_Register_key, roomId))
	if err != nil {
		ziLog.Error(fmt.Sprintf("queryRoomIdToAnchorOpenId,查询%v失败: %v", roomId, err), debug)
		return ""
	}
	if res.Count == 0 { // 没有用户掉线注册
		return ""
	}
	return string(res.Kvs[0].Value)
}

// 通过roomid获取rounid
func queryRoomIdToRoundId(roomId string) (int64, bool) {
	roundId, err := rdb.Get(roomId + "_round")
	if err != nil {
		return 0, false
	}
	roundIdInt, _ := strconv.ParseInt(roundId, 10, 64)
	return roundIdInt, len(roundId) != 0
}

// 通过roomid获取到uid
func queryRoomIdToUid(roomId string) string {
	// 获取房间信息
	result, err := etcdClient.Client.Get(first_ctx, path.Join("/", cfg.Project, common.RoomId_Register_Uid_key, roomId))
	if err != nil {
		ziLog.Error(fmt.Sprintf("queryRoomIdToUid etcdClient.Client.Get err: %v", err), debug)
		return ""
	}
	if len(result.Kvs) == 0 {
		ziLog.Error(fmt.Sprintf("queryRoomIdToUid  etcdClient.Client.Get %v nil", roomId), debug)
		return ""
	}
	return string(result.Kvs[0].Value)
}
