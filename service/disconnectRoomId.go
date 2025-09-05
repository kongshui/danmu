package service

import (
	"errors"
	"fmt"
	"time"

	battlematchv1 "github.com/kongshui/danmu/battlematch/v1"
)

// 将roomid添加到断线列表中
func disconnectRoomIdADD(roomId string) error {
	_, err := rdb.SAdd(disconnect_roomid_db, roomId)
	if err != nil {
		return errors.New("DisconnectRoomId 添加roomId到断线重连列表失败: " + err.Error() + "roomId: " + roomId)
	}
	return nil
}

// 将roomid从断线列表中删除
func disconnectRoomIdDelete(roomId string) error {
	err := rdb.SRem(disconnect_roomid_db, roomId)
	if err != nil {
		return errors.New("ConnectRoomId 从断线重连列表中删除roomId失败: " + err.Error() + "roomId: " + roomId)
	}
	return nil
}

// 查询断线列表中的roomid
func queryDisconnectRoomId() ([]string, error) {
	roomIdList, err := rdb.SMembers(disconnect_roomid_db)
	if err != nil {
		return nil, errors.New("QueryDisconnectRoomId 查询断线重连列表中的roomid失败: " + err.Error())
	}
	return roomIdList, nil
}

// 定期检查断线的roomid是否过期，如果过期则删除
func checkDisconnectRoomIdExpire() {
	t := time.NewTicker(1 * time.Minute) // 每分钟检查一次
	defer t.Stop()
	for {
		<-t.C
		ok, _ := rdb.SetKeyNX(monitor_disconnect_roomid_db, "1", 55*time.Second)
		if !ok {
			continue // 如果已经有监控在运行，则跳过
		}
		roomIdList, err := queryDisconnectRoomId()
		if err != nil {
			ziLog.Error(fmt.Sprintf("checkDisconnectRoomIdExpire 查询断线重连列表中的roomid失败: %v", err), debug)
			continue
		}

		for _, roomId := range roomIdList {
			openId, _ := queryAnchorOpenIdByRoomId(roomId)
			ziLog.Info(fmt.Sprintf("checkDisconnectRoomIdExpire 检查roomId: %v, openId: %v", roomId, openId), debug)
			// 检查roomid是否过期，如果过期则删除
			if rdb.IsExistKey(roomId + "_round") {
				switch platform {
				case "ks":
					if getKsGameInfo(roomId, url_BindUrl) == 1 {
						// 游戏中，不做处理
						continue
					} else {
						if err := battlematchv1.DisconnectMatchRegister(first_ctx, openId); err != nil {
							ziLog.Error(fmt.Sprintf("checkDisconnectRoomIdExpire 匹配组掉线注册失败, roomId : %v,err: %v", roomId, err), debug)
						}
					}
					// 获取roundId
					roundId, ok := queryRoomIdToRoundId(roomId)
					if ok {
						if err := ksSyncGameStatus(SyncGameStatusStruct{
							AnchorOpenId:    openId,
							AppId:           app_id,
							RoomId:          roomId,
							RoundId:         roundId,
							StartTime:       time.Now().UnixMilli() - 1000,
							EndTime:         time.Now().UnixMilli(),
							Status:          2,
							GroupResultList: []GroupResultList{{GroupId: groupid_list[0], Result: 1}},
						}, "stop", false); err != nil {
							ziLog.Error(fmt.Sprintf("checkDisconnectRoomIdExpire 同步对局状态失败, roomId : %v,err: %v", roomId, err), debug)
							endClean(roomId, openId)
						}
					}
				case "dy":
					// 抖音 1 任务不存在 2任务未启动 3任务运行中
					if getDyGameInfo(roomId, url_check_push_url, "live_gift") == 1 &&
						getDyGameInfo(roomId, url_check_push_url, "live_comment") == 1 &&
						getDyGameInfo(roomId, url_check_push_url, "live_like") == 1 {
						endClean(roomId, openId)
					}
				}
			} else {
				// 否则的话清理
				endClean(roomId, openId)
			}
		}
	}
}
