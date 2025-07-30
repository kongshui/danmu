package service

import (
	"fmt"
	"time"

	battlematchv1 "github.com/kongshui/danmu/battlematch/v1"
)

func endConnect(roomId string, openId string) bool {
	// 是否存在roundId
	if rdb.IsExistKey(roomId + "_round") {
		// 获取roundId
		roundId, ok := queryRoomIdToRoundId(roomId)
		if ok {
			switch platform {
			case "ks":
				// 添加到断线列表中，如果有问题后续执行end，思路暂时放这里
				if getKsGameInfo(roomId, url_BindUrl) == 1 {
					disconnectRoomIdADD(roomId)
					return true
				} else {
					if err := battlematchv1.DisconnectMatchRegister(first_ctx, openId); err != nil {
						ziLog.Error(fmt.Sprintf("endConnect 匹配组掉线注册失败, roomId : %v,err: %v", roomId, err), debug)
					}
				}
				if !ksSyncGameStatus(SyncGameStatusStruct{
					AnchorOpenId:    openId,
					AppId:           app_id,
					RoomId:          roomId,
					RoundId:         roundId,
					StartTime:       time.Now().UnixMilli() - 1000,
					EndTime:         time.Now().UnixMilli(),
					Status:          2,
					GroupResultList: []GroupResultList{{GroupId: "绝对中立", Result: 1}},
				}, "stop", false) {
					endClean(roomId, openId)
				}
			default:
				// 抖音 1 任务不存在 2任务未启动 3任务运行中
				if getDyGameInfo(roomId, url_check_push_url, "live_gift") == 1 &&
					getDyGameInfo(roomId, url_check_push_url, "live_comment") == 1 &&
					getDyGameInfo(roomId, url_check_push_url, "live_like") == 1 {
					endClean(roomId, QueryRoomIdInterconvertAnchorOpenId(roomId))
				} else {
					disconnectRoomIdADD(roomId)
				}
			}
		}
	} else {
		// 否则的话清理
		endClean(roomId, openId)
	}
	return true
}

// Connect 暂时没被使用
func connect(roomId string, openId string) bool {
	//清空积分池
	// if err := rdb.HSet(integral_pool_db, roomId, 0); err != nil {
	// 	log.Printf("清空积分池失败, roomId : %v,err: %v", roomId, err)
	// }
	ok, err := rdb.SetKeyNX(integral_pool_Prefix+openId, 0, 0)
	if err != nil {
		ziLog.Error(fmt.Sprintf("connect 设置积分池失败, roomId : %v,err: %v", roomId, err), debug)
	}
	// 设置过期时间
	if ok {
		expireTime := getNextExpireTime()
		rdb.Expire(integral_pool_Prefix+openId, expireTime)
	} else {
		if ttl, _ := rdb.TTL(integral_pool_Prefix + openId); ttl < 0 {
			expireTime := getNextExpireTime()
			rdb.Expire(integral_pool_Prefix+openId, expireTime)
		}
	}
	if _, err := rdb.SAdd(room_id_list_db, roomId); err != nil {
		ziLog.Error(fmt.Sprintf("connect 添加房间id失败, roomId : %v,err: %v", roomId, err), debug)
		return false
	}
	ziLog.Info(fmt.Sprintf("connect 添加房间id成功, roomId : %v", roomId), debug)
	return true
}

// 清除过期房间
func endClean(roomId, openId string) {
	// 查看是否有匹配
	if openId != "" {
		ok, _ := battlematchv1.IsInVs1Group(first_ctx, openId)
		if ok {
			battlematchv1.DisconnectMatchRegister(first_ctx, openId)
		}
	}

	//删除对战信息
	if err := liveCurrentRoundDel(roomId); err != nil {
		ziLog.Error(fmt.Sprintf("endClean 删除对战信息, roomId : %v,err: %v", roomId, err), debug)
	}
	//删除房间id
	// DelRoomIdToAnchorOpenId(roomId)

	// 删除礼物对比消息
	rdb.Del(roomId + "giftSend")

	if err := rdb.SRem(room_id_list_db, roomId); err != nil {
		ziLog.Error(fmt.Sprintf("endClean 删除房间id失败, roomId : %v,err: %v", roomId, err), debug)
	}
	disconnectRoomIdDelete(roomId)
	ziLog.Info(fmt.Sprintf("endClean 删除房间成功, roomId : %v", roomId), debug)
}
